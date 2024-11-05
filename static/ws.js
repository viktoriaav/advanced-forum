import { getState, updateState, resetState } from './state.js';
import { router } from './index.js'
import { navigateTo } from './index.js';


export function connectWebSocket() {
    const socket = new WebSocket('ws://localhost:8090/ws');

    socket.addEventListener('open', async (event) => {
        console.log('WebSocket connection opened:', event);
        
        let state = getState();
        updateUI(state.loggedInUsername);
    });

    socket.addEventListener('close', (event) => {
        console.log('WebSocket connection closed:', event);
    });

    socket.addEventListener('message', (event) => {
        receiveWebSocketMessage(event);
    });

    return socket;
}
// Establish WebSocket connection
export const socket = connectWebSocket();


export function sendMessage(message) {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(message));
        console.log('Sent Message:', message);
    } else {
        console.error('WebSocket is not open. Unable to send message.');
    }
}

export function receiveWebSocketMessage(event) {
    console.log('Raw WebSocket Message:', event.data);

    try {
        let data = JSON.parse(event.data);
        console.log('Parsed Data:', data);

        if (!data) {
            console.error('Invalid or missing data structure:', data);
            return;
        }
        let state;

        switch (data.type) {
            case "Registration":
            case "Login":
                updateState({
                    isAuthenticated: data.data.isAuthenticated,
                    loggedInUsername: data.data.loggedInUsername,
                    allMessagesForUser: data.data.allMessages
                });
                
                state = getState();
                updateUI(state.loggedInUsername);

                const homePage = {
                    message: "homePage"
                };
                sendMessage(homePage);
                router();
                break;

            case "allData":
                updateState({
                    allPosts: data.data.allPosts,
                    AllUsernames: data.data.allUsernames,
                    allCategories: data.data.allCategories,
                    allComments: data.data.allComments,
                    NotifyAllUsersOnlineStatus: data.data.usersOnline
                });
                router();
                break;

            case "createdPost":
                updateState({
                    allPosts: data.data.allPosts
                });
                navigateTo("/");
                router();
                break;

            case "newComment":
                updateState({
                    allComments: data.data.allComments
                });
                router();
                break;

            case "Error":
                updateState({
                    errorMessage: data.message
                });
                navigateTo("/error");
                break;

            case "homePageUpdate":
                state = getState();
                if (state.isAuthenticated) {
                    updateState({
                        allPosts: data.data.allPosts,
                        AllUsernames: data.data.allUsernames,
                        allCategories: data.data.allCategories,
                        allComments: data.data.allComments,
                        NotifyAllUsersOnlineStatus: data.data.usersOnline
                    });
                    router();
                }
                break;

            case "updatAllUsersOnline":                
                state = getState();
                if (state.isAuthenticated) {
                    updateState({
                        NotifyAllUsersOnlineStatus: data.data.usersOnline
                    });
                    router();
                }
                break;

            case "updateAllPosts":
                state = getState();
                if (state.isAuthenticated) {
                    updateState({
                        allPosts: data.data.allPosts
                    });
                    router();
                }
                break;

            case "updateAllComments":
                state = getState();
                if (state.isAuthenticated) {
                    updateState({
                        allComments: data.data.allComments
                    });
                    router();
                }
                break;
            case "updateAllMessages":
                state = getState();
                if (state.isAuthenticated) {
                    updateState({
                        allMessagesForUser: data.data.allMessages
                    });
                    router();
                }
            default:
                console.warn('Unhandled message type:', data.type);
        }
    } catch (error) {
        console.error('Error parsing JSON:', error);
    }
}



// Function to update the user interface
function updateUI(loggedInUsername) {
    const profileDiv = document.getElementById('profile');
    profileDiv.innerHTML = '';

    if (loggedInUsername) {
    
        // Update UI with the global variable
        profileDiv.innerHTML = `
            <div class="dropdown">
                <button class="greeting">
                    Hi, ${loggedInUsername}
                </button>
                <div class="logout" id="logoutButton" data-link>
                    <img class="sign" src="../static/images/signout.png">
                </div>
                <div class="chat-btn" id="chatsButton" data-link>
                    <img class="sign" src="../static/images/chat.png">
                </div>
            </div>
        `;
        const chatsButton = document.getElementById('chatsButton');

        // Add an event listener to the button
        chatsButton.addEventListener('click', function() {
            navigateTo("/chats");
        });

        const logoutButton = document.getElementById('logoutButton');

        // Add an event listener to the button
        logoutButton.addEventListener('click', function() {
            handleLogout();
        });
    } else {
        // Update UI even if no data is present
        profileDiv.innerHTML = `
            <div class="login">
                <a onclick="showPopup('loginPopup')">
                    Login
                    <img class="sign" src="../static/images/sign-in.png">
                </a>
            </div>
            <div class="signup">
                <a onclick="showPopup('signupPopup')">
                    Register
                    <img class="sign" src="../static/images/register.png">
                </a>
            </div>
        `;
    }
}


// Function to handle logout
function handleLogout() {
    let state = getState();
    const userLeft = {
        message: "userLogout",
        username: state.loggedInUsername,
    };
    sendMessage(userLeft);

    resetState();
    state = getState();
    updateUI(state.loggedInUsername);
    navigateTo("/")
    router();
    setTimeout(() => {
        // Close the WebSocket connection
        if (socket.readyState === WebSocket.OPEN) {
            socket.close();
        }
      }, 2000);
    

}