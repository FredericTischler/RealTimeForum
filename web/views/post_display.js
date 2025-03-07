// web/views/post_display.js

import { displayErrorModal } from './error_modal.js';

const postsContainer = document.getElementById("postsContainer");


export let currentCategory = "";
export let currentKeyword = "";

let limit = 2;
let offset = 0;
let isLoading = false;
let allPostsLoaded = false;

// Fonction pour initialiser le chargement des posts et l'écoute du scroll
export function displayPosts() {
    loadMorePosts();
    window.addEventListener("scroll", handleScroll);
}

async function loadMorePosts() {
    if (isLoading || allPostsLoaded) return;
    isLoading = true;
    try {
        // Construire la query string avec les filtres courants
        const queryParams = new URLSearchParams({
            limit: limit,
            offset: offset,
            category: currentCategory,
            keyword: currentKeyword
        });
        const response = await fetch(`/posts?${queryParams.toString()}`, {
            method: "GET",
            credentials: "include"
        });
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText);
        }
        const posts = await response.json();
        if (posts === null) {
            allPostsLoaded = true;
            window.removeEventListener("scroll", handleScroll);
            if (offset === 0) {
                renderNoPosts();
            }
        } else {
            appendPosts(posts);
            offset += limit;
        }
    } catch (error) {

    }
    isLoading = false;
    if (!allPostsLoaded && document.body.offsetHeight < window.innerHeight) {
        loadMorePosts();
    }
}

// Toujours dans post_display.js

export function updateFilters(category, keyword, author = "") {
    currentCategory = category;
    currentKeyword = keyword;
    // Réinitialisation de la pagination et de l'affichage
    offset = 0;
    allPostsLoaded = false;
    window.addEventListener("scroll", handleScroll);
    const postsContainer = document.getElementById("postsContainer");
    if (postsContainer) postsContainer.innerHTML = "";
    // Charger les posts avec les nouveaux filtres
    loadMorePosts();
}

function renderNoPosts() {
    const postsContainer = document.getElementById("postsContainer");
    if (postsContainer) {
        postsContainer.innerHTML = "<p>Aucun post à afficher.</p>";
    }
}

function renderPosts(posts) {
    const postsContainer = document.getElementById("postsContainer");
    if (!postsContainer) return;

    // Si c'est le premier chargement, vider le conteneur
    if (offset === 0) {
        postsContainer.innerHTML = "";
    }

    // On trie le lot actuel (l'API doit renvoyer les posts déjà ordonnés globalement)
    posts.sort((a, b) => new Date(b.CreatedAt) - new Date(a.CreatedAt));

    posts.forEach(post => {
        const postElement = document.createElement("div");
        postElement.classList.add("post");

        postElement.innerHTML = `
      <h3>${post.Title}</h3>
      <p><em>Category : ${post.Category}</em></p>
      <p>${post.Content}</p>
      <p><small>Created on ${new Date(post.CreatedAt).toLocaleString()}</small></p>
    `;
        postsContainer.appendChild(postElement);
    });
}
