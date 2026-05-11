# go-bulkify

A Go tool for sending bulk HTTP requests concurrently against any API endpoint.

This project is designed to demonstrate:
- Concurrent bulk request execution via goroutine workers
- Flexible, file-based configuration with environment variable overrides
- Generic payload builder pattern for any request body shape
- Clean separation between config, runner, and payload logic

---

## ✨ Features

- Sends bulk HTTP requests concurrently using a configurable worker pool.
- Supports any HTTP method: `POST`, `PUT`, `PATCH`, etc.
- Configures URL, headers, workers, and timeout via a central `config.json`.
- Overrides any config value at runtime through environment variables — no file edits needed.
- Supports multiple environment configs: `config.dev.json`, `config.prod.json`, etc.
- Generic `Wrap[T]` helper lets you plug in any typed payload struct without touching other files.

---

## 🧠 Business Rules

1. All requests are distributed across `worker_count` goroutines running in parallel.
2. Each request body is built by `BuildPayload(index)`, where `index` runs from `1` to `total_records`.
3. A request is counted as failed when the server returns HTTP `4xx` or `5xx`.
4. Environment variables always take precedence over values in the JSON config file.

---

## ⚙️ Configuration

Edit `config.json` to set the target API and tuning parameters:

```json
{
  "url": "https://httpbin.org/post",
  "method": "POST",
  "headers": {
    "Content-Type": "application/json",
    "Authorization": "Bearer <token>"
  },
  "worker_count": 10,
  "total_records": 4000,
  "timeout_seconds": 10
}
```

### Environment Variable Overrides

Any field can be overridden at runtime without editing the file:

| Variable | Overrides |
|---|---|
| `BULKIFY_URL` | `url` |
| `BULKIFY_METHOD` | `method` |
| `BULKIFY_WORKER_COUNT` | `worker_count` |
| `BULKIFY_TOTAL_RECORDS` | `total_records` |
| `BULKIFY_TIMEOUT_SECONDS` | `timeout_seconds` |

```bash
BULKIFY_URL=https://prod.api.com BULKIFY_WORKER_COUNT=20 ./go-bulkify
```

---

## 🔧 Custom Payload

Edit `payload/payload.go` to define your own request body.

**Option 1 — edit the default struct directly:**
```go
type Payload struct {
    ProductName string `json:"product_name"`
    Status       string `json:"status"`
}

func BuildPayload(index int) (any, string) {
    name := fmt.Sprintf("Product_%05d", index)
    return Payload{ProductName: name, Status: "active"}, name
}
```

**Option 2 — keep the file untouched and use `Wrap` in `main.go`:**
```go
type Order struct {
    OrderID int    `json:"order_id"`
    Product string `json:"product"`
}

cfg.BuildPayload = payload.Wrap(func(i int) (Order, string) {
    label := fmt.Sprintf("order-%05d", i)
    return Order{OrderID: i, Product: "Widget"}, label
})
```

---

## 🚀 How to Run

### Option 1: Run with Go
```bash
go run main.go
```

Use a specific config file:
```bash
go run main.go -config config.prod.json
```

### Option 2: Build and run binary
```bash
go build -o go-bulkify .
./go-bulkify -config config.json
```

### Option 3: Run with Docker
Build the image:
```bash
docker build -t go-bulkify:dev .
```
Run the bulk requests:
```bash
docker run --rm go-bulkify:dev ./main -config config.json
```
Or keep the container running and exec in:
```bash
docker run -d --name go-bulkify go-bulkify:dev
docker exec go-bulkify ./main -config config.json
```

---

## 📊 Example Output

```
🚀 Starting: url=https://httpbin.org/post  method=POST  workers=5  total=20
2026/05/11 10:00:01 ✅ [Worker 2 | #1] SUCCESS — Product_00001
2026/05/11 10:00:01 ✅ [Worker 1 | #2] SUCCESS — Product_00002
2026/05/11 10:00:02 ❌ [Worker 3 | #5] FAILED  — server returned status 500
...
🏁 Done! Success: 19 | Failed: 1 | Total: 20
```

---

## 📁 Project Structure

```
go-bulkify/
├── config/
│   └── config.go        # Config struct + LoadFromFile + env overrides
├── internal/
│   └── runner/
│       ├── client.go    # HTTP request execution
│       └── worker.go    # Worker pool and Run()
├── payload/
│   └── payload.go       # Payload struct, BuildPayload, Wrap[T]
├── config.json          # Central configuration file
├── main.go
└── Dockerfile
```

---

## 📄 License
MIT License