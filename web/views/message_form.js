// web/views/private_chat.js
let ws;
let offset = 0;

export async function startPrivateChat(myUserId, targetUserId, targetUsername) {
    // Création du modal
    const chatModal = document.getElementById("privateChatModal");
    chatModal.innerHTML = `
        <div class="chat-box">
            <h3>${targetUsername}</h3>
            <div id="messages" style="height: 300px; overflow-y: scroll; display: flex; flex-direction: column-reverse;"></div>
            <input type="text" id="chatInput" placeholder="Message...">
            <button id="sendBtn">Send</button>
        </div>
    `;
    document.body.appendChild(chatModal);

    // Initialiser WebSocket
    ws = new WebSocket("ws://localhost:8443/ws/private");

    ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        appendMessage(msg.from, msg.message, myUserId);
    };

    document.getElementById("sendBtn").addEventListener("click", () => {
        const input = document.getElementById("chatInput");
        const message = input.value.trim();
        if (message !== "") {
            ws.send(JSON.stringify({ to: targetUserId, message }));
            appendMessage(myUserId, message, myUserId);
            input.value = "";
        }
    });

    // Chargement initial des 10 derniers messages
    offset = 0;
    await loadPreviousMessages(targetUserId);

    // Scroll pour charger plus
    const messagesDiv = document.getElementById("messages");
    messagesDiv.addEventListener("scroll", async () => {
        if (messagesDiv.scrollTop === messagesDiv.scrollHeight - messagesDiv.clientHeight) {
            await loadPreviousMessages(targetUserId);
        }
    });
}

// Fonction pour charger les messages précédents
async function loadPreviousMessages(targetUserId) {
    try {
        const response = await fetch(`/message?with=${targetUserId}&offset=${offset}`, {
            credentials: "include"
        });
        if (!response.ok) {
            throw new Error("Erreur lors du chargement des messages");
        }
        
        const messages = await response.json();

        console.log(messages)

        const container = document.getElementById("messages");
        messages.reverse().forEach(msg => {
            appendMessage(msg.from_user, msg.message, msg.current_user_id);
        });

        offset += 10;
    } catch (error) {
        console.error(error);
    }
}

function appendMessage(senderId, message, currentUserId) {
    const messagesDiv = document.getElementById("messages");
    const messageEl = document.createElement("p");
    messageEl.innerHTML = `<strong>${senderId === currentUserId ? "Me" : "Them"}:</strong> ${message}`;
    messagesDiv.prepend(messageEl);
}