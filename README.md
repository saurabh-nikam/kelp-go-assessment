# Go Gin Request Coalescing System

This project demonstrates a high-performance concurrent API system built with **Go** and the **Gin** framework. It solves the problem of redundant resource utilization when handling multiple concurrent requests for the same data.

## Problem Statement

In a scenario where multiple users request data for the same company simultaneously, naive implementations would trigger identical, expensive calculations for every request. This system implements **Request Coalescing** (also known as request collapsing or deduplication) to ensure that:

1.  **Duplicate Requests are Coalesced**: If 100 users request financials for Company X at the same time, only *one* calculation is performed. The result is shared with all 100 users.
2.  **Shared Initial Data**: Different APIs (Financials, Sales, Employee Stats) often rely on the same "Initial Data" for a company. This system ensures this initial data is calculated only once per company, even if different API endpoints are hit simultaneously.

## Features

*   **Concurrent Request Handling**: Built on Go's robust concurrency primitives.
*   **Singleflight Pattern**: Uses `golang.org/x/sync/singleflight` to suppress duplicate function calls.
*   **Optimized Resource Usage**: drastically reduces CPU and I/O by avoiding redundant work.
*   **Clean Architecture**: Separation of concerns with Models, Service Layer, and Handlers.

## Architecture

The system uses a two-layer coalescing strategy in the `Service` layer:

1.  **API Level**: Coalesces requests to the same endpoint (e.g., multiple `GET /api/financials?companyId=123`).
2.  **Data Level**: Coalesces the underlying "Initial Data" calculation, which is shared across different endpoints.

## Getting Started

### Prerequisites

*   Go 1.21 or higher

### Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/saurabh-nikam/kelp-go-assessment.git
    cd kelp-go-assessment
    ```

2.  Install dependencies:
    ```bash
    go mod tidy
    ```

### Running the Application

Start the server:

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`.

### API Endpoints

| Method | Endpoint | Description | Example |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/financials` | Get company financials | `/api/financials?companyId=123` |
| `GET` | `/api/sales` | Get company sales data | `/api/sales?companyId=123` |
| `GET` | `/api/employee` | Get employee statistics | `/api/employee?companyId=123` |

### Testing & Verification

To verify the concurrency logic and request coalescing, run the included test suite:

```bash
go test -v ./test/...
```

**What the tests verify:**
*   **Concurrent Financials**: Simulates 5 concurrent requests. Verifies that the expensive calculation happens only once and others receive the shared result.
*   **Shared Initial Data**: Simulates concurrent requests to *different* endpoints (Financials and Sales). Verifies that the common "Initial Data" calculation runs only once.

## Project Structure

```
.
├── cmd
│   └── server
│       └── main.go          # Entry point
├── internal
│   ├── api
│   │   ├── handlers         # HTTP Handlers
│   │   └── routes           # Router setup
│   ├── models               # Data structures
│   └── service              # Business logic & Singleflight implementation
├── test
│   └── concurrency_test.go  # Verification tests
├── go.mod
└── README.md
```
