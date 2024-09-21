// Leaflet map setup
let map;

function initializeMap(containerId, options = {}) {
    const defaultOptions = {
        crs: L.CRS.Simple,
        minZoom: -1,
    };
    map = L.map(containerId, { ...defaultOptions, ...options });
}

function createEitTaradaleMap() {
    const svgDimensions = { width: 561.568, height: 962.941 };
    const minCoordinates = { x: 128.009, y: 82.331 };
    const imageUrl = "/static/map.svg";
    const bounds = [
        [0, 0],
        [svgDimensions.height, svgDimensions.width],
    ];

    L.imageOverlay(imageUrl, bounds).addTo(map);

    fetchBuildingsData()
        .then((data) => renderBuildings(data, svgDimensions, minCoordinates))
        .catch((error) =>
            console.error("Error fetching building data:", error)
        );

    map.fitBounds(bounds);
}

function fetchBuildingsData() {
    return fetch("static/buildings.json").then((response) => response.json());
}

function renderBuildings(data, svgDimensions, minCoordinates) {
    data.buildings.forEach((building) => {
        const x = building.coordinates.x - minCoordinates.x;
        const y =
            svgDimensions.height - (building.coordinates.y - minCoordinates.y);

        const rectangle = L.rectangle([
            [y - 19, x],
            [y, x + 19],
        ]).addTo(map);

        rectangle.on("click", () => {
            GetAllDevices(building.name);
            console.log("Building clicked:", building.name);
        });
    });
}

function GetFilterOptions() {
    fetchAndPopulateSelect(
        "/api/site",
        "siteFilter",
        "site_name",
        "site_id",
        "All Sites"
    );
    setupBuildingFilter();
    setupRoomFilter();
    fetchAndPopulateSelect(
        "/api/emergency-device-type",
        "deviceTypeFilter",
        "emergency_device_type_name",
        null,
        "All Device Types"
    );
}

function fetchAndPopulateSelect(
    url,
    selectId,
    textKey,
    valueKey,
    defaultOptionText
) {
    fetch(url)
        .then((response) => response.json())
        .then((data) => {
            const select = document.getElementById(selectId);
            select.innerHTML = "";

            // Add the default option
            addDefaultOption(select, defaultOptionText);

            // Check if data is valid and is an array
            if (data && Array.isArray(data) && data.length > 0) {
                data.forEach((item) => {
                    const option = document.createElement("option");
                    option.text = item[textKey];
                    if (valueKey) option.value = item[valueKey];
                    select.add(option);
                });
            } else {
                console.log(`No data available for ${selectId}`);
            }
        })
        .catch((error) => {
            console.error(`Error fetching ${selectId} data:`, error);
        });
}

function addDefaultOption(select, text) {
    const defaultOption = document.createElement("option");
    defaultOption.text = text;
    defaultOption.selected = true;
    select.add(defaultOption);
}

function setupBuildingFilter() {
    document.getElementById("siteFilter").addEventListener("change", () => {
        const selectedSite = document.getElementById("siteFilter").value;
        fetchAndPopulateSelect(
            `/api/building?siteId=${selectedSite}`,
            "buildingFilter",
            "building_code",
            "building_id",
            "All Buildings"
        );
        clearRoomFilter();
    });
}

function setupRoomFilter() {
    document.getElementById("buildingFilter").addEventListener("change", () => {
        const selectedBuilding =
            document.getElementById("buildingFilter").value;
        fetchAndPopulateSelect(
            `/api/room?buildingId=${selectedBuilding}`,
            "roomFilter",
            "room_code",
            null,
            "All Rooms"
        );
    });
}

function clearRoomFilter() {
    const roomSelect = document.getElementById("roomFilter");
    roomSelect.innerHTML = "";
    addDefaultOption(roomSelect, "All Rooms");
}

function FilterBySite() {
    const siteName =
        document.getElementById("siteFilter").selectedOptions[0].text;
    const siteId = document.getElementById("siteFilter").value;
    console.log("Site Name:", siteName, "Site ID:", siteId);

    const mapElement = document.getElementById("map");
    const toggleMapButton = document.getElementById("toggleMap");

    // Clear the table body
    clearTableBody();

    if (siteName === "All Sites") {
        console.log("Filter by site: All Sites");
        mapElement.classList.add("d-none");
        toggleMapButton.classList.add("d-none");
        // Change the device list width to col-xxl-12
        document.querySelector(".device-list").classList.remove("col-xxl-9");
        document.querySelector(".device-list").classList.add("col-xxl-12");
        GetAllDevices();
        return;
    }

    if (siteId === "1") {
        // EIT Taradale should always also be id = 1, as its the first site inserted into the database (see seed.go)
        console.log("Filter by site: EIT Taradale");
        // Clear the map layers
        map.eachLayer((layer) => {
            if (
                layer instanceof L.ImageOverlay ||
                layer instanceof L.Rectangle
            ) {
                map.removeLayer(layer);
            }
        });
        mapElement.classList.remove("d-none");
        toggleMapButton.classList.remove("d-none");
        // Change the device list width to col-xxl-9
        document.querySelector(".device-list").classList.remove("col-xxl-12");
        document.querySelector(".device-list").classList.add("col-xxl-9");
        createEitTaradaleMap();
        GetAllDevices("", siteId);
        return;
    }

    GetAllDevices("", siteId);
    clearRoomFilter();
    updateMapForSite(siteId);
}

function clearTableBody() {
    const tableBody = document.getElementById("emergency-device-body");
    if (tableBody) {
        tableBody.innerHTML = "";
    } else {
        console.error("Table body element not found");
    }
}

function updateMapForSite(siteId) {
    fetch(`/api/site/${siteId}`)
        .then((response) => response.json())
        .then((data) => {
            console.log("Site Data:", data);
            document.getElementById("map").classList.remove("d-none");
            document.getElementById("toggleMap").classList.remove("d-none");
            // Change the device list width to col-xxl-9
            document
                .querySelector(".device-list")
                .classList.remove("col-xxl-12");
            document.querySelector(".device-list").classList.add("col-xxl-9");

            const imageUrl = data.site_map_image_path.String;
            console.log("Image URL:", imageUrl);

            const image = new Image();
            image.src = imageUrl;
            image.onload = function () {
                const imgWidth = this.width;
                const imgHeight = this.height;
                console.log("Image dimensions:", imgWidth, "x", imgHeight);

                const newBounds = [
                    [0, 0],
                    [imgHeight, imgWidth],
                ];

                map.eachLayer((layer) => {
                    if (
                        layer instanceof L.ImageOverlay ||
                        layer instanceof L.Rectangle
                    ) {
                        map.removeLayer(layer);
                    }
                });

                L.imageOverlay(imageUrl, newBounds).addTo(map);
                map.fitBounds(newBounds);
            };
        })
        .catch((error) => console.error("Error updating map:", error));
}

// Update the event listener to include table clearing
document.getElementById("siteFilter").addEventListener("change", () => {
    FilterBySite();
    clearTableBody();
});

// Initialize the map and populate filter options
initializeMap("map");
GetFilterOptions();

console.log(role);

async function GetAllDevices(buildingCode = "", siteId = "") {
    try {
        let url = "/api/emergency-device";
        const params = new URLSearchParams();
        if (buildingCode) params.append("building_code", buildingCode);
        if (siteId) params.append("site_id", siteId);
        if (params.toString()) url += `?${params.toString()}`;

        const response = await fetch(url);
        console.log("Response:", response);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const devices = await response.json();
        console.log("Devices:", devices);

        const tbody = document.getElementById("emergency-device-body");
        if (!tbody) {
            console.error("Table body element not found");
            return;
        }

        if (!Array.isArray(devices) || devices.length === 0) {
            tbody.innerHTML = `<tr><td colspan="11" class="text-center">No devices found for this site.</td></tr>`;
        } else {
            tbody.innerHTML = devices.map(formatDeviceRow).join("");
        }
    } catch (err) {
        console.error("Failed to fetch devices:", err);
        const tbody = document.getElementById("emergency-device-body");
        if (tbody) {
            tbody.innerHTML = `<tr><td colspan="11" class="text-center">Error fetching devices. Please try again.</td></tr>`;
        }
    }
}

function formatDeviceRow(device) {
    if (!device) return "";
    const formatDateMonthYear = (dateString) =>
        formatDate(dateString, { year: "numeric", month: "long" });
    const formatDateFull = (dateString) =>
        formatDate(dateString, {
            year: "numeric",
            month: "long",
            day: "numeric",
        });

    const badgeClass = getBadgeClass(device.status.String);
    const buttons = getActionButtons(device);

    const isAdmin = role === "admin";

    return `
        <tr>
            <td data-label="Device Type">${
                device.emergency_device_type_name
            }</td>
            <td data-label="Extinguisher Type">${
                device.extinguisher_type_name.String
            }</td>
            <td data-label="Room">${device.room_code}</td>
            <td data-label="Serial Number">${device.serial_number.String}</td>
            <td data-label="Manufacture Date">${formatDateMonthYear(
                device.manufacture_date.Time
            )}</td>
            <td data-label="Expire Date">${formatDateMonthYear(
                device.expire_date.Time
            )}</td>
            ${
                isAdmin
                    ? `<td data-label="Last Inspection Date">${formatDateFull(
                          device.last_inspection_date.Time
                      )}</td>`
                    : ""
            }
            ${
                isAdmin
                    ? `<td data-label="Next Inspection Date">${formatDateFull(
                          device.next_inspection_date.Time
                      )}</td>`
                    : ""
            }
            <td data-label="Size">${device.size.String}</td>
            <td data-label="Status">
                <span class="badge ${badgeClass}">${device.status.String}</span>
            </td>
            <td>
                <div class="action-buttons">
                    ${buttons}
                </div>
            </td>
        </tr>
    `;
}

function formatDate(dateString, options) {
    if (!dateString || dateString === "0001-01-01T00:00:00Z") {
        return "N/A";
    }
    return new Date(dateString).toLocaleDateString("en-NZ", options);
}

function getBadgeClass(status) {
    switch (status) {
        case "Active":
            return "text-bg-success";
        case "Expired":
            return "text-bg-danger";
        default:
            return "text-bg-warning";
    }
}

function getActionButtons(device) {
    let buttons = `<button class="btn btn-primary" onclick="DeviceNotes('${device.description.String}')">Notes</button>`;
    if (role === "admin") {
        buttons += `
            <button class="btn btn-secondary" onclick="ViewDeviceInspection(${device.emergency_device_id})">Inspect</button>
            <button class="btn btn-warning" onclick="EditDevice(${device.emergency_device_id})">Edit</button>
            <button class="btn btn-danger" onclick="DeleteDevice(${device.emergency_device_id})">Delete</button>
        `;
    }
    return buttons;
}

// Initial fetch without filtering
GetAllDevices();

function AddDevice() {
    // Fetch the sites and populate the select options
    fetch("/api/site")
        .then((response) => response.json())
        .then((data) => {
            const select = document.getElementById("site");
            // Clear previous options
            select.innerHTML = "";
            // Add a default option and select it
            const defaultOption = document.createElement("option");
            defaultOption.text = "Select a Site";
            defaultOption.value = "";
            defaultOption.selected = true;
            defaultOption.disabled = true;
            select.add(defaultOption);
            data.forEach((item) => {
                const option = document.createElement("option");
                option.text = item.site_name;
                option.value = item.site_id;
                select.add(option);
            });
        })
        .catch((error) => console.error("Error:", error));

    // Function to clear building and room options
    function clearBuildingAndRoom() {
        const buildingSelect = document.getElementById("building");
        const roomSelect = document.getElementById("room");
        buildingSelect.innerHTML =
            "<option value='' selected disabled>Select a Building</option>";
        roomSelect.innerHTML =
            "<option value='' selected disabled>Select a Room</option>";
    }

    // Function to fetch and populate buildings
    function fetchAndPopulateBuildings(siteId) {
        fetch(`/api/building?siteId=${siteId}`)
            .then((response) => response.json())
            .then((data) => {
                const select = document.getElementById("building");
                select.innerHTML =
                    "<option value='' selected disabled>Select a Building</option>";
                data.forEach((item) => {
                    const option = document.createElement("option");
                    option.text = item.building_code;
                    option.value = item.building_id;
                    select.add(option);
                });

                // If there's only one building, select it automatically
                if (data.length === 1) {
                    select.value = data[0].building_id;
                    select.dispatchEvent(new Event("change"));
                }
            })
            .catch((error) => console.error("Error:", error));
    }

    // Function to fetch and populate rooms
    function fetchAndPopulateRooms(buildingId) {
        fetch(`/api/room?buildingId=${buildingId}`)
            .then((response) => response.json())
            .then((data) => {
                const select = document.getElementById("room");
                select.innerHTML =
                    "<option value='' selected disabled>Select a Room</option>";
                data.forEach((item) => {
                    const option = document.createElement("option");
                    option.text = item.room_code;
                    option.value = item.room_id;
                    select.add(option);
                });
            })
            .catch((error) => console.error("Error:", error));
    }

    // Event listener for site change
    document.getElementById("site").addEventListener("change", (event) => {
        const selectedSiteId = event.target.value;
        clearBuildingAndRoom();

        if (selectedSiteId) {
            fetchAndPopulateBuildings(selectedSiteId);
        }
    });

    // Event listener for building change
    document.getElementById("building").addEventListener("change", (event) => {
        const selectedBuildingId = event.target.value;

        if (selectedBuildingId) {
            fetchAndPopulateRooms(selectedBuildingId);
        }
    });

    // Fetch the device types and populate the select options
    fetch("/api/emergency-device-type")
        .then((response) => response.json())
        .then((data) => {
            const select = document.getElementById("status");
            // Clear previous options
            select.innerHTML = "";
            // Add a default option and select it
            const defaultOption = document.createElement("option");
            defaultOption.text = "Select Device type";
            defaultOption.selected = true;
            defaultOption.disabled = true;
            select.add(defaultOption);
            data.forEach((item) => {
                const option = document.createElement("option");
                option.text = item.emergency_device_type_name; // Set the text of the option
                select.add(option);
            });
        })
        .catch((error) => console.error("Error:", error));

    // Fetch the extinguisher types and populate the select options
    fetch("/api/extinguisher-type")
        .then((response) => response.json())
        .then((data) => {
            const select = document.getElementById("extinguisherType");
            // Clear previous options
            select.innerHTML = "";
            // Add a default option and select it
            const defaultOption = document.createElement("option");
            defaultOption.text = "Select Extinguisher Type";
            defaultOption.selected = true;
            defaultOption.disabled = true;
            select.add(defaultOption);
            data.forEach((item) => {
                const option = document.createElement("option");
                option.text = item.extinguisher_type_name; // Set the text of the option
                select.add(option);
            });
        })
        .catch((error) => console.error("Error:", error));

    // Show the modal after populating the select options
    $("#addModal").modal("show");
}

function EditDevice(deviceId) {
    console.log(`Edit device with ID: ${deviceId}`);
    // Add your edit logic here
}

function DeleteDevice(deviceId) {
    console.log(`Delete device with ID: ${deviceId}`);
    // Add your delete logic here
}

// Change to add inspection
function ViewDeviceInspection(deviceId) {
    console.log(`Inspect device with ID: ${deviceId}`);

    // Show the modal
    $("#viewInspectionModal").modal("show");
}

function ViewInspectionDetails(inspectionId) {
    console.log(`View inspection details for inspection ID: ${inspectionId}`);
    // Add your view inspection details logic here
}

function AddInspection() {
    // Close the view inspection modal
    $("#viewInspectionModal").modal("hide");

    // Show the modal
    $("#addInspectionModal").modal("show");
}

function DeviceNotes(description) {
    // Populate the modal with the description
    document.getElementById("notesModalBody").innerText = description;

    // Show the modal
    $("#notesModal").modal("show");
}

// Function to toggle the map visibility
function ToggleMap() {
    var map = document.getElementById("map");
    var deviceList = document.querySelector(".device-list");

    // Check if the map is currently visible
    if (map.classList.contains("d-none")) {
        // Map is hidden, show the map and set device list back to col-xxl-9 width
        map.classList.remove("d-none");
        map.classList.add("col-xxl-3");
        deviceList.classList.remove("col-xxl-12");
        deviceList.classList.add("col-xxl-9");
    } else {
        // Map is visible, hide the map and make device list 100% width
        map.classList.add("d-none");
        deviceList.classList.remove("col-xxl-9");
        deviceList.classList.add("col-xxl-12");
    }
}
