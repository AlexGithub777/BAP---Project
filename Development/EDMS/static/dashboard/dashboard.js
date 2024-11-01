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
    const imageUrl = "/static/site_maps/EIT_Taradale.svg";
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
    return fetch("static/assets/buildings.json").then((response) =>
        response.json()
    );
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
            getAllDevices(building.name, "1");
            console.log("Building clicked:", building.name);
        });
    });
}

function getFilterOptions() {
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

function filterBySite() {
    const siteName =
        document.getElementById("siteFilter").selectedOptions[0].text;
    const siteId = document.getElementById("siteFilter").value;

    // Clear the table body
    clearTableBody();

    if (siteName === "All Sites") {
        hideMap();
        getAllDevices();
        return;
    }

    if (siteId === "1") {
        // Hard coded - EIT Taradale should always also be id = 1, as its the first site inserted into the database (see seed.go)
        // Clear the map layers
        map.eachLayer((layer) => {
            if (
                layer instanceof L.ImageOverlay ||
                layer instanceof L.Rectangle
            ) {
                map.removeLayer(layer);
            }
        });
        showMap();
        createEitTaradaleMap();
        getAllDevices("", siteId);
        return;
    }

    getAllDevices("", siteId);
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
            // Check if site has a map image
            if (!data.site_map_image_path.String) {
                hideMap();
                return;
            }
            showMap();

            const imageUrl = data.site_map_image_path.String;

            const image = new Image();
            image.src = imageUrl;
            image.onload = function () {
                const imgWidth = this.width;
                const imgHeight = this.height;

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

function hideMap() {
    document.getElementById("map").classList.add("d-none");
    document.getElementById("toggleMap").classList.add("d-none");
    // Change the device list width to col-xxl-12
    document.querySelector(".device-list").classList.remove("col-xxl-9");
    document.querySelector(".device-list").classList.add("col-xxl-12");
}

function showMap() {
    document.getElementById("map").classList.remove("d-none");
    document.getElementById("toggleMap").classList.remove("d-none");
    // Change the device list width to col-xxl-9
    document.querySelector(".device-list").classList.remove("col-xxl-12");
    document.querySelector(".device-list").classList.add("col-xxl-9");
}

// Update the event listener to include table clearing
document.getElementById("siteFilter").addEventListener("change", () => {
    filterBySite();
    clearTableBody();
});

// Initialize the map and populate filter options
initializeMap("map");
getFilterOptions();

let currentPage = 1;
let rowsPerPage = 10;
let allDevices = [];

async function getAllDevices(buildingCode = "", siteId = "") {
    try {
        let url = "/api/emergency-device";
        const params = new URLSearchParams();
        if (buildingCode) params.append("building_code", buildingCode);
        if (siteId) params.append("site_id", siteId);
        if (params.toString()) url += `?${params.toString()}`;

        const response = await fetch(url);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        allDevices = await response.json();

        updateTable();
    } catch (err) {
        console.error("Failed to fetch devices:", err);
        const tbody = document.getElementById("emergency-device-body");
        if (tbody) {
            tbody.innerHTML = `<tr><td colspan="11" class="text-center">Error fetching devices. Please try again.</td></tr>`;
        }
    }
}

function updateTable() {
    const tbody = document.getElementById("emergency-device-body");
    if (!tbody) {
        console.error("Table body element not found");
        return;
    }

    const startIndex = (currentPage - 1) * rowsPerPage;
    const endIndex = startIndex + rowsPerPage;
    const pageDevices = allDevices.slice(startIndex, endIndex);

    if (!Array.isArray(pageDevices) || pageDevices.length === 0) {
        tbody.innerHTML = `<tr><td colspan="11" class="text-center">No devices found.</td></tr>`;
    } else {
        tbody.innerHTML = pageDevices.map(formatDeviceRow).join("");
    }

    updatePaginationControls();
}

// JavaScript
function updatePaginationControls() {
    const totalPages = Math.ceil(allDevices.length / rowsPerPage);
    const paginationEl = document.querySelector(".pagination");
    const isMobile = window.innerWidth < 768; // Detect mobile devices

    let paginationHTML = `
        <li class="page-item ${currentPage === 1 ? "disabled" : ""}">
            <a class="page-link" href="#" data-page="${
                currentPage - 1
            }" aria-label="Previous">
                <span aria-hidden="true">&laquo;</span>
            </a>
        </li>
    `;

    function addPageNumber(pageNum) {
        paginationHTML += `
            <li class="page-item ${
                currentPage === pageNum ? "active" : ""
            }" aria-current="page">
                <a class="page-link" href="#" data-page="${pageNum}">${pageNum}</a>
            </li>
        `;
    }

    function addEllipsis() {
        paginationHTML += `
            <li class="page-item disabled">
                <span class="page-link">...</span>
            </li>
        `;
    }

    if (isMobile) {
        // Simplified pagination for mobile
        if (totalPages <= 3) {
            for (let i = 1; i <= totalPages; i++) {
                addPageNumber(i);
            }
        } else {
            addPageNumber(1);
            if (currentPage !== 1 && currentPage !== totalPages) {
                addPageNumber(currentPage);
            }
            addPageNumber(totalPages);
        }
    } else {
        // Desktop pagination (keep your existing logic here)
        if (totalPages <= 7) {
            for (let i = 1; i <= totalPages; i++) {
                addPageNumber(i);
            }
        } else {
            addPageNumber(1);
            if (currentPage > 3) addEllipsis();

            let start = Math.max(2, currentPage - 1);
            let end = Math.min(totalPages - 1, currentPage + 1);

            if (currentPage <= 3) {
                end = 4;
            } else if (currentPage >= totalPages - 2) {
                start = totalPages - 3;
            }

            for (let i = start; i <= end; i++) {
                addPageNumber(i);
            }

            if (currentPage < totalPages - 2) addEllipsis();
            addPageNumber(totalPages);
        }
    }

    paginationHTML += `
        <li class="page-item ${currentPage === totalPages ? "disabled" : ""}">
            <a class="page-link" href="#" data-page="${
                currentPage + 1
            }" aria-label="Next">
                <span aria-hidden="true">&raquo;</span>
            </a>
        </li>
    `;

    paginationEl.innerHTML = paginationHTML;

    function handlePaginationClick(e) {
        e.preventDefault();
        e.stopPropagation();

        let target = e.target.closest(".page-link");

        if (target && target.hasAttribute("data-page")) {
            const newPage = parseInt(target.getAttribute("data-page"), 10);

            if (
                newPage !== currentPage &&
                newPage > 0 &&
                newPage <= totalPages
            ) {
                currentPage = newPage;
                updateTable();
            }
        }
    }

    // Remove existing event listeners
    paginationEl.removeEventListener("click", handlePaginationClick);
    paginationEl.removeEventListener("touchstart", handlePaginationClick);

    // Add event listeners to the pagination container
    paginationEl.addEventListener("click", handlePaginationClick);
    paginationEl.addEventListener("touchstart", handlePaginationClick);
}

// Event listener for rows per page dropdown
document.getElementById("rowsPerPage").addEventListener("change", (e) => {
    rowsPerPage = parseInt(e.target.value);
    currentPage = 1; // Reset to first page when changing rows per page
    updateTable();
});

// Initial fetch without filtering
getAllDevices();

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

    // Declare isAdmin within the function
    let isAdmin = false;

    // Ensure role is defined and check for "Admin"
    if (role === "Admin") {
        isAdmin = true;
    }

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
                <div class="btn-group">
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
            return "text-bg-warning";
        case "Inspection Failed":
            return "text-bg-danger";
        case "Inactive":
            return "text-bg-secondary";
        default:
            return "text-bg-warning";
    }
}

function getActionButtons(device) {
    let buttons = `<button class="btn btn-primary" onclick="deviceNotes('${device.description.String}')">Notes</button>`;

    if (role === "Admin") {
        // Only show Inspect button if device is a Fire Extinguisher
        const isFireExtinguisher =
            device.emergency_device_type_name === "Fire Extinguisher";

        if (isFireExtinguisher) {
            buttons += `
                <button class="btn btn-secondary" onclick="viewDeviceInspections(${device.emergency_device_id})">Inspect</button>`;
        }

        buttons += `
            <button class="btn btn-warning" onclick="editDevice(${device.emergency_device_id})">Edit</button>
            <button class="btn btn-danger" onclick="showDeleteModal(${device.emergency_device_id},'emergency-device', '<br>${device.emergency_device_type_name} - Serial Number: ${device.serial_number.String}')">Delete</button>
        `;
    }
    return buttons;
}

// Function to clear building and room options
function clearBuildingAndRoom() {
    const buildingSelects = document.querySelectorAll(".buildingInput");
    const roomSelects = document.querySelectorAll(".roomInput");
    buildingSelects.forEach((select) => {
        select.innerHTML =
            "<option value='' selected disabled>Select a Building</option>";
    });
    roomSelects.forEach((select) => {
        select.innerHTML =
            "<option value='' selected disabled>Select a Room</option>";
    });
}

function addDevice() {
    document.getElementById("addDeviceForm").reset();
    document.getElementById("addDeviceForm").classList.remove("was-validated");

    const emergencyDeviceTypePromise = populateDropdown(
        ".emergencyDeviceTypeInput",
        "/api/emergency-device-type",
        "Select Device Type",
        "emergency_device_type_id",
        "emergency_device_type_name"
    );

    const extinguisherTypePromise = populateDropdown(
        ".extinguisherTypeInput",
        "/api/extinguisher-type",
        "Select Extinguisher Type",
        "extinguisher_type_id",
        "extinguisher_type_name"
    );

    const sitePromise = populateDropdown(
        ".siteInput",
        "/api/site",
        "Select a Site",
        "site_id",
        "site_name"
    );

    Promise.all([
        emergencyDeviceTypePromise,
        extinguisherTypePromise,
        sitePromise,
    ])
        .then(() => {
            // Event listener for site change
            document
                .querySelector(".siteInput")
                .addEventListener("change", (event) => {
                    const selectedSiteId = event.target.value;
                    clearBuildingAndRoom();

                    if (selectedSiteId) {
                        fetchAndPopulateBuildings(selectedSiteId);
                    }
                });

            // Event listener for building change
            document
                .querySelector(".buildingInput")
                .addEventListener("change", (event) => {
                    const selectedBuildingId = event.target.value;

                    if (selectedBuildingId) {
                        fetchAndPopulateRooms(selectedBuildingId);
                    }
                });

            // Show the modal after all dropdowns are populated
            $("#addDeviceModal").modal("show");
        })
        .catch((error) => {
            console.error("Error loading dropdown data:", error);
        });
}

function editDevice(deviceId) {
    // Function to handle visibility of extinguisher-specific fields
    function updateExtinguisherFields() {
        const selectElement = document.querySelector(
            "#editEmergencyDeviceTypeInput"
        );
        const selectedOption =
            selectElement.options[selectElement.selectedIndex];
        const selectedDeviceType = selectedOption.textContent;
        console.log(selectedDeviceType);

        if (selectedDeviceType !== "Fire Extinguisher") {
            // Hide and clear extinguisher-specific fields if not a Fire Extinguisher
            document
                .querySelector(".editExtinguisherTypeInputDiv")
                .classList.add("d-none");
            document.querySelector("#editExtinguisherTypeInput").value = ""; // Clear selected value
            document
                .querySelector(".editLastInspectionDateInputDiv")
                .classList.add("d-none");
            document.querySelector("#editLastInspectionDateInput").value = ""; // Clear selected value
        } else {
            // Show extinguisher-specific fields and set default
            document
                .querySelector(".editExtinguisherTypeInputDiv")
                .classList.remove("d-none");
            document.querySelector(
                "#editExtinguisherTypeInput"
            ).selectedIndex = 0; // Set to "Select Extinguisher Type"
            document
                .querySelector(".editLastInspectionDateInputDiv")
                .classList.remove("d-none");
        }
    }

    // Clear the form before showing the modal
    document.getElementById("editDeviceForm").reset();
    document.getElementById("editDeviceForm").classList.remove("was-validated");
    document.getElementById("editDeviceID").value = deviceId;

    // Fetch the dropdown data
    const emergencyDeviceTypePromise = populateDropdown(
        ".emergencyDeviceTypeInput",
        "/api/emergency-device-type",
        "Select Device Type",
        "emergency_device_type_id",
        "emergency_device_type_name"
    );

    const extinguisherTypePromise = populateDropdown(
        ".extinguisherTypeInput",
        "/api/extinguisher-type",
        "Select Extinguisher Type",
        "extinguisher_type_id",
        "extinguisher_type_name"
    );

    const sitePromise = populateDropdown(
        ".editSiteInput",
        "/api/site",
        "Select a Site",
        "site_id",
        "site_name"
    );

    // Wait for all dropdowns to be populated before proceeding
    Promise.all([
        emergencyDeviceTypePromise,
        extinguisherTypePromise,
        sitePromise,
    ])
        .then(() => {
            // Event listener for site change
            document
                .querySelector(".editSiteInput")
                .addEventListener("change", (event) => {
                    const selectedSiteId = event.target.value;
                    clearBuildingAndRoom();

                    if (selectedSiteId) {
                        fetchAndPopulateBuildings(selectedSiteId);
                    }
                });

            // Event listener for building change
            document
                .querySelector(".editBuildingInput")
                .addEventListener("change", (event) => {
                    const selectedBuildingId = event.target.value;

                    if (selectedBuildingId) {
                        fetchAndPopulateRooms(selectedBuildingId);
                    }
                });

            // Fetch the device data
            fetch(`/api/emergency-device/${deviceId}`)
                .then((response) => response.json())
                .then((data) => {
                    // Populate the form with the data
                    document.getElementById(
                        "editEmergencyDeviceTypeInput"
                    ).value = data.emergency_device_type_id;
                    document.getElementById("editExtinguisherTypeInput").value =
                        data.extinguisher_type_id.Int64;
                    document.getElementById("editSerialNumberInput").value =
                        data.serial_number.String;
                    document.getElementById("editManufactureDateInput").value =
                        data.manufacture_date.Time.split("T")[0];
                    document.getElementById("editSizeInput").value =
                        data.size.String;
                    document.getElementById("editDescriptionInput").value =
                        data.description.String;
                    document.getElementById("editSiteInput").value =
                        data.site_id;
                    document.getElementById(
                        "editLastInspectionDateInput"
                    ).value = data.last_inspection_date.Time.split("T")[0];
                    document.getElementById("editStatusInput").value =
                        data.status.String;

                    // Populate the building and room dropdowns
                    fetchAndPopulateBuildings(data.site_id)
                        .then(() => fetchAndPopulateRooms(data.building_id))
                        .then(() => {
                            // Set the building and room values
                            document.getElementById("editBuildingInput").value =
                                data.building_id;
                            document.getElementById("editRoomInput").value =
                                data.room_id;
                        });

                    // Check and update visibility of extinguisher fields
                    updateExtinguisherFields();

                    // Now that the data is populated, show the modal
                    $("#editDeviceModal").modal("show");
                });
        })
        .catch((error) => {
            console.error("Error loading dropdown data:", error);
        });

    // Add change event listener for device type dropdown
    document
        .querySelector("#editEmergencyDeviceTypeInput")
        .addEventListener("change", updateExtinguisherFields);
}

function fetchAndPopulateBuildings(siteId) {
    return fetch(`/api/building?siteId=${siteId}`)
        .then((response) => response.json())
        .then((data) => {
            const selects = document.querySelectorAll(
                ".buildingInput, .editBuildingInput"
            );
            selects.forEach((select) => {
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
            });
            return data;
        })
        .catch((error) => console.error("Error:", error));
}

function fetchAndPopulateRooms(buildingId) {
    return fetch(`/api/room?buildingId=${buildingId}`)
        .then((response) => response.json())
        .then((data) => {
            const selects = document.querySelectorAll(
                ".roomInput",
                ".editRoomInput"
            );
            selects.forEach((select) => {
                select.innerHTML =
                    "<option value='' selected disabled>Select a Room</option>";
                data.forEach((item) => {
                    const option = document.createElement("option");
                    option.text = item.room_code;
                    option.value = item.room_id;
                    select.add(option);
                });
            });
            return data;
        })
        .catch((error) => console.error("Error:", error));
}

/// Fetch the form and the submit button
const addDeviceForm = document.querySelector("#addDeviceForm");
const addDeviceButton = document.querySelector("#addDeviceBtn");
/// Fetch the form and the submit button
const editDeviceForm = document.querySelector("#editDeviceForm");
const editDeviceButton = document.querySelector("#editDeviceBtn");

// Function to validate select elements
function validateSelect(selectElement) {
    if (selectElement.value === "") {
        selectElement.setCustomValidity("Please make a selection");
    } else {
        selectElement.setCustomValidity("");
    }
}

// Function to validate dates
function validateDates() {
    const currentDate = new Date().toISOString().split("T")[0];
    let isValid = true;

    const manufactureDateInputs = document.querySelectorAll(
        ".manufactureDateInput"
    );
    const lastInspectionDateInputs = document.querySelectorAll(
        ".lastInspectionDateInput"
    );
    const deviceTypeInputs = document.querySelectorAll(
        ".emergencyDeviceTypeInput"
    );

    manufactureDateInputs.forEach((manufactureDate, index) => {
        const lastInspectionDate = lastInspectionDateInputs[index];
        const deviceTypeSelect = deviceTypeInputs[index];
        const isFireExtinguisher =
            deviceTypeSelect.options[deviceTypeSelect.selectedIndex]
                ?.textContent === "Fire Extinguisher";

        // Validate manufacture date
        if (manufactureDate.value && manufactureDate.value > currentDate) {
            manufactureDate.setCustomValidity(
                "Manufacture date cannot be in the future"
            );
            document.querySelectorAll(".manufactureDateFeedback")[
                index
            ].textContent = "Manufacture date cannot be in the future";
            isValid = false;
        } else {
            manufactureDate.setCustomValidity("");
            document.querySelectorAll(".manufactureDateFeedback")[
                index
            ].textContent = "";
        }

        // Only validate last inspection date if it's a Fire Extinguisher
        if (isFireExtinguisher && lastInspectionDate) {
            if (
                lastInspectionDate.value &&
                lastInspectionDate.value > currentDate
            ) {
                lastInspectionDate.setCustomValidity(
                    "Last inspection date cannot be in the future"
                );
                document.querySelectorAll(".lastInspectionDateFeedback")[
                    index
                ].textContent = "Last inspection date cannot be in the future";
                isValid = false;
            } else if (
                lastInspectionDate.value &&
                lastInspectionDate.value === manufactureDate.value
            ) {
                lastInspectionDate.setCustomValidity(
                    "Last inspection date cannot be the same as manufacture date"
                );
                document.querySelectorAll(".lastInspectionDateFeedback")[
                    index
                ].textContent =
                    "Last inspection date cannot be the same as manufacture date";
                isValid = false;
            } else if (
                lastInspectionDate.value &&
                manufactureDate.value &&
                lastInspectionDate.value < manufactureDate.value
            ) {
                lastInspectionDate.setCustomValidity(
                    "Last inspection date cannot be before manufacture date"
                );
                document.querySelectorAll(".lastInspectionDateFeedback")[
                    index
                ].textContent =
                    "Last inspection date cannot be before manufacture date";
                isValid = false;
            } else {
                lastInspectionDate.setCustomValidity("");
                document.querySelectorAll(".lastInspectionDateFeedback")[
                    index
                ].textContent = "";
            }
        } else if (lastInspectionDate) {
            // Clear validation for non-Fire Extinguisher devices
            lastInspectionDate.setCustomValidity("");
            document.querySelectorAll(".lastInspectionDateFeedback")[
                index
            ].textContent = "";
        }
    });
    return isValid;
}

// Function to validate length for input and textarea elements
function validateLength(element, maxLength) {
    if (element.value.length > maxLength) {
        element.setCustomValidity(
            `This field is too long, maximum ${maxLength} characters.`
        );
    } else {
        element.setCustomValidity("");
    }
}

// Function to handle device type changes
function handleDeviceTypeChange(event) {
    const selectedOption = event.target.options[event.target.selectedIndex];
    const selectedDeviceType = selectedOption.textContent;
    const isEdit = event.target.id === "editEmergencyDeviceTypeInput";
    const prefix = isEdit ? "edit" : "";

    if (selectedDeviceType !== "Fire Extinguisher") {
        // Clear and hide extinguisher type
        const extinguisherTypeInput = document.querySelector(
            `#${prefix}ExtinguisherTypeInput`
        );
        if (extinguisherTypeInput) {
            extinguisherTypeInput.value = "";
            document
                .querySelector(`.${prefix}ExtinguisherTypeInputDiv`)
                .classList.add("d-none");
        }

        // Clear and hide last inspection date
        const lastInspectionDateInput = document.querySelector(
            `#${prefix}LastInspectionDateInput`
        );
        if (lastInspectionDateInput) {
            lastInspectionDateInput.value = "";
            document
                .querySelector(`.${prefix}LastInspectionDateInputDiv`)
                .classList.add("d-none");
        }
    } else {
        // Show fields for Fire Extinguisher
        document
            .querySelector(`.${prefix}ExtinguisherTypeInputDiv`)
            .classList.remove("d-none");
        document
            .querySelector(`.${prefix}LastInspectionDateInputDiv`)
            .classList.remove("d-none");
    }

    // Validate dates after changing device type
    validateDates();
}

document.addEventListener("DOMContentLoaded", function () {
    // Add change event listeners to device type inputs
    const deviceTypeInputs = document.querySelectorAll(
        ".emergencyDeviceTypeInput"
    );
    deviceTypeInputs.forEach((input) => {
        input.addEventListener("change", handleDeviceTypeChange);
    });

    const description = document.querySelector(".descriptionInput");
    const manufactureDate = document.querySelector(".manufactureDateInput");
    const lastInspectionDate = document.querySelector(
        ".lastInspectionDateInput"
    );
    const editDescriptionInput = document.querySelector(
        "#editDescriptionInput"
    );
    const editManufactureDateInput = document.querySelector(
        "#editManufactureDateInput"
    );
    const editLastInspectionDateInput = document.querySelector(
        "#editLastInspectionDateInput"
    );

    // Validate edit description length
    editDescriptionInput.addEventListener("input", function () {
        validateLength(this, 255);
    });

    // Validate edit manufacture date
    editManufactureDateInput.addEventListener("change", validateDates);

    // Validate edit last inspection date
    editLastInspectionDateInput.addEventListener("change", validateDates);

    // Validate description length
    description.addEventListener("input", function () {
        validateLength(this, 255);
    });

    // Add event listeners to select elements
    document
        .querySelectorAll(
            ".emergencyDeviceTypeInput, .siteInput, .buildingInput, .roomInput"
        )
        .forEach((select) => {
            select.addEventListener("change", function () {
                validateSelect(this);
            });
        });

    // Add event listeners for date validation
    manufactureDate.addEventListener("change", validateDates);
    lastInspectionDate.addEventListener("change", validateDates);

    // Add event listener to the add device button
    addDeviceButton.addEventListener("click", function (event) {
        // Validate all select elements before form submission
        document
            .querySelectorAll(
                ".emergencyDeviceTypeInput, .siteInput, .buildingInput, .roomInput"
            )
            .forEach((select) => {
                validateSelect(select);
            });

        // Validate description length
        validateLength(description, 255);

        // Validate dates
        const datesValid = validateDates();

        if (!addDeviceForm.checkValidity() || !datesValid) {
            event.preventDefault();
            event.stopPropagation();
        } else {
            // If the form is valid, submit it
            addDeviceForm.submit();
        }

        addDeviceForm.classList.add("was-validated");
    });

    // Add event listener to the edit device button
    editDeviceButton.addEventListener("click", function (event) {
        event.preventDefault(); // Prevent default form submission

        // Validate all select elements before form submission
        document
            .querySelectorAll(
                ".emergencyDeviceTypeInput, .editSiteInput, .editBuildingInput, .roomInput"
            )
            .forEach((select) => {
                validateSelect(select);
            });

        // Validate description length
        validateLength(description, 255);

        // Validate dates
        const datesValid = validateDates();

        // Check if the form is valid
        if (!editDeviceForm.checkValidity() || !datesValid) {
            event.stopPropagation();
            editDeviceForm.classList.add("was-validated");
        } else {
            // If the form is valid, prepare to send the PUT request
            const formData = new FormData(editDeviceForm);
            const jsonData = {};
            for (const [key, value] of formData.entries()) {
                jsonData[key] = value;
            }
            console.log(jsonData);
            // Send the PUT request
            fetch(
                `/api/emergency-device/${
                    document.getElementById("editDeviceID").value
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

function deleteDevice(deviceId) {
    console.log(`Delete device with ID: ${deviceId}`);
    // Add your delete logic here
}

function viewDeviceInspections(deviceId) {
    console.log(`Inspect device with ID: ${deviceId}`);

    // Clear the inspection table
    document.getElementById("inspectionTable").innerHTML = "";

    // Clear the hidden input field
    document.getElementById("inspect_device_id").value = "";

    // Fetch the inspections for this device
    fetch(`/api/inspection?device_id=${deviceId}`)
        .then((response) => response.json())
        .then((data) => {
            const inspectionTable = document.getElementById("inspectionTable");

            if (!data || !Array.isArray(data) || data.length === 0) {
                inspectionTable.innerHTML = `
                    <tr>
                        <td colspan="4" class="text-center">No inspections found</td>
                    </tr>
                `;
            } else {
                inspectionTable.innerHTML = data
                    .map((inspection) => {
                        const formattedDate = inspection.inspection_date.Valid
                            ? new Date(
                                  inspection.inspection_date.Time
                              ).toLocaleDateString("en-NZ", {
                                  day: "numeric",
                                  month: "long",
                                  year: "numeric",
                              })
                            : "No Date Available";

                        // Determine badge color based on inspection status
                        let badgeClass = "badge text-bg-primary"; // default color
                        if (inspection.inspection_status === "Passed") {
                            badgeClass = "badge text-bg-success";
                        } else if (inspection.inspection_status === "Failed") {
                            badgeClass = "badge text-bg-danger";
                        }

                        return `
                            <tr>
                                <td>${formattedDate}</td>
                                <td>${
                                    inspection.inspector_name || "Unknown"
                                }</td>
                                <td>
                                    <span class="badge ${badgeClass}">${
                            inspection.inspection_status || "Not Set"
                        }</span>
                                </td>
                                <td>
                                    <button class="btn btn-primary" onclick="viewInspectionDetails(${
                                        inspection.emergency_device_inspection_id
                                    })">View</button>
                                </td>
                            </tr>
                        `;
                    })
                    .join("");
            }
            // log the data to the console
            console.log(data);

            if (data) {
                // Set the modal title with the device serial number
                document.getElementById("inspectionModalTitle").innerText =
                    `Extinguisher Inspections - Serial Number: ${data[0].serial_number}` ||
                    "Unknown";
            } else {
                document.getElementById("inspectionModalTitle").innerText =
                    "Extinguisher Inspections";
            }
        })
        .catch((error) => {
            console.error("Error fetching inspection data:", error);
            document.getElementById("inspectionTable").innerHTML = `
                <tr>
                    <td colspan="4" class="text-center">Failed to load inspections</td>
                </tr>
            `;
        });

    // Set the device ID in the hidden input field
    document.getElementById("inspect_device_id").value = deviceId;

    // Show the modal
    $("#viewInspectionModal").modal("show");
}

function viewInspectionDetails(inspectionId) {
    console.log(`View inspection details for inspection ID: ${inspectionId}`);
    // Close the view inspection modal
    $("#viewInspectionModal").modal("hide");

    // Fetch the inspection details
    fetch(`/api/inspection/${inspectionId}`)
        .then((response) => response.json())
        .then((data) => {
            // Populate the modal with the inspection details
            document.getElementById("inspector_username").innerText =
                data.inspector_name || "Unknown";
            // Options for formatting date and date-time
            const dateOptions = {
                day: "numeric",
                month: "long",
                year: "numeric",
            };

            const dateTimeOptions = {
                ...dateOptions,
                hour: "numeric",
                minute: "numeric",
                hour12: true,
            };

            // Update inspection date display
            document.getElementById("ViewInspectionDateInput").innerText = data
                .inspection_date.Valid
                ? formatDate(data.inspection_date.Time, dateOptions)
                : "No Date Available";

            // Update created at date display
            document.getElementById("ViewInspectionCreatedAt").innerText = data
                .created_at.Valid
                ? formatDate(data.created_at.Time, dateTimeOptions)
                : "No Date Available";

            // Create badge for the inspection status
            const statusBadge = document.createElement("span");
            statusBadge.className = "badge";

            // Set badge color based on status
            switch (data.inspection_status) {
                case "Passed":
                    statusBadge.classList.add("bg-success");
                    statusBadge.innerText = "Passed";
                    break;
                case "Failed":
                    statusBadge.classList.add("bg-danger");
                    statusBadge.innerText = "Failed";
                    break;
                default:
                    statusBadge.classList.add("bg-secondary");
                    statusBadge.innerText = "Not Set";
            }

            // Clear previous status display and append the new badge
            const statusContainer = document.getElementById(
                "ViewInspectionStatus"
            );
            statusContainer.innerHTML = ""; // Clear existing content
            statusContainer.appendChild(statusBadge); // Append the badge

            document.getElementById("viewNotes").innerText =
                data.notes.String || "";
            document.getElementById("ViewdeviceSerialNumber").innerText =
                data.serial_number || "Unknown";
            // Check checkboxes based on boolean values
            document.getElementById("ViewIsConspicuous").checked =
                data.is_conspicuous.Bool && data.is_conspicuous.Valid;
            document.getElementById("ViewIsAccessible").checked =
                data.is_accessible.Bool && data.is_accessible.Valid;
            document.getElementById("ViewIsAssignedLocation").checked =
                data.is_assigned_location.Bool &&
                data.is_assigned_location.Valid;
            document.getElementById("ViewIsSignVisible").checked =
                data.is_sign_visible.Bool && data.is_sign_visible.Valid;
            document.getElementById("ViewIsAntiTamperDeviceIntact").checked =
                data.is_anti_tamper_device_intact.Bool &&
                data.is_anti_tamper_device_intact.Valid;
            document.getElementById("ViewIsSupportBracketSecure").checked =
                data.is_support_bracket_secure.Bool &&
                data.is_support_bracket_secure.Valid;
            document.getElementById("ViewWorkOrderRequired").checked =
                data.work_order_required.Bool && data.work_order_required.Valid;
            document.getElementById(
                "ViewAreOperatingInstructionsClear"
            ).checked =
                data.are_operating_instructions_clear.Bool &&
                data.are_operating_instructions_clear.Valid;
            document.getElementById("ViewIsMaintenanceTagAttached").checked =
                data.is_maintenance_tag_attached.Bool &&
                data.is_maintenance_tag_attached.Valid;
            document.getElementById("ViewIsNoExternalDamage").checked =
                data.is_no_external_damage.Bool &&
                data.is_no_external_damage.Valid;
            document.getElementById("ViewIsChargeGaugeNormal").checked =
                data.is_charge_gauge_normal.Bool &&
                data.is_charge_gauge_normal.Valid;
            document.getElementById("ViewIsReplaced").checked =
                data.is_replaced.Bool && data.is_replaced.Valid;
            document.getElementById(
                "ViewAreMaintenanceRecordsComplete"
            ).checked =
                data.are_maintenance_records_complete.Bool &&
                data.are_maintenance_records_complete.Valid;

            // Show the modal
            $("#viewInspectionDetailsModal").modal("show");
        })
        .catch((error) => {
            console.error("Error fetching inspection details:", error);
        });
}

console.log("User ID:", user_id);

function addInspection() {
    const deviceId = document.getElementById("inspect_device_id").value;

    // Set the user ID in the hidden input field
    document.getElementById("inspect_user_id").value = user_id;

    console.log(`Adding inspection for device ID: ${deviceId}`);

    // Close the view inspection modal
    $("#viewInspectionModal").modal("hide");

    // Clear the form and reset validation classes
    const addInspectionForm = document.getElementById("addInspectionForm");
    addInspectionForm.reset();
    addInspectionForm.classList.remove("was-validated");

    // Clear the feedback messages
    const feedbackElements =
        addInspectionForm.querySelectorAll(".invalid-feedback");

    feedbackElements.forEach((element) => {
        element.textContent = "";
    });

    // Add the device ID to the hidden input field
    const deviceIdInput = document.getElementById("add_inspection_device_id");
    deviceIdInput.value = deviceId;

    console.log(`Device ID set in form: ${deviceIdInput.value}`);

    // Show the add inspection modal
    $("#addInspectionModal").modal("show");
}

document.addEventListener("DOMContentLoaded", function () {
    // Select form elements
    const addInspectionButton = document.querySelector("#addInspectionBtn");
    const addInspectionForm = document.querySelector("#addInspectionForm");
    const inspectionDateInput = document.querySelector("#InspectionDateInput");
    const inspectionStatus = document.querySelector("#inspectionStatus");
    const inspectionDateFeedback = document.getElementById(
        "inspectionDateFeedback"
    );
    const inspectionStatusFeedback = inspectionStatus.nextElementSibling;

    // Ensure the elements are found
    if (
        !addInspectionButton ||
        !addInspectionForm ||
        !inspectionDateInput ||
        !inspectionStatus
    ) {
        console.error("Required elements not found in the DOM.");
        return;
    }

    // Get all checkboxes except workOrderRequired and isReplaced
    const checkboxes = Array.from(
        addInspectionForm.querySelectorAll('input[type="checkbox"]')
    ).filter(
        (checkbox) => !["workOrderRequired", "isReplaced"].includes(checkbox.id)
    );

    // Function to validate checkboxes based on inspection status
    function validateInspectionStatus() {
        const allChecked = checkboxes.every((checkbox) => checkbox.checked);

        if (inspectionStatus.value === "Passed" && !allChecked) {
            inspectionStatus.setCustomValidity(
                "All inspection criteria must be met to mark as Passed"
            );
            inspectionStatusFeedback.textContent =
                "All inspection criteria must be met to mark as Passed";
            return false;
        } else if (!inspectionStatus.value) {
            inspectionStatus.setCustomValidity(
                "Please select an inspection status"
            );
            inspectionStatusFeedback.textContent =
                "Please select an inspection status";
            return false;
        } else {
            inspectionStatus.setCustomValidity("");
            inspectionStatusFeedback.textContent = "";
            return true;
        }
    }

    // Add event listener to inspection status dropdown
    inspectionStatus.addEventListener("change", function () {
        validateInspectionStatus();
        // Show validation message immediately on change
        if (addInspectionForm.classList.contains("was-validated")) {
            inspectionStatusFeedback.style.display = this.validationMessage
                ? "block"
                : "none";
            inspectionStatusFeedback.textContent = this.validationMessage || "";
        }
    });

    // Add change event listeners to all checkboxes
    checkboxes.forEach((checkbox) => {
        checkbox.addEventListener("change", () => {
            if (inspectionStatus.value === "Passed") {
                validateInspectionStatus();
                // Update status feedback if form is already validated
                if (addInspectionForm.classList.contains("was-validated")) {
                    inspectionStatusFeedback.style.display =
                        inspectionStatus.validationMessage ? "block" : "none";
                    inspectionStatusFeedback.textContent =
                        inspectionStatus.validationMessage;
                }
            }
        });
    });

    // Event listener for inspection date input change
    inspectionDateInput.addEventListener("input", function () {
        const currentDate = new Date().toISOString().split("T")[0];

        if (
            inspectionDateInput.value &&
            inspectionDateInput.value <= currentDate
        ) {
            inspectionDateInput.setCustomValidity(""); // Clear the custom validity
            inspectionDateFeedback.textContent = ""; // Clear feedback message
        } else if (inspectionDateInput.value > currentDate) {
            inspectionDateInput.setCustomValidity(
                "Inspection date cannot be in the future"
            );
            inspectionDateFeedback.textContent =
                "Inspection date cannot be in the future";
        }
    });

    // Add event listener to the add inspection button
    addInspectionButton.addEventListener("click", function (event) {
        event.preventDefault(); // Prevent default form submission

        // Check if inspection date is provided
        if (inspectionDateInput.value) {
            const currentDate = new Date().toISOString().split("T")[0];
            // Check if the inspection date is in the future
            if (inspectionDateInput.value > currentDate) {
                inspectionDateInput.setCustomValidity(
                    "Inspection date cannot be in the future"
                );
                inspectionDateFeedback.textContent =
                    "Inspection date cannot be in the future";
            } else {
                inspectionDateInput.setCustomValidity(""); // Clear any previous validity message
                inspectionDateFeedback.textContent = ""; // Clear feedback message
            }
        } else {
            inspectionDateInput.setCustomValidity(
                "Please provide an inspection date"
            );
            inspectionDateFeedback.textContent =
                "Please provide an inspection date";
        }

        // Validate inspection status
        validateInspectionStatus();

        // Add was-validated class before checking validity
        addInspectionForm.classList.add("was-validated");

        // Validate the form before submission
        if (!addInspectionForm.checkValidity()) {
            event.stopPropagation();

            // Show custom error messages for invalid fields
            if (inspectionStatus.validationMessage) {
                inspectionStatusFeedback.style.display = "block";
                inspectionStatusFeedback.textContent =
                    inspectionStatus.validationMessage;
            }

            if (inspectionDateInput.validationMessage) {
                inspectionDateFeedback.style.display = "block";
                inspectionDateFeedback.textContent =
                    inspectionDateInput.validationMessage;
            }
        } else {
            // If the form is valid, submit it
            addInspectionForm.submit();
        }
    });
});

function deviceNotes(description) {
    // Populate the modal with the description
    document.getElementById("notesModalBody").innerText = description;

    // Show the modal
    $("#notesModal").modal("show");
}

// Function to toggle the map visibility
function toggleMap() {
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

async function searchDevices() {
    const siteFilter = document.getElementById("siteFilter");
    const searchInput = document.getElementById("searchInput");
    const searchValue = searchInput.value.toLowerCase();

    console.log("siteFilter:", siteFilter.value);

    // If site filter is "All Sites", reload all devices
    if (
        document.getElementById("siteFilter").selectedOptions[0].text ===
        "All Sites"
    ) {
        await getAllDevices();
    } else {
        // Reload devices for the selected site
        await getAllDevices("", siteFilter.value);
    }

    allDevices = allDevices.filter((device) => {
        const baseSearch =
            device.emergency_device_type_name
                .toLowerCase()
                .includes(searchValue) ||
            device.extinguisher_type_name.String.toLowerCase().includes(
                searchValue
            ) ||
            device.room_code.toLowerCase().includes(searchValue) ||
            device.serial_number.String.toLowerCase().includes(searchValue) ||
            device.manufacture_date.Time.toLowerCase().includes(searchValue) ||
            device.expire_date.Time.toLowerCase().includes(searchValue) ||
            device.size.String.toLowerCase().includes(searchValue) ||
            device.status.String.toLowerCase().includes(searchValue) ||
            device.description.String.toLowerCase().includes(searchValue);

        // Add admin-only search fields if user is admin
        if (role === "Admin") {
            // Format the dates for searching
            const lastInspectionFormatted = new Date(
                device.last_inspection_date.Time
            )
                .toLocaleDateString("en-NZ", {
                    day: "numeric",
                    month: "long",
                    year: "numeric",
                })
                .toLowerCase();

            const nextInspectionFormatted = new Date(
                device.next_inspection_date.Time
            )
                .toLocaleDateString("en-NZ", {
                    day: "numeric",
                    month: "long",
                    year: "numeric",
                })
                .toLowerCase();

            return (
                baseSearch ||
                lastInspectionFormatted.includes(searchValue) ||
                nextInspectionFormatted.includes(searchValue)
            );
        }

        return baseSearch;
    });

    updateTable();
}

document.getElementById("searchInput").addEventListener("input", () => {
    searchDevices();
});
