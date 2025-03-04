let postFormVisible = false;

document.getElementById("createPostBtn").addEventListener("click", () => {
    const postFormContainer = document.getElementById("post-form-container");
    // Bascule l'affichage du formulaire
    postFormVisible = !postFormVisible;
    postFormContainer.style.display = postFormVisible ? "block" : "none";

    if (postFormVisible) {
        postFormContainer.innerHTML = `
            <div class="container">
                <h2>Create / Edit Post</h2>
                <form id="postForm" action="/api/posts" method="POST">
                    <label for="title">Titre</label>
                    <input type="text" name="title" id="title" required>
                    
                    <label for="content">Contenu</label>
                    <textarea name="content" id="content" required></textarea>
                    
                    <label for="category">Catégorie</label>
                    <input type="text" name="category" id="category">
                    
                    <!-- Champ caché pour l'édition -->
                    <input type="hidden" name="postId" id="postId" value="">
                    
                    <button type="submit">Submit</button>
                </form>
            </div>
        `;

        // Ajoute l'écouteur de soumission du formulaire
        document.getElementById("postForm").addEventListener("submit", async (event) => {
            event.preventDefault();
            const form = event.target;
            const formData = new FormData(form);
            const postData = Object.fromEntries(formData.entries());
            const isEditMode = postData.postId !== "";
            const url = isEditMode ? `/api/posts/${postData.postId}` : '/api/posts';
            const method = isEditMode ? "PUT" : "POST";

            try {
                const response = await fetch(url, {
                    method: method,
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify(postData)
                });

                if (!response.ok) {
                    const errorText = await response.text();
                    alert(`Erreur: ${errorText}`);
                    return;
                }

                const result = await response.json();
                alert(isEditMode ? "Post mis à jour avec succès!" : "Post créé avec succès!");
                // Vous pouvez ici appeler une fonction pour rafraîchir dynamiquement la liste des posts
                console.log("Résultat:", result);
            } catch (error) {
                console.error("Erreur lors de la soumission du formulaire:", error);
                alert("Une erreur est survenue lors de la soumission du formulaire.");
            }
        });
    }
});
