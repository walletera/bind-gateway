package inbound

import (
    "context"
    "log/slog"

    "github.com/walletera/bind-gateway/pkg/logattr"
    paymentsapi "github.com/walletera/payments-types/api"
    "github.com/walletera/werrors"
)

type EventsHandler interface {
    HandleInboundPaymentReceived(ctx context.Context, inboundPaymentReceived *PaymentReceived) werrors.WError
}

type EventsHandlerImpl struct {
    paymentsApiClient *paymentsapi.Client
    deserializer      *EventsDeserializer
    logger            *slog.Logger
}

func NewEventsHandlerImpl(client *paymentsapi.Client, logger *slog.Logger) *EventsHandlerImpl {
    return &EventsHandlerImpl{
        paymentsApiClient: client,
        deserializer:      NewEventsDeserializer(),
        logger:            logger.With(logattr.Component("gateway.inbound.EventsHandlerImpl")),
    }
}

func (ev *EventsHandlerImpl) HandleInboundPaymentReceived(ctx context.Context, inboundPaymentReceived *PaymentReceived) werrors.WError {
    resp, err := ev.paymentsApiClient.PostPayment(ctx, &inboundPaymentReceived.Payment, paymentsapi.PostPaymentParams{
        XWalleteraCorrelationID: paymentsapi.OptUUID{
            Value: inboundPaymentReceived.correlationId,
            Set:   true,
        },
    })
    if err != nil {
        ev.logger.Error("failed creating payment on payments api", logattr.Error(err.Error()))
        return werrors.NewRetryableInternalError(err.Error())
    }
    switch r := resp.(type) {
    case *paymentsapi.Payment:
        ev.logger.Info("gateway event InboundPaymentReceived processed successfully",
            logattr.PaymentId(r.ID.String()),
            logattr.EventType(inboundPaymentReceived.Type()),
        )
        return nil
    case *paymentsapi.PostPaymentInternalServerError:
        ev.logger.Error("failed creating payment on payments api",
            logattr.CorrelationId(inboundPaymentReceived.CorrelationID()),
            logattr.Error(r.ErrorMessage),
        )
        return werrors.NewRetryableInternalError(r.ErrorMessage)
    case *paymentsapi.PostPaymentBadRequest:
        ev.logger.Error("failed creating payment on payments api",
            logattr.CorrelationId(inboundPaymentReceived.CorrelationID()),
            logattr.Error(r.ErrorMessage))
        return werrors.NewNonRetryableInternalError(r.ErrorMessage)
    case *paymentsapi.PostPaymentConflict:
        ev.logger.Error("failed creating payment on payments api",
            logattr.CorrelationId(inboundPaymentReceived.CorrelationID()),
            logattr.Error(r.ErrorMessage))
        return werrors.NewNonRetryableInternalError(r.ErrorMessage)
    case *paymentsapi.PostPaymentUnauthorized:
        ev.logger.Error("failed creating payment on payments api: unauthorized")
        return werrors.NewNonRetryableInternalError("failed creating payment on payments api",
            logattr.CorrelationId(inboundPaymentReceived.CorrelationID()),
            logattr.Error("unauthorized"),
        )
    default:
        ev.logger.Error("unexpected error creating payment on payments api",
            logattr.CorrelationId(inboundPaymentReceived.CorrelationID()),
        )
        return werrors.NewRetryableInternalError("unexpected error creating payment on payments api")
    }
}
