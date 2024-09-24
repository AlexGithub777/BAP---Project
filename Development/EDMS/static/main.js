function logout() {
    window.location.href = "/logout";
}

function viewNotifications() {
    console.log("View notifications");
    $("#notificationsModal").modal("show");
    // Add your view notifications logic here
}

// hot reload
if (window.EventSource) {
    new EventSource("http://localhost:8090/internal/reload").onmessage = () => {
        setTimeout(() => {
            location.reload();
        });
    };
}

function showDeleteModal(id, entityType) {
    const deleteModal = document.getElementById("deleteModal");
    const deleteForm = document.getElementById("deleteForm");
    const deleteIdInput = document.getElementById("deleteId");
    const modalBody = deleteModal.querySelector(".modal-body p");
    const deleteButton = deleteModal.querySelector(".modal-footer .btn-danger");

    if (
        deleteModal &&
        deleteForm &&
        deleteIdInput &&
        modalBody &&
        deleteButton
    ) {
        // Capitalize the first letter of entityType
        const capitalizedEntityType =
            entityType.charAt(0).toUpperCase() + entityType.slice(1);

        // Update modal text
        modalBody.textContent = `Are you sure you want to delete this ${entityType}?`;
        deleteButton.textContent = `Delete ${capitalizedEntityType}`;

        // Set form action and ID
        deleteForm.action = `/api/${entityType}/${id}`;
        deleteIdInput.value = id;

        console.log(deleteForm.action, deleteIdInput.value);

        // Show the modal
        const modal = new bootstrap.Modal(deleteModal);
        modal.show();
    } else {
        console.error(
            "One or more elements required for delete modal not found"
        );
    }
}

// Event listener for delete form submission
// Event listener for delete form submission
document
    .getElementById("deleteForm")
    .addEventListener("submit", function (event) {
        event.preventDefault();

        fetch(this.action, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
                // Add any other necessary headers here
            },
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    window.location.href = data.redirectURL;
                } else if (data.message) {
                    window.location.href = data.redirectURL;
                } else {
                    console.error("Error:", data);
                    // Handle errors (e.g., show an error message
                }
            })
            .catch((error) => {
                console.error("Error:", error);
                // Handle errors (e.g., show an error message to the user)
            });
    });
