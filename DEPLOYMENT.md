# Deployment Guide - List Manager API

This guide covers deployment procedures for the List Manager API across different environments.

## Table of Contents

- [Local Development](#local-development)
- [Environment Variables](#environment-variables)
- [Docker Deployment](#docker-deployment)
- [Render Deployment](#render-deployment)
- [Health Check Verification](#health-check-verification)
- [Troubleshooting](#troubleshooting)

---

## Local Development

### Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- Make (optional, for automation)

### Step 1: Start MongoDB

```bash
docker-compose up -d mongodb
```

This starts MongoDB on port 27017 with persistent storage in a named volume.

### Step 2: Configure Environment

Create a `.env` file or export environment variables:

```bash
export MONGO_URI=mongodb://localhost:27017/
export MONGO_DB_NAME=listmanager
export PORT=8085
```

### Step 3: Run the Application

**Option A: Using Go run (development)**
```bash
go run cmd/api/main.go
```

**Option B: Using Make**
```bash
make run
```

**Option C: Build and run binary**
```bash
make build
./bin/list-manager-api
```

### Step 4: Verify Deployment

```bash
curl http://localhost:8085/healthz
```

Expected response:
```json
{
  "status": "up",
  "server": "up",
  "database": {
    "status": "connected"
  },
  "timestamp": "2026-02-13T10:30:00Z",
  "checks": {
    "mongodb": {
      "status": "passed"
    }
  }
}
```

### Step 5: Stop Services

```bash
# Stop API
Ctrl+C (if running with go run)

# Stop MongoDB
docker-compose down
```

---

## Environment Variables

### Required Variables

| Variable | Description | Example | Default |
|-----------|-------------|---------|---------|
| `MONGO_URI` | MongoDB connection string | `mongodb://user:pass@host:27017/` | - |
| `MONGO_DB_NAME` | Database name | `listmanager` | - |
| `PORT` | HTTP server port | `8085` | `8080` |

### Optional Variables

| Variable | Description | Example | Default |
|-----------|-------------|---------|---------|
| `ENVIRONMENT` | Environment name | `production` | `development` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` | `info` |
| `CORS_ORIGINS` | Allowed CORS origins | `https://example.com` | `*` |

### Configuration Priority

1. Environment variables (highest priority)
2. `.env` file (development only)
3. Default values (lowest priority)

---

## Docker Deployment

### Build Image

```bash
docker build -t list-manager-api:latest .
```

### Run Container

```bash
docker run -d \
  --name list-manager-api \
  -p 8085:8085 \
  -e MONGO_URI=mongodb://host.docker.internal:27017/ \
  -e MONGO_DB_NAME=listmanager \
  -e PORT=8085 \
  list-manager-api:latest
```

### Docker Compose (Full Stack)

```bash
docker-compose up -d
```

This starts:
- API on `http://localhost:8085`
- MongoDB on `localhost:27017`
- Mongo Express UI on `http://localhost:8081`

### View Logs

```bash
docker-compose logs -f api
```

### Stop All Services

```bash
docker-compose down
```

---

## Render Deployment

The List Manager API is configured for deployment on [Render](https://render.com).

### Prerequisites

- GitHub repository connected to Render
- Render account with free/starter tier
- Render MongoDB instance created

### Step 1: Create MongoDB Instance

1. In Render Dashboard, go to **New → MongoDB**
2. Configure:
   - Name: `list-manager-db`
   - Database: `listmanager`
   - User: `listmanager` (auto-generated password)
3. Deploy and wait for instance to be ready
4. Copy **Internal Database URL** from instance details

### Step 2: Create Web Service

1. In Render Dashboard, go to **New → Web Service**
2. Connect your GitHub repository
3. Configure:
   - Name: `list-manager-api`
   - Region: Oregon (or closest to users)
   - Branch: `main`
   - Runtime: `Go`
   - Build Command: `go build -o bin/list-manager-api cmd/api/main.go`
   - Start Command: `./bin/list-manager-api`

### Step 3: Configure Environment Variables

Add the following in Render Web Service → Environment:

| Key | Value |
|-----|-------|
| `MONGO_URI` | (Your MongoDB Internal Database URL) |
| `MONGO_DB_NAME` | `listmanager` |
| `PORT` | `8085` (Render default, can be omitted) |

⚠️ **Important:** Enable **Sync to Render** for MONGO_URI to keep it in sync with your MongoDB instance.

### Step 4: Deploy

Click **Create Web Service**. Render will:
1. Build your application from source
2. Deploy to Render infrastructure
3. Provide a public URL (`https://list-manager-api.onrender.com`)

### Step 5: Verify Deployment

```bash
curl https://list-manager-api.onrender.com/healthz
```

### Step 6: Configure Custom Domain (Optional)

1. In Render Web Service, go to **Settings → Custom Domains**
2. Add your domain (e.g., `api.yourdomain.com`)
3. Update DNS records as instructed by Render
4. Wait for SSL certificate provisioning

### Auto-Deploy Settings

By default, Render auto-deploys on every push to `main`. To disable:

- Go to **Settings → Deploy Context**
- Uncheck **Auto-Deploy branch**

For manual deploys:

- Go to **Manual Deploy**
- Select branch and click **Deploy commit**

### Render Service Regions

| Region | Code | Location |
|--------|------|----------|
| Oregon | `oregon` | US West |
| Frankfurt | `frankfurt` | EU Central |
| Singapore | `singapore` | Asia Pacific |
| Ohio | `ohio` | US East |
| Virginia | `virginia` | US East |

---

## Health Check Verification

### Health Check Endpoint

**URL:** `GET /healthz`

**Success Response (200 OK):**
```json
{
  "status": "up",
  "server": "up",
  "database": {
    "status": "connected"
  },
  "timestamp": "2026-02-13T10:30:00Z",
  "checks": {
    "mongodb": {
      "status": "passed"
    }
  }
}
```

**Degraded Response (200 OK with degraded status):**
```json
{
  "status": "degraded",
  "server": "up",
  "database": {
    "status": "disconnected"
  },
  "timestamp": "2026-02-13T10:30:00Z",
  "checks": {
    "mongodb": {
      "status": "failed",
      "error": "connection timeout"
    }
  }
}
```

### Monitoring Health Checks

**Cron Job Monitoring:**

```bash
# Add to crontab for monitoring
*/5 * * * * curl -f http://localhost:8085/healthz || alert-admin
```

**Render Health Checks:**

Render automatically health checks your service every 30 seconds. Configure health check path in Render Web Service → Health Check.

---

## Troubleshooting

### Application Won't Start

**Problem:** `listen tcp :8085: bind: address already in use`

**Solutions:**
```bash
# Find process using the port
lsof -i :8085

# Kill the process
kill -9 <PID>

# Or use different port
export PORT=8086
```

### MongoDB Connection Refused

**Problem:** `connection refused` or `no reachable servers`

**Solutions:**

**For Docker Compose:**
```bash
# Verify MongoDB is running
docker-compose ps

# Check MongoDB logs
docker-compose logs mongodb

# Restart MongoDB
docker-compose restart mongodb
```

**For Render:**
- Verify MONGO_URI environment variable
- Check MongoDB instance is not suspended
- Ensure MongoDB and API are in same region

### CORS Errors

**Problem:** `No 'Access-Control-Allow-Origin' header`

**Solution:**
- For development, CORS is configured to allow all origins (`*`)
- For production, configure `CORS_ORIGINS` environment variable
- Check browser console for preflight request failures

### Build Failures

**Problem:** Build fails during `go build`

**Solutions:**
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Verify dependencies
go mod verify

# Try building locally first
make build
```

### High Memory Usage

**Problem:** Service runs out of memory on Render

**Solutions:**
- Upgrade to paid tier with more RAM
- Optimize MongoDB queries with indexes
- Reduce connection pool size in MongoDB client
- Add pagination to `GET /items` endpoint

### Logs Not Appearing

**Problem:** Can't see application logs

**Solutions:**

**For Docker Compose:**
```bash
docker-compose logs -f api
```

**For Render:**
- Go to Logs tab in Render Dashboard
- Select stream type: Server, Application
- Adjust time range if needed

---

## Monitoring and Observability

### Current State

- **Logging:** Zap structured logging to stdout
- **Health Checks:** `/healthz` endpoint with MongoDB verification

### Future Enhancements (Planned)

- **OpenTelemetry Integration:** Distributed tracing and metrics
- **External Monitoring:** Integration with Datadog, New Relic, or Grafana Cloud
- **Alerting:** PagerDuty or Slack notifications for critical failures
- **Log Aggregation:** Centralized logging with Loki, CloudWatch Logs, or ELK

### Viewing Logs

**Docker Compose:**
```bash
docker-compose logs -f --tail=100 api
```

**Render:**
- Logs available in Render Dashboard → Logs
- Real-time log streaming
- 7-day retention on free tier

---

## Security Checklist

Before deploying to production, ensure:

- [ ] MongoDB credentials are in environment variables (not in code)
- [ ] CORS origins are restricted to specific domains
- [ ] HTTPS is enabled (automatic on Render)
- [ ] Health check endpoint does not expose sensitive information
- [ ] Logs do not contain passwords or secrets
- [ ] Rate limiting is configured (if applicable)
- [ ] Input validation is comprehensive

---

## Rollback Procedure

### Render Automatic Rollback

If health checks fail for 5+ minutes, Render automatically rolls back to the last successful deployment.

### Manual Rollback

1. Go to **Deployments** in Render Dashboard
2. Find the last successful deployment
3. Click **Rollback to this deployment**

### Database Rollback

For destructive database changes:

1. Restore from MongoDB backup (if available)
2. Or implement database migrations for version-controlled schema changes

---

## Support

- **GitHub Issues:** [Report bugs or feature requests](https://github.com/lucaspereirasilva0/list-manager-api/issues)
- **Render Documentation:** [Official Render Docs](https://render.com/docs)
- **MongoDB Documentation:** [Official MongoDB Docs](https://www.mongodb.com/docs/)

---

**Last Updated:** 2026-02-13
