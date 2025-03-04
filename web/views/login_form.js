

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
                body: formData,
                credentials: "include"
            });
            if (!response.ok) {
                const errorText = await response.text();
                alert("Erreur: " + errorText);
            }
        } catch (err) {
            console.error("Erreur lors du login", err);
            alert("Le login a échoué, veuillez réessayer.");
        }
        // Après login, on laisse main.js vérifier l'état et mettre à jour l'UI.
        // Vous pouvez recharger la page ou appeler checkAuthStatus() si nécessaire.
        location.reload();
    });
});
