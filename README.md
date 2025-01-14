# Bitcoin LTP Service

This repository contains a backend written in Go and a frontend in React that retrieves the last traded price (LTP) of Bitcoin for specified currency pairs using the Kraken API. The backend also supports caching responses for 60 seconds to optimize performance.

## Features

- Retrieve Bitcoin LTP for the following pairs: BTC/USD, BTC/CHF, BTC/EUR.
- Cache LTP responses for 60 seconds.
- Dockerized for ease of deployment.
- Integration tests included for backend.

---

## Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

---

## Running the Application

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-directory>

2. Build and start the Docker containers:
   ```bash
   docker-compose up --build

3. Access the frontend:

   The React frontend will be available at http://localhost:3000.
   


3. Access the backend:

   The React frontend will be available at http://localhost:8082.
   Example API endpoint: http://localhost:8080/api/v1/ltp?pairs=BTCUSD,BTCCHF
   
