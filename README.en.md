# USD/BRL Quote Project (Go + SQLite)

This project exposes an HTTP endpoint that fetches the USD/BRL quote from an external API, stores it in SQLite, and returns only the `bid` value. A simple client calls the server to demonstrate timeouts.

## Architecture (summary)

- `cmd/server`: HTTP server.
- `cmd/client`: example HTTP client.
- `pkg/types`: shared DTOs.
- `app.db`: SQLite database created automatically.

## Endpoints

- `GET /cotacao`

Response:
```json
{ "bid": "5.1234" }
```

## How to run

Requires Go 1.25.6.

Build:
```bash
make build
```

Run server:
```bash
make run-server
```

Run client:
```bash
make run-client
```

Run server + client (client runs and server is stopped afterwards):
```bash
make run
```

## Key settings

In the server (`cmd/server/main.go`):
- `serverAddr`: `:8080`
- `priceAPIURL`: `https://economia.awesomeapi.com.br/json/last/USD-BRL`
- `callTimeout`: 200ms (external API timeout)
- `dbTimeout`: 10ms (SQLite insert timeout)

In the client (`cmd/client/main.go`):
- `apiAddress`: `http://localhost:8080/cotacao`
- `callTimeout`: 300ms (calls expected to succeed)
- `callTimeoutToFail`: 2ms (calls expected to fail)
- `apiCalls`: 5 (number of calls per batch)

## Database

The server creates the `prices` table on startup and clears it on each start. The database file is `app.db` at the project root.

## Notes

- The endpoint returns only the `bid` field.
- External API or DB errors return `500`.

