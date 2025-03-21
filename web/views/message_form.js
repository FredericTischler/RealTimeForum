// web/views/private_chat.js
let ws;
let offset = 0;

export async function startPrivateChat(myUserId, targetUserId, targetUsername) {
    // Création du modal
    const chatModal = document.getElementById("privateChatModal");
    chatModal.innerHTML = `
        <div class="chat-box">
            <h1>${targetUsername}</h1>
            <div id="messages" style="height: 300px; overflow-y: scroll; display: flex; flex-direction: column-reverse;"></div>
            <input type="text" id="chatInput" placeholder="Message...">
            <div class="btnContainer">
                <button id="sendBtn">Send</button>
            </div>
        </div>
    `;
    document.body.appendChild(chatModal);

    // Initialiser WebSocket
    ws = new WebSocket("ws://localhost:8443/ws/private");

    ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        appendMessage(msg.from, msg.message, myUserId);
    };

    
    document.getElementById("sendBtn").addEventListener("click", async () => {
        const input = document.getElementById("chatInput");
        const message = input.value.trim();
    
        if (message !== "") {
            try {
                // Envoyer le message au serveur
                const response = await fetch("/message/insert", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    credentials: "include", // Inclure les cookies pour l'authentification
                    body: JSON.stringify({
                        receiverId: targetUserId, // ID du destinataire
                        content: message,        // Contenu du message
                    }),
                });
    
                if (!response.ok) {
                    throw new Error("Erreur lors de l'envoi du message");
                }

                ws.send(JSON.stringify({ to: targetUserId, message }));
                appendMessage(myUserId, message, myUserId);

            } catch (error) {
                console.error("Erreur lors de l'envoi du message :", error);
                alert("Erreur lors de l'envoi du message");
            }
        }
    });

    // Chargement initial des 10 derniers messages
    offset = 0;
    await loadPreviousMessages(targetUserId, myUserId);

    // Scroll pour charger plus
    const messagesDiv = document.getElementById("messages");
    messagesDiv.addEventListener("scroll", async () => {
        if (messagesDiv.scrollTop === messagesDiv.scrollHeight - messagesDiv.clientHeight) {
            await loadPreviousMessages(targetUserId);
        }
    });
}

// Fonction pour charger les messages précédents
async function loadPreviousMessages(targetUserId, myUserId) {
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
            console.log(msg.message)
            appendMessage(msg.SenderId, msg.Content, myUserId);
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
    senderId === currentUserId ? messageEl.classList.add('sender') : messageEl.classList.add('receiver')
    messagesDiv.prepend(messageEl);
}

