# Golang Final Test - Post Management API

A RESTful API service built with Go, Gin, PostgreSQL, Redis, and Elasticsearch for managing blog posts with full-text search capabilities.

## Features

- Create, read, update posts
- Search posts by tags
- Full-text search with Elasticsearch
- Redis caching
- PostgreSQL database with UUID primary keys
- Docker containerization
- Database migrations

## Tech Stack

- **Backend**: Go 1.21+, Gin Framework
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **Search**: Elasticsearch 8.15.3
- **Containerization**: Docker & Docker Compose

## How to Run This Project

### Prerequisites

- Docker
- Docker Compose

### Running with Docker Compose

1. **Clone the repository** (if not already done)
   ```bash
   git clone <repository-url>
   cd golang-final-test
   ```

2. **Start all services**
   ```bash
   docker-compose up -d
   ```

   This will start:
   - PostgreSQL database on port 5433
   - Redis on port 6379
   - Elasticsearch on port 9200
   - Go API application on port 8080

3. **Check service health**
   ```bash
   docker-compose ps
   ```

4. **View logs** (optional)
   ```bash
   docker-compose logs -f app
   ```

5. **Stop services**
   ```bash
   docker-compose down
   ```

### Environment Variables

The following environment variables are configured in `docker-compose.yml`:

- `DB_HOST=postgres`
- `DB_PORT=5432`
- `DB_USER=appuser`
- `DB_PASSWORD=apppassword`
- `DB_NAME=appdb`
- `DB_SSLMODE=disable`
- `REDIS_HOST=redis`
- `REDIS_PORT=6379`
- `ELASTICSEARCH_HOST=elasticsearch`
- `ELASTICSEARCH_PORT=9200`

## API Endpoints & Testing with cURL

The API server runs on `http://localhost:8080` and provides the following endpoints:

### 1. Create a Post

**Endpoint**: `POST /posts/`

```bash
curl -X POST http://localhost:8080/posts/ \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Blog Post",
    "content": "This is the content of my first blog post. It contains useful information about Go programming.",
    "tags": ["go", "programming", "tutorial"]
  }'
```

**Expected Response**:
```json
{
  "message": "Post created successfully"
}
```

### 2. Get Post by ID

**Endpoint**: `GET /posts/{id}`

First, you'll need a post ID. After creating a post, you can search for posts to get an ID, or check your database.

```bash
# Replace {post-id} with an actual UUID from your database
curl -X GET http://localhost:8080/posts/{post-id}
```

**Expected Response**:
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "title": "My First Blog Post",
  "content": "This is the content of my first blog post. It contains useful information about Go programming.",
  "tags": ["go", "programming", "tutorial"],
  "created_at": "2025-09-18T10:30:00Z"
}
```

### 3. Search Posts by Tag

**Endpoint**: `GET /posts/search-by-tag?tag={tag}`

```bash
curl -X GET "http://localhost:8080/posts/search-by-tag?tag=go"
```

**Expected Response**:
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "My First Blog Post",
    "content": "This is the content of my first blog post. It contains useful information about Go programming.",
    "tags": ["go", "programming", "tutorial"],
    "created_at": "2025-09-18T10:30:00Z"
  }
]
```

### 4. Update Post

**Endpoint**: `PUT /posts/`

```bash
# Replace the ID with an actual UUID from your database
curl -X PUT http://localhost:8080/posts/ \
  -H "Content-Type: application/json" \
  -d '{
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "title": "Updated Blog Post Title",
    "content": "This is the updated content with more detailed information about Go programming and best practices.",
    "tags": ["go", "programming", "tutorial", "advanced"]
  }'
```

**Expected Response**:
```json
{
  "message": "Post updated successfully"
}
```

### 5. Full-Text Search Posts

**Endpoint**: `GET /posts/search?q={query}&offset={offset}&limit={limit}`

```bash
# Basic search
curl -X GET "http://localhost:8080/posts/search?q=programming"

# Search with pagination
curl -X GET "http://localhost:8080/posts/search?q=programming&offset=0&limit=5"

# Search for multiple terms
curl -X GET "http://localhost:8080/posts/search?q=Go%20tutorial"
```

**Expected Response**:
```json
{
  "posts": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "title": "My First Blog Post",
      "content": "This is the content of my first blog post. It contains useful information about Go programming.",
      "tags": ["go", "programming", "tutorial"],
      "created_at": "2025-09-18T10:30:00Z"
    }
  ],
  "total": 1,
  "offset": 0,
  "limit": 10
}
```

## Testing Workflow Example

Here's a complete workflow to test the API:

```bash
# 1. Create a few posts
curl -X POST http://localhost:8080/posts/ \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Go Basics",
    "content": "Learn the fundamentals of Go programming language including variables, functions, and data structures.",
    "tags": ["go", "basics", "programming"]
  }'

curl -X POST http://localhost:8080/posts/ \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Advanced Go Patterns",
    "content": "Explore advanced design patterns in Go including interfaces, goroutines, and channels.",
    "tags": ["go", "advanced", "patterns", "concurrency"]
  }'

curl -X POST http://localhost:8080/posts/ \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Web Development with Gin",
    "content": "Build REST APIs using the Gin framework in Go with middleware and routing.",
    "tags": ["go", "gin", "web", "api"]
  }'

# 2. Search posts by tag
curl -X GET "http://localhost:8080/posts/search-by-tag?tag=go"

# 3. Full-text search
curl -X GET "http://localhost:8080/posts/search?q=goroutines"

# 4. Search with pagination
curl -X GET "http://localhost:8080/posts/search?q=go&offset=0&limit=2"
```

## Database Access

If you need to access the PostgreSQL database directly:

```bash
# Connect to PostgreSQL container
docker exec -it golang-final-test-postgres-1 psql -U appuser -d appdb

# View posts table
\dt
SELECT * FROM posts;
```

## Troubleshooting

### Common Issues

1. **Port already in use**: Make sure ports 8080, 5433, 6379, and 9200 are not being used by other services.

2. **Services not healthy**: Check service health with:
   ```bash
   docker-compose ps
   docker-compose logs [service-name]
   ```

3. **Database connection issues**: Ensure PostgreSQL is fully started before the app starts. The docker-compose file includes health checks and depends_on configurations.

### Logs

View logs for specific services:
```bash
# Application logs
docker-compose logs -f app

# Database logs
docker-compose logs -f postgres

# Redis logs
docker-compose logs -f redis

# Elasticsearch logs
docker-compose logs -f elasticsearch
```

## Development

### Local Development Setup

If you want to run the application locally (outside Docker):

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Set environment variables**:
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5433
   export DB_USER=appuser
   export DB_PASSWORD=apppassword
   export DB_NAME=appdb
   export DB_SSLMODE=disable
   export REDIS_HOST=localhost
   export REDIS_PORT=6379
   export ELASTICSEARCH_HOST=localhost
   export ELASTICSEARCH_PORT=9200
   ```

3. **Start external services**:
   ```bash
   docker-compose up -d postgres redis elasticsearch
   ```

4. **Run the application**:
   ```bash
   go run cmd/api/main.go
   ```

## API Documentation

### Request/Response Formats

All API endpoints expect and return JSON format. Make sure to include the `Content-Type: application/json` header in your requests.

### Error Responses

The API returns standard HTTP status codes and error messages in JSON format:

```json
{
  "error": "Error description here"
}
```

### Validation Rules

- **Post Title**: Required, non-empty string
- **Post Content**: Required, non-empty string  
- **Post Tags**: Required, array of strings
- **Post ID**: Must be a valid UUID format
- **Search Query**: Required for search endpoints
- **Pagination**: offset ≥ 0, limit between 1-100