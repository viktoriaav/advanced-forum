# Advanced Forum

# Project Objectives

In this project, I focused on developing a new and upgraded forum that includes several key features:

## Key Features:

1. **Registration and Login**
   - I implemented a registration and login system where users must register to access the forum. The registration form collects essential information, including:
     - Nickname
     - Age
     - Gender
     - First Name
     - Last Name
     - E-mail
     - Password
   - Users can log in using either their nickname or email, combined with their password.
   - A logout function allows users to sign out from any page on the forum.

2. **Posts and Comments**
   - Users can create posts categorized similarly to the previous forum. 
   - Commenting functionality enables users to respond to posts.
   - I designed a feed display for posts, where users can view posts and access comments by clicking on them.

3. **Private Messages**
   - I developed a chat feature for users to send private messages to one another. This includes:
     - A section displaying online/offline users, organized by the most recent message. New users with no messages are arranged alphabetically.
     - The ability for users to send messages to those online, with this section always visible.
     - When a user clicks on another user’s name, it loads the past messages, displaying the last 10 messages. Scrolling up retrieves an additional 10 messages, implemented with throttling and debouncing techniques to prevent excessive requests.
   - Messages have a specific format that includes:
     - A timestamp indicating when the message was sent.
     - The username of the sender.

## Real-Time Functionality
To ensure a seamless user experience, I integrated WebSockets for real-time messaging. This allows users to receive notifications of new messages instantly without needing to refresh the page.

## Technical Implementation
- **Database**: I used SQLite for data storage, similar to the previous forum.
- **Backend**: The backend was built using Golang to handle data processing and WebSocket communication.
- **Frontend**: JavaScript managed all client-side events and WebSocket interactions, creating a dynamic single-page application (SPA).
- **HTML**: I organized the page elements within a single HTML file, allowing for easy navigation through JavaScript.
- **CSS**: I styled the elements to enhance the user interface.

Through this project, I created a robust and interactive forum that enhances user engagement and communication.


## Used Packages

- All standard go packages are allowed.
- Gorilla websocket
- sqlite3
- bcrypt
- UUID

## Usage

1. You can start the program by running the following command:
```
go run .
```
2. Open http://localhost:8090
3. To end the server:
```
CTRL + C
```

## Developers

- [Olia Priadkina/Olha_Priadkina](https://01.kood.tech/git/Olha_Priadkina)
- [Viktoriia/vavstanc](https://01.kood.tech/git/vavstanc)