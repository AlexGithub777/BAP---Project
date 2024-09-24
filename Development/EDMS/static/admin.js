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

            console.log("File input changed", fileInput, file); // Debugging log

            if (!file) {
                console.log("No file selected"); // Debugging log
                $("#editSiteModal #imagePreview").attr("src", "");
                $("#editSiteModal #editSiteImagePreviewContainer").hide();
                return;
            }

            // Check if the file is an image
            if (!file.type.startsWith("image/")) {
                console.log("Non-image file selected"); // Debugging log
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
                console.log("File read successfully"); // Debugging log
                $("#editSiteModal #imagePreview").attr("src", e.target.result);
                $("#editSiteModal #editSiteImagePreviewContainer").show();
            };
            reader.onerror = function (error) {
                console.log("Error reading file:", error); // Debugging log
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

    // Debugging: Log when the edit modal is opened
    $("#editSiteModal").on("shown.bs.modal", function () {
        console.log("Edit site modal opened");
        console.log(
            "File input element:",
            $("#editSiteModal #editSiteImgInput")[0]
        );
    });
});

// Fetch users from the server
fetch("/api/user")
    .then((response) => response.json())
    .then((users) => {
        // Create a table row for each user
        const userRows = users.map(
            (user) => `
<tr>
<td data-label="Username">${user.username}</td>
<td data-label="Email">${user.email}</td>
<td data-label="Role">${user.role}</td>
<td>
    <div class="btn-group">
        <button class="btn btn-primary edit-user-button" data-id="${user.user_id}">Edit</button>
        <button class="btn btn-danger delete-button" onclick="showDeleteModal(${user.user_id}, 'user', '${user.username}')" data-id="${user.user_id}">Delete</button>
    </div>
</td>
</tr>
`
        );

        // Add the rows to the users table
        $("#users-table tbody").html(userRows.join(""));

        // Add event listeners to the edit and delete buttons
        $(".edit-user-button").click((event) => {
            const id = $(event.target).data("id");
            console.log("Edit button clicked for user with ID:", id);
            // Handle edit
            // Fetch the user data from the nearest row
            const row = $(event.target).closest("tr");
            const username = row.find("td[data-label=Username]").text();
            const email = row.find("td[data-label=Email]").text();
            const role = row.find("td[data-label=Role]").text();

            console.log("User data:", { id, username, email, role });

            // Fill in the form with the user data
            $("#editUserForm input[name=editUserID]").val(id);
            $("#editUserForm input[name=editUserUsername]").val(username);
            $("#editUserForm input[name=editUserEmail]").val(email);
            $("#editUserForm select[name=editUserRole]").val(role);

            // Show the modal
            $("#editUserModal").modal("show");

            // Add event listener to the submit button
        });
    });

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
                    <button class="btn btn-primary edit-button" onclick="editSite(${site.site_id})" data-id="${site.site_id}">Edit</button>
                    <button class="btn btn-danger delete-button" onclick="showDeleteModal(${site.site_id}, 'site', '${site.site_name}')">Delete</button>
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
        // Create a table row for each building
        const buildingRows = buildings.map(
            (building) => `
<tr>
<td data-label="Building Code">${building.building_code}</td>
<td data-label="Site Name">${building.site_name}</td>
<td>
    <div class="btn-group">
        <button class="btn btn-primary edit-building-button" data-id="${building.building_id}">Edit</button>
        <button class="btn btn-danger delete-button" onclick="showDeleteModal(${building.building_id}, 'building', '${building.building_code}')" data-id="${building.building_id}">Delete</button>
    </div> 
</td>
</tr>
`
        );

        // Add the rows to the buildings table
        $("#buildings-table tbody").html(buildingRows.join(""));

        // Add event listeners to the edit and delete buttons
        $(".edit-building-button").click((event) => {
            const id = $(event.target).data("id");
            console.log("Edit button clicked for building with ID:", id);
            // Handle edit
        });
    });

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
        <button class="btn btn-primary edit-room-button" data-id="${room.room_id}">Edit</button>
        <button class="btn btn-danger delete-button" onclick="showDeleteModal(${room.room_id}, 'room', '${room.room_code}')" data-id="${room.room_id}">Delete</button>
    </div>
</td>
</tr>
`
        );

        // Add the rows to the rooms table
        $("#rooms-table tbody").html(roomRows.join(""));

        // Add event listeners to the edit and delete buttons
        $(".edit-room-button").click((event) => {
            const id = $(event.target).data("id");
            console.log("Edit button clicked for room with ID:", id);
        });
    });

// Fetch device types from the server
fetch("/api/emergency-device-type")
    .then((response) => response.json())
    .then((deviceTypes) => {
        // Create a table row for each device type
        const deviceTypeRows = deviceTypes.map(
            (deviceType) => `
<tr>
<td data-label="Device Type">${deviceType.emergency_device_type_name}</td>
<td>
    <div class="btn-group">
        <button class="btn btn-primary edit-device-type-button" data-id="${deviceType.emergency_device_type_id}">Edit</button>
        <button class="btn btn-danger delete-button" onclick="showDeleteModal(${deviceType.emergency_device_type_id}, 'emergency-device-type', '<br>${deviceType.emergency_device_type_name}')" data-id="${deviceType.emergency_device_type_id}">Delete</button>
    </div>
</td>
</tr>
`
        );

        // Add the rows to the device types table
        $("#device-types-table tbody").html(deviceTypeRows.join(""));

        // Add event listeners to the edit and delete buttons
        $(".edit-device-type-button").click((event) => {
            const id = $(event.target).data("id");
            console.log("Edit button clicked for device type with ID:", id);
            // Handle edit
        });
    });

// Function to add a new site to the database
function addSite() {
    // Show the modal
    $("#addSiteModal").modal("show");
}

// Function to edit a site in the database
function editSite(siteId) {
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
                    "/static/map.svg"
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
