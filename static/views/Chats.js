import AbstractView from "./AbstractView.js";
import { getState, updateState } from '../state.js';
import { navigateTo, router } from "../index.js";
import { sendMessage } from "../ws.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Chats");
    }

    async updateApp() {
        let state = getState();
        
        function createUsernamesParagraphs() {
            console.log("Current state:", getState());
            if (!state || !state.AllUsernames || !state.NotifyAllUsersOnlineStatus) {
                return '';
            }
        
            const onlineUsernames = state.NotifyAllUsersOnlineStatus.map(user => user.username);
        
            // Sort usernames based on last message timestamp (if available)
            const sortedUsernames = state.AllUsernames.sort((a, b) => {
                const lastMessageA = getLastMessageTimestamp(a);
                const lastMessageB = getLastMessageTimestamp(b);
        
                if (lastMessageA !== undefined && lastMessageB !== undefined) {
                    return lastMessageB - lastMessageA;
                } else {
                    // If there are no messages, sort alphabetically
                    return a.localeCompare(b);
                }
            });
        
            const usernames = sortedUsernames.map((username) => {
                const isOnline = onlineUsernames.includes(username);
                const onlineStatusIndicator = isOnline ? '<div class="online-indicator"><p class="circle green"></p></div>' : '<div class="online-indicator"><p class="circle red"></p></div>';
        
                return `
                    <div class="user-wrap" id="user-select" data-username="${username}">
                        <div class="user-list-item">${username}</div>
                        ${onlineStatusIndicator}
                    </div>
                `;
            });
        
            return usernames.join("");
        }
        
        

        function getLastMessageTimestamp(username) {
            const loggedInUsername = state.loggedInUsername;
        
            // Filter messages based on the sender or receiver being the loggedInUsername
            const userMessages = state.allMessagesForUser.filter(message =>
                (message.sender === loggedInUsername && message.receiver === username) ||
                (message.sender === username && message.receiver === loggedInUsername)
            );
        
            // Find the latest message timestamp among the filtered messages
            const lastMessage = userMessages.reduce((latestMessage, currentMessage) => {
                const currentTimestamp = new Date(currentMessage.created_at).getTime();
        
                return (!latestMessage || currentTimestamp > new Date(latestMessage.created_at).getTime()) ?
                    currentMessage :
                    latestMessage;
            }, null);
        
            return lastMessage ? new Date(lastMessage.created_at).getTime() : null;
        }
        
        function displayChat() {
            if (!state || !state.chatOpen) {
                return '<div class="chat-message">Select a chat</div>';
            }
        
            // Get the selected chat's messages
            const selectedUsername = state.selectedChatUsername;
            const selectedChatMessages = state.allMessagesForUser.filter((message) =>
                (message.sender === state.loggedInUsername && message.receiver === selectedUsername) ||
                (message.sender === selectedUsername && message.receiver === state.loggedInUsername)
            );
        
            // Load only the last 10 messages if there are more than 10
            const messagesToDisplay = selectedChatMessages.slice(-10);
        
            const messages = messagesToDisplay.map((message) => {
                const { sender, content, created_at } = message;
        
                // Convert the timestamp to a Date object
                const messageDate = new Date(created_at);
        
                // Format the date in "hour:minute day/month/year" format
                const formattedTime = new Intl.DateTimeFormat('en-US', {
                    hour: 'numeric',
                    minute: 'numeric',
                    hour12: true,
                    day: 'numeric',
                    month: 'short',
                    year: '2-digit'
                }).format(messageDate);
        
                return `
                    <p class="chat-message">
                        <span class="sender">${sender}</span>
                        <span class="content-message">${content}</span>
                        <span class="time">${formattedTime}</span>
                    </p>
                `;
            });
        
        
            // Join messages and scroll to the bottom
            const chatContainer = messages.join("");

            
            // Scroll to the bottom
            setTimeout(() => {
                const chatElement = document.getElementById('all-messages');
                chatElement.scrollTop = chatElement.scrollHeight;
            }, 0);
        
            return chatContainer;
        }

        function sendMessageBox() {
            const selectedUsername = state.selectedChatUsername;
            const isUserOnline = state.NotifyAllUsersOnlineStatus.some(user => user.username === selectedUsername);

            if (isUserOnline) {
                return `
                    <div class="send-message-box">
                        <form id="message-form">
                            <input type="text" id="message-input" placeholder="Type your message" required />
                            <button type="submit" id="send-button">Send</button>
                        </form>
                    </div>
                `;
            } else {
                return '';
            }
        }
        const userChattingDiv = state.selectedChatUsername ? `<div class="userChatting" id="username-Chat">${state.selectedChatUsername}</div>` : '';
        
        return `
            <div class="all-chats">
            <div class="chat-field">
                <div class="users-avaliable">
                    <div class="back-field" id="back-to-home">
                        <div class="arrow-icon">
                            <p class="arrow-back">← Back</p>
                        </div>
                    </div>

                    <div class="users">
                        <div class="user-list" style="max-height: 800px; overflow-y: auto;">
                            ${createUsernamesParagraphs()}
                        </div>
                    </div>
                </div>
                <div class="chat-area" id="chat-area">
                    ${userChattingDiv}
                    <div class="all-messages" id="all-messages" style="max-height: none; overflow-y: auto;">
                        ${displayChat()}
                    </div>
                    <div id="typing-indicator"></div>
                    ${sendMessageBox()}
                </div>
            
            </div>
        </div>    
    `;
    }

    async pageAction() {
        let state = getState();
        const backHome = document.getElementById("back-to-home");
        const sendButton = document.getElementById("send-button");

        const userSelectElements = document.querySelectorAll(".user-wrap");

        userSelectElements.forEach((selectElement) => {
            selectElement.addEventListener("click", function () {
                const selectedUsername = this.dataset.username;
                updateState({
                    chatOpen: true,
                    selectedChatUsername: selectedUsername
                });

                router();

            });
        });

        backHome.addEventListener("click", function () {
            navigateTo("/");
        });
        const messageInput = document.getElementById("message-input");
        const messageForm = document.getElementById("message-form");
        if (messageForm) {
            messageForm.addEventListener("submit", function (event) {
                event.preventDefault(); // Prevents the default form submission behavior

                const loggedInUsername = state.loggedInUsername;
                const selectedChatUsername = state.selectedChatUsername;


                if (messageInput && loggedInUsername && selectedChatUsername) {
                    const newMessage = {
                        message: "newMessage",
                        sender: loggedInUsername,
                        receiver: selectedChatUsername,
                        content: messageInput.value,
                        created_at: new Date().toISOString()
                    };

                    sendMessage(newMessage);

                    messageInput.value = "";
                    console.log("Sending message:", newMessage);
                }
            });
        }
        if (messageInput) {
            messageInput.addEventListener("keydown", function (event) {
                if (event.key === "Enter") {
                    event.preventDefault();
    
                    const loggedInUsername = state.loggedInUsername;
                    const selectedChatUsername = state.selectedChatUsername;
    
                    if (messageInput && loggedInUsername && selectedChatUsername) {
                        const newMessage = {
                            message: "newMessage",
                            sender: loggedInUsername,
                            receiver: selectedChatUsername,
                            content: messageInput.value,
                            created_at: new Date().toISOString()
                        };
    
                        sendMessage(newMessage);
                            messageInput.value = "";
                        console.log("Sending message:", newMessage);
                    }
                }
            });
        }
        
        // Get the selected chat's messages
        const selectedUsername = state.selectedChatUsername;
        const selectedChatMessages = state.allMessagesForUser.filter((message) =>
            (message.sender === state.loggedInUsername && message.receiver === selectedUsername) ||
            (message.sender === selectedUsername && message.receiver === state.loggedInUsername)
        );
        let messagesToShow = 10;

        // Add an onscroll event listener to the chat element with debounce
        const chatElement = document.getElementById('all-messages');
        chatElement.onscroll = debounce(function () {
            // Check if the user has scrolled to the top of the chat
            if (chatElement.scrollTop === 0) {
                // If yes, load 10 more messages
                messagesToShow += 10;

                // Update the messages to display based on the new count
                const newMessagesToDisplay = selectedChatMessages.slice(-messagesToShow);

                // Map and format the new messages
                const newMessages = newMessagesToDisplay.map((message) => {
                    const { sender, content, created_at } = message;

                    // Convert the timestamp to a Date object
                    const messageDate = new Date(created_at);

                    // Format the date in "hour:minute day/month/year" format
                    const formattedTime = new Intl.DateTimeFormat('en-US', {
                        hour: 'numeric',
                        minute: 'numeric',
                        hour12: true,
                        day: 'numeric',
                        month: 'short',
                        year: '2-digit'
                    }).format(messageDate);

                    return `
                        <p class="chat-message">
                            <span class="sender">${sender}</span>
                            <span class="content-message">${content}</span>
                            <span class="time">${formattedTime}</span>
                        </p>
                    `;
                });

                // Insert the new messages at the beginning of the chat container
                chatElement.innerHTML = newMessages.join("") + chatElement.innerHTML;
            }
        }, 200); // 200ms debounce time (adjust as needed)

        // Debounce function
        function debounce(func, delay) {
            let timeoutId;
            return function () {
                clearTimeout(timeoutId);
                timeoutId = setTimeout(func, delay);
            };
        }
    }
}
