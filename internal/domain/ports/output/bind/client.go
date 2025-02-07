package bind

import (
    "context"

    "github.com/walletera/bind-gateway/pkg/bind/api"
)

type Client interface {
    CreateCVUTrasfer(ctx context.Context, req *api.CreateTransferRequest) (api.CreateTransferResponse, error)
}
