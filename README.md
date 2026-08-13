# preissuance-ext-authz

An agentgateway `ext_authz` gRPC service that gates the enterprise-agentgateway controller's built-in OAuth issuer _before_ it mints a token - the controller's `KGW_OAUTH_ISSUER_CONFIG.pre_issuance` hook (see [agentgateway-field-kit](https://github.com/day0ops/agentgateway-field-kit)'s `auth-only-mcp` feature).

For every pre-issuance check it:

1. Reads the caller's principal from `CheckRequest.Attributes.Source.Principal`.
2. Allows the request if that principal is in the configured allowlist, denying (so the controller never issues a token) otherwise.

## Configuration

All configuration is via environment variables:

| Variable             | Required | Default | Description                                                               |
| -------------------- | -------- | ------- | ------------------------------------------------------------------------- |
| `PORT`               | no       | `9002`  | gRPC listen port for the ext_authz service.                               |
| `ALLOWED_PRINCIPALS` | no       | (empty) | Comma-separated list of principals allowed through the pre-issuance gate. |

An empty or missing `ALLOWED_PRINCIPALS` denies every request.

## Running locally

```bash
export ALLOWED_PRINCIPALS=user1,user2

go run .
```

The service listens on `:9002` (gRPC) by default and exposes a standard gRPC health check (`grpc.health.v1.Health`).

## Development

```bash
make test   # go test ./...
make lint   # gofmt -l . && go vet ./...
make clean  # remove bin/
```

## Wiring into agentgateway

Reference this service's Service/port from the controller's `pre_issuance` config:

```json
{
  "pre_issuance": {
    "ext_authz": {
      "grpc_service": {
        "target_uri": "auth-only-mcp-preissuance-authz.agentgateway-system.svc.cluster.local:9002"
      }
    },
    "denied_redirect": "https://example.com/access-denied"
  }
}
```

## License

Apache License 2.0 - see [LICENSE](LICENSE).
