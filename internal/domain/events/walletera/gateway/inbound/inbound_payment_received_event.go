package inbound

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/google/uuid"
    "github.com/walletera/eventskit/events"
    "github.com/walletera/payments-types/api"
    "github.com/walletera/werrors"
)

var _ events.Event[EventsHandler] = &PaymentReceived{}

type PaymentReceived struct {
    correlationId uuid.UUID
    api.Payment
}

func NewPaymentReceived(correlationId uuid.UUID, payment api.Payment) *PaymentReceived {
    return &PaymentReceived{correlationId: correlationId, Payment: payment}
}

func (p *PaymentReceived) ID() string {
    return p.Payment.ID.String()
}

func (p *PaymentReceived) Type() string {
    return "InboundPaymentReceived"
}

func (p *PaymentReceived) DataContentType() string {
    return "application/json"
}

func (p *PaymentReceived) CorrelationID() string {
    return p.correlationId.String()
}

func (p *PaymentReceived) Accept(ctx context.Context, handler EventsHandler) werrors.WError {
    return handler.HandleInboundPaymentReceived(ctx, p)
}

func (p *PaymentReceived) Serialize() ([]byte, error) {
    serializedPayment, err := json.Marshal(p.Payment)
    if err != nil {
        return nil, fmt.Errorf("failed serializing OutbounPaymentCreated event: %w", err)
    }
    envelope := events.EventEnvelope{
        Type:          p.Type(),
        CorrelationID: p.correlationId.String(),
        Data:          serializedPayment,
    }
    return json.Marshal(envelope)
}
