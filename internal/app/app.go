package app

import (
    "context"
    "fmt"
    "log/slog"
    "time"

    "github.com/walletera/bind-gateway/internal/adapters/bind"
    "github.com/walletera/bind-gateway/internal/domain/events/walletera/gateway/outbound"
    "github.com/walletera/bind-gateway/internal/domain/events/walletera/payments"
    "github.com/walletera/bind-gateway/pkg/logattr"
    "github.com/walletera/bind-gateway/pkg/paymentsauth"
    "github.com/walletera/eventskit/eventstoredb"
    "github.com/walletera/eventskit/messages"
    "github.com/walletera/eventskit/rabbitmq"
    paymentsapi "github.com/walletera/payments-types/api"
    paymentsevents "github.com/walletera/payments-types/events"
    "github.com/walletera/werrors"
    "go.uber.org/zap"
    "go.uber.org/zap/exp/zapslog"
    "go.uber.org/zap/zapcore"
)

const (
    RabbitMQPaymentsExchangeName              = "payments.events"
    RabbitMQExchangeType                      = "topic"
    RabbitMQPaymentCreatedRoutingKey          = "payment.created"
    RabbitMQQueueName                         = "bind-gateway"
    ESDB_ByCategoryProjection_OutboundPayment = "$ce-bindGateway-outboundPayment"
    ESDB_ByCategoryProjection_InboundPayment  = "$ce-bindGateway-inboundPayment"
    ESDB_SubscriptionGroupName                = "bind-gateway"
    WebhookServerPort                         = 8686
)

type App struct {
    rabbitmqHost     string
    rabbitmqPort     int
    rabbitmqUser     string
    rabbitmqPassword string
    bindEnv          string
    bindUrl          string
    paymentsUrl      string
    esdbUrl          string
    logHandler       slog.Handler
    logger           *slog.Logger
}

func NewApp(opts ...Option) (*App, error) {
    app := &App{}
    err := setDefaultOpts(app)
    if err != nil {
        return nil, fmt.Errorf("failed setting default options: %w", err)
    }
    for _, opt := range opts {
        opt(app)
    }
    return app, nil
}

func (app *App) Run(ctx context.Context) error {

    // TODO log initialization to main
    appLogger := slog.
        New(app.logHandler).
        With(logattr.ServiceName("bind-gateway"))
    app.logger = appLogger

    err := app.execESDBSetupTasks(ctx)
    if err != nil {
        return err
    }

    paymentsMessageProcessor, err := createPaymentsMessageProcessor(app, appLogger)
    if err != nil {
        return fmt.Errorf("failed creating payments message processor: %w", err)
    }

    err = paymentsMessageProcessor.Start(ctx)
    if err != nil {
        return fmt.Errorf("failed starting payments rabbitmq processor: %w", err)
    }

    appLogger.Info("payments message processor started")

    gatewayMessageProcessor, err := createGatewayMessageProcessor(app, appLogger)
    if err != nil {
        return fmt.Errorf("failed creating dinopay message processor: %w", err)
    }

    err = gatewayMessageProcessor.Start(ctx)
    if err != nil {
        return fmt.Errorf("failed starting dinopay message processor: %w", err)
    }

    appLogger.Info("gateway message processor started")

    appLogger.Info("bind-gateway started")

    return nil
}

func (app *App) Stop(ctx context.Context) {
    // TODO implement processor gracefull shutdown
    app.logger.Info("bind-gateway stopped")
}

func setDefaultOpts(app *App) error {
    zapLogger, err := newZapLogger()
    if err != nil {
        return err
    }
    app.logHandler = zapslog.NewHandler(zapLogger.Core(), nil)
    return nil
}

func newZapLogger() (*zap.Logger, error) {
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
    zapConfig := zap.Config{
        Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
        Development: false,
        Sampling: &zap.SamplingConfig{
            Initial:    100,
            Thereafter: 100,
        },
        Encoding:         "json",
        EncoderConfig:    encoderConfig,
        OutputPaths:      []string{"stderr"},
        ErrorOutputPaths: []string{"stderr"},
    }
    return zapConfig.Build()
}

func (app *App) execESDBSetupTasks(_ context.Context) error {
    err := eventstoredb.CreatePersistentSubscription(app.esdbUrl, ESDB_ByCategoryProjection_OutboundPayment, ESDB_SubscriptionGroupName)
    if err != nil {
        return fmt.Errorf("failed creating persistent subscription for %s: %w", ESDB_ByCategoryProjection_OutboundPayment, err)
    }

    err = eventstoredb.CreatePersistentSubscription(app.esdbUrl, ESDB_ByCategoryProjection_InboundPayment, ESDB_SubscriptionGroupName)
    if err != nil {
        return fmt.Errorf("failed creating persistent subscription for %s: %w", ESDB_ByCategoryProjection_InboundPayment, err)
    }
    return nil
}

func createPaymentsMessageProcessor(app *App, logger *slog.Logger) (*messages.Processor[paymentsevents.Handler], error) {
    bindClient, err := bind.NewClient(app.bindUrl)
    if err != nil {
        return nil, fmt.Errorf("failed parsing dinopay url %s: %w", app.bindUrl, err)
    }

    esdbClient, err := eventstoredb.GetESDBClient(app.esdbUrl)
    if err != nil {
        return nil, fmt.Errorf("failed getting esdb client: %w", err)
    }

    eventsDB := eventstoredb.NewDB(esdbClient)
    handler := payments.NewEventsHandler(bindClient, eventsDB, logger)
    queueName := fmt.Sprintf(RabbitMQQueueName)

    rabbitMQClient, err := rabbitmq.NewClient(
        rabbitmq.WithHost(app.rabbitmqHost),
        rabbitmq.WithPort(uint(app.rabbitmqPort)),
        rabbitmq.WithUser(app.rabbitmqUser),
        rabbitmq.WithPassword(app.rabbitmqPassword),
        rabbitmq.WithExchangeName(RabbitMQPaymentsExchangeName),
        rabbitmq.WithExchangeType(RabbitMQExchangeType),
        rabbitmq.WithConsumerRoutingKeys(RabbitMQPaymentCreatedRoutingKey),
        rabbitmq.WithQueueName(queueName),
    )
    if err != nil {
        return nil, fmt.Errorf("creating rabbitmq client: %w", err)
    }

    paymentsMessageProcessor, err := messages.NewProcessor[paymentsevents.Handler](
        rabbitMQClient,
        paymentsevents.NewDeserializer(logger),
        handler,
        withErrorCallback(
            logger.With(
                logattr.Component("payments.rabbitmq.MessageProcessor")),
        ),
    ), nil
    if err != nil {
        return nil, fmt.Errorf("failed creating payments rabbitmq processor: %w", err)
    }

    return paymentsMessageProcessor, nil
}

func createGatewayMessageProcessor(app *App, logger *slog.Logger) (*messages.Processor[outbound.EventsHandler], error) {

    paymentsClient, err := paymentsapi.NewClient(app.paymentsUrl, paymentsauth.NewSecuritySource())
    if err != nil {
        return nil, fmt.Errorf("failed creating payments api client: %w", err)
    }

    esdbMessagesConsumer, err := eventstoredb.NewMessagesConsumer(
        app.esdbUrl,
        ESDB_ByCategoryProjection_OutboundPayment,
        ESDB_SubscriptionGroupName,
    )
    if err != nil {
        return nil, fmt.Errorf("failed creating esdb messages consumer: %w", err)
    }

    esdbClient, err := eventstoredb.GetESDBClient(app.esdbUrl)
    if err != nil {
        return nil, fmt.Errorf("failed creating esdb client: %w", err)
    }

    eventsDB := eventstoredb.NewDB(esdbClient)

    eventsHandler := outbound.NewEventsHandlerImpl(eventsDB, paymentsClient, logger)
    return messages.NewProcessor[outbound.EventsHandler](
            esdbMessagesConsumer,
            outbound.NewEventsDeserializer(),
            eventsHandler,
            withErrorCallback(
                logger.With(
                    logattr.Component("gateway.esdb.MessageProcessor")),
            ),
        ),
        nil
}

func withErrorCallback(logger *slog.Logger) messages.ProcessorOpt {
    return messages.WithErrorCallback(func(wError werrors.WError) {
        logger.Error(
            "failed processing message",
            logattr.Error(wError.Message()))
    })
}
