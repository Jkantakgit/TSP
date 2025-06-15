# Travelling Salesman Problem (TSP)

A React-based interactive app to add cities on a map and solve the Travelling Salesman Problem (TSP) using a backend optimization service.

## Instalation

### Prerequisites

- Node.js
- Golang (1.24.x)
- Make
- Git

### Download

```bash
git clone https://github.com/Jkantakgit/TSP.git
cd TSP
```

### Build

``` bash
make all
```

### Running the app

Since the backend and frontend run as separate services, start them in two separate terminals:
Terminal 1 — Start backend server:

```bash
make run-backend
```

Terminal 2 — Start frontend development server:

```bash
make run-frontend
```

## Backend

The backend is written in Golang and is responsible for optimalizating the Travelling Salesman Problem (TSP).

### Overview

- Receives list of cities and it's possitions from frontend via REST API
- Runs Ant Colony Optimalisation (ACO) algorithm for finding the shortest route visiting all cities
- Returns Optimalised route to the frontend

By default, the backend listens on http://localhost:8080.

API Endpoint

- POST /optimize: Accepts JSON payload with cities, returns the optimized route.

Example request body:

```plaintext
{
  "cities": [
    { "name": "City 1", "x": 10.5, "y": 20.3 },
    { "name": "City 2", "x": 50.2, "y": 60.1 }
  ]
}
```

## Frontend

The frontend is written in TypeScript and using Vite as build tool.

### Frontend Features

- Clicking on the map adds cities
- Visually displays the citties and optimalised routes
- Shows the route as list of cities

## License
This project is licensed under the [GPL License](./LICENSE). See the LICENSE file for details.
