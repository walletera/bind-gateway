package payments

import (
    "context"
    "log/slog"
    "strconv"

    "github.com/google/uuid"
    "github.com/walletera/bind-gateway/internal/adapters/bind"
    "github.com/walletera/bind-gateway/internal/domain/events/walletera/gateway/outbound"
    "github.com/walletera/bind-gateway/pkg/bind/api"
    "github.com/walletera/bind-gateway/pkg/logattr"
    "github.com/walletera/eventskit/eventsourcing"
    paymentEvents "github.com/walletera/payments-types/events"
    "github.com/walletera/werrors"
)

type EventsHandler struct {
    bindClient *bind.Client
    esDB       eventsourcing.DB
    logger     *slog.Logger
}

var _ paymentEvents.Handler = (*EventsHandler)(nil)

func NewEventsHandler(bindClient *bind.Client, esDB eventsourcing.DB, logger *slog.Logger, ) *EventsHandler {
    return &EventsHandler{
        bindClient: bindClient,
        esDB:       esDB,
        logger:     logger.With(logattr.Component("payments.EventsHandler")),
    }
}

func (ev *EventsHandler) HandlePaymentCreated(ctx context.Context, paymentCreated paymentEvents.PaymentCreated) werrors.WError {
    walleteraPaymentId := paymentCreated.Data.ID
    logger := ev.logger.With(
        logattr.EventType(paymentCreated.Type()),
        logattr.PaymentId(walleteraPaymentId.String()),
    )

    bindResp, err := ev.bindClient.CreateCVUTrasfer(ctx,
        &api.CreateTransferRequest{
            CvuOrigen:     paymentCreated.Data.Debtor.AccountDetails.OneOf.CvuAccountDetails.RoutingInfo.OneOf.CvuCvuRoutingInfo.Cvu,
            CbuCvuDestino: api.NewOptString(paymentCreated.Data.Beneficiary.AccountDetails.OneOf.CvuAccountDetails.RoutingInfo.OneOf.CvuCvuRoutingInfo.Cvu),
            CuitDestino:   api.OptString{},
            AliasDestino:  api.OptString{},
            Importe:       paymentCreated.Data.Amount,
            Referencia:    api.OptString{},
            Concepto:      api.OptString{},
            Emails:        nil,
            IdExterno:     api.NewOptString(paymentCreated.Data.ID.String()),
        },
    )
    if err != nil {
        werr := werrors.NewRetryableInternalError("failed creating payment on bind: %s", err.Error())
        logger.Error(werr.Error())
        return werr
    }
    if bindResp == nil {
        werr := werrors.NewRetryableInternalError("dinopay response is nil")
        logger.Error(werr.Error())
        return werr
    }
    var successfullResp *api.CreateTransferResponse
    switch resp := bindResp.(type) {
    case *api.CreateTransferResponse:
        successfullResp = resp
    case *api.CreateTransferBadRequest:
        werr := werrors.NewNonRetryableInternalError("failed creating transfer on bind %s:", resp.Detalle)
        logger.Error(werr.Error())
        return werr
    case *api.CreateTransferUnprocessableEntity:
        werr := werrors.NewNonRetryableInternalError("failed creating transfer on bind %s:", resp.Detalle)
        logger.Error(werr.Error())
        return werr
    case *api.CreateTransferUnauthorized:
    default:
        werr := werrors.NewNonRetryableInternalError("unexpected bind response type %t:", bindResp)
        logger.Error(werr.Error())
        return werr
    }

    logger.Info(
        "bind transfer created successfully",
        logattr.CorrelationId(paymentCreated.CorrelationID()),
        logattr.PaymentId(walleteraPaymentId.String()),
        logattr.BindOperationId(strconv.Itoa(successfullResp.OperacionId.Value)),
        logattr.BindStatus(int(successfullResp.EstadoId.Value)),
    )

    outboundPaymentCreated := outbound.PaymentCreated{
        Id:                uuid.New(),
        PaymentId:         walleteraPaymentId,
        BindPaymentId:     successfullResp.OperacionId.Value,
        BindPaymentStatus: successfullResp.EstadoId.Value,
    }

    streamName := outbound.BuildOutboundPaymentStreamName(walleteraPaymentId.String())

    werr := ev.esDB.AppendEvents(
        ctx,
        streamName,
        eventsourcing.ExpectedAggregateVersion{IsNew: true},
        outboundPaymentCreated,
    )
    if err != nil {
        werr := werrors.NewWrappedError(
            werr,
            "failed handling outbound PaymentCreated event",
            streamName,
        )
        logger.Error(werr.Error())
        return werr
    }

    logger.Info("PaymentCreated event processed successfully")

    return nil
}

func (ev *EventsHandler) HandlePaymentUpdated(_ context.Context, _ paymentEvents.PaymentUpdated) werrors.WError {
    // Ignore, nothing to do
    return nil
}
