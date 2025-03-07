package bind

import (
    "context"
    "strconv"

    "github.com/walletera/bind-gateway/pkg/wuuid"
    "github.com/walletera/werrors"
)

const InboundTransferReceivedEventType = "InboundTransferReceived"

type InboundTransferReceived struct {
    Id           string `json:"id"`
    TransferType string `json:"type"`
    From         struct {
        BankId    string `json:"bank_id"`
        AccountId string `json:"account_id"`
    } `json:"from"`
    Counterparty struct {
        Id          string `json:"id"`
        Name        string `json:"name"`
        IdType      string `json:"id_type"`
        BankRouting struct {
            Scheme  string `json:"scheme"`
            Address string `json:"address"`
        } `json:"bank_routing"`
        AccountRouting struct {
            Scheme  string `json:"scheme"`
            Address string `json:"address"`
        } `json:"account_routing"`
    } `json:"counterparty"`
    Details struct {
        OriginId    int64 `json:"origin_id"`
        OriginDebit struct {
            Cvu  string `json:"cvu"`
            Cuit int64  `json:"cuit"`
        } `json:"origin_debit"`
        OriginCredit struct {
            Cvu  string `json:"cvu"`
            Cuit int64  `json:"cuit"`
        } `json:"origin_credit"`
    } `json:"details"`
    TransactionIds []string    `json:"transaction_ids"`
    Status         string      `json:"status"`
    StartDate      string      `json:"start_date"`
    EndDate        string      `json:"end_date"`
    Challenge      interface{} `json:"challenge"`
    Charge         struct {
        Summary string `json:"summary"`
        Value   struct {
            Currency string  `json:"currency"`
            Amount   float64 `json:"amount"`
        } `json:"value"`
    } `json:"charge"`
}

func (t InboundTransferReceived) ID() string {
    return strconv.FormatInt(t.Details.OriginId, 10)
}

func (t InboundTransferReceived) Type() string {
    return InboundTransferReceivedEventType
}

func (t InboundTransferReceived) CorrelationID() string {
    return wuuid.NewUUID().String()
}

func (t InboundTransferReceived) DataContentType() string {
    return "application/json"
}

func (t InboundTransferReceived) Serialize() ([]byte, error) {
    return nil, werrors.NewNonRetryableInternalError("bind InboundTransferReceived is not serializable")
}

func (t InboundTransferReceived) Accept(ctx context.Context, handler EventsHandler) werrors.WError {
    return handler.HandleInboundTransferReceived(ctx, t)
}
