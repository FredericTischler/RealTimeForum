// web/views/post_display.js

import { displayErrorModal } from './error_modal.js';

const postsContainer = document.getElementById("postsContainer");


export let currentCategory = "";
export let currentKeyword = "";

let limit = 2;
let offset = 0;
let isLoading = false;
let allPostsLoaded = false;

async function isAuthenticated() {
    try {
        const response = await fetch("/auth/status", {
            credentials: "include"
        });
        if (response.ok) {
            const data = await response.json();
            return data.authenticated;
        } else {
            return false;
        }
    } catch (error) {
        console.error("Erreur lors de la vérification d'authentification:", error);
        return false;
    }
}

// Fonction pour initialiser le chargement des posts et l'écoute du scroll
export async function displayPosts() {
    const userLoggedIn = await isAuthenticated();
    if (!userLoggedIn) {
        renderNoPosts();  // L'utilisateur n'est pas connecté, afficher message approprié
        return;
    }
    loadMorePosts();
    window.addEventListener("scroll", handleScroll);
}

async function loadMorePosts() {
    if (isLoading || allPostsLoaded) return;
    isLoading = true;
    
    const userLoggedIn = await isAuthenticated();
    if (!userLoggedIn) {
        renderNoPosts();  // L'utilisateur n'est pas connecté, afficher message approprié
        return;
    }
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
        postsContainer.remove()
    }


}

// Ajoute un lot de posts au conteneur existant
function appendPosts(posts) {
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
        postElement.addEventListener("click", () => {
            displayModal(post);
        });

        postElement.innerHTML = `
      <h3>${post.Post.Title}</h3>
      <p><em>Category : ${post.Post.Category}</em></p>
      <p>${post.Post.Content}</p>
      <p>${post.Username} <small>${new Date(post.Post.CreatedAt).toLocaleString()}</small></p>
    `;
        postsContainer.appendChild(postElement);
    });
}

let debounceTimer;

function handleScroll() {
    // On annule le timer précédent s'il existe
    clearTimeout(debounceTimer);

    // On démarre un nouveau timer avec un délai de 100ms
    debounceTimer = setTimeout(() => {
        // Vérifier si l'utilisateur est proche du bas de la page (à 200px près)
        if ((window.innerHeight + window.pageYOffset) >= document.body.offsetHeight - 200) {
            loadMorePosts();
        }
    }, 100);
}

// N'oubliez pas d'ajouter l'écouteur de scroll
window.addEventListener("scroll", handleScroll);

async function displayModal(post) {
    const modal = document.getElementById("postModal");
    const modalContent = document.getElementById("modalContent");
    const modalTitle = document.getElementById("modalTitle");

    modalContent.innerHTML = `
      <div id="modal-post">
          <h2>${post.Post.Title}</h2>      
          <p><strong>Category:</strong> ${post.Post.Category}</p>
          <p>${post.Post.Content}</p>
          <p>${post.Username} <small>${new Date(post.Post.CreatedAt).toLocaleString()}</small></p>
      </div>
      <h4>Comments:</h4>
      <div id="commentsContainer"></div> <!-- Section pour afficher les commentaires -->
      <form id="commentForm">
        <input type="text" name="comment" placeholder="New comment...">
        <br>
        <button class="authButton" style="margin-top: 10px;" type="submit">Add Comment</button>
      </form>
    `;

    // Charger les commentaires existants
    await fetchComments(post.Post.PostId);

    // Ecouter l'événement du formulaire pour ajouter un commentaire
    const commentForm = document.getElementById("commentForm");
    commentForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const commentContent = commentForm.comment.value;
        console.log(commentContent)
        if (commentContent.trim() === "") {
            alert("Le commentaire ne peut pas être vide !");
            return;
        }

        try {
            const response = await fetch(`/posts/comment/${post.Post.PostId}`, {
                method: "POST",
                credentials: "include",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    PostId: post.Post.PostId,
                    Content: commentContent
                })
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText);
            }

            commentForm.reset();
            await fetchComments(post.Post.PostId); // Rafraîchir les commentaires après ajout
        } catch (error) {
            displayErrorModal("Erreur lors de l'ajout du commentaire.");
        }
    });

    modal.style.display = "block";
}

// Fonction pour récupérer et afficher les commentaires
async function fetchComments(postId) {
    const commentsContainer = document.getElementById("commentsContainer");
    console.log("Fetching comments from:", `/posts/comment/${postId}`);

    try {
        const response = await fetch(`/posts/comment/${postId}`, {
            method: "GET",
            credentials: "include"
        });

        console.log(response)

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText);
        }

        const comments = await response.json();
        console.log("Commentaires récupérés:", comments);

        commentsContainer.innerHTML = "";

        if (comments === null) {
            commentsContainer.innerHTML = "<p>Aucun commentaire pour le moment.</p>";
        } else {
            comments.forEach(comment => {
                const commentElement = document.createElement("div");
                commentElement.classList.add("comment");
                commentElement.innerHTML = `
                    <h2>${comment.Username}</h2>
                    <p> ${comment.Comment.Content}</p>
                    <p><small>${new Date(comment.Comment.CreatedAt).toLocaleString()}</small></p>
                `;
                commentsContainer.appendChild(commentElement);
            });
        }
    } catch (error) {
        console.error("Erreur lors de la récupération des commentaires:", error);
        displayErrorModal("Erreur lors de la récupération des commentaires.");
    }
}


const postModal = document.getElementById("postModal");

postModal.addEventListener("click", (event) => {
    // Si le clic est sur le conteneur du modal lui-même, fermer le modal
    if (event.target === postModal) {
        postModal.style.display = "none";
    }
});

