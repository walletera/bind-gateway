package outbound

import "fmt"

const (
    OutboundPaymentStreamNamePrefix = "bindGateway-outboundPayment"
    InboundPaymentStreamNamePrefix  = "bindGateway-inboundPayment"
)

func BuildOutboundPaymentStreamName(id string) string {
    return fmt.Sprintf("%s.%s", OutboundPaymentStreamNamePrefix, id)
}

func BuildInboundPaymentStreamName(id string) string {
    return fmt.Sprintf("%s.%s", InboundPaymentStreamNamePrefix, id)
}
