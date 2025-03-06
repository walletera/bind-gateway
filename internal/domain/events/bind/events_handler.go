package bind

import (
    "context"
    "fmt"
    "log/slog"
    "strconv"
    "time"

    "github.com/google/uuid"
    accountsapi "github.com/walletera/accounts/types/api/api"
    "github.com/walletera/bind-gateway/internal/domain/events/walletera/gateway/inbound"
    "github.com/walletera/bind-gateway/pkg/logattr"
    "github.com/walletera/bind-gateway/pkg/wuuid"
    "github.com/walletera/eventskit/eventsourcing"
    paymentsapi "github.com/walletera/payments-types/api"
    "github.com/walletera/werrors"
)

const InboundPaymentStreamNamePrefix = "bindGateway-inboundPayment"

type EventsHandler interface {
    HandleInboundTransferCreated(ctx context.Context, event InboundTransferCreated) werrors.WError
}

type EventsHandlerImpl struct {
    db                eventsourcing.DB
    paymentsApiClient paymentsapi.Invoker
    accountsApiClient accountsapi.Invoker
    logger            *slog.Logger
}

func NewEventsHandlerImpl(
    db eventsourcing.DB,
    paymentsApiClient paymentsapi.Invoker,
    accountsApiClient accountsapi.Invoker,
    logger *slog.Logger,
) *EventsHandlerImpl {
    return &EventsHandlerImpl{
        db:                db,
        paymentsApiClient: paymentsApiClient,
        accountsApiClient: accountsApiClient,
        logger:            logger.With(logattr.Component("bind.EventsHandler")),
    }
}

func (h *EventsHandlerImpl) HandleInboundTransferCreated(ctx context.Context, event InboundTransferCreated) werrors.WError {
    correlationUUID := wuuid.NewUUID()
    h.logger.Info("bind webhook event received",
        logattr.EventType(event.Type()),
        logattr.CorrelationId(correlationUUID.String()),
    )

    account, err := h.getBeneficiaryAccount(ctx, event)
    if err != nil {
        return err
    }

    inboundPaymentReceived := h.buildInboundPaymentReceivedEvent(correlationUUID, event, account)

    streamName := h.buildInboundPaymentStreamName(inboundPaymentReceived.ID())
    werr := h.db.AppendEvents(ctx, streamName, eventsourcing.ExpectedAggregateVersion{IsNew: true}, inboundPaymentReceived)
    if werr != nil {
        h.logger.Error("error appending event to stream",
            logattr.PaymentId(inboundPaymentReceived.ID()),
            logattr.EventType(event.Type()),
            logattr.StreamName(streamName),
            logattr.Error(werr.Error()),
        )
        return werrors.NewWrappedError(werr)
    }

    h.logger.Info("bind event InboundTransferCreated processed successfully",
        logattr.PaymentId(inboundPaymentReceived.ID()),
        logattr.EventType(event.Type()),
        logattr.CorrelationId(correlationUUID.String()),
    )
    return nil
}

func (h *EventsHandlerImpl) buildInboundPaymentReceivedEvent(correlationUUID uuid.UUID, event InboundTransferCreated, account accountsapi.Account) *inbound.PaymentReceived {
    paymentUUID := wuuid.NewUUID()
    inboundPaymentReceived := inbound.NewPaymentReceived(
        correlationUUID,
        paymentsapi.Payment{
            ID:       paymentUUID,
            Amount:   event.Charge.Value.Amount,
            Currency: event.Charge.Value.Currency,
            Direction: paymentsapi.OptPaymentDirection{
                Value: paymentsapi.PaymentDirectionInbound,
                Set:   true,
            },
            CustomerId: paymentsapi.OptUUID{
                Value: account.CustomerId,
                Set:   true,
            },
            ExternalId: paymentsapi.OptString{
                Value: strconv.FormatInt(event.Details.OriginId, 10),
                Set:   true,
            },
            SchemeId: paymentsapi.OptString{
                Value: event.Id,
                Set:   true,
            },
            Beneficiary: paymentsapi.OptAccountDetails{
                Value: paymentsapi.AccountDetails{
                    Currency: paymentsapi.OptString{
                        Value: event.Charge.Value.Currency,
                        Set:   true,
                    },
                    AccountType: paymentsapi.OptAccountDetailsAccountType{
                        Value: paymentsapi.AccountDetailsAccountTypeCvu,
                        Set:   true,
                    },
                    AccountDetails: paymentsapi.OptAccountDetailsAccountDetails{
                        Value: paymentsapi.AccountDetailsAccountDetails{
                            OneOf: paymentsapi.AccountDetailsAccountDetailsSum{
                                Type: paymentsapi.CvuAccountDetailsAccountDetailsAccountDetailsSum,
                                CvuAccountDetails: paymentsapi.CvuAccountDetails{
                                    Cuit: paymentsapi.OptString{
                                        Value: strconv.FormatInt(event.Details.OriginCredit.Cuit, 10),
                                        Set:   true,
                                    },
                                    Cvu: paymentsapi.OptString{
                                        Value: event.Details.OriginCredit.Cvu,
                                        Set:   true,
                                    },
                                },
                            },
                        },
                        Set: true,
                    },
                },
                Set: true,
            },
            Debtor: paymentsapi.OptAccountDetails{
                Value: paymentsapi.AccountDetails{
                    Currency: paymentsapi.OptString{
                        Value: event.Charge.Value.Currency,
                        Set:   true,
                    },
                    AccountType: paymentsapi.OptAccountDetailsAccountType{
                        Value: paymentsapi.AccountDetailsAccountTypeCvu,
                        Set:   true,
                    },
                    AccountDetails: paymentsapi.OptAccountDetailsAccountDetails{
                        Value: paymentsapi.AccountDetailsAccountDetails{
                            OneOf: paymentsapi.AccountDetailsAccountDetailsSum{
                                Type: paymentsapi.CvuAccountDetailsAccountDetailsAccountDetailsSum,
                                CvuAccountDetails: paymentsapi.CvuAccountDetails{
                                    Cuit: paymentsapi.OptString{
                                        Value: strconv.FormatInt(event.Details.OriginDebit.Cuit, 10),
                                        Set:   true,
                                    },
                                    Cvu: paymentsapi.OptString{
                                        Value: event.Details.OriginDebit.Cvu,
                                        Set:   true,
                                    },
                                },
                            },
                        },
                        Set: true,
                    },
                },
                Set: true,
            },
            Status: paymentsapi.OptPaymentStatus{
                Value: paymentsapi.PaymentStatusConfirmed,
                Set:   true,
            },
            CreatedAt: paymentsapi.OptDateTime{
                Value: time.Now(),
                Set:   true,
            },
            UpdatedAt: paymentsapi.OptDateTime{
                Value: time.Now(),
                Set:   true,
            },
        },
    )
    return inboundPaymentReceived
}

func (h *EventsHandlerImpl) getBeneficiaryAccount(ctx context.Context, event InboundTransferCreated) (accountsapi.Account, werrors.WError) {
    getAccountResp, err := h.accountsApiClient.GetAccount(
        ctx,
        accountsapi.GetAccountParams{
            AccountType: accountsapi.OptGetAccountAccountType{
                Value: accountsapi.GetAccountAccountTypeCvu,
                Set:   true,
            },
            CvuAccountDetails: accountsapi.OptCvuAccountDetails{
                Value: accountsapi.CvuAccountDetails{
                    Cuit: accountsapi.OptString{
                        Value: strconv.FormatInt(event.Details.OriginCredit.Cuit, 10),
                        Set:   true,
                    },
                    Cvu: accountsapi.OptString{
                        Value: event.Details.OriginCredit.Cvu,
                        Set:   true,
                    },
                    Alias: accountsapi.OptString{},
                },
                Set: true,
            },
        },
    )
    if err != nil {
        return accountsapi.Account{}, werrors.NewRetryableInternalError("failed getting account from accounts service: %s", err.Error())
    }
    switch resp := getAccountResp.(type) {
    case *accountsapi.GetAccountOKApplicationJSON:
        accounts := *resp
        if len(accounts) > 1 {
            return accountsapi.Account{}, werrors.NewNonRetryableInternalError("expected only one account but got more")
        }
        return accounts[0], nil
    case *accountsapi.GetAccountNotFound:
        return accountsapi.Account{}, werrors.NewNonRetryableInternalError("no account matched the provided account details")
    default:
        return accountsapi.Account{}, werrors.NewNonRetryableInternalError("unexpected error when getting account from account service: %s", err.Error())
    }
}

func (h *EventsHandlerImpl) buildInboundPaymentStreamName(paymentId string) string {
    return fmt.Sprintf("%s.%s", InboundPaymentStreamNamePrefix, paymentId)
}
