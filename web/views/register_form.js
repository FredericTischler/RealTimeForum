// web/views/register_form.js

document.getElementById("registerView").addEventListener("click", () => {
    const viewForm = document.getElementById("register");
    viewForm.style.display = "block";
    viewForm.innerHTML = `
    <div class="container">
      <h2>Register</h2>
      <form action="register" method="POST" id="registerForm">
        <label for="user_name">User name</label>
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
  
        <button type="submit">Register</button>
      </form>
    </div>
  `;
});
