let isClicked = false

document.getElementById("loginView").addEventListener("click", () => {
    const viewForm = document.getElementById("login")
    isClicked = true
    viewForm.style.display = isClicked ? "block" : "none"
    viewForm.innerHTML = `
        <div class="container">
            <h2>Login</h2>
            <form action="login" method="POST">
                <label for="user_name">User name or email</label>
                <input type="text" name="user_name" id="user_name" required>
        
                <label for="password">Password</label>
                <input type="password" name="password" id="password" required>
                
                <button type="submit">Login</button>
            </form>
        </div>
    `
})