# Docker Deployment Guide for decodeBot

## 🐳 Building the Docker Image

### Local Build
```bash
cd decodeBot
docker build -t decodebot:latest .
```

### Build with a specific tag
```bash
docker build -t decodebot:v1.0.0 .
```

## 🚀 Running the Container

### Basic Run
```bash
docker run --env-file .env decodebot:latest
```

### Run with environment variables
```bash
docker run \
  -e BOT_TOKEN=your_bot_token \
  -e SERVER_URL=http://your-server:8081 \
  -e BOT_SECRET=your_bot_secret \
  -e ADMIN_ID=your_telegram_id \
  decodebot:latest
```

### Run in detached mode with restart policy
```bash
docker run -d \
  --name decodebot \
  --restart unless-stopped \
  --env-file .env \
  decodebot:latest
```

### Run with custom timezone
```bash
docker run \
  --env-file .env \
  -e TZ=America/New_York \
  decodebot:latest
```

## 🔍 Managing the Container

### View logs
```bash
docker logs decodebot
```

### Follow logs in real-time
```bash
docker logs -f decodebot
```

### Stop the container
```bash
docker stop decodebot
```

### Start the container
```bash
docker start decodebot
```

### Restart the container
```bash
docker restart decodebot
```

### Remove the container
```bash
docker rm decodebot
```

## 🌐 Docker Compose (Optional)

Create a `docker-compose.yml` in the project root:

```yaml
version: '3.8'

services:
  bot:
    build:
      context: ./decodeBot
      dockerfile: Dockerfile
    container_name: decodebot
    restart: unless-stopped
    env_file:
      - ./decodeBot/.env
    environment:
      - TZ=Europe/Warsaw
    depends_on:
      - server
    networks:
      - decode-network

  server:
    build:
      context: ./decodeServer
      dockerfile: Dockerfile
    container_name: decodeserver
    restart: unless-stopped
    env_file:
      - ./decodeServer/.env
    ports:
      - "8081:8081"
    networks:
      - decode-network

networks:
  decode-network:
    driver: bridge
```

### Run with Docker Compose
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f bot

# Stop all services
docker-compose down

# Rebuild and restart
docker-compose up -d --build
```

## 📦 GitHub Container Registry

The GitHub Actions workflow automatically builds and publishes the Docker image to GitHub Container Registry (ghcr.io) on every push to the `main` branch.

### Pull from GitHub Container Registry
```bash
docker pull ghcr.io/ziks29/ciphercore-bot:main
```

### Run from GHCR
```bash
docker run --env-file .env ghcr.io/ziks29/ciphercore-bot:main
```

## 🔧 Environment Variables

Required environment variables (create a `.env` file):

```env
BOT_TOKEN=your_telegram_bot_token
SERVER_URL=http://decodeserver:8081
BOT_SECRET=your_shared_secret
ADMIN_ID=your_telegram_user_id
DEBUG=false
```

## 🏗️ Multi-Architecture Builds

To build for multiple architectures (useful for deployment on different platforms):

```bash
docker buildx create --use
docker buildx build --platform linux/amd64,linux/arm64 -t decodebot:latest .
```

## 📊 Image Size

The final Docker image is approximately **20MB** thanks to:
- Multi-stage builds
- Alpine Linux base image
- Static binary compilation

## 🔐 Security Best Practices

1. **Never commit `.env` files** - Use `.env.example` as a template
2. **Use secrets management** - For production, use Docker secrets or environment variable injection
3. **Keep the base image updated** - Regularly rebuild to get security updates
4. **Run as non-root** - Consider adding a non-root user in the Dockerfile for production

## 🐛 Troubleshooting

### Container exits immediately
```bash
# Check logs for errors
docker logs decodebot

# Run interactively to see output
docker run --rm -it --env-file .env decodebot:latest
```

### Cannot connect to server
- Ensure `SERVER_URL` is correct and accessible from the container
- If using Docker Compose, use the service name (e.g., `http://server:8081`)
- Check network connectivity between containers

### Timezone issues with scheduler
- Verify `TZ` environment variable is set correctly
- Check that `tzdata` package is installed in the container (already included)
