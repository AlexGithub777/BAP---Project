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
        <button class="btn btn-primary edit-button" data-id="${user.user_id}">Edit</button>
        <button class="btn btn-danger delete-button" data-id="${user.user_id}">Delete</button>
    </div>
</td>
</tr>
`
        );

        // Add the rows to the users table
        $("#users-table tbody").html(userRows.join(""));

        // Add event listeners to the edit and delete buttons
        $(".edit-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle edit
        });
        $(".delete-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle delete
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
        <button class="btn btn-primary edit-button" data-id="${site.site_id}">Edit</button>
        <button class="btn btn-danger delete-button" data-id="${site.site_id}">Delete</button>
    </div>
</td>
</tr>
`
        );

        // Add the rows to the sites table
        $("#sites-table tbody").html(siteRows.join(""));

        // Add event listeners to the edit and delete buttons
        $(".edit-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle edit
        });
        $(".delete-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle delete
        });
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
        <button class="btn btn-primary edit-button" data-id="${building.building_id}">Edit</button>
        <button class="btn btn-danger delete-button" data-id="${building.building_id}">Delete</button>
    </div> 
</td>
</tr>
`
        );

        // Add the rows to the buildings table
        $("#buildings-table tbody").html(buildingRows.join(""));

        // Add event listeners to the edit and delete buttons
        $(".edit-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle edit
        });
        $(".delete-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle delete
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
        <button class="btn btn-primary edit-button" data-id="${room.room_id}">Edit</button>
        <button class="btn btn-danger delete-button" data-id="${room.room_id}">Delete</button>
    </div>
</td>
</tr>
`
        );

        // Add the rows to the rooms table
        $("#rooms-table tbody").html(roomRows.join(""));

        // Add event listeners to the edit and delete buttons
        $(".edit-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle edit
        });
        $(".delete-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle delete
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
        <button class="btn btn-primary edit-button" data-id="${deviceType.emergency_device_type_id}">Edit</button>
        <button class="btn btn-danger delete-button" data-id="${deviceType.emergency_device_type_id}">Delete</button>
    </div>
</td>
</tr>
`
        );

        // Add the rows to the device types table
        $("#device-types-table tbody").html(deviceTypeRows.join(""));

        // Add event listeners to the edit and delete buttons
        $(".edit-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle edit
        });
        $(".delete-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle delete
        });
    });

// function to add site dropdown items to navbar
function AddSiteOptions() {
    fetch("/api/site")
        .then((response) => response.json())
        .then((data) => {
            const dropdownMenu = document.getElementById("siteDropdown");
            data.forEach((item) => {
                const listItem = document.createElement("li");
                const anchor = document.createElement("a");
                anchor.classList.add("dropdown-item");
                anchor.href = "#"; // replace with the actual link
                anchor.textContent = item.site_name;
                anchor.dataset.siteId = item.site_id; // store site_id in data attribute
                listItem.appendChild(anchor);
                dropdownMenu.appendChild(listItem);
            });
        })
        .catch((error) => console.error("Error:", error));
}

AddSiteOptions();
