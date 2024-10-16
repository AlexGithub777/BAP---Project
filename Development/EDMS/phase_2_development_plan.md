# EDMS Phase 2 Development Plan

## All Tasks Due - TBA

### Testing & Intergation - TBA

### Prototype 2.0 Release Date & Presentation - TBA

## Development Workflow:

### 1. Create a New Branch

-   Start by creating a new branch from the remote develop branch.
-   Name the branch using the format `feature/task_name` (e.g., `feature/device_management`).

### 2. Develop the Feature

-   Implement the assigned feature(s) in your branch.
-   Regularly commit your changes to the feature branch with clear commit messages.

### 3. Sync with Remote Develop Branch

-   Before merging your feature branch, ensure it’s up-to-date by fetching the latest changes from the remote develop branch.
-   Resolve any merge conflicts carefully to ensure nothing breaks.

### 4. Submit a Pull Request (PR)

-   Once the feature is complete and tested, push all changes to your feature branch.
-   Go to the GitHub repository and submit a [pull request](https://github.com/AlexGithub777/BAP---Project/pulls) from your feature branch to the develop branch. (set base as develop & compare as your branch)
-   Provide a detailed description of the changes in your PR.

### 5. Review and Merge

-   After your pull request is submitted, I will review the changes.
-   Once approved, the changes will be merged into the develop branch.

| **Team Member** | **Task**                         | **Subtasks**                                                                                                                                                                                                                                                                                                                                         |
| --------------- | -------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Liam**        | Filter devices                   | - Filter device list using dropdown filters<br>- Update the table when a filter option is selected<br>- Update filter UI on on mobile                                                                                                                                                                                                                |
| **Alex**        | Search devices<br>Update Main UI | - Search devices by all properties, update device list with results on keypress:<br>- Device Type <br>- Extinguisher Type<br>- Room<br>- Serial Number<br>- Manufacture Date<br>- Expire Date<br>- Size<br>- Status<br>If current user role is Admin, then also allow search by thse properties:<br>- Last Inspection Date<br>- Next Inspection Date |
| **Aidan**       | Notifications                    | - Generate all Notifications on login (expired device or device inspection due)<br>-Generate and display HTML in notifcations modal <br>- Dsiplay Notification count in NavBar<br>- Update Notification Count in NavBar<br>- Clear Notification<br>- Clear All Notifications                                                                         |
| **James**       | Sort Devices                     | - Sort the device list by clicking any column heading in the table, update the table after each click, toggle between ascending and descending order.<br>-Update the icon arrow in the table heading to reflect the state of ordering. (▲ for asc, ▲ for desc)                                                                                       |
| **Joe**         | Inspection Management            | - Create Inspection (UI modal already complete, develop backend and link to frontend.)<br>- Read all inspections (single device) (UI done, develop backend and link to frontend)<br>- View inspection details (Single device)<br> Display inspection details (use non-editable version of "add inspection" modal)                                    |
