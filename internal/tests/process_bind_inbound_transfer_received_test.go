package tests

import (
    "bytes"
    "context"
    "fmt"
    "net/http"
    "testing"

    "github.com/cucumber/godog"
    "github.com/walletera/bind-gateway/internal/app"
)

const (
    rawDinopayPaymentCreatedEventKey            = "rawDinopayPaymentCreatedEventKey"
    paymentsCreateDepositEndpointExpectationKey = "paymentsCreateDepositEndpointExpectationKey"
    getAccountEndpointExpectationKey            = "getAccountEndpointExpectationKey"
)

func TestBindInboundTransferReceivedProcessing(t *testing.T) {

    suite := godog.TestSuite{
        ScenarioInitializer: InitializeProcessBindInboundTransferReceivedScenario,
        Options: &godog.Options{
            Format:   "pretty",
            Paths:    []string{"features/bind_inbound_transfer_received.feature"},
            TestingT: t, // Testing instance that will run subtests.
        },
    }

    if suite.Run() != 0 {
        t.Fatal("non-zero status returned, failed to run feature tests")
    }
}

func InitializeProcessBindInboundTransferReceivedScenario(ctx *godog.ScenarioContext) {
    ctx.Before(beforeScenarioHook)
    ctx.Step(`^a running bind-gateway$`, aRunningBindGateway)
    ctx.Step(`^a Bind transfer.cvu.received event:$`, aBindTransferReceivedEvent)
    ctx.Step(`^an accounts endpoint to get accounts:$`, anAccountsEndpointToGetAccounts)
    ctx.Step(`^a payments endpoint to create payments:$`, aPaymentsEndpointToCreatePayments)
    ctx.When(`^the webhook event is received$`, theWebhookEventIsReceived)
    ctx.Step(`^the bind-gateway creates the corresponding payment on the Payments API$`, theBindGatewayCreatesTheCorrespondingPaymentOnThePaymentsAPI)
    ctx.Step(`^the bind-gateway produces the following log:$`, theBindGatewayProducesTheFollowingLog)
    ctx.After(afterScenarioHook)
}

func aBindTransferReceivedEvent(ctx context.Context, event *godog.DocString) (context.Context, error) {
    if event == nil || len(event.Content) == 0 {
        return ctx, fmt.Errorf("the WithdrawalCreated event is empty or was not defined")
    }
    return context.WithValue(ctx, rawDinopayPaymentCreatedEventKey, []byte(event.Content)), nil
}

func anAccountsEndpointToGetAccounts(ctx context.Context, mockserverExpectation *godog.DocString) (context.Context, error) {
    return createMockServerExpectation(ctx, mockserverExpectation, getAccountEndpointExpectationKey)
}

func aPaymentsEndpointToCreatePayments(ctx context.Context, mockserverExpectation *godog.DocString) (context.Context, error) {
    return createMockServerExpectation(ctx, mockserverExpectation, paymentsCreateDepositEndpointExpectationKey)
}

func theWebhookEventIsReceived(ctx context.Context) (context.Context, error) {
    rawEvent := ctx.Value(rawDinopayPaymentCreatedEventKey).([]byte)
    url := fmt.Sprintf("http://127.0.0.1:%d/webhooks", app.WebhookServerPort)
    httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(rawEvent))
    if err != nil {
        return ctx, fmt.Errorf("failed sending webhook event: %w", err)
    }
    resp, err := http.DefaultClient.Do(httpReq)
    if err != nil {
        return ctx, fmt.Errorf("failed sending request to payments api: %w", err)
    }
    if resp.StatusCode != http.StatusCreated {
        return ctx, fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
    }
    return ctx, nil
}

func theBindGatewayCreatesTheCorrespondingPaymentOnThePaymentsAPI(ctx context.Context) (context.Context, error) {
    id := expectationIdFromCtx(ctx, paymentsCreateDepositEndpointExpectationKey)
    err := verifyExpectationMetWithin(ctx, id, expectationTimeout)
    return ctx, err
}
