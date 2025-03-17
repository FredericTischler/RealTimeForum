import { PostForm } from "./post_form.js";
import { displayPosts, updateFilters } from "./post_display.js";
import { displayLoginForm } from "./auth_form.js";

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

    // Affichage du contenu principal (les aside et la section postsSection restent dans le HTML)
    const contentContainer = document.querySelector(".content-container");
    if (contentContainer) {
        contentContainer.style.display = "flex";
    }

    // Appel à la fonction d'affichage des posts
    displayPosts();

    const ws = new WebSocket("ws://localhost:8443/ws");

    ws.onopen = () => {
        console.log("Connexion WebSocket établie");
    };

    ws.onmessage = (event) => {
        const onlineUsers = JSON.parse(event.data);
        const usersListSection = document.getElementById("usersList");
        if (usersListSection) {
            let html = "<ul>";
            onlineUsers.forEach(user => {
                html += `
              <li>
                <span class="username">${user.Username}</span>
                <span class="gender" style="display: none;">${user.Gender}</span>
                <span class="age" style="display: none;">${user.Age}</span>
              </li>`;
            });
            html += "</ul>";
            usersListSection.innerHTML = html;
        }
    };



    ws.onerror = (error) => {
        console.error("Erreur WebSocket :", error);
    };

    ws.onclose = () => {
        console.log("Connexion WebSocket fermée");
    };
}

function filterUsers(username, gender, age) {
    const usersListSection = document.getElementById("usersList");
    if (usersListSection) {
        const users = Array.from(usersListSection.querySelectorAll("li"));
        users.forEach(user => {
            const userUsername = user.querySelector(".username").textContent;
            const userGender = user.querySelector(".gender").textContent;
            const userAge = user.querySelector(".age").textContent;
            const shouldDisplay = 
                (username === "" || userUsername.includes(username)) &&
                (gender === "" || userGender === gender) && 
                (age === "" || userAge === age);
            user.style.display = shouldDisplay ? "block" : "none";
        });
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

