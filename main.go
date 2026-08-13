// Command preissuance-ext-authz runs an Envoy/agentgateway ext_authz gRPC
// service that gates the enterprise-agentgateway controller's OAuth issuer
// before it mints a token (KGW_OAUTH_ISSUER_CONFIG.pre_issuance), based on
// the caller's source principal.
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/day0ops/preissuance-ext-authz/pkg/authz"
	"github.com/day0ops/preissuance-ext-authz/pkg/config"
	"github.com/day0ops/preissuance-ext-authz/pkg/healthz"
	"github.com/day0ops/preissuance-ext-authz/pkg/logger"
	"github.com/day0ops/preissuance-ext-authz/pkg/version"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "preissuance-ext-authz: %v\n", err)
		os.Exit(1)
	}
}

// run builds and serves the ext_authz gRPC server until the process
// receives SIGINT/SIGTERM, then shuts down gracefully.
func run() error {
	cfg := config.Load()

	log, err := logger.New(false)
	if err != nil {
		return fmt.Errorf("build logger: %w", err)
	}
	defer log.Sync() //nolint:errcheck

	authzServer := authz.NewServer(cfg, log)

	grpcServer := grpc.NewServer()
	authv3.RegisterAuthorizationServer(grpcServer, authzServer)

	healthServer := healthz.New()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	addr := ":" + cfg.Port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		log.Info("starting preissuance-ext-authz",
			zap.String("addr", addr),
			zap.Int("allowed_principals", len(cfg.AllowedPrincipals)),
			zap.String("version", version.String()),
		)
		if err := grpcServer.Serve(lis); err != nil {
			return fmt.Errorf("serve grpc: %w", err)
		}
		return nil
	})
	group.Go(func() error {
		<-groupCtx.Done()
		log.Info("shutting down preissuance-ext-authz")
		grpcServer.GracefulStop()
		return nil
	})

	if err := group.Wait(); err != nil {
		return fmt.Errorf("run server: %w", err)
	}
	return nil
}
