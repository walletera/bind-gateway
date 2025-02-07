package outbound

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/google/uuid"
    "github.com/walletera/bind-gateway/internal/domain/events/walletera/gateway"
    "github.com/walletera/bind-gateway/pkg/bind/api"
    "github.com/walletera/eventskit/events"
    "github.com/walletera/werrors"
)

var _ events.Event[EventsHandler] = PaymentCreated{}

type PaymentCreated struct {
    Id                uuid.UUID                          `json:"id,omitempty"`
    PaymentId         uuid.UUID                          `json:"withdrawal_id,omitempty"`
    BindPaymentId     int                                `json:"bind_payment_id"`
    BindPaymentStatus api.CreateTransferResponseEstadoId `json:"bind_payment_status"`
    CreatedAt         int64                              `json:"created_at,omitempty"`
}

func (o PaymentCreated) ID() string {
    return fmt.Sprintf("%s-%s", o.Type(), o.Id)
}

func (o PaymentCreated) Type() string {
    return "OutboundPaymentCreated"
}

func (o PaymentCreated) DataContentType() string {
    return "application/json"
}

func (o PaymentCreated) CorrelationID() string {
    panic("not implemented yet")
}

func (o PaymentCreated) Accept(ctx context.Context, handler EventsHandler) werrors.WError {
    return handler.HandleOutboundPaymentCreated(ctx, o)
}

func (o PaymentCreated) Serialize() ([]byte, error) {
    data, err := json.Marshal(o)
    if err != nil {
        return nil, fmt.Errorf("failed serializing OutbounPaymentCreated event: %w", err)
    }
    envelope := gateway.EventEnvelope{
        Type: "OutboundPaymentCreated",
        Data: data,
    }
    return json.Marshal(envelope)
}
