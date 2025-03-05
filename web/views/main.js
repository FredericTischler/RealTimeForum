import { PostForm } from "./post_form.js";
import { displayPosts } from './post_display.js';
import { displayLoginForm } from './auth_form.js'; // Import ajouté

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
    logoutButton.textContent = "Logout";
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
    createPostButton.textContent = "Créer un post";
    PostForm(createPostButton);

    // Utiliser le conteneur authentification dans le header
    let authContainer = document.getElementById("authContainer");
    if (authContainer) {
        authContainer.innerHTML = "";
        authContainer.appendChild(logoutButton);
        authContainer.appendChild(createPostButton);
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
