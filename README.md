
# Weather Query Application

This application provides a weather query service that aggregates requests for the same location to minimize calls to external weather APIs. It fetches data from two different weather services and returns the average temperature.

## Features

- **Request Aggregation**: Groups incoming requests by location and holds them for up to 5 seconds to reduce redundant API calls.
- **Concurrent Processing**: Handles multiple requests concurrently and efficiently.
- **Asynchronous Database Logging**: Saves query results to a database without blocking user requests.
- **External API Integration**: Fetches weather data from `weatherapi.com` and `weatherstack.com`.
- **Configurable**: Application settings can be managed through a `config.yml` file.
- **Tested**: Includes unit and integration tests to ensure reliability.

## Requirements

- Go 1.18 or higher
- A running instance of SQLite, PostgreSQL, MySQL, or MariaDB (the project is configured to use SQLite by default).

## Installation

1.  **Clone the repository:**
    ```sh
    git clone <repository-url>
    cd weather-query-application
    ```

2.  **Install dependencies:**
    ```sh
    go mod tidy
    ```

## Configuration

The application uses a `config.yml` file for configuration. You can modify this file to change the server port, database connection, and API keys.

```yaml
server:
  port: 8080

database:
  dsn: "weather.sqlite"

weather_api:
  weather_api_key: "your_weather_api_key"
  weather_stack_key: "your_weather_stack_key"
```

## Running the Application

To run the application, use the following command:

```sh
go run cmd/server/main.go
```

The server will start on the port specified in the `config.yml` file (default is `8080`).

## API Endpoint

### Get Weather

- **Endpoint**: `/weather`
- **Method**: `GET`
- **Query Parameters**:
  - `q`: The location you want to query (e.g., `Istanbul`).
- **Success Response** (`200 OK`):
  ```json
  {
    "location": "<location>",
    "temperature": <average-temp>
  }
  ```
- **Error Response** (`400 Bad Request`):
  ```
  location is required
  ```

### `curl` Examples

Here are some `curl` examples to test the endpoint:

**1. Single request for a location:**

```sh
curl "http://localhost:8080/weather?q=Istanbul"
```

**2. Multiple requests for the same location (run in separate terminals):**

```sh
# Terminal 1
curl "http://localhost:8080/weather?q=London"

# Terminal 2
curl "http://localhost:8080/weather?q=London"
```

**3. Request for a different location:**

```sh
curl "http://localhost:8080/weather?q=Paris"
```

**4. Request without a location (will return an error):**

```sh
curl "http://localhost:8080/weather"
```

## Postman Collection

A Postman collection is available in the `postman_collection.json` file. You can import this file into Postman to easily test the API endpoints. 

## Database

The application logs every successful weather query to a database. By default, it uses SQLite, and the database file is named `weather.sqlite`.

### Table Schema

The data is stored in the `weather_queries` table with the following schema:

-   `id` (INTEGER, PRIMARY KEY AUTOINCREMENT)
-   `location` (TEXT)
-   `service_1_temperature` (REAL)
-   `service_2_temperature` (REAL)
-   `request_count` (INTEGER) - The number of grouped requests for this query.
-   `created_at` (DATETIME)

### Verifying Logs

You can check the logs by querying the database directly. For SQLite, you can use the following command in your terminal:

```sh
sqlite3 -header weather.sqlite "SELECT * FROM weather_queries;"