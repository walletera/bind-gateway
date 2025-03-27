package tests

import (
    "bytes"
    "context"
    "fmt"
    "net/http"
    "regexp"
    "testing"

    "github.com/cucumber/godog"
    "github.com/walletera/bind-gateway/internal/app"
    "github.com/walletera/bind-gateway/internal/domain/events/walletera/gateway/inbound"
    "github.com/walletera/eventskit/eventstoredb"
)

const (
    rawBindTransferReceivedEventKey             = "rawBindTransferReceivedEventKey"
    paymentsCreateDepositEndpointExpectationKey = "paymentsCreateDepositEndpointExpectationKey"
    accountsEndpointFailsAccountNotFound        = "accountEndpointFailsAccountNotFound"
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
    ctx.Step(`^an accounts endpoint to get accounts that fails with account not found:$`, anAccountsEndpointToGetAccountsThatFailsWithAccountNotFound)
    ctx.Step(`^the InboundPaymentReceived event is parked$`, theInboundPaymentReceivedEventIsParked)
    ctx.After(afterScenarioHook)
}

func aBindTransferReceivedEvent(ctx context.Context, event *godog.DocString) (context.Context, error) {
    if event == nil || len(event.Content) == 0 {
        return ctx, fmt.Errorf("the WithdrawalCreated event is empty or was not defined")
    }
    return context.WithValue(ctx, rawBindTransferReceivedEventKey, []byte(event.Content)), nil
}

func anAccountsEndpointToGetAccounts(ctx context.Context, mockserverExpectation *godog.DocString) (context.Context, error) {
    return createMockServerExpectation(ctx, mockserverExpectation, getAccountEndpointExpectationKey)
}

func aPaymentsEndpointToCreatePayments(ctx context.Context, mockserverExpectation *godog.DocString) (context.Context, error) {
    return createMockServerExpectation(ctx, mockserverExpectation, paymentsCreateDepositEndpointExpectationKey)
}

func theWebhookEventIsReceived(ctx context.Context) (context.Context, error) {
    rawEvent := ctx.Value(rawBindTransferReceivedEventKey).([]byte)
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

func anAccountsEndpointToGetAccountsThatFailsWithAccountNotFound(ctx context.Context, mockserverExpectation *godog.DocString) (context.Context, error) {
    return createMockServerExpectation(ctx, mockserverExpectation, accountsEndpointFailsAccountNotFound)
}

func theInboundPaymentReceivedEventIsParked(ctx context.Context) (context.Context, error) {
    rawEvent := ctx.Value(rawBindTransferReceivedEventKey).([]byte)
    reg := regexp.MustCompile("\"origin_id\": (\\d+),")
    found := reg.FindSubmatch(rawEvent)
    if found == nil || len(found) != 2 {
        return ctx, fmt.Errorf("couldn't get origin id from raw bind transfer")
    }
    client, err := eventstoredb.GetESDBClient(eventStoreDBUrl)
    if err != nil {
        return ctx, err
    }
    db := eventstoredb.NewDB(client)
    parkedEvents, err := db.ReadEvents(
        ctx,
        fmt.Sprintf("$persistentsubscription-%s::%s-parked", "$ce-bindGateway-inboundPayment", app.ESDB_SubscriptionGroupName),
    )
    if err != nil {
        return ctx, err
    }
    if len(parkedEvents) == 0 {
        return ctx, fmt.Errorf("no parked events originId")
    }
    if len(parkedEvents) > 1 {
        return ctx, fmt.Errorf("multiple parked events originId")
    }
    event, err := inbound.NewEventsDeserializer().Deserialize(parkedEvents[0].RawEvent)
    if err != nil {
        return ctx, err
    }
    if event.Type() != inbound.PaymentReceivedType {
        return ctx, fmt.Errorf("unexpected parked event type: %s", event.Type())
    }
    return ctx, nil
}
