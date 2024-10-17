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
        5;
    });
});

console.log(is_current_user_default_admin);

// whenever is_current_user_default_admin is false, hide the actions buttons from any row with user.default_admin = true

// Fetch users from the server
fetch("/api/user")
    .then((response) => response.json())
    .then((users) => {
        console.log("Original users:", users);
        console.log(
            "Current user ID:",
            current_user_id,
            "Type:",
            typeof current_user_id
        );

        // Convert current_user_id to a number
        const currentUserIdNumber = parseInt(current_user_id, 10);

        // Log the converted current user ID
        console.log(
            "Converted current user ID:",
            currentUserIdNumber,
            "Type:",
            typeof currentUserIdNumber
        );

        // Sort the users array to put the current user first
        users.sort((a, b) => {
            if (a.user_id === currentUserIdNumber) return -1;
            if (b.user_id === currentUserIdNumber) return 1;
            return a.username.localeCompare(b.username); // Sort others alphabetically
        });

        console.log("Sorted users:", users);

        // Create a table row for each user
        const userRows = users.map((user) => {
            console.log("Processing user:", user);
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
                    : `<button class="btn btn-primary edit-user-button" data-id="${
                          user.user_id
                      }">Edit</button>
                       ${
                           hideDelete
                               ? ""
                               : `<button class="btn btn-danger delete-button" onclick="showDeleteModal(${user.user_id}, 'user', '${user.username}', '${currentUserIdNumber}')" data-id="${user.user_id}">Delete</button>`
                       }`
            }
        </div>
    </td>
</tr>
`;
        });

        // Add the rows to the users table
        $("#users-table tbody").html(userRows.join(""));

        // Add event listeners to the edit and delete buttons
        $(".edit-user-button").click(async (event) => {
            const id = $(event.target).data("id");
            console.log("Edit button clicked for user with ID:", id);
            // Handle edit
            // Fetch the user data from the nearest row
            const row = $(event.target).closest("tr");
            const username = row.find("td[data-label=Username]").text();
            const email = row.find("td[data-label=Email]").text();
            const role = row.find("td[data-label=Role]").text();

            const default_admin = await fetch(`/api/user/${username}`)
                .then((response) => response.json())
                .then((user) => {
                    console.log("User data:", user);
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

            console.log("Updated user ID:", updatedUserId);
            console.log("Current user ID:", current_user_id);

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
                    console.log("JSON data:", jsonData);
                    fetch(
                        `/api/user/${
                            document.getElementById("editUserID").value
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
                            console.log("Success:", data);
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

            var editDeviceTypeForm =
                document.getElementById("editDeviceTypeForm");

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
                    console.log("JSON data:", jsonData);
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
                            console.log("Success:", data);
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

/*
(function () {
    "use strict";

    // Fetch the form and the submit button
    var form = document.querySelector("#editUserForm");
    var submitButton = document.querySelector("#editUserBtn");

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
*/
// Function to add a new building to the database
function addBuilding() {
    // Show the modal
    $("#addBuildingModal").modal("show");
}

// Function to edit a building in the database
function editBuilding(buildingId) {
    // Clear the form
    $("#editBuildingForm")[0].reset();
    $("#currentSiteMapContainer").hide();
    $("#editSiteImagePreviewContainer").hide();

    // Fetch the building data from the server
    fetch(`/api/building/${buildingId}`)
        .then((response) => response.json())
        .then((building) => {
            // Fill in the form with the site data
            $("#editBuildingForm input[name=editBuildingID]").val(building.building_id);
            $("#editBuildingForm input[name=editSiteID]").val(building.site_id);
            $("#editBuildingForm input[name=editBuildingCode]").val(building.building_code);

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