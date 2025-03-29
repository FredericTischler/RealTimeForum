// web/views/error_modal.js

export function displayErrorModal(message) {
    // Création de l'overlay du modal
    const overlay = document.createElement("div");
    overlay.classList.add("error-modal-overlay");

    // Création du contenu du modal
    const modal = document.createElement("div");
    modal.classList.add("error-modal");
    modal.innerHTML = `<p>${message}</p>`;

    // Bouton pour fermer le modal
    const closeButton = document.createElement("button");
    closeButton.textContent = "Close";
    closeButton.addEventListener("click", () => {
        overlay.remove();
    });

    modal.appendChild(closeButton);
    overlay.appendChild(modal);

    // Ajout du modal dans le document
    document.body.appendChild(overlay);
}
