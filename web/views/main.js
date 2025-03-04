import { PostForm } from "./post_form.js";

// Fonction qui interroge l'endpoint d'authentification
async function checkAuthStatus() {
    try {
        // Le cookie HTTPOnly "session_token" est automatiquement envoyé
        const response = await fetch('/auth/status', {
            credentials: 'include'
        });
        if (response.ok) {
            const data = await response.json();
            if (data.authenticated) {
                updateUIAfterLogin();
            } else {
                updateUIAfterLogout();
            }
        } else {
            updateUIAfterLogout();
        }
    } catch (error) {
        console.error("Erreur lors de la vérification d'authentification:", error);
        updateUIAfterLogout();
    }
}

// Mise à jour de l'interface pour un utilisateur connecté
function updateUIAfterLogin() {
    // Supprimer les boutons "Login" et "Register"
    const loginRegisterContainer = document.getElementById("login-register-container");
    loginRegisterContainer.innerHTML = "";

    // Vider les conteneurs de formulaires
    document.getElementById("login").innerHTML = "";
    document.getElementById("register").innerHTML = "";

    // Créer le bouton Logout
    const logoutButton = document.createElement("button");
    logoutButton.id = "logoutBtn";
    logoutButton.textContent = "Logout";
    logoutButton.addEventListener("click", async () => {
        // Appeler l'API logout pour détruire la session côté serveur
        try {
            const response = await fetch('/logout', { method: 'POST', credentials: 'include' });
            if (!response.ok) {
                alert("Erreur lors du logout");
                return;
            }
        } catch (err) {
            console.error("Erreur lors du logout", err);
        }
        location.reload();
    });

    // Créer le bouton "Créer un post"
    const createPostButton = document.createElement("button");
    createPostButton.id = "createPostBtn";
    createPostButton.textContent = "Créer un post";
    PostForm(createPostButton)

    // Créer ou récupérer le conteneur d'authentification dans le header
    let authContainer = document.getElementById("authContainer");
    authContainer.innerHTML = "";
    authContainer.appendChild(logoutButton);
    authContainer.appendChild(createPostButton);
}

// Mise à jour de l'interface pour un utilisateur non authentifié
function updateUIAfterLogout() {
    // Assurez-vous que le conteneur login-register contient bien les boutons de connexion
    const loginRegisterContainer = document.getElementById("login-register-container");
    // Si les boutons ne sont pas présents, vous pouvez les recréer (selon votre logique)
    if (!document.getElementById("loginView") && !document.getElementById("registerView")) {
        loginRegisterContainer.innerHTML = `
      <button id="loginView">Login</button>
      <button id="registerView">Register</button>
    `;
        // Vous pouvez ré-attacher les écouteurs si nécessaire ou recharger la page.
    }
}

// Au chargement du DOM, on vérifie l'état d'authentification
document.addEventListener("DOMContentLoaded", () => {
    checkAuthStatus();
});
