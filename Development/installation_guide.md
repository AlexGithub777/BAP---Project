# EDMS App Installation Guide

This guide will help you install and set up the project on your local machine.

## Prerequisites

Before starting, ensure the following are installed on your machine:

1. **Go 1.22.5**: [Download and Install Go](https://go.dev/doc/install).
2. **PostgreSQL Portable**: [Download PostgreSQL Portable](https://drive.google.com/file/d/14JKK4coDqtd-SqW5QGn4VizIklcd4thP/view?usp=sharing).
3. **Goose**: A database migration tool.
4. **Air**: A hot-reloading tool for Go.

## Installation Steps

### 1. Install Goose

Goose is used for managing database migrations. Install it by running:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### 2. Install Air (Hot Reloading)

Install Air for live reloading during development:

```bash
go install github.com/air-verse/air@latest
```

### 3. Clone the Repository

Create a new folder and clone the project from GitHub:

```bash
git clone https://github.com/AlexGithub777/BAP---Project.git
```

### 4. Set Up PostgreSQL Database

1. Extract the downloaded `PostgreSQL.zip` file from the provided Google Drive link.

2. Open the extracted folder.

3. Run `startdb.bat` to start the PostgreSQL server. **(don't close the terminal window)**

### 5. Open the Project in VSCode

Open the cloned repository in Visual Studio Code.

Then in the terminal navigate to the project root.

```bash
cd .\Development\EDMS\
```

### 6. Create the `.env` File

In the root directory of the project `BAP---Project\Development\EDMS`, create a new file named `.env`.

Get the necessary environment variable values (such as `DB_USER`, `DB_PASSWORD`, and `JWT_SECRET`) from the `#dotenv` channel in the project's Discord server.

The structure of the `.env` file should look like this:

```bash
DB_USER="your_postgres_username"
DB_PASSWORD="your_postgres_password"
DB_NAME="your_database_name"
DB_HOST="localhost"
DB_PORT=5432
JWT_SECRET="your_jwt_secret"
```

### 7. Run Database Migrations

To initialize the database tables, run the `goose_up.bat` script:

```bash
./goose_up.ps1
```

This will create the necessary tables in your PostgreSQL database.

### 8. Start the Application with Air

Run the application using Air.
Air will automatically rebuild the project when changes are detected.

If you're using Air for the first time, you may need to specify the config file:

Run in project root:

```bash
air -c .air.toml
```

Then after you should be able to use:

```bash
air
```

### 9. Access the Application

Once the application is running, the terminal should display a link:

```bash
http://localhost:3000
```

`Ctrl + Click` or open this link in your browser to access the application. You can create an account and log in.

### 10. Troubleshooting

GOPATH Environment Variable
If you encounter errors related to Go paths, ensure your GOPATH is set correctly. [Follow this guide to set your GOPATH.](https://go.dev/wiki/SettingGOPATH)

Additionally, add `C:\Users\%USERPROFILE%\go\bin` to your system PATH:

1. Open **System Properties** (`Windows Key + Pause`).

2. Click **Advanced System Settings**.

3. Click **Environment Variables**.

4. Under User Variables, click **New** and add `C:\Users\%USERPROFILE%\go\bin`.

After updating your GOPATH, restart your terminal and try again.

<hr>

For any issues or further assistance, feel free to ask in the project's Discord channel!
