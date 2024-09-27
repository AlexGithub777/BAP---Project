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

function formatEntityType(entityType) {
    return entityType
        .split("-")
        .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
        .join(" ");
}

function showDeleteModal(id, entityType, entityName, currentUserId) {
    const deleteModal = document.getElementById("deleteModal");
    const deleteForm = document.getElementById("deleteForm");
    const currentUserIdInput = document.getElementById("deleteCurrentUserID");
    const deleteIdInput = document.getElementById("deleteId");
    const modalBody = deleteModal.querySelector(".modal-body p");
    const deleteButton = deleteModal.querySelector(".modal-footer .btn-danger");

    if (
        deleteModal &&
        deleteForm &&
        deleteIdInput &&
        currentUserIdInput &&
        modalBody &&
        deleteButton
    ) {
        // Format and capitalize the entityType
        const formattedEntityType = formatEntityType(entityType);

        // Update modal text
        modalBody.innerHTML = `Are you sure you want to delete ${formattedEntityType}: ${entityName}?`;
        deleteButton.textContent = `Delete ${formattedEntityType}`;

        // Add any additional logic here

        // if the current user id is passed, add to delete form action
        if (currentUserId) {
            deleteForm.action = `/api/${entityType}/${id}?currentUserId=${currentUserId}`;
        } else {
            deleteForm.action = `/api/${entityType}/${id}`;
        }

        deleteIdInput.value = id;
        currentUserIdInput.value = currentUserId;

        console.log(
            deleteForm.action,
            deleteIdInput.value,
            currentUserIdInput.value
        );

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

        // Create a FormData object
        const formData = new FormData(this);

        // Convert FormData to a plain object
        const plainFormData = Object.fromEntries(formData.entries());

        fetch(this.action, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
                // Add any other necessary headers here
            },
            // add cur
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    window.location.href = data.redirectURL;
                } else if (data.message) {
                    window.location.href = data.redirectURL;
                } else {
                    console.error("Error:", data);
                    // Handle errors (e.g., show an error message)
                }
            })
            .catch((error) => {
                console.error("Error:", error);
                // Handle errors (e.g., show an error message to the user)
            });
    });
