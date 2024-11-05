import AbstractView from "./AbstractView.js";
import { getState } from '../state.js';
import { sendMessage } from "../ws.js";
import { navigateTo } from "../index.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Create post");
    }

    async updateApp() {
        let state = getState();


        // Define a function to create category checkboxes
        function createCategoryCheckboxes() {
            const categories = state.allCategories.map((category) => {
                return `
                <div class="checkbox-rect">
                    <input class="checkbox-spin" type="checkbox" id="${category.category}" name="categories[]" value="${category.category}">
                    <label for="${category.category}">
                        ${category.category}
                    </label>
                </div>`
            });

            return categories.join("");
        }
        

        return `
            <div class="create-post" id="create-post">
                <div class="create-form">
                    <form method="POST" class="create-post-form">
                        <div class="back-home-wrap" id="back-home">
                            <div class="back-home">
                                <a href="/" class="back-home-btn" id="back-home-btn" data-link>Back on Home Page</a>
                            </div>
                        </div>
                        <div class="start-discussion">
                            <span>What's on your mind?</span>
                        </div>
                        <div class="category-choose">
                            ${createCategoryCheckboxes()}
                        </div>
                        <input type="hidden" id="createdBy" name="createdyBy" value="${state.loggedInUsername}">
                        <input type="text" id="title" name="title" placeholder="Post title ..." required> <br>
                        <input type="text" id="content" name="content" placeholder="Post content ..." required> <br>
                        <div class="submit-post">
                            <input class="submit" type="submit" value="Submit">
                        </div>
                    </form>
                </div>
            </div>
        `;
    }
    async pageAction() {
        let state = getState();
        if (state.isAuthenticated) {
            const createPostForm = document.querySelector(".create-post-form");

            createPostForm.addEventListener("submit", function (e) {
                e.preventDefault(); // Prevent the default form submission
                const titleInput = document.getElementById("title");
                const contentInput = document.getElementById("content");
                const createdByInput = document.getElementById("createdBy");
                const selectedCategoriesInput = Array.from(document.querySelectorAll('input[name="categories[]"]:checked'));

                const title = titleInput.value;
                const content = contentInput.value;
                const createdBy = createdByInput.value;
                const selectedCategories = selectedCategoriesInput.map(checkbox => checkbox.value);

                // Prepare the post data
                const postData = {
                    message: "createPost",
                    createdBy: createdBy,
                    title: title,
                    content: content,
                    categories: selectedCategories,
                };
                console.log("WebSocket Message:", postData);

                // Send the post data as a JSON string to the WebSocket
                sendMessage(postData);
                navigateTo("/");
            });
        }
    }
}
