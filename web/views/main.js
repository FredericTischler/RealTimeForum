import { PostForm } from "./post_form.js";
import { displayPosts, updateFilters } from "./post_display.js";
import { displayLoginForm } from "./auth_form.js";

let ws
// Fonction qui interroge l'endpoint d'authentification
async function checkAuthStatus() {
    try {
        const response = await fetch("/auth/status", {
            credentials: "include"
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
    // Injection dynamique du menu déroulant des catégories dans le header
    const categoriesDiv = document.getElementById("categories");
    if (categoriesDiv) {
        categoriesDiv.innerHTML = `
      <select id="categoryMenu">
        <option value="">-- All categories --</option>
        <option value="Graffiti">Graffiti</option>
        <option value="Rap">Rap</option>
        <option value="Beatmaking">Beatmaking</option>
        <option value="Breakdance">Breakdance</option>
        <option value="Streetwear">Streetwear</option>
      </select>
    `;
    }

    // Injection dynamique de la search bar dans le header
    const searchbarDiv = document.getElementById("searchbar");
    if (searchbarDiv) {
        searchbarDiv.innerHTML = `
      <form id="searchForm">
        <input type="text" id="searchInput" placeholder="Search..." />
        <button type="submit">Search</button>
      </form>
    `;

        const searchForm = document.getElementById("searchForm");
        searchForm.addEventListener("submit", (e) => {
            e.preventDefault();
            // Récupération des valeurs
            const keyword = document.getElementById("searchInput").value.trim();
            const category = document.getElementById("categoryMenu").value;
            updateFilters(category, keyword);
        });
    }

    // Nettoyage des zones d'authentification
    const loginDiv = document.getElementById("login");
    if (loginDiv) loginDiv.innerHTML = "";
    const registerDiv = document.getElementById("register");
    if (registerDiv) registerDiv.innerHTML = "";

    // Création du bouton Logout
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
            const response = await fetch("/logout", { method: "POST", credentials: "include" });
            if (!response.ok) {
                alert("Erreur lors du logout");
                return;
            }
        } catch (err) {
            console.error("Erreur lors du logout", err);
        }
        location.reload();
    });

    // Création du bouton "Créer un post"
    const createPostButton = document.createElement("button");
    createPostButton.id = "createPostBtn";
    createPostButton.textContent = "+";
    document.body.appendChild(createPostButton);
    PostForm(createPostButton);

    // Injection du bouton Logout dans le header
    const authContainer = document.getElementById("authContainer");
    if (authContainer) {
        authContainer.innerHTML = "";
        authContainer.appendChild(logoutButton);
    }

    // Masquer la section d'authentification
    const authSection = document.getElementById("authSection");
    if (authSection) {
        authSection.style.display = "none";
    }

    const usersSidebar = document.getElementById("usersSidebar");
    if (usersSidebar) {
        usersSidebar.innerHTML = `
      <section id="userFilters" class="filters">
        <form id="userFiltersForm">
          <div class="filter-group">
            <label for="userGender">Gender :</label>
            <select id="userGender" name="gender">
              <option value="">Tous</option>
              <option value="male">Male</option>
              <option value="female">Female</option>
              <option value="other">Other</option>
            </select>
          </div>
          <div class="filter-group">
            <label for="userName">Username :</label>
            <input type="text" id="userName" name="username" placeholder="Search by username">
          </div>
          <button type="submit">Filter</button>
        </form>
      </section>
      <section id="usersList"></section>
    `;
    }

    // Affichage du contenu principal (les aside et la section postsSection restent dans le HTML)
    const contentContainer = document.querySelector(".content-container");
    if (contentContainer) {
        contentContainer.style.display = "flex";
    }

    // Appel à la fonction d'affichage des posts
    displayPosts();

    // Connexion WebSocket
    ws = new WebSocket("ws://localhost:8443/ws");

    ws.onopen = () => {
        console.log("Connexion WebSocket établie");
    };

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        console.log("Données reçues :", data.type); // Log pour vérifier les données
        
        if (data.type === "users") {
            const usersListSection = document.getElementById("usersList");
            if (usersListSection) {
                let html = "<ul>";
                data.users.forEach((user) => {
                    html += `<li><span class="user">${user.Username}</span></li>`;
                });
                html += "</ul>";
                usersListSection.innerHTML = html;
            }
        } else if (data.type === "message") {
            console.log(data.content)
            const messageElement = document.createElement("div");
            messageElement.textContent = `${data.sender}: ${data.content}`;
            messageContainer.appendChild(messageElement);
        }
    };

    ws.onerror = (error) => {
        console.error("Erreur WebSocket :", error);
    };

    ws.onclose = () => {
        console.log("Connexion WebSocket fermée");
    };
}

// Fonction pour charger les messages historiques
async function loadMessages(sender, recipient) {
    try {
        const response = await fetch(`/api/messages?sender=${sender}&recipient=${recipient}`);
        if (!response.ok) {
            throw new Error("Failed to fetch messages");
        }
        const messages = await response.json();

        // Afficher les messages dans le modal
        const messageContainer = document.getElementById("messageContainer");
        messageContainer.innerHTML = ""; // Vider les messages précédents

        messages.forEach((message) => {
            const messageElement = document.createElement("div");
            messageElement.textContent = `${message.sender}: ${message.content}`;
            messageContainer.appendChild(messageElement);
        });
    } catch (error) {
        console.error("Error loading messages:", error);
    }
}

function updateUIAfterLogout() {
    const contentContainer = document.getElementById("contentContainer");
    if (contentContainer) {
        contentContainer.style.display = "none";
    }
    const authSection = document.getElementById("authSection");
    if (authSection) {
        authSection.style.display = "block";
    }
    const usersSideBar = document.getElementById("usersSidebar")
    usersSideBar.style.display = "none";
    displayLoginForm(); // Affiche le formulaire de connexion
}

document.addEventListener("DOMContentLoaded", () => {
    checkAuthStatus();
});

// Sélectionnez le modal et les éléments associés
const messageModal = document.getElementById("messageModal");
const closeModal = document.querySelector("#messageModal .close");
const recipientUsernameElement = document.getElementById("recipientUsername");
const messageContainer = document.getElementById("messageContainer");
const messageInput = document.getElementById("messageInput");
const sendMessageButton = document.getElementById("sendMessageButton");

// Fonction pour ouvrir le modal
function openMessageModal(username) {
    recipientUsernameElement.textContent = username; // Afficher le nom de l'utilisateur
    messageModal.style.display = "block"; // Afficher le modal
    messageContainer.innerHTML = ""; // Vider les messages précédents
    messageInput.value = ""; // Vider l'input

    // Charger les messages historiques
    const currentUser = " "; // Remplacez par le nom de l'utilisateur actuel
    loadMessages(currentUser, username);
}

// Fonction pour fermer le modal
function closeMessageModal() {
    messageModal.style.display = "none";
}

// Ajouter un événement de clic sur les utilisateurs connectés
document.getElementById("usersSidebar").addEventListener("click", (event) => {
    if (event.target.classList.contains("user")) {
        const username = event.target.textContent;
        openMessageModal(username);
    }
});

// Fermer le modal lorsque l'utilisateur clique sur la croix
closeModal.addEventListener("click", closeMessageModal);

// Fermer le modal lorsque l'utilisateur clique en dehors du modal
window.addEventListener("click", (event) => {
    if (event.target === messageModal) {
        closeMessageModal();
    }
});

// Exemple d'envoi de message
sendMessageButton.addEventListener("click", () => {
    const message = messageInput.value.trim();
    if (message) {
        const recipient = recipientUsernameElement.textContent;
        const sender = " "; // Remplacez par le nom de l'utilisateur actuel

        // Envoyer le message via WebSocket
        const messageData = {
            type: "message",
            sender: sender,
            recipient: recipient,
            content: message,
        };
        ws.send(JSON.stringify(messageData));

        // Ajouter le message à l'affichage
        const messageElement = document.createElement("div");
        messageElement.textContent = `${sender}: ${message}`;
        messageContainer.appendChild(messageElement);

        // Vider l'input
        messageInput.value = "";
    }
});