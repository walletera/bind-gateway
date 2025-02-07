package bind

import (
    "context"
    "fmt"

    "github.com/walletera/bind-gateway/pkg/bind/api"
)

type Client struct {
    client *api.Client
}

func NewClient(url string) (*Client, error) {
    client, err := api.NewClient(url, NewSecuritySource())
    if err != nil {
        return nil, fmt.Errorf("failed creating dinopay api client: %w", err)
    }
    return &Client{
        client: client,
    }, nil
}

func (c *Client) CreateCVUTrasfer(ctx context.Context, req *api.CreateTransferRequest) (api.CreateTransferRes, error) {
    return c.client.CreateTransfer(ctx, req)
}
