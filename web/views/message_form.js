export async function loadMessages(sender, recipient) {
    try {
        const response = await fetch(`/api/messages?sender=${sender}&recipient=${recipient}`);
        if (!response.ok) {
            throw new Error("Failed to fetch messages");
        }
        const messages = await response.json();

        // Afficher les messages dans le modal
        const messageContainer = document.getElementById("messageContainer");
        messageContainer.innerHTML = ""; // Vider les messages précédents

        messages.forEach((message) => {
            const messageElement = document.createElement("div");
            messageElement.textContent = `${message.sender}: ${message.content}`;
            messageContainer.appendChild(messageElement);
        });
    } catch (error) {
        console.error("Error loading messages:", error);
    }
}