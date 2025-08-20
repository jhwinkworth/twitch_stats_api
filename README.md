# Twitch Stats API

This service provides Twitch video statistics through a simple REST API.

## Endpoints

### Get Video Stats
GET /streamers/{channel_id}/videos?n={n}

- **Path Parameter**
  - `channel_id`: Twitch channel ID.
- **Query Parameter**
  - `n`: Number of videos to fetch.

**Example:**
GET http://localhost:8080/streamers/12345/videos?n=5

---

## Running with Docker

### 1. Prerequisites
- [Docker](https://docs.docker.com/get-docker/) installed on your system.
- A valid Twitch **Client ID** and **Client Secret**.

### 2. Prepare Environment Variables
Create a `.env` file in the project root:

```env
PORT=8080
TWITCH_CLIENT_ID=your-twitch-client-id
TWITCH_CLIENT_SECRET=your-twitch-client-secret
TWITCH_CHANNEL_ID=12826
```
- 12826 is the channel ID of the official Twitch channel

### 3. Build the Docker image
```
docker build -t twitch-stats-api .
```

### 4. Run the container
```
docker run --env-file .env -p 8080:8080 twitch-stats-api
```

### 5. Test the API
curl "http://localhost:8080/streamers/12826/videos?n=5"
- --env-file .env injects your environment variables at runtime.
- -p 8080:8080 maps the container port to your host.

Expected response
```
{
  "total_views": 587021,
  "average_views": 117404.2,
  "total_duration_minutes": 452.783333333333,
  "views_per_minute": 1296.47219052527,
  "most_viewed_title": "Twitch Public Access (August 1, 2025) | w/ @merrykish @snackless @unsanitylive @ajlive3",
  "most_viewed_view_count": 287856
}
```

## Development notes
Tests run during Docker build:
```
RUN go test -v -tags=integration ./...
```
The build will fail if tests fail.
