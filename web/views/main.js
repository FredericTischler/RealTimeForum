import { PostForm } from "./post_form.js";
import { displayPosts, updateFilters } from "./post_display.js";
import { displayLoginForm } from "./auth_form.js";
import { startPrivateChat } from "./message_form.js";

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

async function updateUIAfterLogin() {
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
    await setupWebSocket();
    await fetchUsers();

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
        // Appeler updateUsersList directement avec les utilisateurs récupérés
        await updateUsersList(users);
    } catch (error) {
        console.error(error);
    }
}

async function updateUsersList(users) {
    const userListContainer = document.getElementById("usersList");
    if (!userListContainer) return;

    userListContainer.innerHTML = "";
    userElements = {}; // réinitialiser les références

    try {
        const authResponse = await fetch("/auth/status", { credentials: "include" });
        if (!authResponse.ok) return;
        
        const authData = await authResponse.json();
        if (!authData.authenticated) return;
        
        const currentUserID = authData.userId;
        
        // Trier les utilisateurs selon les critères
        const sortedUsers = await sortUsers(users, currentUserID);
        
        // Afficher les utilisateurs triés
        sortedUsers.forEach(user => {
            if (user.UserId !== currentUserID) {
                const userCard = createUserCard(user, currentUserID);
                userListContainer.appendChild(userCard);
                
                // Stocke les références dans userElements pour mises à jour ultérieures
                userElements[user.UserId] = {
                    statusText: userCard.querySelector(".status"),
                    indicator: userCard.querySelector(".status-indicator")
                };
            }
        });
    } catch (error) {
        console.error("Erreur lors de la récupération des données:", error);
    }
}

async function sortUsers(users, currentUserId) {
    try {
        // Récupérer les derniers messages pour chaque conversation
        const lastMessages = await Promise.all(
            users.map(async user => {
                if (user.UserId === currentUserId) return null;
                
                try {
                    const response = await fetch(`/message?with=${user.UserId}&offset=0&limit=1`, {
                        credentials: "include"
                    });
                    
                    if (!response.ok) {
                        console.warn(`Erreur lors de la récupération des messages pour l'utilisateur ${user.UserId}`);
                        return { user, lastMessage: null };
                    }
                    
                    const messages = await response.json();
                    
                    if (!Array.isArray(messages)) {
                        console.warn(`Réponse inattendue pour les messages de l'utilisateur ${user.UserId}`, messages);
                        return { user, lastMessage: null };
                    }
                    
                    return {
                        user,
                        lastMessage: messages.length > 0 ? messages[0] : null,
                        isOnline: onlineUserIds.has(user.UserId) // Ajouter le statut en ligne
                    };
                } catch (error) {
                    console.error(`Erreur lors du traitement des messages pour ${user.UserId}:`, error);
                    return { user, lastMessage: null, isOnline: false };
                }
            })
        );
        
        // Filtrer les résultats null (pour l'utilisateur courant)
        const filtered = lastMessages.filter(item => item !== null);
        
        // Séparer les utilisateurs connectés et non connectés
        const onlineUsers = filtered.filter(item => item.isOnline);
        const offlineUsers = filtered.filter(item => !item.isOnline);
        
        // Trier les utilisateurs connectés par date du dernier message (plus récent d'abord)
        onlineUsers.sort((a, b) => {
            // Si les deux ont des messages, trier par date
            if (a.lastMessage && b.lastMessage) {
                return new Date(b.lastMessage.SentAt) - new Date(a.lastMessage.SentAt);
            }
            // Si seul a a un message, il passe avant
            if (a.lastMessage) return -1;
            // Si seul b a un message, il passe avant
            if (b.lastMessage) return 1;
            // Sinon trier par nom d'utilisateur
            return a.user.Username.localeCompare(b.user.Username);
        });
        
        // Trier les utilisateurs non connectés par nom d'utilisateur
        offlineUsers.sort((a, b) => a.user.Username.localeCompare(b.user.Username));
        
        // Combiner les deux listes (connectés d'abord, puis non connectés)
        return [...onlineUsers.map(item => item.user), ...offlineUsers.map(item => item.user)];
    } catch (error) {
        console.error("Erreur lors du tri des utilisateurs:", error);
        // En cas d'erreur, retourner simplement la liste des utilisateurs (sans l'utilisateur courant)
        return users.filter(user => user.UserId !== currentUserId);
    }
}

function createUserCard(user, currentUserId) {
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
    const isOnline = onlineUserIds.has(user.UserId);
    userInfo.innerHTML = `<strong>${user.Username}</strong><br><span class="status">${isOnline ? "Online" : "Offline"}</span>`;
    
    const statusIndicator = document.createElement("div");
    statusIndicator.classList.add("status-indicator", isOnline ? "online" : "offline");
    
    userCard.appendChild(avatar);
    userCard.appendChild(userInfo);
    userCard.appendChild(statusIndicator);
    
    userCard.addEventListener("click", () => {
        startPrivateChat(currentUserId, user.UserId, user.Username);
    });
    console.log(userCard)
    return userCard;
}

function setupWebSocket() {
    return new Promise((resolve) => {
        const ws = new WebSocket("ws://localhost:8443/ws");

        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            onlineUserIds = new Set(data.map(item => item.UserId));
            
            // Mettre à jour les indicateurs de statut
            for (const userId in userElements) {
                const isOnline = onlineUserIds.has(userId);
                userElements[userId].statusText.textContent = isOnline ? "Online" : "Offline";
                userElements[userId].indicator.className = `status-indicator ${isOnline ? "online" : "offline"}`;
            }
            
            // Résoudre la Promise seulement après avoir reçu la première liste
            resolve();
        };

        ws.onopen = () => console.log("WebSocket connecté");
        ws.onerror = (error) => console.error("WebSocket erreur :", error);
        ws.onclose = () => console.log("WebSocket fermé");
    });
}
