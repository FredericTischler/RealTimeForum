import { PostForm } from "./post_form.js";
import { displayPosts } from './post_display.js';
import { displayLoginForm } from './auth_form.js';
import { createFilterForm} from "./post_display.js";

// Fonction qui interroge l'endpoint d'authentification
async function checkAuthStatus() {
    try {
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

function updateUIAfterLogin() {
    // Vérifier et nettoyer le conteneur "login-register-container" s'il existe
    const loginRegisterContainer = document.getElementById("login-register-container");
    if (loginRegisterContainer) {
        loginRegisterContainer.innerHTML = "";
    }

    // Nettoyer les formulaires d'authentification
    const loginDiv = document.getElementById("login");
    if (loginDiv) {
        loginDiv.innerHTML = "";
    }
    const registerDiv = document.getElementById("register");
    if (registerDiv) {
        registerDiv.innerHTML = "";
    }

    // Mise à jour de l'interface pour un utilisateur connecté : injection des boutons Logout et "Créer un post"
    const logoutButton = document.createElement("button");
    logoutButton.id = "logoutBtn";
    logoutButton.innerHTML = `
<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" fill="#BB86FC" class="bi bi-box-arrow-in-right" viewBox="0 0 16 16">
  <path fill-rule="evenodd" d="M6 3.5a.5.5 0 0 1 .5-.5h8a.5.5 0 0 1 .5.5v9a.5.5 0 0 1-.5.5h-8a.5.5 0 0 1-.5-.5v-2a.5.5 0 0 0-1 0v2A1.5 1.5 0 0 0 6.5 14h8a1.5 1.5 0 0 0 1.5-1.5v-9A1.5 1.5 0 0 0 14.5 2h-8A1.5 1.5 0 0 0 5 3.5v2a.5.5 0 0 0 1 0z"/>
  <path fill-rule="evenodd" d="M11.854 8.354a.5.5 0 0 0 0-.708l-3-3a.5.5 0 1 0-.708.708L10.293 7.5H1.5a.5.5 0 0 0 0 1h8.793l-2.147 2.146a.5.5 0 0 0 .708.708z"/>
</svg>
`;
    logoutButton.addEventListener("click", async () => {
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

    const createPostButton = document.createElement("button");
    createPostButton.id = "createPostBtn";
    createPostButton.textContent = "+";
    document.body.appendChild(createPostButton);
    PostForm(createPostButton);

    // Utiliser le conteneur authentification dans le header
    let authContainer = document.getElementById("authContainer");
    if (authContainer) {
        authContainer.innerHTML = "";
        authContainer.appendChild(logoutButton);
    }

    // Masquer la section d'authentification puisque l'utilisateur est connecté
    const authSection = document.getElementById("authSection");
    if (authSection) {
        authSection.style.display = "none";
    }

    // Afficher le contenu principal (posts, etc.)
    const contentContainer = document.querySelector(".content-container");
    if (contentContainer) {
        contentContainer.style.display = "block";
    }
    createFilterForm()
    displayPosts()
}


function updateUIAfterLogout() {
    const contentContainer = document.querySelector(".content-container");
    if(contentContainer) {
        contentContainer.style.display = "none";
    }
    const authSection = document.getElementById("authSection");
    if(authSection) {
        authSection.style.display = "block";
    }
    displayLoginForm(); // Affiche le formulaire de connexion
}

document.addEventListener("DOMContentLoaded", () => {
    checkAuthStatus();
});
