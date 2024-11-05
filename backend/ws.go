package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Declare a mutex to safely manage the list of connected clients
var clientsMutex sync.Mutex
var clients = make(map[*websocket.Conn]struct{})

func HandleWebSocket(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		clientsMutex.Lock()
		delete(clients, conn)
		clientsMutex.Unlock()
		conn.Close()
	}()

	// Add the new connection to the list
	clientsMutex.Lock()
	clients[conn] = struct{}{}
	clientsMutex.Unlock()

	// Handle WebSocket messages
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("Received message: %s\n", p)
		// Handle different message types (text, binary, etc.) if needed
		if messageType == websocket.TextMessage {
			// Parse and handle the message as needed
			var message map[string]interface{}
			if err := json.Unmarshal(p, &message); err != nil {
				log.Println(err)
				continue
			}

			messageType, ok := message["message"].(string)
			log.Print(messageType)
			if !ok {
				log.Println("Invalid message format")
				continue
			}

			// Route the message to the appropriate handler
			switch messageType {
			case "register":
				RegisterHandler(conn, r, db, message)

			case "login":
				LoginHandler(conn, r, db, message)
			case "homePage":
				// Get needed data for the home page
				allUsernames, err := GetAllUsernames(db)
				if err != nil {
					SendWebSocketMessage(conn, Response{Success: false, Message: "Failed to get all usernames"})
					return
				}

				allPosts, err := GetAllPosts(db)
				if err != nil {
					SendWebSocketMessage(conn, Response{Success: false, Message: "Failed to get all posts"})
					return
				}
				allCategories, err := GetCategories(db)
				if err != nil {
					SendWebSocketMessage(conn, Response{Success: false, Message: "Failed to get all categories"})
					return
				}
				allComments, err := GetAllComments(db)
				if err != nil {
					SendWebSocketMessage(conn, Response{Success: false, Message: "Failed to get all comments"})
					return
				}
				usersOnline, err := GetAllOnlineUsers(db)
				if err != nil {
					SendWebSocketMessage(conn, Response{Success: false, Message: "Failed to get all online users"})
					return
				}
				// Send the data back to the frontend
				responseData := struct {
					AllUsernames   []string     `json:"allUsernames"`
					AllPosts       []Post       `json:"allPosts"`
					AllCategories  []Category   `json:"allCategories"`
					AllComments    []Comment    `json:"allComments"`
					AllUsersOnline []OnlineUser `json:"usersOnline"`
				}{
					AllUsernames:   allUsernames,
					AllPosts:       allPosts,
					AllCategories:  allCategories,
					AllComments:    allComments,
					AllUsersOnline: usersOnline,
				}

				SendWebSocketMessageSuccess(conn, SuccessResponse{Type: "allData", Success: true, Message: "Home Page Data", Data: responseData})
				BroadcastChanges(conn, "homePageUpdate", responseData)

			case "createPost":
				CreatePostHandler(conn, r, db, message)
			case "submitComment":
				SubmitCommentHandler(conn, r, db, message)

			case "userLogout":
				DeleteSessionHandler(conn, r, db, message)
			case "newMessage":
				SubmitMessageHandler(conn, r, db, message)

			}
		}
	}
}

type WebSocketUpdate struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Broadcast changes to all connected clients
func BroadcastChanges(conn *websocket.Conn, messagetype string, data interface{}) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	// Create a WebSocketUpdate message
	updateMessage := WebSocketUpdate{
		Type: messagetype,
		Data: data,
	}

	// Iterate over all connected clients and send them the update
	for client := range clients {
		err := client.WriteJSON(updateMessage)
		if err != nil {
			log.Println("Error sending update to client:", err)
			// Handle error as needed (e.g., remove the client from the list)
		}
	}
}
