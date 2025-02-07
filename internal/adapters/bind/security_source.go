package bind

import (
    "context"

    "github.com/walletera/bind-gateway/pkg/bind/api"
)

type SecuritySource struct {
}

func NewSecuritySource() *SecuritySource {
    return &SecuritySource{}
}

func (s *SecuritySource) BearerAuth(ctx context.Context, operationName api.OperationName) (api.BearerAuth, error) {
    //TODO implement me
    return api.BearerAuth{Token: "fakeToken"}, nil
}
