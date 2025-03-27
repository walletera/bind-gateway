package inbound

import (
    "context"
    "encoding/json"
    "fmt"
    "strconv"

    "github.com/google/uuid"
    "github.com/walletera/eventskit/events"
    "github.com/walletera/werrors"
)

const PaymentReceivedType = "InboundPaymentReceived"

var _ events.Event[EventsHandler] = &PaymentReceived{}

type PaymentReceived struct {
    CorrelationId      uuid.UUID       `json:"correlation_id"`
    OriginCreditCuit   int64           `json:"origin_credit_cuit"`
    OriginCreditCvu    string          `json:"origin_credit_cvu"`
    OriginDebitCuit    int64           `json:"origin_debit_cuit"`
    OriginDebitCvu     string          `json:"origin_debit_cvu"`
    ChargeValueAmount  float64         `json:"charge_value_amount"`
    Currency           string          `json:"currency"`
    OriginId           int64           `json:"origin_id"`
    CoelsaId           string          `json:"coelsa_id"`
    RawInboundTransfer json.RawMessage `json:"raw_inbound_transfer"`
}

func (p PaymentReceived) ID() string {
    return strconv.FormatInt(p.OriginId, 10)
}

func (p PaymentReceived) Type() string {
    return PaymentReceivedType
}

func (p PaymentReceived) DataContentType() string {
    return "application/json"
}

func (p PaymentReceived) CorrelationID() string {
    return p.CorrelationId.String()
}

func (p PaymentReceived) Accept(ctx context.Context, handler EventsHandler) werrors.WError {
    return handler.HandleInboundPaymentReceived(ctx, p)
}

func (p PaymentReceived) Serialize() ([]byte, error) {
    serializedPayment, err := json.Marshal(p)
    if err != nil {
        return nil, fmt.Errorf("failed serializing OutbounPaymentCreated event: %w", err)
    }
    envelope := events.EventEnvelope{
        Type:          p.Type(),
        CorrelationID: p.CorrelationId.String(),
        Data:          serializedPayment,
    }
    return json.Marshal(envelope)
}
