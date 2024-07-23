# ZeroTier Webhook Service

This application is demonstrating webhook service for ZeroTier Central events. It includes a server for receiving webhook events, storing them in a database, and providing an endpoint for searching stored events.

## Getting Started

### Prerequisites

- Go (version 1.19 or later recommended)
- Docker (optional for containerization)

### Installing and Running

1. Clone the repository:

    ```bash
    git clone https://github.com/akhilesharora/zerotier-webhook
    cd zerotier-webhook
    ```

2. Build and run the server application:

    ```bash
    make build
    make run
    ```

3. Alternatively, you can use Docker to build and run the application:

    ```bash
    make docker-build
    make docker-run
    ```

### Application Overview

The application provides a server for handling ZeroTier Central webhook events and storing them in a SQLite database.

### Server

The server handles:
- Receiving webhook events from ZeroTier Central.
- Storing events in a SQLite database.
- Providing an endpoint for searching stored events.

### Flow of Control

- ZeroTier Central sends webhook events to the server.
- The server processes these events and stores them in the database.
- Clients can query the server to search for stored events.

### API Endpoints

### POST /webhook
Receives webhook payloads from ZeroTier Central.

### Request Body:
The request body should contain the webhook payload in JSON format.

### Response:
- `200 OK`: Webhook processed successfully
- `400 Bad Request`: Error reading request body
- `500 Internal Server Error`: Error processing payload

### GET /search
Searches for stored events based on provided query parameters.

### Query Parameters:
- `network_id` (optional): Filter events by network ID
- `member_id` (optional): Filter events by member ID
- `user_id` (optional): Filter events by user ID

### Response:
- `200 OK`: Returns a JSON array of matching events
- `500 Internal Server Error`: Error occurred while fetching events

### Example:
```bash
curl -X GET "http://localhost:8080/search?network_id=8056c2e21c000001&member_id=12345&user_id=user@example.com"
```

### Project Structure

- `cmd/`: Contains the main application code.
- `pkg/`: Library code for database operations and HTTP handlers.
- `Dockerfile`: Dockerfile for building the application.
- `Makefile`: Automates build and run tasks.

### Caveats and Limitations

- **Concurrency Handling**: Basic concurrency handling is implemented. Future versions could aim to improve this for high-load scenarios.
- **Error Handling**: Basic error handling is implemented; more contextual errors can be added like Graceful shutdown of the service.
- **Testing Coverage**: Basic coverage for major functionalities. Edge cases and stress conditions can be added for improvement.
- **Data Persistence**: Currently uses SQLite for data storage.
- **Authentication**: The current implementation does not include signature verification for the search endpoint. This should be added for production use.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
