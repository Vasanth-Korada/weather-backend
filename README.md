# Weather Wonder Backend (GoLang, Postgres)

This project serves as a backend service for weather-wonder app functionalities. It provides various APIs to fetch weather information, maintain search history, and perform CRUD operations on the weather data.

<img width="50%" src="https://github.com/Vasanth-Korada/weather-app-frontend/assets/50695446/d5e774fa-38f1-4b7e-9cc5-e1ff1282a228">

[Weather Wonder FullStack App YouTube Demo Link](https://youtu.be/df1Wx9Bwt10)

## API Endpoints and Features

- **POST /create_weather**: Creates a new weather entry.
- **DELETE /delete_weather/:id**: Deletes the weather entry with the specified ID.
- **PUT /update_weather/:id**: Updates the weather entry with the specified ID.
- **GET /history**: Retrieves the history of weather entries.
- **POST /login**: Authenticates a user and starts a session.
- **POST /register**: Registers a new user account.
- **POST /logout**: Ends the user's session and logs them out.

- **Fetch Weather**: Retrieve weather information for a specific city.
- **Search History**: Maintain a history of the user's weather searches.
- **CRUD Operations**: Allows creating, reading, updating, and deleting weather records.
- **JWT Authentication**: Used Token-based JWT Authentication for Login, Register and Logout.
- **Error Handling and Validation**: Added proper error handling and validation.
- **DB Constraints**: Added appropriate DB constraints like Primary Key, Unique Key, Foreign Key

## Prerequisites

- Go 
- PostgreSQL

## Setup & Installation

1. **Clone the Repository**:
    ```bash
    git clone https://github.com/Vasanth-Korada/weather-backend.git
    cd weather-backend
    ```

2. **Setup Database**:
    - Create a database in MySQL named `weathers`.
    - Import the SQL schema provided in the `models` directory.

3. **Environment Variables**: 
    Rename `.env.example` to `.env` and update the variables accordingly.
    ```bash
    DB_USER=your_db_user
    DB_PASSWORD=your_db_password
    ...
    ```

4. **Install Dependencies**:
    ```bash
    go mod download
    ```

5. **Run the Server**:
    ```bash
    go run *.go
    ```

The server will start and listen on the port specified in the `.env` file.
