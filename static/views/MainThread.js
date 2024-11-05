import AbstractView from "./AbstractView.js";
import { getState } from '../state.js';
import { sendMessage } from "../ws.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Real-Time-Forum");
    }
    
    async updateApp() {
        
        let state = getState();

       function createUsernamesParagraphs() {
            // Check if state or required properties are null or undefined
            if (!state || !state.AllUsernames || !state.NotifyAllUsersOnlineStatus) {
                return '';  // Return an empty string or handle the error in an appropriate way
            }

            const onlineUsernames = state.NotifyAllUsersOnlineStatus.map(user => user.username);

            const usernames = state.AllUsernames.map((username) => {
                // Check if the user is online
                const isOnline = onlineUsernames.includes(username);

                // Add a green or red circle based on online status
                const onlineStatusIndicator = isOnline ? '<div class="online-indicator"><p class="circle green"></p></div>' : '<div class="online-indicator"><p class="circle red"></p></div>';

                return `
                    <div class="user-with-indicator">
                        <p class="category-button category">${username}</p>
                        ${onlineStatusIndicator}
                    </div>
                `;
            });

            return usernames.join("");
        }

    
        // Define a function to create the posts
        function createPosts() {
            const posts = state.allPosts.map((post) => {
                const truncatedContent = post.content.slice(0, 100); // Take only the first 100 characters
                return `
                    <div class="post" data-username="${post.username}" data-category="${post.post_category}" id="post">
                        <div class="post-category">
                            <span>${post.post_category}</span>
                        </div>
                        <a href="/post/${post.post_id}" class="title" data-link>${post.title} by: ${post.username}</a>
                        <p class="content">${truncatedContent}...</p>
                        <div class="reactions">
                            <a href="/post/${post.post_id}" class="comments" data-link></a>
                        </div>
                    </div>
                `;
            });
            return posts.join("");
        }

        // Generate the HTML based on the user's login status
        if (state.isAuthenticated) {
            const generatedHTML = `
                <div class="body" id="body">
                    <div class="discussion">
                        <a href="/create-post" id="create-post-btn" data-link>
                            Create a post
                        </a>
                    </div>
                    <div class="post-wrapper">
                        <div class="category-buttons categories">
                            ${createUsernamesParagraphs()}
                        </div>
                        <div class="all-posts" id="all-posts">
                            <h2 class="posts">Posts</h2>
                            ${createPosts()}
                        </div>
                    </div>
                </div>
            `;

            // Append the generated HTML to the "app" element
            return generatedHTML;
        } else {
            const generatedHTML = `
                <div class="login-to-continue-wrap">
                    <div class="login-to-continue">
                        <p>Login or register to continue!</p>
                    </div>
                </div>
            `;

            // Append the generated HTML to the "app" element
            return generatedHTML;
        }
    }

    async pageAction(){
        let state = getState();
        if (!state.isAuthenticated) {
            const loginForm = document.querySelector(".form");
            const identifierInput = document.getElementById("identifier");
            const passwordInput = document.getElementById("password");


            loginForm.addEventListener("submit", function (e) {
                e.preventDefault(); // Prevent the default form submission

                const identifier = identifierInput.value;
                const password = passwordInput.value;
                // Prepare the login data
                const loginData = {
                    message: "login",
                    identifier: identifier,
                    password: password,
                };
                console.log("WebSocket Message:", loginData);
                // Send the login data as a JSON string to the WebSocket
                sendMessage(loginData);
                closePopup("loginPopup");
            });

            const registrationForm = document.getElementById("form-registration");
            const emailInput = document.getElementById("email");
            const firstNameInput = document.getElementById("first-name");
            const lastNameInput = document.getElementById("last-name");
            const ageInput = document.getElementById("age");
            const usernameInput = document.getElementById("username");
            const passwordRegInput = document.getElementById("password-reg");
            const genderInput = document.getElementById("gender");

            registrationForm.addEventListener("submit", function (e) {
                e.preventDefault(); // Prevent the default form submission

                const email = emailInput.value;
                const firstName = firstNameInput.value;
                const lastName = lastNameInput.value;
                const age = ageInput.value;
                const gender = genderInput.value;
                const username = usernameInput.value;
                const password = passwordRegInput.value;

                // Prepare the registration data
                const registrationData = {
                    message: "register",
                    email: email,
                    "first-name": firstName,
                    "last-name": lastName,
                    age: age,
                    gender: gender,
                    username: username,
                    password: password,
                };
                console.log("WebSocket Message:", registrationData);

                // Send the registration data as a JSON string to the WebSocket
                sendMessage(registrationData);

                // Optionally, you can close the popup here if needed
                closePopup("signupPopup");
            });
        } 
    }   
}

