package tests

import (
    "context"
    "fmt"
    "testing"
    "time"

    "github.com/cucumber/godog"
    "github.com/walletera/bind-gateway/internal/app"
    "github.com/walletera/eventskit/events"

    "github.com/walletera/eventskit/rabbitmq"
)

const (
    rawWithdrawalCreatedEventKey                     = "rawWithdrawalCreatedEvent"
    bindEndpointCreatePaymentsExpectationIdKey       = "bindEndpointCreatePaymentsExpectationId"
    paymentsEndpointUpdateWithdrawalExpectationIdKey = "paymentsEndpointUpdateWithdrawalExpectationId"
    expectationTimeout                               = 5 * time.Second
)

func TestPaymentCreatedEventProcessing(t *testing.T) {

    suite := godog.TestSuite{
        ScenarioInitializer: InitializeProcessWithdrawalCreatedScenario,
        Options: &godog.Options{
            Format:   "pretty",
            Paths:    []string{"features/payment_created.feature"},
            TestingT: t, // Testing instance that will run subtests.
        },
    }

    if suite.Run() != 0 {
        t.Fatal("non-zero status returned, failed to run feature tests")
    }
}

func InitializeProcessWithdrawalCreatedScenario(ctx *godog.ScenarioContext) {
    ctx.Before(beforeScenarioHook)
    ctx.Given(`^a running bind-gateway$`, aRunningBindGateway)
    ctx.Given(`^a PaymentCreated event:$`, aPaymentCreatedEvent)
    ctx.Given(`^a bind endpoint to create transfers:$`, aBindEndpointToCreateTransfers)
    ctx.Given(`^a payments endpoint to update payments:$`, aPaymentsEndpointToUpdatePayments)
    ctx.When(`^the event is published$`, theEventIsPublished)
    ctx.Then(`^the bind-gateway creates the corresponding payment on the Bind API$`, theDinopayGatewayCreatesTheCorrespondingPaymentOnTheDinoPayAPI)
    ctx.Then(`^the bind-gateway updates the payment on payments service$`, theDinopayGatewayUpdatesThePaymentOnPaymentsService)
    ctx.Then(`the bind-gateway fails creating the corresponding payment on the DinoPay API$`, theDinoPayGatewayFailsCreatingTheCorrespondingPayment)
    ctx.Then(`^the bind-gateway produces the following log:$`, theBindGatewayProducesTheFollowingLog)
    ctx.After(afterScenarioHook)
}

func aPaymentCreatedEvent(ctx context.Context, event *godog.DocString) (context.Context, error) {
    if event == nil || len(event.Content) == 0 {
        return ctx, fmt.Errorf("the WithdrawalCreated event is empty or was not defined")
    }
    return context.WithValue(ctx, rawWithdrawalCreatedEventKey, []byte(event.Content)), nil
}

func aBindEndpointToCreateTransfers(ctx context.Context, mockserverExpectation *godog.DocString) (context.Context, error) {
    return createMockServerExpectation(ctx, mockserverExpectation, bindEndpointCreatePaymentsExpectationIdKey)
}

func aPaymentsEndpointToUpdatePayments(ctx context.Context, mockserverExpectation *godog.DocString) (context.Context, error) {
    return createMockServerExpectation(ctx, mockserverExpectation, paymentsEndpointUpdateWithdrawalExpectationIdKey)
}

func theEventIsPublished(ctx context.Context) (context.Context, error) {
    publisher, err := rabbitmq.NewClient(
        rabbitmq.WithExchangeName(app.RabbitMQPaymentsExchangeName),
        rabbitmq.WithExchangeType(app.RabbitMQExchangeType),
    )
    if err != nil {
        return nil, fmt.Errorf("error creating rabbitmq client: %s", err.Error())
    }

    rawEvent := ctx.Value(rawWithdrawalCreatedEventKey).([]byte)
    err = publisher.Publish(ctx, publishable{rawEvent: rawEvent}, events.RoutingInfo{
        Topic:      app.RabbitMQPaymentsExchangeName,
        RoutingKey: app.RabbitMQPaymentCreatedRoutingKey,
    })
    if err != nil {
        return nil, fmt.Errorf("error publishing WithdrawalCreated event to rabbitmq: %s", err.Error())
    }

    return ctx, nil
}

func theDinopayGatewayCreatesTheCorrespondingPaymentOnTheDinoPayAPI(ctx context.Context) (context.Context, error) {
    id := expectationIdFromCtx(ctx, bindEndpointCreatePaymentsExpectationIdKey)
    err := verifyExpectationMetWithin(ctx, id, expectationTimeout)
    return ctx, err
}

func theDinopayGatewayUpdatesThePaymentOnPaymentsService(ctx context.Context) (context.Context, error) {
    id := expectationIdFromCtx(ctx, paymentsEndpointUpdateWithdrawalExpectationIdKey)
    err := verifyExpectationMetWithin(ctx, id, expectationTimeout)
    return ctx, err
}

func theDinoPayGatewayFailsCreatingTheCorrespondingPayment(ctx context.Context) (context.Context, error) {
    id := expectationIdFromCtx(ctx, bindEndpointCreatePaymentsExpectationIdKey)
    err := verifyExpectationMetWithin(ctx, id, expectationTimeout)
    return ctx, err
}
