# API_CONTRACT.md

# Autonomous AI Creator - API Contract

This document defines the communication between the React frontend and Go backend.

**Purpose**
- Prevent confusion while development.
- Allow frontend and backend to be developed independently.
- Keep request and response formats consistent.

---

# Base URL

Development

http://localhost:8080/api

Production

https://your-domain.com/api

---

# Response Format

## Success

```json
{
  "success": true,
  "data": {}
}
```

## Error

```json
{
  "success": false,
  "message": "Something went wrong"
}
```

---

# 1. Initialize Agent

## Endpoint

POST /api/agent/init

## Purpose

Creates a new autonomous AI agent.

This endpoint is REQUIRED by the hackathon.

---

## Request

```json
{
  "persona": {
    "name": "Nova",
    "domain": "AI Security"
  }
}
```

### Fields

| Field | Type | Required | Description |
|--------|------|----------|-------------|
| name | string | ✅ | Name of AI persona |
| domain | string | ✅ | AI niche/domain |

---

## Success Response

```json
{
  "agentId": "abc123"
}
```

---

## Frontend Flow

User clicks

Create Agent

↓

Frontend sends POST request

↓

Backend creates agent

↓

Returns agentId

↓

Frontend redirects to Dashboard

---

# 2. Feed

## Endpoint

GET /api/agent/feed?agentId=abc123

## Purpose

Returns all published posts.

This endpoint is REQUIRED by the hackathon.

---

## Query Parameters

| Name | Required | Description |
|------|----------|-------------|
| agentId | ✅ | Agent ID returned during initialization |

---

## Response

```json
{
  "posts": [
    {
      "id": "post_1",
      "createdAt": "2026-08-07T10:30:00Z",
      "text": "Prompt Injection is becoming a serious concern...",
      "rationale": "Chosen because...",
      "sources": [
        "https://..."
      ]
    }
  ]
}
```

---

# 3. Agent Details

## Endpoint

GET /api/agent/:agentId

## Purpose

Returns general information about the agent.

---

## Response

```json
{
    "id":"abc123",
    "name":"Nova",
    "domain":"AI Security",
    "status":"ACTIVE",
    "createdAt":"2026-08-07T08:00:00Z"
}
```

---

# 4. Current Activity

## Endpoint

GET /api/agent/:agentId/activity

## Purpose

Returns the current task of the AI.

Used for live dashboard.

---

## Response

```json
{
    "status":"Searching",
    "currentTask":"Reading OpenAI Blog",
    "progress":65
}
```

---

# 5. Topic Queue

## Endpoint

GET /api/agent/:agentId/topics

## Purpose

Returns accepted and rejected topics.

---

## Response

```json
{
    "accepted":[
        {
            "title":"MCP Security",
            "score":92
        }
    ],
    "rejected":[
        {
            "title":"GPT-6 Launch",
            "reason":"Already covered"
        }
    ]
}
```

---

# 6. Memory

## Endpoint

GET /api/agent/:agentId/memory

## Purpose

Returns important memories used by the AI.

---

## Response

```json
{
    "interests":[
        "Prompt Injection",
        "MCP",
        "AI Security"
    ],
    "recentTopics":[
        "Claude",
        "Gemini",
        "OpenAI"
    ]
}
```

---

# 7. Analytics

## Endpoint

GET /api/agent/:agentId/stats

## Purpose

Returns dashboard statistics.

---

## Response

```json
{
    "published":15,
    "rejected":38,
    "memoryNodes":122,
    "sources":7
}
```

---

# 8. Logs

## Endpoint

GET /api/agent/:agentId/logs

## Purpose

Returns recent agent activity.

---

## Response

```json
[
    {
        "time":"10:45",
        "action":"Searching OpenAI Blog"
    },
    {
        "time":"10:47",
        "action":"Comparing Topics"
    },
    {
        "time":"10:48",
        "action":"Rejected GPT News"
    },
    {
        "time":"10:49",
        "action":"Publishing MCP Security"
    }
]
```

---

# Frontend Routes

These are React routes.

| Route | Purpose |
|--------|----------|
| / | Landing Page |
| /create | Create Agent |
| /dashboard/:agentId | Dashboard |
| /dashboard/:agentId/feed | Feed |
| /dashboard/:agentId/topics | Topic Queue |
| /dashboard/:agentId/memory | Memory |
| /dashboard/:agentId/analytics | Analytics |
| /dashboard/:agentId/settings | Settings |

---

# Development Workflow

User opens app

↓

Creates Persona

↓

POST /api/agent/init

↓

Receives agentId

↓

Redirects to Dashboard

↓

Dashboard fetches

- Agent Details
- Activity
- Feed
- Memory
- Topics
- Analytics

↓

Auto refresh every few seconds

---

# Notes

- Follow JSON formats exactly.
- Keep field names consistent.
- Always return ISO 8601 timestamps.
- Never expose API keys to the frontend.
- Validate agentId before processing requests.
- The two required hackathon endpoints are:
  - POST /api/agent/init
  - GET /api/agent/feed
- All other endpoints are for frontend convenience and can be adjusted if both frontend and backend agree.

---

# Frontend & Backend Responsibilities

## Frontend (React)

- UI/UX
- Forms
- Dashboard
- API calls
- Loading states
- Error handling
- Responsive design

## Backend (Go)

- Agent initialization
- AI orchestration
- Topic discovery
- Memory management
- Feed generation
- Scheduler
- Database
- HTTP APIs