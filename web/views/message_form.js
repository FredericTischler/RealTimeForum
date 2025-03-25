// web/views/private_chat.js
let ws;
let offset = 0;
const chatModal = document.getElementById("privateChatModal");

export async function startPrivateChat(myUserId, targetUserId, targetUsername) {
    // Création du modal
    chatModal.innerHTML = `
        <div class="chat-box">
            <h1>${targetUsername}</h1>
            <div id="messages" style="height: 300px; overflow-y: scroll; display: flex; flex-direction: column;"></div>
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
                input.value = "";


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
        if (messagesDiv.scrollTop <= 50) { // Déclenche quand on arrive en haut
            await loadPreviousMessages(targetUserId, myUserId);
        }
    });


   // Rendre le modal visible
    chatModal.style.display = "flex";

    // Supprimer tout ancien écouteur pour éviter des doublons
    document.removeEventListener("click", closeChatOnOutsideClick);

    // Ajouter un nouvel écouteur
    document.addEventListener("click", closeChatOnOutsideClick);
}

function closeChatOnOutsideClick(event) {
    const chatBox = document.querySelector(".chat-box");

    // Vérifie si le clic est en dehors du chat
    if (chatBox && !chatBox.contains(event.target)) {
        chatModal.style.display = "none";
        
        // Supprime l'écouteur pour éviter les conflits lors du prochain affichage
        document.removeEventListener("click", closeChatOnOutsideClick);
    }
}

// Fonction pour charger les messages précédents
async function loadPreviousMessages(targetUserId, myUserId) {
    try {
        const messagesDiv = document.getElementById("messages");
        const previousScrollHeight = messagesDiv.scrollHeight; // Sauvegarder la hauteur avant ajout

        const response = await fetch(`/message?with=${targetUserId}&offset=${offset}`, {
            credentials: "include"
        });

        if (!response.ok) {
            throw new Error("Erreur lors du chargement des messages");
        }
        
        const messages = await response.json();

        if (messages.length === 0) return; // Stop si plus de messages

        // Ajout des nouveaux messages EN HAUT
        const fragment = document.createDocumentFragment();
        messages.forEach(msg => {
            const messageEl = document.createElement("p");
            let date = new Date(msg.SentAt);
            let hours = date.getHours()
            let minutes = date.getMinutes()
            let year = date.getFullYear()
            let mounth = date.getMonth()
            let day = date.getDate()

            messageEl.innerHTML = `<strong>${msg.Content}</strong><br><small>${hours <= 9 ? "0" + hours : hours}:${minutes <= 9 ? "0" + minutes : minutes} 
            ${year}-${mounth <= 9 ? "0" + mounth : mounth}-${day <= 9 ? "0" + day : day}</small>`;
            msg.SenderId === myUserId ? messageEl.classList.add('sender') : messageEl.classList.add('receiver');
            fragment.prepend(messageEl);
        });

        messagesDiv.prepend(fragment); // Ajouter tout d'un coup pour éviter le reflow

        offset += 10; // Incrémenter l'offset pour éviter de recharger les mêmes messages

        // Ajuster le scroll pour éviter un saut brutal
        setTimeout(() => {
            messagesDiv.scrollTop = messagesDiv.scrollHeight - previousScrollHeight;
        }, 100);
    } catch (error) {
        console.error(error);
    }
}


function appendMessage(senderId, message, currentUserId) {
    const messagesDiv = document.getElementById("messages");
    const messageEl = document.createElement("p");
    let date = new Date();
    let hours = date.getHours()
    let minutes = date.getMinutes()
    let year = date.getFullYear()
    let mounth = date.getMonth()
    let day = date.getDate()

    messageEl.innerHTML = `<strong>${message}</strong><br><small>${hours <= 9 ? "0" + hours : hours}:${minutes <= 9 ? "0" + minutes : minutes} 
            ${year}-${mounth <= 9 ? "0" + mounth : mounth}-${day <= 9 ? "0" + day : day}</small>`;
    senderId === currentUserId ? messageEl.classList.add('sender') : messageEl.classList.add('receiver');

    messagesDiv.appendChild(messageEl); // Ajouter en bas

    // Scroller vers le bas pour voir le nouveau message
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}



