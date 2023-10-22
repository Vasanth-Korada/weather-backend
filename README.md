# Weather Backend Service

This project serves as a backend service for weather-related functionalities. It provides various APIs to fetch weather information, maintain search history, and perform CRUD operations on the weather data.

## Features

- **Fetch Weather**: Retrieve weather information for a specific city.
- **Search History**: Maintain a history of user's weather searches.
- **CRUD Operations**: Allows creating, reading, updating, and deleting weather records.

## Prerequisites

- Go (Version: X.XX) [Provide your Go version here]
- MySQL (Version: X.XX) [Provide your MySQL version here]
- Any other dependencies...

## Setup & Installation

1. **Clone the Repository**:
    ```bash
    git clone https://github.com/Vasanth-Korada/weather-backend.git
    cd weather-backend
    ```

2. **Setup Database**:
    - Create a database in MySQL named `weather_db`.
    - Import the SQL schema provided in the `sql` directory.

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
    go run main.go
    ```

The server will start and listen on the port specified in the `.env` file.

## API Endpoints

- **Fetch Weather**: `GET /api/weather/:city`
- **Create Weather Record**: `POST /api/weather`
- **Get Weather Search History**: `GET /api/history`
- **Delete Weather Record**: `DELETE /api/weather/:id`
- **Update Weather Record**: `PUT /api/weather/:id`