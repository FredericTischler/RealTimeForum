document.getElementById("loginView").addEventListener("click", () => {
    const viewForm = document.getElementById("login");
    viewForm.style.display = "block";
    viewForm.innerHTML = `
        <div class="container">
            <h2>Login</h2>
            <form action="login" method="POST">
                <label for="user_name">User name</label>
                <input type="text" name="identifier" id="identifier" required>
                
                <label for="password">Password</label>
                <input type="password" name="password" id="password" required>
                
                <button type="submit">Login</button>
            </form>
        </div>
    `;
});
