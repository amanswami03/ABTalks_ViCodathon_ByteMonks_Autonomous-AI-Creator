# Autonomous AI Creator — Backend Skeleton

Go backend for the Abtalks "Autonomous AI Creator" hackathon challenge.

## What's already built
- `POST /api/agent/init` — creates a persona, starts the autonomous loop
- `GET /api/agent/feed?agentId=...` — returns posts, newest first
- Topic discovery via Hacker News (free, no key)
- Grok API client for editorial decisions + writing
- In-memory store with per-agent memory (seen topics + published posts)
- A background scheduler that runs the full cycle every 20 minutes, forever, with no further human input

## Setup

```bash
cp .env.example .env
# paste your Grok API key into .env

go mod init ai-persona-agent   # if go.mod isn't already resolved
go get github.com/google/uuid
go mod tidy

go run cmd/server/main.go
```

Test locally:
```bash
curl -X POST localhost:8080/api/agent/init \
  -H "Content-Type: application/json" \
  -d '{"persona":{"name":"Ada","domain":"AI Security"}}'

curl "localhost:8080/api/agent/feed?agentId=PASTE_ID_HERE"
```

## Where to get things

| Need | Source |
|---|---|
| Grok API key | https://docs.x.ai/docs/api-overview |
| Topic source (already wired) | Hacker News public API, no key needed |
| Optional 2nd topic source | Tavily / NewsAPI / arXiv API — add a new file in `internal/topics/` |

## TODOs left for you / Copilot (search "TODO(copilot)" in the code)
1. **Persona system prompt** (`internal/api/handlers.go`) — make it more opinionated and specific, add 2-3 example sentences in the voice you want.
2. **Topic filtering** (`internal/topics/hackernews.go`) — currently pulls raw top stories; add keyword filtering or a second source for better relevance.
3. **Scheduler interval** (`internal/api/handlers.go`, the `StartScheduler` call) — tune how often it runs.
4. **Persistence** — currently everything is in-memory and resets on restart. Swap `internal/store/store.go` for PostgreSQL if your deploy target restarts the process.
5. **Logging of skipped topics** — currently skip decisions aren't stored anywhere visible; consider logging them somewhere you can show judges "look, it rejects things too."
6. **Deploy** — pick a host that stays alive for 48h continuously (Railway, Render, Fly.io — avoid serverless functions that sleep, since the scheduler needs to keep running in the background).

## Important: don't add a manual trigger
There is intentionally no "generate post now" endpoint. The evaluator only calls `/init` once and `/feed` repeatedly — publishing must happen entirely on its own via the scheduler.
