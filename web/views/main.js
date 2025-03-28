import { PostForm } from "./post_form.js";
import { displayPosts, updateFilters } from "./post_display.js";
import { displayLoginForm } from "./auth_form.js";
import { startPrivateChat, checkNewMessages, showMessagesModal } from "./message_form.js";

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

    const messages = document.getElementById("messagesNotifs");
    if (messages) {
        messages.innerHTML = `
        <div class="notification-container">
            <div id="messageNotification" class="notification-badge hidden"></div>
            <svg id="messageIcon" xmlns="http://www.w3.org/2000/svg" width="32" height="32" fill="#BB86FC" viewBox="0 0 16 16">
                <path d="M.05 3.555A2 2 0 0 1 2 2h12a2 2 0 0 1 1.95 1.555L8 8.414.05 3.555ZM0 4.697v7.104l5.803-3.558L0 4.697ZM6.761 8.83l-6.57 4.026A2 2 0 0 0 2 14h12a2 2 0 0 0 1.808-1.144l-6.57-4.026L8 9.586l-1.239-.757Zm3.436-.586L16 11.801V4.697l-5.803 3.546Z"/>
            </svg>
        </div>
        `;
        
    }
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

    document.getElementById("messageIcon").addEventListener("click", showMessagesModal);

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
        //Séparer les utilisateurs connectés et non connectés
        const onlineUsers = [];
        const offlineUsers = [];
        
        for (const user of users) {
            if (user.UserId === currentUserId) continue;
            
            if (onlineUserIds.has(user.UserId)) {
                onlineUsers.push({ user, isOnline: true });
            } else {
                offlineUsers.push({ user, isOnline: false });
            }
        }

        //Pour les utilisateurs connectés, récupérer les derniers messages
        const onlineUsersWithMessages = await Promise.all(
            onlineUsers.map(async ({ user }) => {
                try {
                    const response = await fetch(`/message?with=${user.UserId}&offset=0&limit=1`, {
                        credentials: "include"
                    });

                    if (!response.ok) {
                        return { user, lastMessage: null };
                    }

                    const messages = await response.json();
                    return {
                        user,
                        lastMessage: Array.isArray(messages) && messages.length > 0 ? messages[0] : null
                    };
                } catch (error) {
                    console.error(`Erreur lors de la récupération des messages pour ${user.UserId}:`, error);
                    return { user, lastMessage: null };
                }
            })
        );

        //Trier les utilisateurs connectés
        onlineUsersWithMessages.sort((a, b) => {
            // D'abord ceux avec messages (triés par date récente)
            if (a.lastMessage && b.lastMessage) {
                return new Date(b.lastMessage.SentAt) - new Date(a.lastMessage.SentAt);
            }
            if (a.lastMessage) return -1;
            if (b.lastMessage) return 1;
            
            // Ensuite ceux sans messages (triés par pseudo)
            return a.user.Username.localeCompare(b.user.Username);
        });

        //Trier les utilisateurs non connectés par pseudo
        offlineUsers.sort((a, b) => a.user.Username.localeCompare(b.user.Username));

        //Combiner les listes
        return [
            ...onlineUsersWithMessages.map(item => item.user),
            ...offlineUsers.map(item => item.user)
        ];
    } catch (error) {
        console.error("Erreur lors du tri des utilisateurs:", error);
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
    return userCard;
}

function setupWebSocket() {
    return new Promise((resolve) => {
        const ws = new WebSocket("ws://localhost:8443/ws");

        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            
            // Mise à jour des utilisateurs en ligne
            onlineUserIds = new Set(data.map(item => item.UserId));
            
            // Mettre à jour les indicateurs de statut
            for (const userId in userElements) {
                const isOnline = onlineUserIds.has(userId);
                userElements[userId].statusText.textContent = isOnline ? "Online" : "Offline";
                userElements[userId].indicator.className = `status-indicator ${isOnline ? "online" : "offline"}`;
            }
            
            // Vérifier les nouveaux messages (simplifié - en vrai il faudrait un système plus robuste)
            checkNewMessages();
            
            resolve();
        };

        ws.onopen = () => console.log("WebSocket connecté");
        ws.onerror = (error) => console.error("WebSocket erreur :", error);
        ws.onclose = () => console.log("WebSocket fermé");
    });
}

// Vérifier les nouveaux messages périodiquement
setInterval(checkNewMessages, 1000); 

