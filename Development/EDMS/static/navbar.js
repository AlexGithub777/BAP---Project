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
