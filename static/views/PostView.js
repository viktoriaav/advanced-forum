import AbstractView from "./AbstractView.js";
import { getState } from '../state.js';
import { sendMessage } from "../ws.js";
import { navigateTo } from "../index.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.postId = params.id;
        this.setTitle("Viewing Post");
    }

    async updateApp() {
        let state = getState();
        
        // Helper function to find the post by postId
        function findPostById(allPosts, postId) {
            console.log("All Posts:", allPosts);
            console.log("Searching for post with postId:", postId);
        
            const matchingPosts = allPosts.filter(post => post.post_id === Number(postId));
        
            if (matchingPosts.length > 0) {
                const foundPost = matchingPosts[0];
                console.log("Found post:", foundPost);
                return foundPost;
            } else {
                console.log("Post not found");
                return undefined;
            }
        }
        
        

        let selectedPost = findPostById(state.allPosts, this.postId);
        console.log(selectedPost.post_id)

        // Helper function to filter comments for the selected post
        function filterCommentsForPost(allComments, postId) {
            console.log("All Comments:", allComments);

            return allComments.filter(comment => comment.post_comment_id === Number(postId));
        }

        // Filter comments for the selected post
        let commentsForPost = filterCommentsForPost(state.allComments, this.postId);
        console.log("Comments for post:", commentsForPost)

        // Define a function to create comments
        function createComments(commentsForPost) {
            const commentsHTML = commentsForPost.map((comment) => `
                <div class="comment">
                    <p class="title"> by: ${comment.username}</p>
                    <p class="content">${comment.content}</p>
                </div>
            `).join("");

            return `
                <div class="post-comments" id="post-comments">
                    <p class="all-comments">Comments</p>
                    ${commentsHTML}
                </div>
            `;
        }

        // Define a function to create the comment form
        function createCommentForm() {
            return `
                <div class="comment-form">
                    <p class="login-to">Leave a Comment</p>
                    <form method="POST" class="comment-form-submit">
                        <input type="hidden" name="postID" id="postId" value="${selectedPost.post_id}">
                        <input type="text" id="comment" name="comment" placeholder="Your comment here ..." required> <br>
                        <input type="submit" value="Submit" class="submit">
                    </form>
                </div>
            `;
        }
        function createPostInfo(selectedPost) {
            return `
                <div class="info-post">
                    <div class="post-category">
                        <span>${selectedPost.post_category}</span>
                    </div>
                    <p class="title">${selectedPost.title} by: ${selectedPost.username}</p>
                    <p class="content">${selectedPost.content}</p>
                </div>
            `;
        }

        return `
            <div class="post-page" id="post-page">
                <div class="back-home-wrap" id="back-home">
                    <div class="back-home">
                        <a href="/" class="back-home-btn" id="back-home-btn" data-link>Back on Home Page</a>
                    </div>
                </div>
                <div class="post-container">
                    ${createPostInfo(selectedPost)}
                    ${createComments(commentsForPost)}
                    ${createCommentForm()}
                </div>
            </div>
        `;
    }
    async pageAction() {
        let state = getState();
        if (state.isAuthenticated) {
    
            const commentForm = document.querySelector(".comment-form-submit");
    
            commentForm.addEventListener("submit", function (e) {
                e.preventDefault(); // Prevent the default form submission
                const commentInput = document.getElementById("comment");
                const postIDInput = document.getElementById("postId");
    
                const comment = commentInput.value;
                const postID = postIDInput.value;
                const username = state.loggedInUsername;
    
                // Prepare the comment data
                const commentData = {
                    message: "submitComment",
                    postID: postID,
                    comment: comment,
                    username: username
                };
                console.log("WebSocket Message:", commentData);
    
                // Send the comment data as a JSON string to the WebSocket
                sendMessage(commentData);
                navigateTo(`/post/${postID}`);
                
            });
        }
    }
}
