

let currentUserID;

export function sendMessage(receiverId, content, ws) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        const messageData = {
            type: "message",
            receiver_id: receiverId,
            sender_id: currentUserID,
            content: content,
            sent_at: new Date().toISOString(),
        };
        console.log(messageData)
        ws.send(JSON.stringify(messageData));
        displayMessage(messageData);
    } else {
        console.error("WebSocket est déconnecté");
    }
}

export async function openMessageModal(userId, username, ws) {
    const recipientUsername = document.getElementById("recipientUsername");
    const messageContainer = document.getElementById("messageContainer");
    const messageInput = document.getElementById("messageInput");
    const sendMessageButton = document.getElementById("sendMessageButton");

    recipientUsername.textContent = username;
    messageContainer.innerHTML = "";

    try {
        const response = await fetch(`/message/${userId}`, { credentials: "include" });
        if (!response.ok) {
            throw new Error("Erreur lors de la récupération des messages");
        }
        const messages = await response.json();
        messages.forEach((message) => {
            displayMessage(message);
        });
    } catch (error) {
        console.error("Erreur :", error);
    }

    sendMessageButton.onclick = () => {
        const content = messageInput.value.trim();
        if (content) {
            sendMessage(userId, content, ws);
            messageInput.value = "";
        }
    };

    document.getElementById("messageModal").style.display = "block";
}

export function closeMessageModal() {
    document.getElementById("messageModal").style.display = "none";
}

export function displayMessage(message) {
    const messageContainer = document.getElementById("messageContainer");
    const messageElement = document.createElement("div");
    console.log(message.sender_id)
    messageElement.classList.add("message", message.sender_id === currentUserID ? "sent" : "received");
    messageElement.innerHTML = `
        <strong>${message.content}</strong>
        <span class="timestamp">${new Date(message.sent_at).toLocaleString("fr-FR", { timeZone: "UTC" })}</span>
    `;
    messageContainer.appendChild(messageElement);
    messageContainer.scrollTop = messageContainer.scrollHeight;
}



export function setCurrentUserId(userId) {
    currentUserID = userId;
}


