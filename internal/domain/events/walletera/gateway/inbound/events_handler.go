package inbound

import (
    "context"
    "fmt"
    "log/slog"
    "strconv"
    "time"

    "github.com/google/uuid"
    accountsapi "github.com/walletera/accounts/types/api/api"
    "github.com/walletera/bind-gateway/pkg/logattr"
    "github.com/walletera/bind-gateway/pkg/wuuid"
    paymentsapi "github.com/walletera/payments-types/api"
    "github.com/walletera/werrors"
)

type EventsHandler interface {
    HandleInboundPaymentReceived(ctx context.Context, inboundPaymentReceived PaymentReceived) werrors.WError
}

type EventsHandlerImpl struct {
    accountsApiClient accountsapi.Invoker
    paymentsApiClient paymentsapi.Invoker
    deserializer      *EventsDeserializer
    logger            *slog.Logger
}

func NewEventsHandlerImpl(
    accountsApiClient accountsapi.Invoker,
    paymentsApiClient paymentsapi.Invoker,
    logger *slog.Logger) *EventsHandlerImpl {
    return &EventsHandlerImpl{
        accountsApiClient: accountsApiClient,
        paymentsApiClient: paymentsApiClient,
        deserializer:      NewEventsDeserializer(),
        logger:            logger.With(logattr.Component("gateway.inbound.EventsHandlerImpl")),
    }
}

func (h *EventsHandlerImpl) HandleInboundPaymentReceived(ctx context.Context, inboundPaymentReceived PaymentReceived) werrors.WError {
    correlationId := inboundPaymentReceived.CorrelationId

    account, err := h.getBeneficiaryAccount(ctx, inboundPaymentReceived)
    if err != nil {
        h.logger.Error(
            wrappedErrMsg("failed getting beneficiary account"),
            logattr.CorrelationId(correlationId.String()),
            logattr.Error(err.Error()),
        )
        return err
    }

    payment := h.buildPayment(correlationId, inboundPaymentReceived, account)

    resp, postPaymentErr := h.paymentsApiClient.PostPayment(ctx, &payment, paymentsapi.PostPaymentParams{
        XWalleteraCorrelationID: paymentsapi.OptUUID{
            Value: inboundPaymentReceived.CorrelationId,
            Set:   true,
        },
    })
    if postPaymentErr != nil {
        h.logger.Error(
            wrappedErrMsg("failed creating payment on payments api"),
            logattr.CorrelationId(correlationId.String()),
            logattr.Error(postPaymentErr.Error()),
        )
        return werrors.NewRetryableInternalError(postPaymentErr.Error())
    }
    switch r := resp.(type) {
    case *paymentsapi.Payment:
        h.logger.Info("gateway event InboundPaymentReceived processed successfully",
            logattr.CorrelationId(correlationId.String()),
            logattr.PaymentId(r.ID.String()),
            logattr.EventType(inboundPaymentReceived.Type()),
        )
        return nil
    case *paymentsapi.PostPaymentInternalServerError:
        h.logger.Error(wrappedErrMsg("failed creating payment on payments api"),
            logattr.CorrelationId(correlationId.String()),
            logattr.Error(r.ErrorMessage),
        )
        return werrors.NewRetryableInternalError(r.ErrorMessage)
    case *paymentsapi.PostPaymentBadRequest:
        h.logger.Error(wrappedErrMsg("failed creating payment on payments api"),
            logattr.CorrelationId(correlationId.String()),
            logattr.Error(r.ErrorMessage))
        return werrors.NewNonRetryableInternalError(r.ErrorMessage)
    case *paymentsapi.PostPaymentConflict:
        h.logger.Error(wrappedErrMsg("failed creating payment on payments api"),
            logattr.CorrelationId(correlationId.String()),
            logattr.Error(r.ErrorMessage))
        return werrors.NewNonRetryableInternalError(r.ErrorMessage)
    case *paymentsapi.PostPaymentUnauthorized:
        h.logger.Error(wrappedErrMsg("failed creating payment on payments api: unauthorized"))
        return werrors.NewNonRetryableInternalError("failed creating payment on payments api",
            logattr.CorrelationId(correlationId.String()),
            logattr.Error("unauthorized"),
        )
    default:
        h.logger.Error("unexpected error creating payment on payments api",
            logattr.CorrelationId(correlationId.String()),
        )
        return werrors.NewRetryableInternalError(wrappedErrMsg("unexpected error creating payment on payments api"))
    }
}

func (h *EventsHandlerImpl) getBeneficiaryAccount(ctx context.Context, event PaymentReceived) (accountsapi.Account, werrors.WError) {
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
                        Value: strconv.FormatInt(event.OriginCreditCuit, 10),
                        Set:   true,
                    },
                    Cvu: accountsapi.OptString{
                        Value: event.OriginCreditCvu,
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

func (h *EventsHandlerImpl) buildPayment(correlationUUID uuid.UUID, event PaymentReceived, account accountsapi.Account) paymentsapi.Payment {
    paymentUUID := wuuid.NewUUID()
    payment := paymentsapi.Payment{
        ID:       paymentUUID,
        Amount:   event.ChargeValueAmount,
        Currency: event.Currency,
        Direction: paymentsapi.OptPaymentDirection{
            Value: paymentsapi.PaymentDirectionInbound,
            Set:   true,
        },
        CustomerId: paymentsapi.OptUUID{
            Value: account.CustomerId,
            Set:   true,
        },
        ExternalId: paymentsapi.OptString{
            Value: strconv.FormatInt(event.OriginId, 10),
            Set:   true,
        },
        SchemeId: paymentsapi.OptString{
            Value: event.CoelsaId,
            Set:   true,
        },
        Beneficiary: paymentsapi.OptAccountDetails{
            Value: paymentsapi.AccountDetails{
                Currency: paymentsapi.OptString{
                    Value: event.Currency,
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
                                    Value: strconv.FormatInt(event.OriginCreditCuit, 10),
                                    Set:   true,
                                },
                                Cvu: paymentsapi.OptString{
                                    Value: event.OriginCreditCvu,
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
                    Value: event.Currency,
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
                                    Value: strconv.FormatInt(event.OriginDebitCuit, 10),
                                    Set:   true,
                                },
                                Cvu: paymentsapi.OptString{
                                    Value: event.OriginDebitCvu,
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
    }
    return payment
}

func wrappedErrMsg(msg string) string {
    return fmt.Sprintf("failed processing InboundPaymentReceived event: %s", msg)
}
