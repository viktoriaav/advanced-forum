package forum

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	Type    string `json:"type"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Type    string      `json:"type"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterHandler handles user registration over WebSocket
func RegisterHandler(conn *websocket.Conn, r *http.Request, db *sql.DB, message map[string]interface{}) {
	log.Println("RegisterHandler called.")

	// Extract registration data from the message map
	email, ok := message["email"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid email format"})
		return
	}

	firstName, ok := message["first-name"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid first name format"})
		return
	}

	lastName, ok := message["last-name"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid last name format"})
		return
	}

	username, ok := message["username"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid username format"})
		return
	}

	password, ok := message["password"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid password format"})
		return
	}

	age, ok := message["age"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid age format"})
		return
	}

	gender, ok := message["gender"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid gender format"})
		return
	}

	// Convert email and username to lowercase
	lowercaseEmail := strings.ToLower(email)
	lowercaseUsername := strings.ToLower(username)

	// Check if the user already exists in the database
	var existingUser int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE LOWER(email) = ? OR LOWER(username) = ?", lowercaseEmail, lowercaseUsername).Scan(&existingUser)
	log.Printf("Query: SELECT COUNT(*) FROM users WHERE LOWER(email) = %s OR LOWER(username) = %s\n", lowercaseEmail, lowercaseUsername)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Database error"})
		log.Println("Database error:", err)
		return
	}
	if existingUser > 0 {
		log.Print("user already exists")
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "User already exists"})
		return
	}

	// Registration logic
	createdAt := time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Password hashing error"})
		log.Println("Password hashing error:", err)
		return
	}

	// Continue with user registration
	query := "INSERT INTO users (email, first_name, last_name, username, password, age, gender, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(query, lowercaseEmail, firstName, lastName, lowercaseUsername, hashedPassword, age, gender, createdAt)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Database error"})
		log.Println("Database error:", err)
		return
	}

	// Generate session token for the newly registered user
	token := GenerateSessionToken()

	// Calculate session duration
	sessionDuration := time.Minute * 15

	// Calculate expiration time
	expirationTime := time.Now().Add(sessionDuration)

	// Get the user ID of the newly registered user
	userID, err := getUserID(lowercaseUsername, db)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Failed to get user ID"})
		log.Println("Failed to get user ID:", err)
		return
	}

	// Create a new session record in the database
	err = createSession(conn, userID, token, expirationTime, db)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Database error"})
		log.Println("Database error:", err)
		return
	}

	loggedInUsername, err := GetLoggedInUsername(token, db)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Failed to get logged-in username"})
		return
	}
	isAuthenticated := true

	// Fetch all messages associated with the logged-in user
	allMessages, err := GetAllMessages(db)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Failed to fetch messages"})
		log.Println("Failed to fetch messages:", err)
		return
	}

	// Prepare data to be sent over WebSocket
	responseData := map[string]interface{}{
		"loggedInUsername": loggedInUsername,
		"isAuthenticated":  isAuthenticated,
		"allMessages":      allMessages,
	}

	// Send data over WebSocket
	SendWebSocketMessageSuccess(conn, SuccessResponse{Type: "Registration", Success: true, Message: "Registration successful", Data: responseData})

}

// LoginHandler handles user login over WebSocket
func LoginHandler(conn *websocket.Conn, r *http.Request, db *sql.DB, message map[string]interface{}) {
	log.Println("LoginHandler called.")
	// Ensure that the WebSocket message contains the necessary fields
	identifier, ok := message["identifier"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid message format"})
		return
	}

	password, ok := message["password"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid message format"})
		return
	}

	lowercaseIdentifier := strings.ToLower(identifier)

	var existingUser int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE LOWER(email) = ? OR LOWER(username) = ?", lowercaseIdentifier, lowercaseIdentifier).Scan(&existingUser)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Database error"})
		log.Println("Database error:", err)
		return
	}

	// Check for an existing user
	if existingUser == 0 {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "User not found"})
		return
	}

	// Use the provided identifier to retrieve user from the database
	var user User
	query := "SELECT user_ID, email, username, password FROM users WHERE LOWER(email) = ? OR LOWER(username) = ?"
	err = db.QueryRow(query, lowercaseIdentifier, lowercaseIdentifier).Scan(&user.ID, &user.Email, &user.Username, &user.Password)

	if err == sql.ErrNoRows {
		log.Println("User not found in the database")
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "User not found"})
		return
	} else if err != nil {
		log.Println("Database error while retrieving user:", err)
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Database error"})
		return
	}

	// Check the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid password"})
		return
	}
	token := GenerateSessionToken()

	// Calculate session duration
	sessionDuration := time.Minute * 15

	// Calculate expiration time
	expirationTime := time.Now().Add(sessionDuration)

	// Create a new session record in the database
	err = createSession(conn, user.ID, token, expirationTime, db)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Database error"})
		log.Println("Database error:", err)
		return
	}

	loggedInUsername, err := GetLoggedInUsername(token, db)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Failed to get logged-in username"})
		return
	}
	isAuthenticated := true

	// Fetch all messages associated with the logged-in user
	allMessages, err := GetAllMessages(db)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Failed to fetch messages"})
		log.Println("Failed to fetch messages:", err)
		return
	}

	// Prepare data to be sent over WebSocket
	responseData := map[string]interface{}{
		"loggedInUsername": loggedInUsername,
		"isAuthenticated":  isAuthenticated,
		"allMessages":      allMessages,
	}

	// Send data over WebSocket
	SendWebSocketMessageSuccess(conn, SuccessResponse{Type: "Login", Success: true, Message: "Login successful", Data: responseData})

}

// SendWebSocketMessage sends a JSON-encoded message to the WebSocket client
func SendWebSocketMessage(conn *websocket.Conn, response Response) {
	err := conn.WriteJSON(response)
	if err != nil {
		log.Println("Error writing JSON to WebSocket connection:", err)
	}
}

func SendWebSocketMessageSuccess(conn *websocket.Conn, response SuccessResponse) {
	err := conn.WriteJSON(response)
	if err != nil {
		log.Println("Error writing JSON to WebSocket connection:", err)
	}
}

func GenerateSessionToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func createSession(conn *websocket.Conn, userID int, token string, expirationTime time.Time, db *sql.DB) error {
	log.Print("Trying to create session")

	// Delete expired sessions before creating a new one
	err := DeleteExpiredSessions(db)
	if err != nil {
		log.Println("Error deleting expired sessions:", err)
		return err
	}

	// Check if an active session already exists for the user
	var existingSessionID int
	query := "SELECT session_ID FROM sessions WHERE user_ID = ? AND expires_at > ?"
	err = db.QueryRow(query, userID, time.Now()).Scan(&existingSessionID)

	if err == sql.ErrNoRows {
		// No active session found, create a new one
		insertQuery := "INSERT INTO sessions (token, user_ID, created_at, expires_at) VALUES (?, ?, ?, ?)"
		_, err = db.Exec(insertQuery, token, userID, time.Now(), expirationTime)

		if err != nil {
			log.Println("Error creating a new session:", err)
			return err
		}
		log.Println("New session created successfully")
	} else if err != nil {
		log.Println("Error checking for an existing session:", err)
		return err
	} else {
		// Update the existing session with new token and expiration time
		updateQuery := "UPDATE sessions SET token = ?, expires_at = ? WHERE session_ID = ?"
		_, err = db.Exec(updateQuery, token, expirationTime, existingSessionID)

		if err != nil {
			log.Println("Error updating an existing session:", err)
			return err
		}

		log.Println("Session updated successfully")
	}
	return nil
}

func DeleteExpiredSessions(db *sql.DB) error {
	deleteQuery := "DELETE FROM sessions WHERE expires_at <= ?"
	_, err := db.Exec(deleteQuery, time.Now())
	return err
}

func NotifyAllUsersOnlineStatus(conn *websocket.Conn, userID int, db *sql.DB) error {
	log.Print("Notifying all users about online status change")

	// Fetch online users
	onlineUsers, err := GetAllOnlineUsers(db)
	if err != nil {
		log.Println("Error fetching online users:", err)
		return err
	}

	// Send the online users information to the WebSocket
	responseData := struct {
		UsersOnlineStatus []OnlineUser `json:"usersOnlineStatus"`
	}{
		UsersOnlineStatus: onlineUsers,
	}

	SendWebSocketMessageSuccess(conn, SuccessResponse{
		Type:    "onlineStatusChange",
		Success: true,
		Message: "Update Users Online Status",
		Data:    responseData,
	})
	BroadcastChanges(conn, "updatAllUsersOnline", responseData)

	return nil
}

// CreatePostHandler handles post creation over WebSocket
func CreatePostHandler(conn *websocket.Conn, r *http.Request, db *sql.DB, message map[string]interface{}) {
	log.Println("CreatePostHandler called.")

	// Ensure that the WebSocket message contains the necessary fields
	createdBy, ok := message["createdBy"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid message format lol"})
		return
	}

	title, ok := message["title"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid message format kek"})
		return
	}

	content, ok := message["content"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid message format haha"})
		return
	}

	// Type assert the field to []interface{}
	categoriesInterface, ok := message["categories"].([]interface{})
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid type for categories field"})
		return
	}

	// Convert each interface{} element to string
	var categories []string
	for _, cat := range categoriesInterface {
		category, ok := cat.(string)
		if !ok {
			SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid type for category element"})
			return
		}
		categories = append(categories, category)
	}

	// Insert the new post into the database
	userID, err := getUserID(createdBy, db)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Database error"})
		log.Println("Database error:", err)
		return
	}
	createdAt := time.Now()
	// Insert the new post into the database
	err = createPost(userID, title, content, categories, createdAt, db)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Database error"})
		log.Println("Database error:", err)
		return
	}

	// Prepare data to be sent over WebSocket
	allPosts, err := GetAllPosts(db)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Failed to get all posts"})
		return
	}
	// Send the data back to the frontend
	responseData := struct {
		AllPosts []Post `json:"allPosts"`
	}{
		AllPosts: allPosts,
	}

	SendWebSocketMessageSuccess(conn, SuccessResponse{Type: "createdPost", Success: true, Message: "Update Posts Data", Data: responseData})
	BroadcastChanges(conn, "updateAllPosts", responseData)
}

// Function to insert a new post into the database
func createPost(userID int, title, content string, categories []string, createdAt time.Time, db *sql.DB) error {
	// Insert the post into the posts table
	_, err := db.Exec("INSERT INTO posts (user_ID, title, content, created_at) VALUES (?, ?, ?, ?)", userID, title, content, createdAt)
	if err != nil {
		return err
	}

	// Associate the post with categories in the post_categories table
	for _, category := range categories {
		_, err := db.Exec("INSERT INTO post_categories (post_ID, category_ID) VALUES ((SELECT post_ID FROM posts WHERE user_ID = ? AND title = ? AND content = ? AND created_at = ?), (SELECT category_ID FROM categories WHERE category = ?))", userID, title, content, createdAt, category)
		if err != nil {
			return err
		}
	}

	return nil
}

// SubmitCommentHandler handles comment submission over WebSocket
func SubmitCommentHandler(conn *websocket.Conn, r *http.Request, db *sql.DB, message map[string]interface{}) {
	log.Println("Submit Comment Handler called.")

	username, ok := message["username"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid message format lol"})
		return
	}

	comment, ok := message["comment"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid message format um"})
		return
	}

	postIDStr, ok := message["postID"].(string)
	if !ok {
		// Handle the case where postID is not a string
		log.Println("postID is not a string")
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid message format ok"})
		log.Println("Invalid postID:", postIDStr)
		return
	}

	userID, err := getUserID(username, db)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid message format sure"})
		log.Println("Database error:", err)
		return
	}

	_, err = db.Exec("INSERT INTO comments (post_ID, user_ID, content, created_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP)", postID, userID, comment)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid message format lol"})
		log.Println("Database error:", err)
		return
	}

	// Prepare data to be sent over WebSocket
	allComments, err := GetAllComments(db)
	if err != nil {
		SendWebSocketMessage(conn, Response{Success: false, Message: "Failed to get all posts"})
		return
	}
	// Send the data back to the frontend
	responseData := struct {
		AllComments []Comment `json:"allComments"`
	}{
		AllComments: allComments,
	}

	SendWebSocketMessageSuccess(conn, SuccessResponse{Type: "newComment", Success: true, Message: "Update Posts Data", Data: responseData})
	BroadcastChanges(conn, "updateAllComments", responseData)
}

// DeleteSessionHandler handles session deletion over WebSocket
func DeleteSessionHandler(conn *websocket.Conn, r *http.Request, db *sql.DB, message map[string]interface{}) {
	log.Println("Delete Session Handler called.")

	username, ok := message["username"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid message format"})
		return
	}

	userID, err := getUserID(username, db)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid message format"})
		log.Println("Database error:", err)
		return
	}

	// Delete the session for the given user
	_, err = db.Exec("DELETE FROM sessions WHERE user_ID = ?", userID)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Failed to delete session"})
		log.Println("Database error:", err)
		return
	}
	usersOnline, err := GetAllOnlineUsers(db)
	if err != nil {
		SendWebSocketMessage(conn, Response{Success: false, Message: "Failed to get all online users"})
		return
	}
	// Send the data back to the frontend
	responseData := struct {
		AllUsersOnline []OnlineUser `json:"usersOnline"`
	}{
		AllUsersOnline: usersOnline,
	}

	// Send a success message back to the frontend
	SendWebSocketMessageSuccess(conn, SuccessResponse{Type: "userLogout", Success: true, Message: "User logout successful"})
	BroadcastChanges(conn, "updatAllUsersOnline", responseData)
}

// SubmitMessageHandler handles message submission over WebSocket
func SubmitMessageHandler(conn *websocket.Conn, r *http.Request, db *sql.DB, message map[string]interface{}) {
	log.Println("Submit Message Handler called.")

	sender, ok := message["sender"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid message format: sender"})
		return
	}

	receiver, ok := message["receiver"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid message format: receiver"})
		return
	}

	content, ok := message["content"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid message format: content"})
		return
	}

	created_at, ok := message["created_at"].(string)
	if !ok {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Invalid message format: created_at"})
		return
	}

	_, err := db.Exec("INSERT INTO private_messages (sender, receiver, content, created_at) VALUES (?, ?, ?, ?)", sender, receiver, content, created_at)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Failed to insert message into the database"})
		log.Println("Database error:", err)
		return
	}

	// Fetch all messages associated with the sender
	allMessages, err := GetAllMessages(db)
	if err != nil {
		SendWebSocketMessage(conn, Response{Type: "Error", Success: false, Message: "Failed to fetch messages"})
		log.Println("Failed to fetch messages:", err)
		return
	}

	// Prepare data to be sent over WebSocket
	responseData := map[string]interface{}{
		"allMessages": allMessages,
	}

	// Send data over WebSocket
	SendWebSocketMessageSuccess(conn, SuccessResponse{Type: "newMessageAdd", Success: true, Message: "Message sent successfully", Data: responseData})
	BroadcastChanges(conn, "updateAllMessages", responseData)
}
