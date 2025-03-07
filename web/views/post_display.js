// web/views/post_display.js

import { displayErrorModal } from './error_modal.js';

const postsContainer = document.getElementById("postsContainer");


// Fonction pour récupérer et afficher les posts
export async function displayPosts() {
    try {
        const response = await fetch("/posts", {
            method: "GET",
            credentials: "include"
        });
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText);
        }
        const posts = await response.json();
        renderPosts(posts);
    } catch (error) {
        displayErrorModal("Erreur lors de la récupération des posts.");
    }
}

function renderPosts(posts) {
    if (!postsContainer) return;

    postsContainer.innerHTML = "";

    if (posts.length === 0) {
        postsContainer.innerHTML = "<p>Aucun post à afficher.</p>";
        return;
    }

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
    postElement.addEventListener("click", () => {
        displayModal(post);
    });

    postsContainer.appendChild(postElement);
    
    });

   
}

async function displayModal(post) {
    const modal = document.getElementById("postModal");
    const modalContent = document.getElementById("modalContent");
    const modalTitle = document.getElementById("modalTitle");

    modalTitle.textContent = post.Title;
    modalContent.innerHTML = `
      <p><strong>Category:</strong> ${post.Category}</p>
      <p>${post.Content}</p>
      <p><small>Created on ${new Date(post.CreatedAt).toLocaleString()}</small></p>
      <h4>Comments:</h4>
      <div id="commentsContainer"></div> <!-- Section pour afficher les commentaires -->
      <form id="commentForm">
        <input type="text" name="comment" placeholder="New comment...">
        <br>
        <button class="authButton" style="margin-top: 10px;" type="submit">Add Comment</button>
      </form>
    `;

    // Charger les commentaires existants
    await fetchComments(post.PostId);

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
            const response = await fetch(`/posts/comment/${post.PostId}`, {
                method: "POST",
                credentials: "include",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    PostId: post.PostId,
                    Content: commentContent
                })
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText);
            }

            commentForm.reset();
            await fetchComments(post.PostId); // Rafraîchir les commentaires après ajout
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
                    <p> ${comment.Comment.Content}</p>
                    <p>Posted By: ${comment.Username}<small>${new Date(comment.Comment.CreatedAt).toLocaleString()}</small></p>
                `;
                commentsContainer.appendChild(commentElement);
            });
        }
    } catch (error) {
        console.error("Erreur lors de la récupération des commentaires:", error);
        displayErrorModal("Erreur lors de la récupération des commentaires.");
    }
}


document.getElementById("closeModal").addEventListener("click", () => {
    document.getElementById("postModal").style.display = "none";
});