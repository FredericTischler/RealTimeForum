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
                return;
            }
            const data = await response.json();
            // Stocker les informations de connexion
            localStorage.setItem("token", data.Token);
            localStorage.setItem("userId", data.UserId);
            // Mettre à jour l'interface
            updateUIAfterLogin();
        } catch (err) {
            console.error("Erreur lors du login", err);
            alert("Le login a échoué, veuillez réessayer.");
        }
    });
});

// Vérifier l'état connecté au chargement de la page
document.addEventListener("DOMContentLoaded", () => {
    if (localStorage.getItem("token")) {
        updateUIAfterLogin();
    }
});

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
        localStorage.removeItem("token");
        localStorage.removeItem("userId");
        location.reload();
    });

    // Créer le bouton "Créer un post"
    const createPostButton = document.createElement("button");
    createPostButton.id = "createPostBtn";
    createPostButton.textContent = "Créer un post";
    createPostButton.addEventListener("click", () => {
        // Implémentez ici la logique pour afficher l'interface de création de post
        alert("Interface de création de post à implémenter.");
    });

    // Pour éviter d'empiler plusieurs boutons si updateUIAfterLogin est appelée plusieurs fois,
    // on vérifie si un conteneur dédié existe déjà sinon on le crée.
    let authContainer = document.getElementById("authContainer");
    if (!authContainer) {
        authContainer = document.createElement("div");
        authContainer.id = "authContainer";
        // Ajoutons ce conteneur à l'intérieur du header pour garder la cohérence visuelle
        const header = document.querySelector("header");
        header.appendChild(authContainer);
    }

    // On vide le conteneur pour s'assurer qu'il n'y a pas de doublons
    authContainer.innerHTML = "";
    authContainer.appendChild(logoutButton);
    authContainer.appendChild(createPostButton);
}

