package forum

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// CATEGORIES
type Category struct {
	CategoryID int    `json:"category_id"`
	Category   string `json:"category"`
}

func GetCategories(db *sql.DB) ([]Category, error) {
	if db == nil {
		return nil, errors.New("nil database connection")
	}
	var categories []Category
	query := "SELECT category FROM categories"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category Category
		err := rows.Scan(&category.Category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// POSTS
type Post struct {
	PostID       int        `json:"post_id"`
	Username     string     `json:"username"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	Categories   []Category `json:"categories"`
	PostCategory string     `json:"post_category"`
	CreatedAt    string     `json:"created_at"`
}

func GetAllPosts(db *sql.DB) ([]Post, error) {
	var posts []Post
	query := `
        SELECT p.post_ID, u.username, p.title, p.content, p.created_at,
               c.category
        FROM posts AS p
        INNER JOIN users AS u ON p.user_ID = u.user_ID
        INNER JOIN post_categories AS pc ON p.post_ID = pc.post_ID
        INNER JOIN categories AS c ON pc.category_ID = c.category_ID
        GROUP BY p.post_ID, u.username, p.title, p.content, p.created_at
    `
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		var category Category
		err := rows.Scan(
			&post.PostID, &post.Username, &post.Title, &post.Content, &post.CreatedAt,
			&category.Category,
		)
		if err != nil {
			return nil, err
		}

		// Fetch categories for the current post
		categories, err := GetCategoriesForPost(db, post.PostID)
		if err != nil {
			return nil, err
		}
		post.PostCategory = strings.Join(categories, " ") // Join the categories into a single string

		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostByID(db *sql.DB, postID int) (Post, error) {
	var post Post
	query := `
        SELECT p.post_ID, u.username, p.title, p.content, p.created_at,
               c.category
        FROM posts AS p
        INNER JOIN users AS u ON p.user_ID = u.user_ID
        INNER JOIN post_categories AS pc ON p.post_ID = pc.post_ID
        INNER JOIN categories AS c ON pc.category_ID = c.category_ID
        WHERE p.post_ID = ?
        GROUP BY p.post_ID, u.username, p.title, p.content, p.created_at
    `
	rows, err := db.Query(query, postID)
	if err != nil {
		return Post{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var category Category
		err := rows.Scan(
			&post.PostID, &post.Username, &post.Title, &post.Content, &post.CreatedAt,
			&category.Category,
		)
		if err != nil {
			return Post{}, err
		}

		// Fetch categories for the current post
		categories, err := GetCategoriesForPost(db, post.PostID)
		if err != nil {
			return Post{}, err
		}
		post.PostCategory = strings.Join(categories, " ") // Join the categories into a single string
	}

	return post, nil
}

func GetCategoriesForPost(db *sql.DB, postID int) ([]string, error) {
	categories := []string{}
	query := `
		SELECT c.category
		FROM categories AS c
		INNER JOIN post_categories AS pc ON c.category_ID = pc.category_ID
		WHERE pc.post_ID = ?
	`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

type Comment struct {
	CommentID int    `json:"comment_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	PostID    int    `json:"post_comment_id"`
}

func GetAllComments(db *sql.DB) ([]Comment, error) {
	var comments []Comment
	query := `
		SELECT com.comment_ID, u.username, com.content, com.post_ID
		FROM comments AS com
		INNER JOIN users AS u ON com.user_ID = u.user_ID
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.CommentID, &comment.Username, &comment.Content, &comment.PostID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// USERname
func GetLoggedInUsername(sessionToken string, db *sql.DB) (string, error) {
	// Use sessionToken in your SQL query or any other logic as needed
	query := `
        SELECT u.username
        FROM sessions s
        INNER JOIN users u ON s.user_ID = u.user_ID
        WHERE s.token = ? AND s.expires_at > strftime('%s', 'now')
        LIMIT 1
    `
	var username string
	err := db.QueryRow(query, sessionToken).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
func GetAllUsernames(db *sql.DB) ([]string, error) {
	var usernames []string
	query := "SELECT username FROM users"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}

	return usernames, nil
}

// getUserID retrieves the user ID based on the username
func getUserID(username string, db *sql.DB) (int, error) {
	query := "SELECT user_ID FROM users WHERE LOWER(username) = ?"
	var userID int
	err := db.QueryRow(query, strings.ToLower(username)).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

type OnlineUser struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
}

func GetAllOnlineUsers(db *sql.DB) ([]OnlineUser, error) {
	query := `
        SELECT u.user_ID, u.username
        FROM users AS u
        JOIN sessions AS s ON u.user_ID = s.user_ID
        WHERE s.expires_at > strftime('%s', 'now')  -- Check if session is still active
    `
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var onlineUsers []OnlineUser

	for rows.Next() {
		var user OnlineUser
		if err := rows.Scan(&user.UserID, &user.Username); err != nil {
			return nil, err
		}
		onlineUsers = append(onlineUsers, user)
	}

	return onlineUsers, nil
}

type Message struct {
	ID        int       `json:"message_ID"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// GetAllMessagesForUser fetches all messages associated with a user (either as a sender or receiver)
func GetAllMessages(db *sql.DB) ([]Message, error) {
	query := "SELECT message_ID, sender, receiver, content, created_at FROM private_messages ORDER BY created_at ASC"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]Message, 0)

	for rows.Next() {
		var message Message
		err := rows.Scan(&message.ID, &message.Sender, &message.Receiver, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
