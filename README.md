# Meme API

> **Fork notice:** This project is a fork of [D3vd/Meme_Api](https://github.com/D3vd/Meme_Api).

JSON API for random memes scraped from Reddit.

**Live API & interactive docs:** [https://meme-api.aelx.de](https://meme-api.aelx.de)

## Endpoints

| Endpoint | Description |
|---|---|
| `GET /gimme` | Random meme from default subreddits |
| `GET /gimme/{count}` | N random memes (max 50) |
| `GET /gimme/{subreddit}` | Random meme from a specific subreddit |
| `GET /gimme/{subreddit}/{count}` | N random memes from a specific subreddit (max 50) |

Full request/response schemas and live try-it-out: [`/`](https://meme-api.aelx.de)

## Deployment

### Docker Compose

Create a `docker-compose.yml` and `.env` (see `.env.example`), then run:

```bash
docker compose up -d
```

**`docker-compose.yml`**

```yaml
services:
  meme-api:
    image: ghcr.io/tilalx/meme-api:rolling
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - REDDIT_CLIENT_ID=${REDDIT_CLIENT_ID}
      - REDDIT_CLIENT_SECRET=${REDDIT_CLIENT_SECRET}
      - REDISCLOUD_URL=redis://redis:6379
      - SENTRY_DSN=${SENTRY_DSN}
    depends_on:
      - redis

  redis:
    image: redis:8
    restart: unless-stopped
    volumes:
      - redis-data:/data

volumes:
  redis-data:
```

**`.env`**

```env
REDDIT_CLIENT_ID=your_client_id
REDDIT_CLIENT_SECRET=your_client_secret
SENTRY_DSN=           # optional
```

> Get Reddit API credentials at <https://www.reddit.com/prefs/apps> — create a **script** type app.
