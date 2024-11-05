// Define the initial state
const initialState = {
    loggedInUsername: null,
    allPosts: { post: {} },
    AllUsernames: null,
    isAuthenticated: false,
    allCategories: { category: {} },
    allComments: { comment: {} },
    NotifyAllUsersOnlineStatus: { onlineUsers: {} },
    errorMessage: null,
    allMessagesForUser : { message: {} },
    chatOpen : false,
    selectedChatUsername: null,
    sendTypingNotification: false
};

// Set the initial state
let state = { ...initialState };

// Function to reset the state to the initial state
function resetState() {
    state = { ...initialState };
}



function updateState(changes) {
    state = { ...state, ...changes };
    console.log(state)
    
}

function getState() {
    return { ...state };
}
export { updateState, getState, resetState };