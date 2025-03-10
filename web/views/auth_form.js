// web/views/auth_form.js

import { displayErrorModal } from './error_modal.js';

export function displayLoginForm() {
    const loginContainer = document.getElementById("login");
    loginContainer.style.display = "block";
    loginContainer.innerHTML = `
    <div id="login-container">
      <div class="container">
        <h2>Sign in</h2>
        <form id="loginForm">
          <label for="identifier">Username</label>
          <input type="text" name="identifier" id="identifier" required>
          
          <label for="password">Password</label>
          <input type="password" name="password" id="password" required>
          
          <button class="authButton" type="submit">Sign in</button>
        </form>
        <p>You don't have any account ? <a href="#" id="switchToRegister">Sign up</a></p>
      </div>
   </div>
  `;
    const loginForm = document.getElementById("loginForm");
    loginForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const formData = new FormData(loginForm);
        try {
            const response = await fetch("login", {
                method: "POST",
                body: formData,
                credentials: "include"
            });
            if (!response.ok) {
                const errorText = await response.text();
                displayErrorModal("Erreur: " + errorText);
                return;
            }
            // Après login, recharger la page pour afficher l'interface connectée
            location.reload();
        } catch (err) {
            console.error("Erreur lors du login", err);
            displayErrorModal("Le login a échoué, veuillez réessayer.");
        }
    });

    // Gestion du switch vers le formulaire d'inscription
    document.getElementById("switchToRegister").addEventListener("click", (e) => {
        e.preventDefault();
        displayRegisterForm();
    });
}

export function displayRegisterForm() {
    const registerContainer = document.getElementById("register");
    // Masquer le formulaire de connexion
    document.getElementById("login").style.display = "none";
    registerContainer.style.display = "block";
    registerContainer.innerHTML = `
      <div id="login-container">
        <div class="container">
          <h2>Sign up</h2>
          <form id="registerForm">
            <label for="user_name">Username</label>
            <input type="text" name="user_name" id="user_name" required>
      
            <label for="email">Email</label>
            <input type="email" name="email" id="email" required>
      
            <label for="password">Password</label>
            <input type="password" name="password" id="password" required>
      
            <label for="first_name">Firstname</label>
            <input type="text" name="first_name" id="first_name" required>
      
            <label for="last_name">Lastname</label>
            <input type="text" name="last_name" id="last_name" required>
      
            <label for="age">Age</label>
            <input type="number" name="age" id="age" required>
      
            <label for="gender">Gender</label>
            <select name="gender" id="gender">
              <option value="">Select your gender</option>
              <option value="Female">Female</option>
              <option value="Male">Male</option>
              <option value="Other">Other</option>
            </select>
      
            <button class="authButton" type="submit">Sign up</button>
          </form>
          <p>You already have an account ? <a href="#" id="switchToLogin">Sign in</a></p>
        </div>
      </div>
  `;

  const registerForm = document.getElementById("registerForm");
    registerForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const registerFormData = new FormData(registerForm);
        try {
            const registerResponse = await fetch("register", {
                method: "POST",
                body: registerFormData,
                credentials: "include"
            });
            console.log(registerResponse)
            if (!registerResponse.ok) {
                const errorText = await registerResponse.text();
                displayErrorModal("Erreur: " + errorText);
                return;
            }
            location.reload();
        } catch (err) {
            console.error("Erreur lors du register", err);
            displayErrorModal("Le register a échoué, veuillez réessayer.");
        }
    });

  

    document.getElementById("switchToLogin").addEventListener("click", (e) => {
        e.preventDefault();
        registerContainer.style.display = "none";
        displayLoginForm();
    });
}
