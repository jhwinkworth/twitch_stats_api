# Twitch Stats API

A simple REST API service that fetches and aggregates Twitch video statistics (views, duration, top video, etc.) for a given Twitch channel.

---

## Table of Contents

1. [Overview](#overview)  
2. [Features](#features)  
3. [Getting Started](#getting-started)  
   - [Prerequisites](#prerequisites)  
   - [Environment Variables](#environment-variables)  
   - [Running with Docker](#running-with-docker)  
4. [Usage](#usage)  
   - [API Endpoint](#api-endpoint)  
   - [Example Request](#example-request)  
   - [Example Response](#example-response)  
5. [Testing](#testing) 
6. [Development Notes](#development-notes)  
7. [Roadmap](#roadmap)  

---

## Overview

This service provides aggregated statistics for Twitch channel videos via a RESTful API.  
It can be used to:

- Quickly fetch total view count, average views, and total duration of a channel’s recent videos  
- Identify the most viewed video  
- Build dashboards, reports, or analytics tools around Twitch content  

---

## Features

- Fetch statistics for the **last _n_ videos** of a channel  
- Aggregate data: total views, average views, total duration (in minutes), views per minute  
- Identify most viewed video and its title  
- Dockerized for easy deployment  
- Integration with Twitch API using Client ID / Secret

---

## Getting Started

### Prerequisites

- Docker (for the Docker-based setup)  
- Twitch Developer credentials:  
  - `TWITCH_CLIENT_ID`  
  - `TWITCH_CLIENT_SECRET`  

### Environment Variables

Create a `.env` file in the root of the project with:

```env
PORT=8080
TWITCH_CLIENT_ID=your-twitch-client-id
TWITCH_CLIENT_SECRET=your-twitch-client-secret
TWITCH_CHANNEL_ID=12826
```

- `TWITCH_CHANNEL_ID`: the numeric Twitch channel ID (e.g. “12826” for a specific channel)  
- `PORT`: port for the API server  

---

### Running with Docker

1. Build the Docker image:  
   ```bash
   docker build -t twitch-stats-api .
   
2. Run the container:
   ```bash
   docker run --env-file .env -p 8080:8080 twitch-stats-api

3. Test the API
   ```bash
   curl "http://localhost:8080/streamers/12826/videos?n=5"

## Usage
### API Endpoint
```bash
GET /streamers/{channel_id}/videos?n={n}
```

Path parameter: channel_id — Twitch channel numeric ID
Query parameter: n — number of recent videos to fetch

### Example Request
```bash
curl "http://localhost:8080/streamers/12826/videos?n=5"
```

### Example Response
```bash
{
  "total_views": 587021,
  "average_views": 117404.2,
  "total_duration_minutes": 452.783333333333,
  "views_per_minute": 1296.47219052527,
  "most_viewed_title": "Twitch Public Access (August 1, 2025) | w/ @merrykish @snackless @unsanitylive @ajlive3",
  "most_viewed_view_count": 287856
}
```

## Testing

This project includes integration tests that run during the Docker build:
```bash
RUN go test -v -tags=integration ./...
```

If any tests fail, the Docker build will also fail — ensuring that builds are always tested.

## Development Notes

The build contains go test for the integration suite

Structure:
- cmd/app: main application entry point
- internal/: core business logic
- Env vars are required for authentication with Twitch API

## Roadmap

Some ideas for future improvements:
- Support pagination or cursor-based fetches
- Add caching of video stats to reduce Twitch API calls
- Add more endpoints (e.g. for live streams, followers, clips)
- Add OpenAPI / Swagger documentation
- Add user authentication (if making this a “client” service)
- Add detailed metrics (e.g. views per day, growth rates)
