export function PostForm(createPostButton) {
    let postFormVisible = false;

    createPostButton.addEventListener("click", () => {
        const postFormContainer = document.getElementById("post-form-container");
        // Bascule l'affichage du formulaire
        postFormVisible = !postFormVisible;
        postFormContainer.style.display = postFormVisible ? "block" : "none";

        if (postFormVisible) {
            postFormContainer.innerHTML = `
        <div class="container">
          <h2>Create / Edit Post</h2>
          <form id="postForm">
            <label for="title">Titre</label>
            <input type="text" name="title" id="title" required>
            
            <label for="content">Contenu</label>
            <textarea name="content" id="content" required></textarea>
            
            <label for="category">Catégorie</label>
            <input type="text" name="category" id="category">
            
            <button type="submit">Submit</button>
          </form>
        </div>
      `;

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
                        credentials: "include", // Assure l'envoi des cookies
                        body: JSON.stringify(postData)
                    });

                    if (!response.ok) {
                        const errorText = await response.text();
                        alert(`Erreur: ${errorText}`);
                        return;
                    }

                    // Ici, on peut mettre à jour dynamiquement l'interface ou rediriger en JavaScript
                    alert("Post créé avec succès!");
                    // Par exemple, rediriger via JS vers la page d'accueil :
                    window.location.href = "/";
                } catch (error) {
                    console.error("Erreur lors de la soumission du formulaire:", error);
                    alert("Une erreur est survenue lors de la soumission du formulaire.");
                }
            });
        }
    });
}
