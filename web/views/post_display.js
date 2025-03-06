// web/views/post_display.js

import { displayErrorModal } from './error_modal.js';

let limit = 2;
let offset = 0;
let isLoading = false;
let allPostsLoaded = false;

// Fonction pour initialiser le chargement des posts et l'écoute du scroll
export function displayPosts() {
    loadMorePosts();
    window.addEventListener("scroll", handleScroll);
}

// Fonction de chargement d'un lot de posts
async function loadMorePosts() {
    if (isLoading || allPostsLoaded) return;
    isLoading = true;
    try {
        const response = await fetch(`/posts?limit=${limit}&offset=${offset}`, {
            method: "GET",
            credentials: "include"
        });
        // Si la réponse n'est pas OK, on vérifie si c'est parce qu'il n'y a plus de posts
        if (!response.ok) {
            // Si c'est un 404 (ou un autre code que vous utilisez pour signaler l'absence de posts), on considère que c'est la fin
            if (response.status === 404) {
                allPostsLoaded = true;
                return;
            }
            const errorText = await response.text();
            throw new Error(errorText);
        }
        const posts = await response.json();
        if (posts.length === 0) {
            allPostsLoaded = true;
            if (offset === 0) {
                renderNoPosts();
            }
        } else {
            appendPosts(posts);
            offset += limit;
        }
    } catch (error) {
        // On n'affiche l'erreur que si ce n'est pas simplement la fin des posts
        if (offset === 0) {
            displayErrorModal("Erreur lors de la récupération des posts.");
        }
    }
    isLoading = false;

    // Vérifie si le contenu ne remplit pas encore la fenêtre et charge plus de posts si nécessaire
    if (!allPostsLoaded && document.body.offsetHeight < window.innerHeight) {
        loadMorePosts();
    }
}

function renderNoPosts() {
    const postsContainer = document.getElementById("postsContainer");
    if (postsContainer) {
        postsContainer.innerHTML = "<p>Aucun post à afficher.</p>";
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

        postElement.innerHTML = `
      <h3>${post.Title}</h3>
      <p><em>Category : ${post.Category}</em></p>
      <p>${post.Content}</p>
      <p><small>Created on ${new Date(post.CreatedAt).toLocaleString()}</small></p>
    `;
        postsContainer.appendChild(postElement);
    });
}

// Détecte quand l'utilisateur approche du bas de page pour charger plus de posts
function handleScroll() {
    // Ne déclenche loadMorePosts que si l'utilisateur a effectivement scrollé (pageYOffset > 0)
    if (window.pageYOffset === 0) return;
    if ((window.innerHeight + window.pageYOffset) >= document.body.offsetHeight - 200) {
        loadMorePosts();
    }
}

// Fonction qui crée dynamiquement le formulaire de filtres
export function createFilterForm() {
    const postFilters = document.getElementById("postFilters");
    if (!postFilters) return;
    postFilters.innerHTML = `
    <form id="filtersForm">
      <div class="filter-group">
        <label for="filterCategory">Catégorie :</label>
        <select id="filterCategory" name="category">
          <option value="">--Toutes--</option>
          <option value="News">News</option>
          <option value="Sport">Sport</option>
          <option value="Culture">Culture</option>
          <!-- Ajoutez d'autres catégories selon vos besoins -->
        </select>
      </div>
      <div class="filter-group">
        <label for="filterKeyword">Mot-clé :</label>
        <input type="text" id="filterKeyword" name="keyword" placeholder="Rechercher...">
      </div>
      <div class="filter-group">
        <label for="filterAuthor">Auteur :</label>
        <input type="text" id="filterAuthor" name="author" placeholder="Nom d'utilisateur">
      </div>
      <button type="submit">Filtrer</button>
    </form>
  `;

    // Attacher l'écouteur d'événement pour gérer la soumission des filtres
    const filtersForm = document.getElementById("filtersForm");
    filtersForm.addEventListener("submit", (e) => {
        e.preventDefault();
        const category = document.getElementById("filterCategory").value;
        const keyword = document.getElementById("filterKeyword").value;
        const author = document.getElementById("filterAuthor").value;

        // Ici, vous pouvez rafraîchir l'affichage des posts en fonction des filtres sélectionnés
        // Par exemple, réinitialiser la pagination et appeler une fonction reloadPosts() :
        offset = 0;
        allPostsLoaded = false;
        const postsContainer = document.getElementById("postsContainer");
        if (postsContainer) postsContainer.innerHTML = "";

        // On construit l'URL avec les nouveaux paramètres de filtre
        const queryParams = new URLSearchParams({
            limit: limit,
            offset: offset,
            category: category,
            keyword: keyword,
            author: author
        });

        fetch(`/posts?${queryParams.toString()}`, {
            method: "GET",
            credentials: "include"
        })
            .then(response => response.json())
            .then(posts => {
                if (posts.length === 0) {
                    renderNoPosts();
                    allPostsLoaded = true;
                } else {
                    appendPosts(posts);
                    offset += limit;
                }
            })
            .catch(err => {
                console.error("Erreur lors du filtrage des posts", err);
            });
    });
}

// Par exemple, dans updateUIAfterLogin(), après avoir affiché les boutons, appelez createFilterForm():
function updateUIAfterLogin() {
    // ... votre code existant pour le logout, etc.
    // Afficher le contenu principal (posts, etc.)
    const contentContainer = document.querySelector(".content-container");
    if (contentContainer) {
        contentContainer.style.display = "block";
    }

    // Créer dynamiquement les filtres
    createFilterForm();

    // Lancer l'affichage initial des posts
    displayPosts();
}

