// web/views/post_form.js

import { displayErrorModal } from './error_modal.js';

export function PostForm(createPostButton) {
    let postFormVisible = false;

    createPostButton.addEventListener("click", () => {
        const postFormContainer = document.getElementById("post-form-container");
        // Bascule l'affichage du formulaire
        postFormVisible = !postFormVisible;
        postFormContainer.style.display = postFormVisible ? "flex" : "none";

        if (postFormVisible) {
            postFormContainer.innerHTML = `
        <div class="container">
          <h2>Create / Edit Post</h2>
          <form id="postForm">
            <label for="title">Titre</label>
            <input type="text" name="title" id="title" required>
            
            <label for="category">Catégorie</label>
            <input type="text" name="category" id="category">
            
            <label for="content">Contenu</label>
            <textarea name="content" id="content" required></textarea>
           
            <button type="submit">Submit</button>
          </form>
        </div>
      `;

            // Fermer le modal si l'utilisateur clique à l'extérieur de la boîte du formulaire
            postFormContainer.addEventListener("click", function(e) {
                if (e.target === postFormContainer) {
                    postFormContainer.style.display = "none";
                    postFormVisible = false;
                }
            });

            const postForm = document.getElementById("postForm");
            postForm.addEventListener("submit", async (event) => {
                event.preventDefault();
                const formData = new FormData(postForm);
                const postData = Object.fromEntries(formData.entries());

                try {
                    const response = await fetch("/posts", {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json"
                        },
                        credentials: "include",
                        body: JSON.stringify(postData)
                    });

                    if (!response.ok) {
                        const errorText = await response.text();
                        displayErrorModal(`Erreur: ${errorText}`);
                        return;
                    }
                    window.location.href = "/";
                } catch (error) {
                    console.error("Erreur lors de la soumission du formulaire:", error);
                    displayErrorModal("Une erreur est survenue lors de la soumission du formulaire.");
                }
            });
        }
    });
}
