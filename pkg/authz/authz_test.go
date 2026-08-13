package authz

import (
	"context"
	"testing"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"

	"github.com/day0ops/preissuance-ext-authz/pkg/config"
)

func request(principal string) *authv3.CheckRequest {
	return &authv3.CheckRequest{
		Attributes: &authv3.AttributeContext{
			Source: &authv3.AttributeContext_Peer{
				Principal: principal,
			},
		},
	}
}

func TestCheck(t *testing.T) {
	cfg := &config.Config{AllowedPrincipals: map[string]bool{"user1": true}}
	s := NewServer(cfg, zap.NewNop())

	cases := []struct {
		name      string
		principal string
		wantCode  codes.Code
	}{
		{"allowed principal", "user1", codes.OK},
		{"denied principal", "user2", codes.PermissionDenied},
		{"missing principal", "", codes.PermissionDenied},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := s.Check(context.Background(), request(tc.principal))
			if err != nil {
				t.Fatalf("Check() returned unexpected error: %v", err)
			}
			if got := codes.Code(resp.GetStatus().GetCode()); got != tc.wantCode {
				t.Errorf("Check(%q) status = %v, want %v", tc.principal, got, tc.wantCode)
			}
		})
	}
}
