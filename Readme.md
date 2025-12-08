# User management

User management system is a simple go API built on top of the [chi router](https://github.com/go-chi/chi) with basic CRUD functionalities.

| ![/assets/UserManagementSystem.drawio.png](./assets/UserManagementSystem.drawio.png) | 
|:--:| 
| *Architecture diagram* |

## Features

### 1. User Management Rest API
- Users CRUD operations
- Persistant data storage

### 2. Instrument Management Rest API
- Instrument CRUD operations
- Persistant data storage

### 4. Supports three levels of configuration
- Supports `--config config.yaml`
- Environment variable overrides (`USRM_*`)
- Dynamic subscriptions via YAML config

### 5. Retry
- Database connection Retry handling

### 6. Graceful Shutdown
- OS signal handling
- HTTP server graceful stop
- Database connection termination

## Installation

**1. Clone the repository:**
```bash
git clone https://github.com/ChethiyaNishanath/user-management.git
cd user-management
```

**2. Install dependencies:**
```bash
go mod tidy
```

## Configuration

Create a `config.yaml`:
```yaml
server:
  port: 8080
  shutdownTimeout: 10

database:
  dsn: "postgres://postgres:Test1234@localhost:5432/usermanagementdb?sslmode=disable"

logging:
  level: INFO
```

## Running the Server

The system uses Cobra commands.

Start the HTTP server:
```bash
user-management serve --config config.yaml
```

Console output:
```
Configuration initialized. Using config file: config.yaml
Server starting on port 8080
```

## User Management API Usage

### Create User
`[POST] /users`

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
        "firstName" : "Chethiya",
        "lastName" : "Viharagama",
        "email" : "chethiya@example.com",
        "phone" : "+941234352",
        "age" : 11,
        "status" : "Active"
    }'
```

### Get All Users
`[GET] /users`

```bash
curl -X GET "http://localhost:8080/users?page=1&limit=10" \
  -H "Content-Type: application/json"
```

### Get User by Id
`[GET] /users/{userId}`

```bash
curl -X GET http://localhost:8080/users/{userId} \
  -H "Content-Type: application/json"
```

### Update User
`[PATCH] /users/{userId}`

```bash
curl -X PATCH http://localhost:8080/users/{userId} \
  -H "Content-Type: application/json" \
  -d '{
        "firstName" : "Chethiya",
        "lastName" : "Viharagama",
        "email" : "chethiya@example.com",
        "phone" : "+941234352",
        "age" : 11,
        "status" : "Active"
    }'
```

### Delete User
`[DELETE] /users/{userId}`

```bash
curl -X DELETE http://localhost:8080/users/{userId} \
  -H "Content-Type: application/json"
```

## Instrument API Usage

### Create Instrument
`[POST] /instruments`

```bash
curl -X POST http://localhost:8080/instruments \
  -H "Content-Type: application/json" \
  -d '{
        "Symbol": "AAPL",
        "Name": "Apple Inc.",
        "Instrument_Type": "Equity",
        "Exchange": "NASDAQ",
        "Last_Price": 226.43
    }'
```

### Get All Instruments
`[GET] /instruments`

```bash
curl -X GET http://localhost:8080/instruments \
  -H "Content-Type: application/json"
```

### Get Instrument by Id
`[GET] /instruments/{instrumentId}`

```bash
curl -X GET http://localhost:8080/instruments/{instrumentId} \
  -H "Content-Type: application/json"
```

### Update Instrument
`[PATCH] /instruments/{instrumentId}`

```bash
curl -X PATCH http://localhost:8080/instruments/{instrumentId} \
  -H "Content-Type: application/json" \
  -d '{
        "Symbol": "AAPL",
        "Name": "Apple Inc.",
        "Instrument_Type": "Equity",
        "Exchange": "NASDAQ",
        "Last_Price": 226.43
    }'
```

### Delete Instrument
`[DELETE] /instruments/{instrumentId}`

```bash
curl -X DELETE http://localhost:8080/instruments/{instrumentId} \
  -H "Content-Type: application/json"
```

## CLI

List all commands
```bash
user-management --help
```

Server
```bash
user-management serve --config config.yaml
```

## Testing

Unit tests
```bash
go test ./unit/...
```

Integration tests
```bash
go test ./it/...
```