// Package authz implements the Envoy/agentgateway ext_authz gRPC service
// (envoy.service.auth.v3.Authorization) that gates the controller's OAuth
// issuer before it mints a token, based on the caller's source principal.
package authz

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	"go.uber.org/zap"

	"github.com/day0ops/preissuance-ext-authz/pkg/config"
)

// Server implements authv3.AuthorizationServer, deciding whether the
// controller may continue issuing a token based on the incoming
// CheckRequest's source principal.
//
// NOTE: how CheckRequest.Attributes.Source.Principal gets populated for a
// browser OAuth login through the controller's pre_issuance hook - as
// opposed to the mTLS peer certificate this field normally carries - is
// unverified against a real controller as of this writing. Live-smoke-test
// this against a real deploy (check controller logs, and what actually ends
// up in Source.Principal for a Keycloak-federated login) before relying on
// it; see features/auth-only-mcp/index.js in agentgateway-field-kit.
type Server struct {
	authv3.UnimplementedAuthorizationServer

	cfg *config.Config
	log *zap.Logger
}

// NewServer returns a Server that authorizes requests against cfg's
// allowlist, logging via log.
func NewServer(cfg *config.Config, log *zap.Logger) *Server {
	return &Server{cfg: cfg, log: log}
}

// Check implements the ext_authz Authorization service. It allows the
// request only if the caller's source principal is on the configured
// allowlist.
func (s *Server) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	principal := req.GetAttributes().GetSource().GetPrincipal()
	if principal == "" || !s.cfg.AllowedPrincipals[principal] {
		s.log.Info("pre-issuance check denied", zap.String("principal", principal))
		return deniedResponse(), nil
	}

	s.log.Info("pre-issuance check allowed", zap.String("principal", principal))
	return okResponse(), nil
}

// deniedResponse builds a CheckResponse that denies the request.
func deniedResponse() *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: &rpcstatus.Status{Code: int32(codes.PermissionDenied)},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{Code: typev3.StatusCode_Forbidden},
			},
		},
	}
}

// okResponse builds a CheckResponse that allows the request.
func okResponse() *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: &rpcstatus.Status{Code: int32(codes.OK)},
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{},
		},
	}
}
