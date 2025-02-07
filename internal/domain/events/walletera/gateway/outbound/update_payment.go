package outbound

import (
    "context"
    "fmt"
    "strconv"

    "github.com/google/uuid"
    "github.com/walletera/bind-gateway/pkg/bind/api"
    paymentsapi "github.com/walletera/payments-types/api"
    "github.com/walletera/werrors"
)

func updatePaymentStatus(
    ctx context.Context,
    client *paymentsapi.Client,
    paymentId uuid.UUID,
    bindOperationId int,
    bindOperationStatus string,
) werrors.WError {
    status, err := bindStatus2PaymentsStatus(bindOperationStatus)
    if err != nil {
        return err
    }
    patchPaymentResp, patchPaymentErr := client.PatchPayment(
        ctx,
        &paymentsapi.PaymentUpdate{
            PaymentId: paymentId,
            ExternalId: paymentsapi.OptString{
                Value: strconv.Itoa(bindOperationId),
                Set:   true,
            },
            Status: status,
        },
        paymentsapi.PatchPaymentParams{
            PaymentId: paymentId,
        })
    return handlePatchPaymentResponse(patchPaymentResp, patchPaymentErr)
}

func bindStatus2PaymentsStatus(dinopayStatus string) (paymentsapi.PaymentStatus, werrors.WError) {
    var status paymentsapi.PaymentStatus
    switch dinopayStatus {
    case string(api.CreateTransferResponseEstadoId1):
        // A procesar
        status = paymentsapi.PaymentStatusDelivered
    case string(api.CreateTransferResponseEstadoId2):
        // Aprobada
        status = paymentsapi.PaymentStatusConfirmed
    case string(api.CreateTransferResponseEstadoId3):
        // Rechazada
        status = paymentsapi.PaymentStatusRejected
    case string(api.CreateTransferResponseEstadoId4):
        // A consultar
        status = paymentsapi.PaymentStatusDelivered
    case string(api.CreateTransferResponseEstadoId5):
        // Auditar
        status = paymentsapi.PaymentStatusDelivered
    case string(api.CreateTransferResponseEstadoId6):
        // Devuelta
        // TODO check how should we handle this status?
        status = paymentsapi.PaymentStatusConfirmed
    case string(api.CreateTransferResponseEstadoId7):
        // Devuelta parcialmente
        // TODO check how should we handle this status?
        status = paymentsapi.PaymentStatusConfirmed
    default:
        return "", werrors.NewNonRetryableInternalError(fmt.Sprintf("unknown bind payment status %s", dinopayStatus))
    }
    return status, nil
}

func handlePatchPaymentResponse(patchPaymentResp paymentsapi.PatchPaymentRes, patchPaymentErr error) werrors.WError {
    errMsg := "failed updating payment in payments service: %s"
    if patchPaymentErr != nil {
        return werrors.NewRetryableInternalError(errMsg, patchPaymentErr.Error())
    }
    switch resp := patchPaymentResp.(type) {
    case *paymentsapi.PatchPaymentOK:
    case *paymentsapi.PatchPaymentBadRequest:
        return werrors.NewNonRetryableInternalError(errMsg, resp.ErrorMessage)
    case *paymentsapi.PatchPaymentInternalServerError:
        return werrors.NewNonRetryableInternalError(errMsg, resp.ErrorMessage)
    case *paymentsapi.PatchPaymentUnauthorized:
        return werrors.NewNonRetryableInternalError(errMsg, resp.ErrorMessage)
    default:
        return werrors.NewNonRetryableInternalError("failed updating payment in payments service: unknown error")
    }
    return nil
}
