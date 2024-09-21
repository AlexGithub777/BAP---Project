# EDMS Phase 1 Development Plan

## All Tasks Due - 27/09/2024

### Testing & Intergation - 27/09/2024 - 29/09/2024

### Prototype 1.0 Release Date & Presentation - 29/09/2024

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

| **Team Member** | **Task**                 | **Subtasks**                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |
| --------------- | ------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Liam**        | Device Management        | - Create Device: Develop Modal UI for creating devices.<br>- Connect form to backend to store device.<br>- Display success/error message.<br><br>- Update Device: Develop modal UI for updating device details.<br>- Prefill form with existing backend data.<br>- Display success/error message.<br><br>- Delete Device: Add "Confirm delete" modal.<br>- Link delete button to backend for device deletion.<br>- Display success/error message.                                                                                                                                                                                      |
| **Joe**         | Building Management      | - Create Building: Develop modal UI for creating buildings.<br>- Ensure building is assigned to a site.<br>- Populate site dropdown from backend.<br>- Connect UI form to backend.<br>- Display success/error message.<br><br>- Update Building: Develop modal UI for updating buildings.<br>- Fetch existing data from backend for form prefill.<br>- Display success/error message.<br><br>- Delete Building: Add "Confirm delete" modal.<br>- Handle foreign key constraint errors.<br>- Display success/error message.                                                                                                             |
| **Aidan**       | Room Management          | - Create Room: Develop modal UI for creating rooms.<br>- Ensure room is assigned to a site and building.<br>- Populate site/building dropdowns dynamically from backend.<br>- Link form to backend for room creation.<br>- Display success/error message.<br><br>- Update Room: Develop modal UI for updating rooms.<br>- Prefill form with backend data.<br>- Display success/error message.<br><br>- Delete Room: Add "Confirm delete" modal.<br>- Handle foreign key constraint errors.<br>- Display success/error message.                                                                                                         |
| **James**       | Device Type Management   | - Create Device Type: Ensure unique device type name.<br>- Handle duplication errors.<br>- Use modal UI for creating device types.<br>- Connect UI form to backend.<br>- Display success/error message.<br><br>- Update Device Type: Implement update functionality using modal UI.<br>- Prefill form with backend data.<br>- Display success/error message.<br><br>- Delete Device Type: Add "Confirm delete" modal.<br>- Handle foreign key constraint errors.<br>- Display success/error message.                                                                                                                                   |
| **Alex**        | User and Site Management | - Update User: Ensure only admins can update users.<br>- Add role-based access in UI and backend.<br><br>- Delete User: Role-based access: only admins can delete users.<br>- Add "Confirm delete" modal.<br><br>- Create Site: Implement site creation with map upload.<br>- Connect frontend to backend API.<br><br>- Update Site: Implement site update functionality.<br>- Prefill form with backend data.<br><br>- Delete Site: Add "Confirm delete" modal.<br>- Handle foreign key constraint errors.<br><br>- Add Dashboard UI: Create UI section for displaying different sites.<br>- Dynamically load site data from backend. |
|                 |
