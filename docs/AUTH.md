# Authentication

The observability platform supports optional API key-based authentication for securing telemetry endpoints.

## Configuration

Authentication is controlled via environment variables:

```bash
# Enable authentication
export OBS_PLATFORM_AUTH_ENABLED=true

# Set API keys (comma-separated for multiple keys)
export OBS_PLATFORM_AUTH_API_KEYS="<your-key-1>,<your-key-2>,<your-key-3>"
```

## Behavior

### When Authentication is Disabled (Default)

- All endpoints accept requests without authentication
- Suitable for local development and testing
- Default state: `OBS_PLATFORM_AUTH_ENABLED=false`

### When Authentication is Enabled

- All `/api/v1/*` endpoints require valid API key in `X-API-Key` header
- `/api/health` endpoint remains unauthenticated for monitoring
- Returns `401 Unauthorized` for missing or invalid API keys

## Usage

### Client Requests (with auth enabled)

```bash
# Send traces with API key
curl -X POST http://localhost:8080/api/v1/traces \
  -H "X-API-Key: <your-api-key>" \
  -H "Content-Type: application/json" \
  -d '[{"trace_id": "abc123", ...}]'

# Health check (no auth required)
curl http://localhost:8080/api/health
```

### Client SDK Configuration

When using the observability client SDK, set the API key in requests:

```go
client := observability.NewClient("http://localhost:8080")
```

If authentication is enabled on the server, the client will need to be modified
to include the X-API-Key header in all requests to `/api/v1/* endpoints`

## Testing

Tests use auth-disabled configuration by default for simplicity:

```go
func createTestService(t *testing.T) *Service {
    cfg := config.ServiceSettings{
        Auth: config.AuthSettings{
            Enabled: false, // Auth disabled for tests
        },
    }
}
```

For testing with authentication enabled, use the helper:

```go
func TestWithAuth(t *testing.T) {
    service := createTestServiceWithAuth(t, []string{"EXAMPLE_KEY"})

    // Make request with API key header
    req.Header.Set("X-API-Key", "EXAMPLE_KEY")
    // ...
}
```

## Security Considerations

1. **API Key Storage**: Store API keys securely (e.g., AWS Secrets Manager, Vault)
2. **Key Rotation**: Implement regular key rotation policies
3. **Multiple Keys**: Use the comma-separated list to enable zero-downtime key rotation
4. **HTTPS**: Always use HTTPS in production to encrypt API keys in transit
5. **Key Length**: Use sufficiently long, random keys (minimum 32 characters recommended)

## Example Production Deployment

```bash
# Generate secure API key
API_KEY=$(openssl rand -base64 32)

# Set environment variables
export OBS_PLATFORM_AUTH_ENABLED=true
export OBS_PLATFORM_AUTH_API_KEYS="$API_KEY"
export OBS_PLATFORM_HOST=0.0.0.0
export OBS_PLATFORM_PORT=8080

# Start server
./observability-platform
```

## Key Rotation Strategy

To rotate keys without downtime:

1. Add new key to the comma-separated list:

   ```bash
   export OBS_PLATFORM_AUTH_API_KEYS="<old-key>,<new-key>"
   ```

2. Update clients to use new key

3. Remove old key after migration complete:

   ```bash
   export OBS_PLATFORM_AUTH_API_KEYS="<new-key>"
   ```
