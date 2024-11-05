export default class {
    constructor(params) {
        this.params = params;
    }

    setTitle(title) {
        document.title = title;
    }

    async updateApp() {
        return ""; // Default implementation, you should override this method in each specific view
    }
    async pageAction() {
        return;
    }
}
