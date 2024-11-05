import MainThread from "./views/MainThread.js";
import CreatePost from "./views/CreatePost.js";
import PostView from "./views/PostView.js";
import Chats from "./views/Chats.js";
import ErrorPage from "./views/ErrorPage.js";

const pathToRegex = (path) =>
  new RegExp("^" + path.replace(/\//g, "\\/").replace(/:\w+/g, "(.+)") + "$");

const getParams = (match) => {
  const values = match.result.slice(1);
  const keys = Array.from(match.route.path.matchAll(/:(\w+)/g)).map(
    (result) => result[1]
  );

  return Object.fromEntries(
    keys.map((key, i) => {
      return [key, values[i]];
    })
  );
};

export const navigateTo = (url) => {
  history.pushState(null, null, url);
  router();
};



export const router = async () => {
  const routes = [
    { path: "/", view: MainThread },
    { path: "/create-post", view: CreatePost },
    { path: "/post/:id", view: PostView },
    { path: "/chats", view: Chats },
    { path: "/error", view: ErrorPage },

  ];

  // Test each route for potential match
  const potentialMatches = routes.map((route) => {
    return {
      route: route,
      result: location.pathname.match(pathToRegex(route.path)),
    };
  });

  let match = potentialMatches.find(
    (potentialMatch) => potentialMatch.result !== null
  );

  if (!match) {
    match = {
      route: routes[0],
        result: [location.pathname],
    };
  }

  let view = new match.route.view(getParams(match)); // Assign the view here

  document.querySelector("#app").innerHTML = await view.updateApp();
  await view.pageAction(); 
};

window.addEventListener("popstate", router);

document.addEventListener("DOMContentLoaded", () => {
  document.body.addEventListener("click", (e) => {
    if (e.target.matches("[data-link]")) {
      e.preventDefault();
      navigateTo(e.target.href);
    }
  });

  router();
});
