package bind

import (
    "context"
    "fmt"
    "log/slog"

    "github.com/google/uuid"
    "github.com/walletera/bind-gateway/internal/domain/events/walletera/gateway/inbound"
    "github.com/walletera/bind-gateway/pkg/logattr"
    "github.com/walletera/bind-gateway/pkg/wuuid"
    "github.com/walletera/eventskit/eventsourcing"
    paymentsapi "github.com/walletera/payments-types/privateapi"
    "github.com/walletera/werrors"
)

const InboundPaymentStreamNamePrefix = "bindGateway-inboundPayment"

type EventsHandler interface {
    HandleInboundTransferReceived(ctx context.Context, event InboundTransferReceived) werrors.WError
}

type EventsHandlerImpl struct {
    db                eventsourcing.DB
    paymentsApiClient paymentsapi.Invoker
    logger            *slog.Logger
}

func NewEventsHandlerImpl(
    db eventsourcing.DB,
    paymentsApiClient paymentsapi.Invoker,
    logger *slog.Logger,
) *EventsHandlerImpl {
    return &EventsHandlerImpl{
        db:                db,
        paymentsApiClient: paymentsApiClient,
        logger:            logger.With(logattr.Component("bind.EventsHandler")),
    }
}

func (h *EventsHandlerImpl) HandleInboundTransferReceived(ctx context.Context, event InboundTransferReceived) werrors.WError {
    correlationUUID := wuuid.NewUUID()
    h.logger.Info("bind webhook event received",
        logattr.EventType(event.Type()),
        logattr.BindOperationId(fmt.Sprintf("%d", event.Details.OriginId)),
        logattr.CorrelationId(correlationUUID.String()),
    )

    inboundPaymentReceived := h.buildInboundPaymentReceivedEvent(correlationUUID, event)

    streamName := BuildInboundPaymentStreamName(inboundPaymentReceived.ID())
    werr := h.db.AppendEvents(ctx, streamName, eventsourcing.ExpectedAggregateVersion{IsNew: true}, inboundPaymentReceived)
    if werr != nil {
        h.logger.Error("error appending event to stream",
            logattr.CorrelationId(correlationUUID.String()),
            logattr.PaymentId(inboundPaymentReceived.ID()),
            logattr.EventType(event.Type()),
            logattr.StreamName(streamName),
            logattr.Error(werr.Error()),
        )
        return werrors.NewWrappedError(werr)
    }

    h.logger.Info("bind event InboundTransferReceived processed successfully",
        logattr.CorrelationId(correlationUUID.String()),
        logattr.PaymentId(inboundPaymentReceived.ID()),
        logattr.EventType(event.Type()),
        logattr.CorrelationId(correlationUUID.String()),
    )
    return nil
}

func (h *EventsHandlerImpl) buildInboundPaymentReceivedEvent(correlationUUID uuid.UUID, event InboundTransferReceived) inbound.PaymentReceived {
    inboundPaymentReceived := inbound.PaymentReceived{
        CorrelationId:      correlationUUID,
        OriginCreditCuit:   event.Details.OriginCredit.Cuit,
        OriginCreditCvu:    event.Details.OriginCredit.Cvu,
        OriginDebitCuit:    event.Details.OriginDebit.Cuit,
        OriginDebitCvu:     event.Details.OriginDebit.Cvu,
        ChargeValueAmount:  event.Charge.Value.Amount,
        Currency:           event.Charge.Value.Currency,
        OriginId:           event.Details.OriginId,
        CoelsaId:           event.Id,
        RawInboundTransfer: nil,
    }
    return inboundPaymentReceived
}

func BuildInboundPaymentStreamName(id string) string {
    return fmt.Sprintf("%s.%s", InboundPaymentStreamNamePrefix, id)
}
