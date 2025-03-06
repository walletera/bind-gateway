package bind

import (
    "encoding/json"
    "fmt"

    "github.com/walletera/eventskit/events"
)

type WebhookEvent struct {
    Id           string          `json:"id"`
    Object       string          `json:"object"`
    Created      string          `json:"created"`
    Data         json.RawMessage `json:"data"`
    Type         string          `json:"type"`
    Redeliveries int             `json:"redeliveries"`
}

type EventsDeserializer struct{}

func NewEventsDeserializer() *EventsDeserializer {
    return &EventsDeserializer{}
}

func (e EventsDeserializer) Deserialize(rawEvent []byte) (events.Event[EventsHandler], error) {
    var webhookEvent WebhookEvent
    err := json.Unmarshal(rawEvent, &webhookEvent)
    if err != nil {
        return nil, fmt.Errorf("failed unmarshalling bind webhook event: %w", err)
    }
    switch webhookEvent.Type {
    case "transfer.cvu.received":
        var inboundTransferCreated InboundTransferCreated
        err = json.Unmarshal(webhookEvent.Data, &inboundTransferCreated)
        if err != nil {
            return nil, fmt.Errorf("failed unmarshalling bind webhook event of type transfer.cvu.received: %w", err)
        }
        return inboundTransferCreated, nil
    default:
        return nil, fmt.Errorf("unexpected webhook event")
    }
}
