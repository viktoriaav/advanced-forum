import AbstractView from "./AbstractView.js";
import { getState } from "../state.js";

export default class extends AbstractView {
    constructor(params) {
        super(params);
        this.setTitle("Error");
    }

    async updateApp() {
        const state = getState();
        console.log(state.errorMessage)
        return `  
        <div class="error-wrap"
            <div class="back-home-wrap" id="back-home">
                <div class="back-home-error">
                    <a href="/" class="back-home-btn" id="back-home-btn" data-link>Go back</a>
                </div>
                <div class="error-page">
                    <div class="error-message">${state.errorMessage}</div>
                </div>
            </div>
        </div>
    ` 
}
}