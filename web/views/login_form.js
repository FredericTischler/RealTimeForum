import { PostFrom } from './post_form.js'

document.getElementById("loginView").addEventListener("click", () => {
    const viewForm = document.getElementById("login");
    viewForm.style.display = "block";
    viewForm.innerHTML = `
    <div class="container">
      <h2>Login</h2>
      <form id="loginForm">
        <label for="identifier">User name</label>
        <input type="text" name="identifier" id="identifier" required>
        
        <label for="password">Password</label>
        <input type="password" name="password" id="password" required>
        
        <button type="submit">Login</button>
      </form>
    </div>
  `;

    const loginForm = document.getElementById("loginForm");
    loginForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const formData = new FormData(loginForm);
        try {
            const response = await fetch("login", {
                method: "POST",
                body: formData
            });
            if (!response.ok) {
                const errorText = await response.text();
                alert("Erreur: " + errorText);
            }
        } catch (err) {
            console.error("Erreur lors du login", err);
            alert("Le login a échoué, veuillez réessayer.");
        }
        updateUIAfterLogin();
    });
});

// Vous pouvez supprimer le code de vérification de connexion via localStorage,
// car la session est gérée via le cookie côté serveur.

function updateUIAfterLogin() {
    // Retirer complètement les boutons de connexion et d'inscription s'ils existent
    const loginBtn = document.getElementById("loginView");
    const registerBtn = document.getElementById("registerView");
    if (loginBtn) loginBtn.remove();
    if (registerBtn) registerBtn.remove();

    // Vider les containers de formulaire
    document.getElementById("login").innerHTML = "";
    document.getElementById("register").innerHTML = "";

    // Créer le bouton Logout
    const logoutButton = document.createElement("button");
    logoutButton.id = "logoutBtn";
    logoutButton.textContent = "Logout";
    logoutButton.addEventListener("click", () => {
        // Logique de déconnexion, par exemple appeler l'API logout
        // Puis recharger la page
        location.reload();
    });

    // Créer le bouton "Créer un post" sans ajouter d'écouteur ici
    const createPostButton = document.createElement("button");
    createPostButton.id = "createPostBtn";
    createPostButton.textContent = "Créer un post";
    PostFrom(createPostButton)
    // Aucune fonction n'est attachée ici ; le script post_form.js se chargera d'ajouter son propre écouteur

    // Créer ou récupérer le conteneur d'authentification dans le header
    let authContainer = document.getElementById("authContainer");
    if (!authContainer) {
        authContainer = document.createElement("div");
        authContainer.id = "authContainer";
        const header = document.querySelector("header");
        header.appendChild(authContainer);
    }

    // Vider le conteneur pour éviter les doublons
    authContainer.innerHTML = "";
    authContainer.appendChild(logoutButton);
    authContainer.appendChild(createPostButton);
}

