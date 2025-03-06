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
    bindOperationStatus api.CreateTransferResponseEstadoId,
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

func bindStatus2PaymentsStatus(bindStatus api.CreateTransferResponseEstadoId) (paymentsapi.PaymentStatus, werrors.WError) {
    var status paymentsapi.PaymentStatus
    switch bindStatus {
    case api.CreateTransferResponseEstadoId1:
        // A procesar
        status = paymentsapi.PaymentStatusDelivered
    case api.CreateTransferResponseEstadoId2:
        // Aprobada
        status = paymentsapi.PaymentStatusConfirmed
    case api.CreateTransferResponseEstadoId3:
        // Rechazada
        status = paymentsapi.PaymentStatusRejected
    case api.CreateTransferResponseEstadoId4:
        // A consultar
        status = paymentsapi.PaymentStatusDelivered
    case api.CreateTransferResponseEstadoId5:
        // Auditar
        // No se pudo resolver el estado definitivo
        // de la transferencia en línea durante cierto
        // tiempo (30 seg). El estado definitivo
        // sera resuelto manualmente.
        status = paymentsapi.PaymentStatusDelivered
    case api.CreateTransferResponseEstadoId6:
        // Devuelta
        // TODO check how should we handle this status?
        status = paymentsapi.PaymentStatusConfirmed
    case api.CreateTransferResponseEstadoId7:
        // Devuelta parcialmente
        // TODO check how should we handle this status?
        status = paymentsapi.PaymentStatusConfirmed
    default:
        return "", werrors.NewNonRetryableInternalError(fmt.Sprintf("unknown bind payment status %d", bindStatus))
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
