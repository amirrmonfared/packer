# Packer Service

A service that calculates the optimal distribution of packs (boxes) to fulfill a given order size.  

This service:

1. Stores a configurable set of pack sizes (e.g., 250, 500, 1000, etc.).  
2. Calculates how many total boxes (and in which sizes) are needed to ship a specified number of items, minimizing leftover items as well as the number of boxes used.

## Table of Contents

1. [Features](#features)
2. [Installation and Requirements](#installation-and-requirements)
3. [Usage](#usage)
4. [Docker Instructions](#docker-instructions)
5. [Endpoints](#endpoints)
6. [Testing](#testing)
7. [Swagger Documentation](#swagger-documentation)

---

## Features

- **Memory-based pack size storage**: Pack sizes can be fetched and updated at runtime using the API.  
- **Optimal calculation**: Uses a dynamic programming approach to find the minimal leftover, and among equal leftovers, the minimal total box count.  
- **Prometheus metrics**: Exposes metrics via `/metrics`.  
- **Health check**: Simple `/healthz` endpoint to verify service uptime.  
- **Swagger documentation**: API documentation available at `/swagger/index.html`.  

---

## Installation and Requirements

1. **Go 1.23.4+** (for building and running the service).
2. **Make** (optional, if you want to use the provided Makefile).
3. **Docker** (optional, if you want to build and run the Docker image).

To install required modules locally, run:

```bash
go mod download
```

---

## Usage

### Running Locally

1. **Clone** this repository.  
2. **Build** and run:

   ```bash
   cd packer
   go build -o packer ./cmd/packer
   ./packer --port=8080
   ```

   Alternatively, you can run it without building:

   ```bash
   cd packer
   go run ./cmd/packer/main.go --port=8080
   ```

   By default, the service listens on port `8080`. You can override that by passing the `--port` flag.

### Accessing the API

Once the service is running (e.g., on `http://localhost:8080`), you can:

- Update pack sizes,
- Get current pack sizes,
- Calculate needed packs for an order.

A simple cURL command to check health:

```bash
curl http://localhost:8080/healthz
```

Should return:

```json
{
  "message": "OK"
}
```

---

## Docker Instructions

### Building the Image

A Dockerfile is located in `images/packer/Dockerfile` (or you can use the main Dockerfile in the repository root). For example, to build from the root `Dockerfile`, run:

```bash
docker build -t packer-service .
```

### Running the Container

```bash
docker run -d -p 8080:8080 packer-service
```

You can then visit `http://localhost:8080/healthz` to ensure the service is responding.

---

## Endpoints

Below are the primary endpoints exposed under `/api/v1`:

1. **Update pack sizes**  
   - **Endpoint**: `POST /api/v1/packs`  
   - **Request Body** (JSON):

     ```json
     {
       "packs": [250, 500, 1000, 2000, 5000]
     }
     ```

   - **Response Body** (JSON):

     ```json
     {
       "packs": [250, 500, 1000, 2000, 5000]
     }
     ```

   - **Description**: Updates the entire list of available pack sizes in memory.

2. **Get pack sizes**  
   - **Endpoint**: `GET /api/v1/packs`  
   - **Response Body** (JSON):

     ```json
     {
       "packs": [250, 500, 1000, 2000, 5000]
     }
     ```

   - **Description**: Retrieves the currently available pack sizes.

3. **Calculate required packs**  
   - **Endpoint**: `POST /api/v1/calculate`  
   - **Request Body** (JSON):

     ```json
     {
       "items": 1200
     }
     ```

   - **Response Body** (JSON):

     ```json
     {
       "order": 1200,
       "leftover": 0,
       "total_packs": 2,
       "distribution": {
         "500": 1,
         "700": 1
       },
       "total_items_shipped": 1200
     }
     ```

   - **Description**: Calculates how many packs are needed (and which sizes) to fulfill the requested number of items. Minimizes leftover first, then uses the fewest packs possible.

4. **Prometheus metrics**  
   - **Endpoint**: `GET /metrics`  
   - **Description**: Exposes default Go and custom metrics for Prometheus monitoring.

5. **Health check**  
   - **Endpoint**: `GET /healthz`  
   - **Description**: Returns a simple `{"message": "OK"}` to confirm the service is up.

6. **Swagger docs**  
   - **Endpoint**: `GET /swagger/index.html` (generally accessible at `/swagger/` in your browser).  
   - **Description**: Interactive API documentation.

---

## Testing

This repository includes unit tests for various components (Fulfillment logic, Store concurrency, etc.). To run tests:

```bash
go test ./...
```

If you have a `Makefile`, you could also run:

```bash
make test
```

---

## Swagger Documentation

The Swagger files are located in `docs/packer/`. They may be generated using [Swag](https://github.com/swaggo/swag) if you need to update them:

1. Install `swag` tool if not already installed:

   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. Generate the docs:

   ```bash
   swag init -g cmd/packer/main.go --output docs/packer
   ```

3. Re-build or re-run the application to serve the updated documentation.
