# API Usage Guide

Mind-OS provides a RESTful API to interact with the artificial brain. This guide demonstrates how to use the core endpoints using `curl`.

**Base URL**: `http://localhost:8080`
**Swagger UI**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## 1. Sensory Input (Chat & Signals)

Send text or physical signals to the brain. This triggers the **Amygdala** (emotion), **Hippocampus** (memory retrieval/encoding), and **PFC** (response generation).

**Endpoint**: `POST /api/v1/sensory-inputs`

### Chat Example
```bash
curl -X POST http://localhost:8080/api/v1/sensory-inputs \
  -H "Content-Type: application/json" \
  -d '{
    "type": "chat",
    "text": "Hello, how are you feeling today?"
  }'
```

### Response Example
```json
{
  "mindState": {
    "currentReaction": [
      { "code": "JOY", "value": 50 }
    ],
    "moodStability": 0.8,
    "motivation": 65,
    "sanity": 100,
    "replyText": "I am feeling quite productive and balanced."
  },
  "reply": "I am feeling quite productive and balanced.",
  "debug": {
    "cortisol": 10,
    "oxytocin": 40,
    "predictedReward": 15
  }
}
```

---

## 2. Check Brain State

Get a snapshot of the current internal state (hormones, motivation, sanity). This endpoint supports **ETag** caching.

**Endpoint**: `GET /api/v1/brain-states/current`

### Request
```bash
curl -v http://localhost:8080/api/v1/brain-states/current
```

### Response Example
```json
{
  "motivation": 65,
  "motivationLevel": "normal",
  "sanity": 100,
  "sanityLevel": "stable",
  "stmCount": 5,
  "ltmCount": 120
}
```

**Note**: If you send the `If-None-Match` header with the previous ETag, the server will return `304 Not Modified` if the state hasn't changed.

---

## 3. Sleep Cycle (Memory Consolidation)

Trigger the sleep process to move Short-Term Memories (STM) to Long-Term Memories (LTM) and perform forgetting/cleanup.

**Endpoint**: `POST /api/v1/sleep-cycles`

### Request
```bash
curl -X POST http://localhost:8080/api/v1/sleep-cycles
```

### Response Example
```json
{
  "message": "Sleep consolidation completed",
  "consolidatedCount": 3,
  "forgottenCount": 1,
  "stmCount": 0,
  "ltmCount": 123
}
```

---

## 4. Daydreaming (DMN Activation)

Trigger the Default Mode Network (DMN) to simulate mind-wandering or random memory recall.

**Endpoint**: `POST /api/v1/daydreams`

### Request
```bash
curl -X POST http://localhost:8080/api/v1/daydreams
```

### Response Example
```json
{
  "message": "Daydreaming... (Stub)"
}
```

---

## 5. Motivation & Feedback

Manage the reward system (Basal Ganglia).

### Positive Feedback (Reward)
```bash
curl -X POST http://localhost:8080/api/motivation/feedback \
  -H "Content-Type: application/json" \
  -d '{ "isPositive": true }'
```

### Apply Motivation Decay (Time passing)
```bash
curl -X POST http://localhost:8080/api/motivation/decay
```
