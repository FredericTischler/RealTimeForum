// web/views/private_chat.js
let ws;
let offset = 0;
const chatModal = document.getElementById("privateChatModal");
export let isChatOpen = false;

export async function startPrivateChat(myUserId, targetUserId, targetUsername) {
    // Marquer les messages comme lus
    isChatOpen = true;
    await fetch(`/message/mark-as-read?sender_id=${myUserId}`, {
        method: 'POST',
        credentials: 'include'
    });
        
    // Recharger le compteur de notifications
    await checkNewMessages();

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

    const input = document.getElementById("chatInput");
    input.addEventListener("keypress", (e) => {
        if (e.key === "Enter") {
            e.preventDefault();
            document.getElementById("sendBtn").click();
        }
    });
    document.getElementById("sendBtn").addEventListener("click", async () => {
        const message = input.value.trim();
    
        if (message !== "") {
            try {
                // Envoyer le message au serveur
                isChatOpen = true;
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
        isChatOpen = false;
        // Supprime l'écouteur pour éviter les conflits lors du prochain affichage
        document.removeEventListener("click", closeChatOnOutsideClick);
    }
}

// Fonction pour charger les messages précédents
async function loadPreviousMessages(targetUserId, myUserId) {
    isChatOpen = true;
    
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
            let mounth = date.getMonth() + 1
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
    let mounth = date.getMonth() + 1
    let day = date.getDate()

    console.log(mounth)

    messageEl.innerHTML = `<strong>${message}</strong><br><small>${hours <= 9 ? "0" + hours : hours}:${minutes <= 9 ? "0" + minutes : minutes} 
            ${year}-${mounth <= 9 ? "0" + mounth : mounth}-${day <= 9 ? "0" + day : day}</small>`;
    senderId === currentUserId ? messageEl.classList.add('sender') : messageEl.classList.add('receiver');

    messagesDiv.appendChild(messageEl); // Ajouter en bas

    // Scroller vers le bas pour voir le nouveau message
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}

let unreadMessagesCount = 0;

function updateUnreadMessagesCount(count) {
    if (isChatOpen) return;
    unreadMessagesCount = count;
    const badge = document.getElementById("messageNotification");
    if (count > 0) {
        badge.textContent = count;
        badge.classList.remove("hidden");
    } else {
        badge.classList.add("hidden");
    }
}

export async function checkNewMessages() {
        try {
            const response = await fetch("/message/unread", { 
                credentials: "include" 
            });
            if (response.ok) {
                const data = await response.json();
                updateUnreadMessagesCount(data.count);
            }
        } catch (error) {
            console.error("Error checking new messages:", error);
        }
}

export async function showMessagesModal() {

    let modal = document.getElementById("messagesModal");
    
    if (modal) {
        modal.remove();
    }

    modal = document.createElement("div");
    modal.id = "messagesModal";
    modal.innerHTML = `
        <h2>${unreadMessagesCount > 0 ?  unreadMessagesCount == 1 ? `${unreadMessagesCount} new message recorded` : `${unreadMessagesCount} new messages recorded` : 'No new messages'}</h2>
        <div id="conversationsList"></div>
    `;
    document.body.appendChild(modal);
    
    // Fermer la modal en cliquant à l'intérieur
    modal.addEventListener("click", (e) => {
        if (e.target === modal) {
            modal.style = "display: none;";
        }
    });
    
    // Charger les conversations avec un filtre pour les nouveaux messages
    await loadConversations(unreadMessagesCount > 0);
    
    modal.style.display = "block";
}

async function loadConversations(onlyUnread = false) {
    
    try {
        const response = await fetch('/users', { credentials: 'include' });
        if (!response.ok) {
            console.error('Failed to fetch users:', response.status);
            return;
        }

        const users = await response.json();
        if (!users || !Array.isArray(users)) {
            console.error('Invalid users data:', users);
            return;
        }

        const authResponse = await fetch('/auth/status', { credentials: 'include' });
        if (!authResponse.ok) {
            console.error('Not authenticated:', authResponse.status);
            return;
        }

        const authData = await authResponse.json();
        const currentUserID = authData.userId;

        const list = document.getElementById('conversationsList');
        if (!list) {
            console.error('Conversations list element not found');
            return;
        }

        // Filtrer et traiter chaque conversation
        const validUsers = users.filter(user => user.UserId && user.UserId !== currentUserID);
        
        const conversations = await Promise.all(
            validUsers.map(async user => {
                try {
                    let url = `/message?with=${user.UserId}&offset=0&limit=1`;
                    if (onlyUnread) {
                        url += '&unread=true';
                    } else {
                        return
                    }

                    const messagesResponse = await fetch(url, {
                        credentials: 'include'
                    });

                    if (!messagesResponse.ok) {
                        console.error(`Failed to fetch messages for user ${user.UserId}`);
                        return null;
                    }

                    const messages = await messagesResponse.json();

                    
                    // Si on ne veut que les non lus et qu'il n'y en a pas, on ignore
                    if (onlyUnread && (!Array.isArray(messages) || messages.length === 0)) {
                        return null;
                    }

                    return {
                        user,
                        lastMessage: Array.isArray(messages) && messages.length > 0 ? messages[0] : null
                    };
                } catch (error) {
                    console.error(`Error processing user ${user.UserId}:`, error);
                    return null;
                }
            })
        );

        // Filtrer les conversations valides et les trier
        const validConversations = conversations.filter(c => c !== null)
            .sort((a, b) => {
                if (a.lastMessage && b.lastMessage) {
                    return new Date(b.lastMessage.SentAt) - new Date(a.lastMessage.SentAt);
                }
                if (a.lastMessage) return -1;
                if (b.lastMessage) return 1;
                return 0;
            });

        // Afficher les conversations
        if (validConversations.length === 0) {
            list.innerHTML = '<div class="no-conversations">' + 
                (onlyUnread ? 'No new messages...' : 'No messages...') + 
                '</div>';
            return;
        }


        validConversations.forEach(conv => {
            const item = document.createElement('div');
            item.className = 'conversation-item';
            
            // Ajouter une classe si le message est non lu
            const isUnread = conv.lastMessage && !conv.lastMessage.Is_read && 
                             conv.lastMessage.SenderId !== currentUserID;
            if (isUnread) {
                item.classList.add('unread-message');
            }
            
            item.innerHTML = `
                <div class="user-info">${conv.user.Username}</div>
                <div class="message-preview">${conv.lastMessage?.Content || ''}</div>
                <div class="message-time">${
                    conv.lastMessage ? new Date(conv.lastMessage.SentAt).toLocaleString() : ''
                }</div>
            `;
            item.addEventListener('click', () => {
                startPrivateChat(currentUserID, conv.user.UserId, conv.user.Username);
                document.getElementById('messagesModal').style.display = 'none';
            });
            list.appendChild(item);
        });

    } catch (error) {
        console.error('Error loading conversations:', error);
        const list = document.getElementById('conversationsList');
        if (list) {
            list.innerHTML = '<div class="error-message"></div>';
        }
    }
}
