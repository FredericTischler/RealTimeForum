# Real-Time Forum

## Introduction
Welcome to the **Real-Time Forum** project! This project is a more advanced version of a previous forum implementation, introducing real-time interactions through WebSockets, private messaging, and enhanced user functionalities.

## Objectives
This forum will allow users to:
- Register and log in.
- Create and comment on posts.
- Send and receive private messages in real time.

The project is built using the following technologies:
- **SQLite** for database management.
- **Golang** for backend logic and WebSockets handling.
- **JavaScript** for frontend interactions and WebSockets clients.
- **HTML & CSS** for structuring and styling the UI.

## Features
### Registration & Login
- Users must register using:
  - Nickname
  - Age
  - Gender
  - First Name
  - Last Name
  - Email
  - Password
- Users can log in using either their nickname or email combined with their password.
- Users can log out from any page.

### Posts & Comments
- Users can create posts categorized by topic.
- Users can comment on posts.
- Posts are displayed in a feed format.
- Comments are only visible when a post is opened.

### Private Messages
- Users can send private messages to each other.
- A chat section displays online/offline users, sorted by the last message sent.
- Users can access their previous messages.
- The last 10 messages load initially, with more messages loading on scroll.
- Messages contain:
  - Timestamp
  - Sender's name
- Real-time messaging without requiring a page refresh (via WebSockets).

## Project Structure
The forum consists of a **Single Page Application (SPA)** with one HTML file. All page transitions are managed by JavaScript.

### Backend
- Built with **Golang**.
- Manages data operations and WebSocket connections.
- Uses **SQLite** for data storage.

### Frontend
- Built with **JavaScript**.
- Manages user interactions and WebSockets communication.
- Updates the UI dynamically without full page reloads.

## Allowed Packages
The following Go packages are permitted:
- **Standard Go libraries**
- **Gorilla WebSocket** (for real-time communication)
- **sqlite3** (for database operations)
- **bcrypt** (for password hashing)
- **gofrs/uuid** or **google/uuid** (for unique identifiers)

### Restrictions
- **No frontend frameworks** like React, Angular, or Vue are allowed.

## Learning Outcomes
By completing this project, you will gain experience in:
- Web development fundamentals (HTML, CSS, JavaScript, HTTP, sessions, and cookies).
- Backend and frontend communication.
- Managing WebSockets in both Go and JavaScript.
- Database manipulation using SQL.




