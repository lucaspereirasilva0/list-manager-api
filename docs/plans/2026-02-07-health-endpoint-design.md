# Health Endpoint Design Document

**Date**: 2026-02-07
**Author**: Auto-generated design
**Status**: Implemented

## Overview

This document describes the implementation of a health check endpoint (`/healthz`) for the list-manager-api. The endpoint allows automated deployment workflows and monitoring systems to verify the application status and connectivity to MongoDB.

## Requirements

- **URL Path**: `/healthz`
- **HTTP Method**: `GET`
- **Authentication**: Not required
- **Response Format**: JSON
- **HTTP Status**: Always `200 OK`
- **Verification**: Server HTTP status + MongoDB connectivity
- **MongoDB Ping Timeout**: 10 seconds
- **Middlewares**: CORS and logging applied (no authentication middleware)

## Architecture

### Components

1. **HealthHandler** (`cmd/api/handlers/health.go`)
   - Interface for health check operations
   - Implementation with MongoDB client and logger dependencies
   - No service layer dependency (infrastructure-only check)

2. **ClientWrapper.Ping** (`internal/database/mongodb/client.go`)
   - New method to verify MongoDB connectivity
   - Uses `mongo.Client.Ping(ctx, readpref.Primary())`
   - Returns error if connection fails

3. **Server Integration** (`cmd/api/server/server.go`)
   - Health handler registered as a field in Server struct
   - Route `/healthz` registered in `setupRoutes()`
   - No error handling middleware applied

4. **Initialization** (`cmd/api/main.go`)
   - Health handler created with MongoDB client and logger
   - Passed to `NewServer()` constructor

## Data Structures

### Response JSON

```json
{
  "status": "up|degraded|down",
  "server": "up",
  "database": "connected|disconnected",
  "timestamp": "2026-02-07T21:30:00Z",
  "checks": {
    "database": {
      "status": "passed|failed",
      "error": "error message if failed"
    }
  }
}
```

## Testing

### Unit Tests

**File**: `cmd/api/handlers/health_test.go`

Test Cases:
1. MongoDB connected → returns "up" status
2. MongoDB disconnected → returns "degraded" status
3. MongoDB timeout → returns "degraded" status with timeout error
4. HTTP status is always 200
5. JSON response is valid and complete
6. All response fields are present

### Success Criteria

1. ✅ Endpoint `/healthz` responds to GET requests
2. ✅ Verifies both server and MongoDB connectivity
3. ✅ 10-second timeout applied to MongoDB ping
4. ✅ Always returns HTTP 200 status code
5. ✅ Response JSON is valid with all required fields
6. ✅ CORS middleware is applied
7. ✅ Logging middleware is applied
8. ✅ No authentication required
9. ✅ All unit tests pass
