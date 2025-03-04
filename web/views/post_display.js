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

    // Pour chaque post, créer un élément et l'ajouter au conteneur
    posts.forEach(post => {
        console.log(post)
        const postElement = document.createElement("div");
        postElement.classList.add("post");

        postElement.innerHTML = `
      <h3>${post.Title}</h3>
      <p>${post.Content}</p>
      <p><em>Catégorie : ${post.Category}</em></p>
      <p><small>Créé le ${new Date(post.CreatedAt).toLocaleString()}</small></p>
    `;
        postsContainer.appendChild(postElement);
    });
}
