package inbound

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/google/uuid"
    "github.com/walletera/eventskit/events"
    paymentsapi "github.com/walletera/payments-types/api"
)

type EventsDeserializer struct {
}

func NewEventsDeserializer() *EventsDeserializer {
    return &EventsDeserializer{}
}

func (d *EventsDeserializer) Deserialize(rawPayload []byte) (events.Event[EventsHandler], error) {
    var event events.EventEnvelope
    err := json.Unmarshal(rawPayload, &event)
    if err != nil {
        return nil, fmt.Errorf("error deserializing message with payload %s: %w", rawPayload, err)
    }
    switch event.Type {
    case "InboundPaymentReceived":
        var payment paymentsapi.Payment
        err := json.Unmarshal(event.Data, &payment)
        if err != nil {
            log.Printf("error deserializing InboundPaymentReceived event data %s: %s", event.Data, err.Error())
        }
        return NewPaymentReceived(uuid.MustParse(event.CorrelationID), payment), nil
    default:
        return nil, fmt.Errorf("unexpected event type: %s", event.Type)
    }
}
