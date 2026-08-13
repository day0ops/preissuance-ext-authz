// Package healthz provides the gRPC health checking service used by
// orchestrators (e.g. Kubernetes gRPC liveness/readiness probes) to observe
// whether the ext_authz server is ready to serve traffic.
package healthz

import "google.golang.org/grpc/health"

// New returns a standard gRPC health server. Callers register it with
// grpc_health_v1.RegisterHealthServer and mark it serving via
// SetServingStatus once the service is ready to accept traffic.
func New() *health.Server {
	return health.NewServer()
}
