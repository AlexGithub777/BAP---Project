// Define priorityOrder globally
const priorityOrder = {
    "Inspection Failed": 0,
    Expired: 1,
    "Inspection Due": 2,
    "Expiring Soon": 3,
    "Inspection Due Soon": 4,
};

export async function getAllDevices(buildingCode = "", siteId = "") {
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

        const devices = await response.json();
        return devices; // Return the devices instead of storing in global variable
    } catch (err) {
        console.error("Failed to fetch devices:", err);
        return []; // Return empty array in case of error
    }
}

export async function updateDeviceStatus(deviceId, status) {
    try {
        const response = await fetch(
            `/api/emergency-device/${deviceId}/status`,
            {
                method: "PUT",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ status: status }),
            }
        );

        const data = await response.json();

        if (data.error) {
            console.error("Error:", data.error);
            window.location.href = data.redirectURL;
            throw new Error(data.error);
        } else if (data.message) {
            console.log("Success:", data.message);
            // Only redirect if redirectURL is provided
            if (data.redirectURL) {
                window.location.href = data.redirectURL;
            }
            return true; // Indicate success
        } else {
            console.error("Unexpected response:", data);
            throw new Error("Unexpected response");
        }
    } catch (error) {
        console.error(`Failed to update status for device ${deviceId}:`, error);
        return false; // Indicate failure
    }
}

export async function generateNotifications() {
    const allDevices = await getAllDevices(); // Store the returned devices

    const currentDate = new Date();
    const thirtyDaysFromNow = new Date();
    thirtyDaysFromNow.setDate(currentDate.getDate() + 30);

    // Helper function to calculate days difference
    const calculateDaysOverdue = (date) => {
        const diffTime = currentDate - new Date(date);
        return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    };

    // Helper function to check if a date is today or in the past
    const isDateDueOrPast = (date) => {
        const targetDate = new Date(date);
        // Reset time parts to compare just the dates
        targetDate.setHours(0, 0, 0, 0);
        const today = new Date();
        today.setHours(0, 0, 0, 0);
        return targetDate <= today;
    };

    // Update device statuses based on dates
    for (const device of allDevices) {
        let statusUpdated = false;

        // Check inspection date
        if (
            device.next_inspection_date.Valid &&
            isDateDueOrPast(device.next_inspection_date.Time) &&
            device.status.String !== "Inspection Failed" &&
            device.status.String !== "Inspection Due"
        ) {
            const success = await updateDeviceStatus(
                device.emergency_device_id,
                "Inspection Due"
            );
            if (success) {
                device.status.String = "Inspection Due"; // Update local state
                statusUpdated = true;
                // reload devices
                await getAllDevices();
            }
        }

        // Check expire date
        if (
            device.expire_date.Valid &&
            isDateDueOrPast(device.expire_date.Time) &&
            device.status.String !== "Expired" &&
            !statusUpdated
        ) {
            // Only update if no other status change

            const success = await updateDeviceStatus(
                device.emergency_device_id,
                "Expired"
            );
            if (success) {
                device.status.String = "Expired"; // Update local state
                // reload devices
                await getAllDevices();
            }
        }
    }

    // Create a Map to store unique devices by their ID
    const notificationMap = new Map();

    // Helper function to add device to map with reason and days
    const addDeviceWithReason = (device, reason, days = null) => {
        if (!notificationMap.has(device.emergency_device_id)) {
            notificationMap.set(device.emergency_device_id, {
                ...device,
                notification_details: [
                    {
                        reason: reason,
                        days: days,
                    },
                ],
            });
        } else {
            const existing = notificationMap.get(device.emergency_device_id);
            if (
                !existing.notification_details.some(
                    (detail) => detail.reason === reason
                )
            ) {
                existing.notification_details.push({
                    reason: reason,
                    days: days,
                });
            }
        }
    };

    // Check each device
    allDevices.forEach((device) => {
        // Skip inactive devices
        if (device.status.String === "Inactive") {
            return;
        }

        // Check failed inspection status
        if (device.status.String === "Inspection Failed") {
            addDeviceWithReason(device, "Inspection Failed");
        }

        // Check inspection due status
        if (
            device.status.String === "Inspection Due" &&
            device.next_inspection_date.Valid
        ) {
            const daysOverdue = calculateDaysOverdue(
                device.next_inspection_date.Time
            );
            addDeviceWithReason(device, "Inspection Due", daysOverdue);
        }

        // Check expired status
        if (device.status.String === "Expired" && device.expire_date.Valid) {
            const daysOverdue = calculateDaysOverdue(device.expire_date.Time);
            addDeviceWithReason(device, "Expired", daysOverdue);
        }

        // Check expire date within 30 days
        if (device.expire_date.Valid) {
            const expireDate = new Date(device.expire_date.Time);
            if (expireDate > currentDate && expireDate <= thirtyDaysFromNow) {
                const daysUntil = Math.ceil(
                    (expireDate - currentDate) / (1000 * 60 * 60 * 24)
                );
                addDeviceWithReason(device, "Expiring Soon", daysUntil);
            }
        }

        // Check next inspection date within 30 days
        if (device.next_inspection_date.Valid) {
            const inspectionDate = new Date(device.next_inspection_date.Time);
            if (
                inspectionDate > currentDate &&
                inspectionDate <= thirtyDaysFromNow
            ) {
                const daysUntil = Math.ceil(
                    (inspectionDate - currentDate) / (1000 * 60 * 60 * 24)
                );
                addDeviceWithReason(device, "Inspection Due Soon", daysUntil);
            }
        }
    });

    // Convert map to array
    const notifications = Array.from(notificationMap.values());

    // Sort notifications by priority
    const priorityOrder = {
        "Inspection Failed": 0,
        Expired: 1,
        "Inspection Due": 2,
        "Expiring Soon": 3,
        "Inspection Due Soon": 4,
    };

    notifications.sort((a, b) => {
        const aPriority = Math.min(
            ...a.notification_details.map((d) => priorityOrder[d.reason])
        );
        const bPriority = Math.min(
            ...b.notification_details.map((d) => priorityOrder[d.reason])
        );

        if (aPriority !== bPriority) {
            return aPriority - bPriority;
        }

        // If same priority, sort by days overdue (highest first)
        const aDetails = a.notification_details.find(
            (d) => priorityOrder[d.reason] === aPriority
        );
        const bDetails = b.notification_details.find(
            (d) => priorityOrder[d.reason] === bPriority
        );
        return (bDetails?.days || 0) - (aDetails?.days || 0);
    });

    return notifications;
}

export function generateNotificationHTML(notifications) {
    const getStatusBadge = (detail) => {
        const { reason, days } = detail;
        let badgeClass = "";
        let icon = "";
        let text = "";

        switch (reason) {
            case "Inspection Failed":
                badgeClass = "bg-danger text-light";
                icon = '<i class="text-danger fa fa-exclamation-circle"></i>';
                text = "Inspection Failed";
                break;
            case "Expired":
                badgeClass = "bg-danger text-light";
                icon = '<i class="text-danger fa fa-exclamation-circle"></i>';
                text = `Expired (${days} days ago)`;
                break;
            case "Inspection Due":
                badgeClass = "bg-danger text-light";
                icon = '<i class="text-danger fa fa-exclamation-circle"></i>';
                text = `Inspection Due (${days} days ago)`;
                break;
            case "Expiring Soon":
                badgeClass = "bg-warning text-black";
                icon =
                    '<i class="text-warning fa-solid fa-exclamation-triangle"></i>';
                text = `Expires (In ${days} days)`;
                break;
            case "Inspection Due Soon":
                badgeClass = "bg-warning text-black";
                icon =
                    '<i class="text-warning fa-solid fa-exclamation-triangle"></i>';
                text = `Inspection Due (In ${days} days)`;
                break;
        }

        return { badgeClass, icon, text };
    };

    let html = "";

    notifications.forEach((device) => {
        // Get highest priority notification detail
        const mainDetail = device.notification_details.reduce((a, b) =>
            priorityOrder[a.reason] < priorityOrder[b.reason] ? a : b
        );

        const { badgeClass, icon, text } = getStatusBadge(mainDetail);

        html += `
            <div class="card mb-3">
                <div class="card-body">
                    <h5 class="card-title">
                        ${device.emergency_device_type_name}
                        ${icon}
                    </h5>
                    <div class="card-text">
                        <div class="d-flex justify-content-between">
                            <div>
                                <span>Serial Number: ${
                                    device.serial_number.String
                                }</span><br />
                                <span>Room: ${device.room_code}</span><br />
                                <span>Status: 
                                    <span class="badge ${badgeClass}">
                                        ${text}
                                    </span>
                                </span>
                            </div>
                            <div>
                                ${
                                    mainDetail.reason.includes("Inspection")
                                        ? `<button class="btn btn-primary" onclick="viewDeviceInspections(${device.emergency_device_id})">
                                        Inspect
                                    </button>`
                                        : ""
                                }
                                <button class="btn btn-secondary" onclick="">
                                    Clear
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        `;
    });

    return html;
}

export async function updateNotificationsUI() {
    try {
        const notifications = await generateNotifications();
        const html = generateNotificationHTML(notifications);

        // Update the notifications section
        const notificationsElement = document.getElementById(
            "deviceNotificationsCards"
        );
        if (notificationsElement) {
            notificationsElement.innerHTML = html;
        } else {
            console.error("Notifications element not found");
        }

        // Update the notification count
        const notificationCountElement = document.querySelector(
            ".notification-count"
        );
        if (notificationCountElement) {
            notificationCountElement.textContent = notifications.length;
        } else {
            console.error("Notification count element not found");
        }
    } catch (error) {
        console.error("Failed to generate notifications:", error);
        // Optionally show user-friendly error message
        const notificationsElement = document.getElementById(
            "deviceNotificationsCards"
        );
        if (notificationsElement) {
            notificationsElement.innerHTML = `
                <div class="alert alert-danger" role="alert">
                    Failed to load notifications. Please try refreshing the page.
                </div>
            `;
        }

        // Reset notification count to 0 if there's an error
        const notificationCountElements = document.querySelectorAll(
            ".notification-count"
        );
        if (notificationCountElements) {
            notificationCountElements.forEach((element) => {
                element.textContent = "0";
            });
        }
    }
}
