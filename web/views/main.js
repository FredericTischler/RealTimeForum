import { PostForm } from "./post_form.js";
import { displayPosts, updateFilters } from "./post_display.js";
import { displayLoginForm } from "./auth_form.js";

// Pour stocker les IDs des utilisateurs en ligne
let onlineUserIds = new Set();
// Pour conserver les références aux éléments DOM de chaque utilisateur
let userElements = {};

// Fonction d'authentification (inchangée)
async function checkAuthStatus() {
    try {
        const response = await fetch("/auth/status", { credentials: "include" });
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
    // Injection dynamique du menu des catégories et de la search bar
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

    // Configuration de la sidebar des utilisateurs
    const usersSidebar = document.getElementById("usersSidebar");
    if (usersSidebar) {
        usersSidebar.innerHTML = `
      <section id="userFilters" class="filters">
        <form id="userFiltersForm">
          <div class="filter-group">
            <label for="userGender">Gender :</label>
            <select id="userGender" name="gender">
              <option value="">Tous</option>
              <option value="Male">Male</option>
              <option value="Female">Female</option>
              <option value="Other">Other</option>
            </select>
          </div>
          <div class="filter-group">
            <label for="userAge">Age :</label>
            <input type="number" id="userAge" name="userage" placeholder="Search by age">
            <label for="userName">Username :</label>
            <input type="text" id="userName" name="username" placeholder="Search by username">
          </div>
          <button type="submit">Filter</button>
        </form>
      </section>
      <section id="usersList"></section>
    `;
    }

    const userFiltersForm = document.getElementById("userFiltersForm");
    if (userFiltersForm) {
        userFiltersForm.addEventListener("submit", (e) => {
            e.preventDefault();
            const username = document.getElementById("userName").value.trim();
            const gender = document.getElementById("userGender").value;
            const age = document.getElementById("userAge").value;
            filterUsers(username, gender, age);
        });
    }

    // Affichage du contenu principal
    const contentContainer = document.querySelector(".content-container");
    if (contentContainer) {
        contentContainer.style.display = "flex";
    }

    displayPosts();
    fetchUsers().then(() => {
        setupWebSocket();
    });
}

function filterUsers(username, gender, age) {
    const userCards = Array.from(document.querySelectorAll(".user-card"));
    userCards.forEach(card => {
        const cardUsername = card.dataset.username.toLowerCase();
        const cardGender = card.dataset.gender.toLowerCase();
        const cardAge = card.dataset.age; // On suppose que l'age est stocké en chaîne

        const matchesUsername = username === "" || cardUsername.includes(username.toLowerCase());
        const matchesGender = gender === "" || cardGender === gender.toLowerCase();
        const matchesAge = age === "" || cardAge === age;
        card.style.display = (matchesUsername && matchesGender && matchesAge) ? "flex" : "none";
    });
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
    const usersSideBar = document.getElementById("usersSidebar");
    if (usersSideBar) {
        usersSideBar.style.display = "none";
    }
    displayLoginForm();
}

document.addEventListener("DOMContentLoaded", () => {
    checkAuthStatus();
});

async function fetchUsers() {
    try {
        const response = await fetch("/users", { credentials: "include" });
        if (!response.ok) {
            throw new Error("Erreur lors de la récupération des utilisateurs");
        }
        const users = await response.json();
        updateUsersList(users);
    } catch (error) {
        console.error(error);
    }
}

function updateUsersList(users) {
    const userListContainer = document.getElementById("usersList");
    if (!userListContainer) return;

    userListContainer.innerHTML = "";
    userElements = {}; // réinitialiser les références

    users.forEach(user => {
        const userCard = document.createElement("div");
        userCard.classList.add("user-card");

        userCard.dataset.username = user.Username;
        userCard.dataset.gender = user.Gender;
        userCard.dataset.age = user.Age;

        const avatar = document.createElement("div");
        avatar.classList.add("user-avatar");
        avatar.textContent = user.Username.charAt(0).toUpperCase();

        const userInfo = document.createElement("div");
        userInfo.classList.add("user-info");
        // Utilise onlineUserIds pour déterminer le statut initial
        const isOnline = onlineUserIds.has(user.UserId);
        userInfo.innerHTML = `<strong>${user.Username}</strong><br><span class="status">${isOnline ? "Online" : "Offline"}</span>`;

        const statusIndicator = document.createElement("div");
        statusIndicator.classList.add("status-indicator", isOnline ? "online" : "offline");

        userCard.appendChild(avatar);
        userCard.appendChild(userInfo);
        userCard.appendChild(statusIndicator);
        userListContainer.appendChild(userCard);

        // Stocke les références dans userElements pour mises à jour ultérieures
        userElements[user.UserId] = {
            statusText: userInfo.querySelector(".status"),
            indicator: statusIndicator
        };

    });
}

function setupWebSocket() {
    const ws = new WebSocket("ws://localhost:8443/ws");

    ws.onmessage = (event) => {
        // Supposons que le serveur envoie un tableau d'IDs (strings) des utilisateurs en ligne
        const data = JSON.parse(event.data);
        onlineUserIds = new Set(data.map(item => item.UserId));
        // Met à jour le DOM pour chaque utilisateur présent dans userElements
        for (const userId in userElements) {
            if (onlineUserIds.has(userId)) {
                userElements[userId].statusText.textContent = "Online";
                userElements[userId].indicator.classList.remove("offline");
                userElements[userId].indicator.classList.add("online");
            } else {
                userElements[userId].statusText.textContent = "Offline";
                userElements[userId].indicator.classList.remove("online");
                userElements[userId].indicator.classList.add("offline");
            }
        }
    };

    ws.onopen = () => console.log("WebSocket connecté");
    ws.onerror = (error) => console.error("WebSocket erreur :", error);
    ws.onclose = () => console.log("WebSocket fermé");
}
