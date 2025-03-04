export function PostFrom(createPostButton) {
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
                <form id="postForm" action="posts" method="POST">
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
        }
    });
}