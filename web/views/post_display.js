// web/views/post_display.js

// Fonction pour récupérer et afficher les posts
export async function displayPosts() {
    try {
        const response = await fetch("/posts", {
            method: "GET",
            credentials: "include"  // Pour envoyer les cookies si nécessaire
        });
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText);
        }
        const posts = await response.json();
        renderPosts(posts);
    } catch (error) {
        console.error("Erreur lors de la récupération des posts :", error);
    }
}

// Fonction qui prend un tableau de posts et met à jour le DOM
function renderPosts(posts) {
    const postsContainer = document.getElementById("postsContainer");
    if (!postsContainer) return;

    // Vider le conteneur
    postsContainer.innerHTML = "";

    if (posts.length === 0) {
        postsContainer.innerHTML = "<p>Aucun post à afficher.</p>";
        return;
    }

    // Trier les posts par date décroissante (les plus récents en premier)
    posts.sort((a, b) => new Date(b.CreatedAt) - new Date(a.CreatedAt));

    // Pour chaque post, créer un élément et l'ajouter au conteneur
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
