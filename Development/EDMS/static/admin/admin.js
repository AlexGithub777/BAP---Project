// admin.js
import {
    updateNotificationsUI,
    refreshAfterChange,
} from "/static/main/notifications.js";
import {
    viewDeviceInspections,
    viewInspectionDetails,
    addInspection,
    initializeInspectionForm,
} from "/static/main/inspections.js";

initializeInspectionForm();

document.addEventListener("DOMContentLoaded", async function () {
    if (role === "Admin") {
        // Check for the refresh flag
        const shouldRefresh = sessionStorage.getItem(
            "shouldRefreshNotifications"
        );
        if (shouldRefresh === "true") {
            sessionStorage.removeItem("shouldRefreshNotifications");
            await refreshAfterChange();
        } else {
            // Normal initial load
            await updateNotificationsUI();
        }
    }
});

window.viewDeviceInspections = viewDeviceInspections;
window.viewInspectionDetails = viewInspectionDetails;
window.addInspection = addInspection;

$(document).ready(function () {
    // Image Preview
    $("#siteMapImgInput").change(function () {
        var file = this.files[0];
        // Check if the file is an image
        if (!file.type.startsWith("image/")) {
            $("#siteMapImgInput").val("");
            Toastify({
                text: "Please upload an image file.",
                duration: 6000,
                close: true,
                gravity: "top", // `top` or `bottom`
                position: "center", // `left`, `center` or `right`
                backgroundColor: "linear-gradient(to right, #ff5f6d, #ffc371)",
            }).showToast();
            return;
        }

        var reader = new FileReader();
        reader.onload = function (e) {
            $("#imagePreview").attr("src", e.target.result);
            $("#imagePreviewContainer").show();
        };
        reader.readAsDataURL(file);
    });
    // Use event delegation to handle dynamically loaded content
    $(document).on(
        "change",
        "#editSiteModal #editSiteImgInput",
        function (event) {
            var fileInput = event.target;
            var file = fileInput.files && fileInput.files[0];
            if (!file) {
                $("#editSiteModal #imagePreview").attr("src", "");
                $("#editSiteModal #editSiteImagePreviewContainer").hide();
                return;
            }

            // Check if the file is an image
            if (!file.type.startsWith("image/")) {
                fileInput.value = ""; // Clear the file input
                $("#editSiteModal #imagePreview").attr("src", "");
                $("#editSiteModal #editSiteImagePreviewContainer").hide();
                Toastify({
                    text: "Please upload an image file.",
                    duration: 6000,
                    close: true,
                    gravity: "top",
                    position: "center",
                    backgroundColor:
                        "linear-gradient(to right, #ff5f6d, #ffc371)",
                }).showToast();
                return;
            }

            var reader = new FileReader();
            reader.onload = function (e) {
                $("#editSiteModal #imagePreview").attr("src", e.target.result);
                $("#editSiteModal #editSiteImagePreviewContainer").show();
            };
            reader.onerror = function (error) {
                fileInput.value = ""; // Clear the file input
                $("#editSiteModal #imagePreview").attr("src", "");
                $("#editSiteModal #editSiteImagePreviewContainer").hide();
                Toastify({
                    text: "Error reading the file. Please try again.",
                    duration: 6000,
                    close: true,
                    gravity: "top",
                    position: "center",
                    backgroundColor:
                        "linear-gradient(to right, #ff5f6d, #ffc371)",
                }).showToast();
            };
            reader.readAsDataURL(file);
        }
    );
});

// Fetch users from the server
fetch("/api/user")
    .then((response) => response.json())
    .then((users) => {
        // Convert current_user_id to a number
        const currentUserIdNumber = parseInt(current_user_id, 10);

        // Sort the users array to put the current user first
        users.sort((a, b) => {
            if (a.user_id === currentUserIdNumber) return -1;
            if (b.user_id === currentUserIdNumber) return 1;
            return a.username.localeCompare(b.username); // Sort others alphabetically
        });

        // Create a table row for each user
        const userRows = users.map((user) => {
            // Convert user.default_admin to a boolean
            var isAdmin = JSON.parse(user.default_admin);
            // Convert is_current_user_default_admin to a boolean
            var current_default_admin = JSON.parse(
                is_current_user_default_admin
            );

            const hideDelete = current_default_admin && isAdmin;

            // Determine whether to hide action buttons based on conditions
            const hideActions = !current_default_admin && isAdmin;

            // Generate the row HTML
            return `
<tr${user.user_id === currentUserIdNumber ? ' class="table-primary"' : ""}>
    <td data-label="Username">${user.username}</td>
    <td data-label="Email">${user.email}</td>
    <td data-label="Role">${user.role}</td>
    <td>
    <div class="btn-group">
    ${
        hideActions
            ? "<span class='text-muted'>No actions available</span>"
            : `
            <button class="btn btn-warning p-2 edit-user-button" onclick="editUser(${
                user.user_id
            })"
                    title="Edit User">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" 
                    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/>
                    <path d="m15 5 4 4"/>
                </svg>
            </button>
            ${
                hideDelete
                    ? ""
                    : `<button class="btn btn-danger p-2 delete-button" 
                            onclick="showDeleteModal(${user.user_id}, 'user', '${user.username}', '${currentUserIdNumber}')" 
                            data-id="${user.user_id}" 
                            title="Delete User">
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" 
                            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M3 6h18"/>
                            <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/>
                            <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>
                            <line x1="10" y1="11" x2="10" y2="17"/>
                            <line x1="14" y1="11" x2="14" y2="17"/>
                        </svg>
                    </button>`
            }
            `
    }
</div>
    </td>
</tr>
`;
        });

        // Add the rows to the users table
        $("#users-table tbody").html(userRows.join(""));
    });

export async function editUser(userId) {
    const id = userId;
    console.log("Edit user ID:", id);
    // Handle edit
    // Fetch the user data from the nearest row
    const row = $(event.target).closest("tr");
    const username = row.find("td[data-label=Username]").text();
    const email = row.find("td[data-label=Email]").text();
    const role = row.find("td[data-label=Role]").text();

    const default_admin = await fetch(`/api/user/${username}`)
        .then((response) => response.json())
        .then((user) => {
            return user.default_admin.toString();
        });

    // Fill in the form with the user data
    $("#editUserForm")[0].reset();
    $("#editUserForm input[name=current_user_id]").val(current_user_id);
    $("#editUserForm input[name=user_id]").val(id);
    $("#editUserForm input[name=username]").val(username);
    $("#editUserForm input[name=email]").val(email);
    $("#editUserForm select[name=role]").val(role);
    $("#editUserForm input[name=default_admin]").val(default_admin);

    // Set the form action to the update endpoint for this user
    $("#editUserForm").attr("action", `/api/user/${id}`);

    // Get the user ID of the user being updated
    const updatedUserId = $("#editUserForm input[name=user_id]").val();

    // If the current user ID is equal to the user being updated, display the password field
    if (current_user_id === updatedUserId) {
        $("#passwordField").show();
    } else {
        $("#passwordField").hide();
    }

    // Show the modal
    $("#editUserModal").modal("show");
    // Clear validation classes
    $("#editUserForm").removeClass("was-validated");

    var editUserForm = document.getElementById("editUserForm");

    // Add event listener to the submit button
    $("#editUserBtn").click(function (event) {
        // Check if the form is valid
        if (!editUserForm.checkValidity()) {
            event.stopPropagation();
            editUserForm.classList.add("was-validated");
        } else {
            // If the form is valid, prepare to send the PUT request
            const formData = new FormData(editUserForm);
            const jsonData = {};
            for (const [key, value] of formData.entries()) {
                jsonData[key] = value;
            }
            fetch(`/api/user/${document.getElementById("editUserID").value}`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(jsonData),
            })
                .then((response) => response.json())
                .then((data) => {
                    if (data.error) {
                        window.location.href = data.redirectURL;
                    } else if (data.message) {
                        window.location.href = data.redirectURL;
                    } else {
                        console.error("Unexpected response:", data);
                        // Handle unexpected responses (e.g., show an error message)
                        throw new Error("Unexpected response");
                    }
                })
                .catch((error) => {
                    console.error("Fetch error:", error);
                    // Optionally display a user-friendly error message
                });
        }
    });
}

// Fetch site data from the server
fetch("/api/site")
    .then((response) => response.json())
    .then((sites) => {
        // Create a table row for each site
        const siteRows = sites.map(
            (site) => `
                <tr>
                    <td data-label="Site Name">${site.site_name}</td>
                    <td data-label="Site Address">${site.site_address}</td>
                    <td data-label="Actions">
                        <div class="btn-group">
                            <button class="btn btn-warning p-2 edit-button" 
                                    onclick="editSite(${site.site_id})" 
                                    data-id="${site.site_id}" 
                                    title="Edit Site">
                                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" 
                                    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                    <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/>
                                    <path d="m15 5 4 4"/>
                                </svg>
                            </button>
                            <button class="btn btn-danger p-2 delete-button" 
                                    onclick="showDeleteModal(${site.site_id}, 'site', '${site.site_name}')" 
                                    title="Delete Site">
                                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" 
                                    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                                    <path d="M3 6h18"/>
                                    <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/>
                                    <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>
                                    <line x1="10" y1="11" x2="10" y2="17"/>
                                    <line x1="14" y1="11" x2="14" y2="17"/>
                                </svg>
                            </button>
                        </div>
                    </td>
                </tr>
            `
        );

        // Add the rows to the sites table
        $("#sites-table tbody").html(siteRows.join(""));
    });

// Fetch buildings from the server
fetch("/api/building")
    .then((response) => response.json())
    .then((buildings) => {
        const buildingRows = buildings.map(
            (building) => `
        <tr>
            <td data-label="Building Code">${building.building_code}</td>
            <td data-label="Site Name">${building.site_name}</td>
            <td>
                <div class="btn-group">
                    <button class="btn btn-warning p-2 edit-building-button" onclick="editBuilding(${building.building_id})"
                            title="Edit Building">
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" 
                            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/>
                            <path d="m15 5 4 4"/>
                        </svg>
                    </button>
                    <button class="btn btn-danger p-2 delete-button" 
                            onclick="showDeleteModal(${building.building_id}, 'building', '${building.building_code}')" 
                            title="Delete Building">
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" 
                            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M3 6h18"/>
                            <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/>
                            <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>
                            <line x1="10" y1="11" x2="10" y2="17"/>
                            <line x1="14" y1="11" x2="14" y2="17"/>
                        </svg>
                    </button>
                </div> 
            </td>
        </tr>
        `
        );

        // Add the rows to the buildings table
        $("#buildings-table tbody").html(buildingRows.join(""));
    });

export function editBuilding(buildingId) {
    const id = buildingId;

    populateDropdown(
        ".siteInput",
        "/api/site",
        "Select a Site",
        "site_id",
        "site_name"
    );

    // Handle edit
    // Fetch the building data from the server
    fetch(`/api/building/${id}`)
        .then((response) => response.json())
        .then((building) => {
            // Populate the form with the data
            document.getElementById("editBuildingID").value =
                building.building_id;
            document.getElementById("editBuildingCode").value =
                building.building_code;
            document.getElementById("editBuildingSite").value =
                building.site_id;
        })
        .catch((error) => {
            console.error("Fetch error: ", error);
        });

    // Clear validation classes
    $("#editBuildingForm").removeClass("was-validated");

    // Show the modal
    $("#editBuildingModal").modal("show");

    // Set the form action to the update endpoint for this building
    $("#editBuildingForm").attr("action", `/api/building/${id}`);

    var editBuildingForm = document.getElementById("editBuildingForm");

    // Add event listener to the submit button
    $("#editBuildingBtn").click(function (event) {
        // Check if the form is valid
        if (!editBuildingForm.checkValidity()) {
            event.stopPropagation();
            editBuildingForm.classList.add("was-validated");
        } else {
            // If the form is valid, prepare to send the PUT request
            const formData = new FormData(editBuildingForm);
            const jsonData = {};
            for (const [key, value] of formData.entries()) {
                jsonData[key] = value;
            }
            fetch(
                `/api/building/${
                    document.getElementById("editBuildingID").value
                }`,
                {
                    method: "PUT",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(jsonData),
                }
            )
                .then((response) => response.json())
                .then((data) => {
                    if (data.error) {
                        window.location.href = data.redirectURL;
                    } else if (data.message) {
                        window.location.href = data.redirectURL;
                    } else {
                        console.error("Unexpected response:", data);
                        // Handle unexpected responses (e.g., show an error message)
                        throw new Error("Unexpected response");
                    }
                })
                .catch((error) => {
                    console.error("Fetch error:", error);
                    // Optionally display a user-friendly error message
                });
        }
    });
}

// Fetch rooms from the server
fetch("/api/room")
    .then((response) => response.json())
    .then((rooms) => {
        // Create a table row for each room
        const roomRows = rooms.map(
            (room) => `
        <tr>
            <td data-label="Room Code">${room.room_code}</td>
            <td data-label="Building Code">${room.building_code}</td>
            <td data-label="Site Name">${room.site_name}</td>
            <td>
                <div class="btn-group">
                    <button class="btn btn-warning p-2 edit-room-button" onclick="editRoom(${room.room_id})" 
                            title="Edit Room">
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" 
                            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/>
                            <path d="m15 5 4 4"/>
                        </svg>
                    </button>
                    <button class="btn btn-danger p-2 delete-button" 
                            onclick="showDeleteModal(${room.room_id}, 'room', '${room.room_code}')" 
                            title="Delete Room">
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" 
                            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M3 6h18"/>
                            <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/>
                            <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>
                            <line x1="10" y1="11" x2="10" y2="17"/>
                            <line x1="14" y1="11" x2="14" y2="17"/>
                        </svg>
                    </button>
                </div>
            </td>
        </tr>
        `
        );

        // Add the rows to the rooms table
        $("#rooms-table tbody").html(roomRows.join(""));
    });

export function editRoom(roomId) {
    const id = roomId;

    // Clear custom validation and invalid classes from building select
    const buildingInput = document.getElementById("editRoomBuildingCode");
    buildingInput.setCustomValidity("");
    buildingInput.classList.remove("is-invalid");

    // Clear the form and validation states
    document.getElementById("editRoomForm").reset();
    $("#editRoomForm").removeClass("was-validated");

    // First fetch the room data
    fetch(`/api/room/${id}`)
        .then((response) => response.json())
        .then((room) => {
            // Populate sites dropdown
            fetch("/api/site")
                .then((response) => response.json())
                .then((sites) => {
                    const siteSelect = $("#editRoomSite");
                    const siteOptions = sites.map(
                        (site) =>
                            `<option value="${site.site_id}" 
                             ${site.site_id === room.site_id ? "selected" : ""}>
                             ${site.site_name}
                             </option>`
                    );
                    siteSelect.html(
                        `<option value="">Select a Site</option>` +
                            siteOptions.join("")
                    );

                    // After setting site, load buildings for that site
                    return fetch(`/api/building?siteId=${room.site_id}`);
                })
                .then((response) => response.json())
                .then((buildings) => {
                    const buildingSelect = $(".buildingInput");

                    if (!buildings || buildings.length === 0) {
                        buildingSelect.html(
                            `<option value="" disabled selected>No buildings for site</option>`
                        );
                        buildingSelect.prop("disabled", true);
                    } else {
                        buildingSelect.prop("disabled", false);
                        const buildingOptions = buildings.map(
                            (building) =>
                                `<option value="${building.building_id}" 
                                 ${
                                     building.building_id === room.building_id
                                         ? "selected"
                                         : ""
                                 }>
                                 ${building.building_code}
                                 </option>`
                        );
                        buildingSelect.html(
                            `<option value="">Select a Building</option>` +
                                buildingOptions.join("")
                        );
                    }

                    // Set other form fields after dropdowns are populated
                    document.getElementById("editRoomID").value = room.room_id;
                    document.getElementById("editRoomCode").value =
                        room.room_code;
                });
        })
        .catch((error) => {
            console.error("Fetch error: ", error);
        });

    // Show the modal
    $("#editRoomModal").modal("show");

    // Clear validation classes
    $("#editRoomForm").removeClass("was-validated");

    // Event handler for site changes
    $(".siteInput")
        .off("change")
        .on("change", function () {
            const siteId = $(this).val();
            const buildingSelect = $(".buildingInput");

            if (!siteId) {
                buildingSelect.html(
                    `<option value="">Select a Building</option>`
                );
                buildingSelect.prop("disabled", true);
                return;
            }

            fetch(`/api/building?siteId=${siteId}`)
                .then((response) => response.json())
                .then((buildings) => {
                    if (!buildings || buildings.length === 0) {
                        buildingSelect.html(
                            `<option value="" disabled selected>No buildings for site</option>`
                        );
                        buildingSelect.prop("disabled", true);
                    } else {
                        buildingSelect.prop("disabled", false);
                        const buildingOptions = buildings.map(
                            (building) =>
                                `<option value="${building.building_id}">
                             ${building.building_code}
                             </option>`
                        );
                        buildingSelect.html(
                            `<option value="">Select a Building</option>` +
                                buildingOptions.join("")
                        );
                    }
                })
                .catch((error) => {
                    console.error("Error loading buildings:", error);
                    buildingSelect.html(
                        `<option value="" disabled selected>Error loading buildings</option>`
                    );
                    buildingSelect.prop("disabled", true);
                });
        });

    var editRoomForm = document.getElementById("editRoomForm");

    // Add change event listener to building select
    buildingInput.addEventListener("change", function () {
        // Clear validation state when user selects a valid building
        if (this.value !== "0" && this.value !== "") {
            this.setCustomValidity("");
            this.classList.remove("is-invalid");
            editRoomForm.classList.remove("was-validated");
        }
    });

    // Add change event listener to site select to reset building validation
    document
        .querySelector(".siteInput")
        .addEventListener("change", function () {
            // Reset building validation state when site changes
            buildingInput.setCustomValidity("");
            buildingInput.classList.remove("is-invalid");
            editRoomForm.classList.remove("was-validated");
        });

    // Event handler for edit room button
    $("#editRoomBtn").click(function (event) {
        // check if the form is valid
        event.preventDefault();
        event.stopPropagation();

        // Get the building select element
        var isDisabled = buildingInput.disabled;
        var hasNoOptions = buildingInput.options.length <= 1; // Only has default option

        // Clear previous validation state
        buildingInput.setCustomValidity("");

        if (isDisabled || hasNoOptions) {
            buildingInput.setCustomValidity(
                "No buildings available for selected site"
            );
            buildingInput.classList.add("is-invalid");
            var invalidFeedback = buildingInput.nextElementSibling;
            invalidFeedback.textContent =
                "No buildings available for selected site";
        } else if (buildingInput.value === "0" || buildingInput.value === "") {
            buildingInput.setCustomValidity(
                "Room must be assigned to a building"
            );
            buildingInput.classList.add("is-invalid");
            var invalidFeedback = buildingInput.nextElementSibling;
            invalidFeedback.textContent = "Room must be assigned to a building";
        }

        editRoomForm.classList.add("was-validated");

        // If form is valid (including our custom validation), p it
        if (editRoomForm.checkValidity() && buildingInput.validity.valid) {
            // If the form is valid, prepare to send the PUT request
            const formData = new FormData(editRoomForm);
            const jsonData = {};
            for (const [key, value] of formData.entries()) {
                jsonData[key] = value;
            }
            fetch(`/api/room/${document.getElementById("editRoomID").value}`, {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(jsonData),
            })
                .then((response) => response.json())
                .then((data) => {
                    if (data.error) {
                        window.location.href = data.redirectURL;
                    } else if (data.message) {
                        window.location.href = data.redirectURL;
                    } else {
                        console.error("Unexpected response:", data);
                        // Handle unexpected responses (e.g., show an error message)
                        throw new Error("Unexpected response");
                    }
                })
                .catch((error) => {
                    console.error("Fetch error:", error);
                    // Optionally display a user-friendly error message
                });
        }
    });
}

// Fetch device types from the server
fetch("/api/emergency-device-type")
    .then((response) => response.json())
    .then((deviceTypes) => {
        const deviceTypeRows = deviceTypes.map(
            (deviceType) => `
        <tr>
            <td data-label="Device Type">${deviceType.emergency_device_type_name}</td>
            <td>
                <div class="btn-group">
                    <button class="btn btn-warning p-2 edit-device-type-button" onclick="editDeviceType(${deviceType.emergency_device_type_id})"
                            title="Edit Device Type">
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" 
                            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M17 3a2.85 2.83 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/>
                            <path d="m15 5 4 4"/>
                        </svg>
                    </button>
                    <button class="btn btn-danger p-2 delete-button" 
                            onclick="showDeleteModal(${deviceType.emergency_device_type_id}, 'emergency-device-type', '<br>${deviceType.emergency_device_type_name}')" 
                            data-id="${deviceType.emergency_device_type_id}" 
                            title="Delete Device Type">
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" 
                            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M3 6h18"/>
                            <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/>
                            <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/>
                            <line x1="10" y1="11" x2="10" y2="17"/>
                            <line x1="14" y1="11" x2="14" y2="17"/>
                        </svg>
                    </button>
                </div>
            </td>
        </tr>
        `
        );

        // Add the rows to the device types table
        $("#device-types-table tbody").html(deviceTypeRows.join(""));
    });

export function editDeviceType(deviceTypeId) {
    const id = deviceTypeId;
    document.getElementById("editDeviceTypeID").value = id;

    //Fetch device type by id and autofill form
    fetch(`/api/emergency-device-type/${id}`)
        .then((response) => response.json())
        .then((data) => {
            //Populate the form with the data
            document.getElementById("editDeviceTypeName").value =
                data.emergency_device_type_name;
        })
        .catch((error) => {
            console.error("Fetch error: ", error);
        });

    // Handle edit
    $("#editDeviceTypeModal").modal("show");

    // Clear validation classes
    $("#editDeviceTypeForm").removeClass("was-validated");

    var editDeviceTypeForm = document.getElementById("editDeviceTypeForm");

    // Function to handle form submission
    function handleSubmit(event) {
        event.preventDefault(); // Prevent actual form submission

        // Check if the form is valid
        if (!editDeviceTypeForm.checkValidity()) {
            event.stopPropagation();
            editDeviceTypeForm.classList.add("was-validated");
        } else {
            // If the form is valid, prepare to send the PUT request
            const formData = new FormData(editDeviceTypeForm);
            const jsonData = {};
            for (const [key, value] of formData.entries()) {
                jsonData[key] = value;
            }
            fetch(
                `/api/emergency-device-type/${
                    document.getElementById("editDeviceTypeID").value
                }`,
                {
                    method: "PUT",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(jsonData),
                }
            )
                .then((response) => response.json())
                .then((data) => {
                    if (data.error) {
                        window.location.href = data.redirectURL;
                    } else if (data.message) {
                        window.location.href = data.redirectURL;
                    } else {
                        console.error("Unexpected response:", data);
                        // Handle unexpected responses (e.g., show an error message)
                        throw new Error("Unexpected response");
                    }
                })
                .catch((error) => {
                    console.error("Fetch error:", error);
                    // Optionally display a user-friendly error message
                });
        }
    }

    // Add event listener to the form for submit event (triggered by Enter key)
    $(editDeviceTypeForm).off("submit").on("submit", handleSubmit);

    // Add event listener to the submit button
    $("#editDeviceTypeBtn").off("click").on("click", handleSubmit);
}

// Function to edit a site in the database
export function editSite(siteId) {
    // Clear the form
    $("#editSiteForm")[0].reset();
    $("#currentSiteMapContainer").hide();
    $("#editSiteImagePreviewContainer").hide();

    // Fetch the site data from the server
    fetch(`/api/site/${siteId}`)
        .then((response) => response.json())
        .then((site) => {
            // Fill in the form with the site data
            $("#editSiteForm input[name=editSiteID]").val(site.site_id);
            $("#editSiteForm input[name=editSiteName]").val(site.site_name);
            $("#editSiteForm input[name=editSiteAddress]").val(
                site.site_address
            );

            // Set the form action to the update endpoint for this site
            $("#editSiteForm").attr("action", `/api/site/${site.site_id}`);

            // Set the form method to POST
            $("#editSiteForm").attr("method", "POST");

            // Check if the site is the main site
            // If it is, hide the image input
            if (site.site_id != 1) {
                $("#editSiteImgInput").show();
            }

            // Check if the image path is valid
            if (site.site_map_image_path.Valid) {
                // Set the image source to the image path
                $("#editSiteForm img[name=currentSiteMap]").attr(
                    "src",
                    site.site_map_image_path.String
                );
                $("#currentSiteMapContainer").show();
            }
            if (site.site_id == 1) {
                $("#editSiteForm img[name=currentSiteMap]").attr(
                    "src",
                    "/static/site_maps/EIT_Taradale.svg"
                );
                $("#currentSiteMapContainer").show();
                $("#editSiteImgInput").hide();
            }
        });
    // Show the modal

    $("#editSiteModal").modal("show");
}

(function () {
    "use strict";

    // Fetch the form and the submit button
    var form = document.querySelector("#addSiteForm");
    var submitButton = document.querySelector("#addSiteBtn");

    // Add event listener to the submit button
    submitButton.addEventListener(
        "click",
        function (event) {
            if (!form.checkValidity()) {
                event.preventDefault();
                event.stopPropagation();
            } else {
                // If the form is valid, submit it
                form.submit();
            }

            form.classList.add("was-validated");
        },
        false
    );
})();

// Validate and submit the form for Device Type Name
(function () {
    "use strict";

    // Fetch the form and the submit button
    var form = document.querySelector("#addDeviceTypeForm");
    var submitButton = document.querySelector("#addDeviceTypeBtn");

    // Add event listener to the submit button
    submitButton.addEventListener(
        "click",
        function (event) {
            if (!form.checkValidity()) {
                event.preventDefault();
                event.stopPropagation();
            } else {
                // If the form is valid, submit it
                form.submit();
            }

            form.classList.add("was-validated");
        },
        false
    );
})();

export function AddBuilding() {
    // Clear the form before showing it
    document.getElementById("addBuildingForm").reset();
    document
        .getElementById("addBuildingForm")
        .classList.remove("was-validated");

    // populate the site select dropdown
    populateDropdown(
        ".siteInput",
        "/api/site",
        "Select a Site",
        "site_id",
        "site_name"
    );
}

export function AddRoom() {
    // Clear the form before showing it
    document.getElementById("addRoomForm").reset();
    document.getElementById("addRoomForm").classList.remove("was-validated");

    // Clear custom validation and invalid classes from building select
    const buildingInput = document.getElementById("addRoomBuildingCode");
    buildingInput.setCustomValidity("");
    buildingInput.disabled = false;
    buildingInput.classList.remove("is-invalid");

    // Clear the building options
    $(".buildingInput").html(
        `<option value="" disabled selected>Select a Building</option>`
    );

    // populate the site select dropdown
    populateDropdown(
        ".siteInput",
        "/api/site",
        "Select a Site",
        "site_id",
        "site_name"
    );

    // Filter buildings based on the selected site
    $(".siteInput").change(function () {
        // populate the building select dropdown
        populateDropdown(
            ".buildingInput",
            "/api/building",
            "Select a Building",
            "building_id",
            "building_code"
        );

        var siteId = document.getElementById("addRoomSite").value;

        // Fetch the buildings for the selected site
        fetch(`/api/building?siteId=${siteId}`)
            .then((response) => response.json())
            .then((buildings) => {
                // Check if buildings is null
                if (buildings === null) {
                    // Add a default disabled option for "No buildings"
                    $(".buildingInput").html(
                        `<option value="0" disabled selected>No buildings for site</option>`
                    );

                    // Disable the building dropdown
                    $(".buildingInput").prop("disabled", true);
                } else {
                    // Enable the building dropdown
                    $(".buildingInput").prop("disabled", false);
                    // Create a dropdown option for each building
                    const buildingOptions = buildings.map(
                        (building) =>
                            `<option value="${building.building_id}">${building.building_code}</option>`
                    );

                    // Add the options to the building dropdown
                    $(".buildingInput").html(
                        `<option value="">Select a Building</option>` +
                            buildingOptions.join("")
                    );
                }
            });
    });
}

(function () {
    "use strict";

    var form = document.querySelector("#addRoomForm");
    var submitButton = document.querySelector("#addRoomBtn");
    var buildingInput = document.querySelector(".buildingInput");

    // Add change event listener to building select
    buildingInput.addEventListener("change", function () {
        // Clear validation state when user selects a valid building
        if (this.value !== "0" && this.value !== "") {
            this.setCustomValidity("");
            this.classList.remove("is-invalid");
            form.classList.remove("was-validated");
        }
    });

    // Add change event listener to site select to reset building validation
    document
        .querySelector(".siteInput")
        .addEventListener("change", function () {
            // Reset building validation state when site changes
            buildingInput.setCustomValidity("");
            buildingInput.classList.remove("is-invalid");
            form.classList.remove("was-validated");
        });

    submitButton.addEventListener(
        "click",
        function (event) {
            event.preventDefault();
            event.stopPropagation();

            // Get the building select element
            var isDisabled = buildingInput.disabled;
            var hasNoOptions = buildingInput.options.length <= 1; // Only has default option

            // Clear previous validation state
            buildingInput.setCustomValidity("");

            if (isDisabled || hasNoOptions) {
                buildingInput.setCustomValidity(
                    "No buildings available for selected site"
                );
                buildingInput.classList.add("is-invalid");
                var invalidFeedback = buildingInput.nextElementSibling;
                invalidFeedback.textContent =
                    "No buildings available for selected site";
            } else if (
                buildingInput.value === "0" ||
                buildingInput.value === ""
            ) {
                buildingInput.setCustomValidity(
                    "Room must be assigned to a building"
                );
                buildingInput.classList.add("is-invalid");
                var invalidFeedback = buildingInput.nextElementSibling;
                invalidFeedback.textContent =
                    "Room must be assigned to a building";
            }

            form.classList.add("was-validated");

            // If form is valid (including our custom validation), submit it
            if (form.checkValidity() && buildingInput.validity.valid) {
                form.submit();
            }
        },
        false
    );
})();

(function () {
    "use strict";

    var form = document.querySelector("#editRoomForm");
    var submitButton = document.querySelector("#editRoomBtn");
    var buildingInput = document.querySelector(".buildingInput");

    // Add change event listener to building select
    buildingInput.addEventListener("change", function () {
        // Clear validation state when user selects a valid building
        if (this.value !== "0" && this.value !== "") {
            this.setCustomValidity("");
            this.classList.remove("is-invalid");
            form.classList.remove("was-validated");
        }
    });

    // Add change event listener to site select to reset building validation
    document
        .querySelector(".siteInput")
        .addEventListener("change", function () {
            // Reset building validation state when site changes
            buildingInput.setCustomValidity("");
            buildingInput.classList.remove("is-invalid");
            form.classList.remove("was-validated");
        });

    submitButton.addEventListener(
        "click",
        function (event) {
            event.preventDefault();
            event.stopPropagation();

            // Get the building select element
            var isDisabled = buildingInput.disabled;
            var hasNoOptions = buildingInput.options.length <= 1; // Only has default option

            // Clear previous validation state
            buildingInput.setCustomValidity("");

            if (isDisabled || hasNoOptions) {
                buildingInput.setCustomValidity(
                    "No buildings available for selected site"
                );
                buildingInput.classList.add("is-invalid");
                var invalidFeedback = buildingInput.nextElementSibling;
                invalidFeedback.textContent =
                    "No buildings available for selected site";
            } else if (
                buildingInput.value === "0" ||
                buildingInput.value === ""
            ) {
                buildingInput.setCustomValidity(
                    "Room must be assigned to a building"
                );
                buildingInput.classList.add("is-invalid");
                var invalidFeedback = buildingInput.nextElementSibling;
                invalidFeedback.textContent =
                    "Room must be assigned to a building";
            }

            form.classList.add("was-validated");

            // If form is valid (including our custom validation), submit it
            if (form.checkValidity() && buildingInput.validity.valid) {
                form.submit();
            }
        },
        false
    );
})();

(function () {
    "use strict";

    // Fetch the form and the submit button
    var form = document.querySelector("#addBuildingForm");
    var submitButton = document.querySelector("#addBuildingBtn");

    // Add event listener to the submit button
    submitButton.addEventListener(
        "click",
        function (event) {
            if (!form.checkValidity()) {
                event.preventDefault();
                event.stopPropagation();
            } else {
                // If the form is valid, submit it
                form.submit();
            }

            form.classList.add("was-validated");
        },
        false
    );
})();

(function () {
    "use strict";

    // Fetch the form and the submit button
    var form = document.querySelector("#editSiteForm");
    var submitButton = document.querySelector("#editSiteBtn");

    // Add event listener to the submit button
    submitButton.addEventListener(
        "click",
        function (event) {
            if (!form.checkValidity()) {
                event.preventDefault();
                event.stopPropagation();
            } else {
                // If the form is valid, submit it
                form.submit();
            }

            form.classList.add("was-validated");
        },
        false
    );
})();

// Make functions available globally
window.editDeviceType = editDeviceType;
window.editUser = editUser;
window.editBuilding = editBuilding;
window.editRoom = editRoom;
window.editSite = editSite;
window.AddBuilding = AddBuilding;
window.AddRoom = AddRoom;
