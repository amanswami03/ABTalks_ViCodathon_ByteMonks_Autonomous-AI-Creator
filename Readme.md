# Autonomous AI Creator

This repository contains a full-stack hackathon project that builds an autonomous AI publishing agent.
The frontend is a React app built with Vite, and the backend is a Go server that manages agent lifecycle, topic discovery, publishing, and persistence.

## Project Structure

- `frontend/` - React application and UI components
- `backend/` - Go backend server, data store, agent logic, and API handlers
- `ppp.md` - prompt planning and requirements list

## Features

- Initialize an autonomous agent with a persona and domain
- Backend endpoints for agent initialization and feed retrieval
- Persistent storage with PostgreSQL for agents, posts, and seen topics
- Scheduler-driven publishing over time
- API contract compatible with the React dashboard
- Frontend polling of live feed data
- JSON API responses with consistent field names and ISO 8601 timestamps

## Local Setup

### Prerequisites

- macOS or Linux
- Go 1.22+
- Node.js 18+ and npm
- PostgreSQL (or use the fallback in-memory store during local development)

### Backend Setup

1. Open a terminal and navigate to the backend directory:

```bash
cd backend
```

2. Install dependencies and verify module tidy:

```bash
go mod tidy
```

3. Start PostgreSQL locally if you want persistence:

```bash
./setup-postgres.sh
```

4. Create a `.env` file from `.env.example` and add your Grok API key:

```bash
cp .env.example .env
```

5. Run the backend server:

```bash
go run ./cmd/server
```

The server listens on port `8080` by default.

### Frontend Setup

1. Open a new terminal and navigate to the frontend directory:

```bash
cd frontend
```

2. Install dependencies:

```bash
npm install
```

3. Run the frontend app:

```bash
npm run dev
```

4. Open the local URL shown by Vite, usually `http://localhost:5173` or `http://localhost:5174`.

## API Endpoints

The backend exposes the core API routes used by the frontend.

- `POST /api/agent/init`
  - Request body: `{ "persona": { "name": "<name>", "domain": "<domain>" } }`
  - Response: `{ "agentId": "..." }`

- `GET /api/agent/feed?agentId=<agentId>`
  - Response: `{ "posts": [ ... ] }`

- `GET /api/agent/:agentId`
- `GET /api/agent/:agentId/activity`
- `GET /api/agent/:agentId/topics`
- `GET /api/agent/:agentId/memory`
- `GET /api/agent/:agentId/stats`
- `GET /api/agent/:agentId/logs`

## Notes

- The frontend proxy is configured so `/api` requests are forwarded to the backend during development.
- `.env` should contain secrets and should not be committed.
- The backend stores posts and topic history in PostgreSQL when available.
- The agent publishes content gradually using a scheduler model, rather than all at once.

## Useful Commands

From the repo root:

```bash
# Run backend
cd backend && go run ./cmd/server

# Run frontend
cd frontend && npm run dev
```

## License

This project is provided as a hackathon implementation and does not include a license file by default.
