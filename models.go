package match

import (
	"github.com/quagmt/udecimal"

	"github.com/miliu2cc/match-engine/protocol"
)

const (
	EngineVersion = "v1.0.0"

	SnapshotSchemaVersion = 1
)

type Side = protocol.Side

const (
	Buy  Side = protocol.SideBuy
	Sell Side = protocol.SideSell
)

type OrderType = protocol.OrderType

const (
	Market   OrderType = protocol.OrderTypeMarket
	Limit    OrderType = protocol.OrderTypeLimit
	FOK      OrderType = protocol.OrderTypeFOK
	IOC      OrderType = protocol.OrderTypeIOC
	PostOnly OrderType = protocol.OrderTypePostOnly
	Cancel   OrderType = protocol.OrderTypeCancel
)

type Order struct {
	ID        string           `json:"id"`
	Side      Side             `json:"side"`
	Price     udecimal.Decimal `json:"price"`
	Size      udecimal.Decimal `json:"size"`
	Type      OrderType        `json:"type"`
	UserID    uint64           `json:"user_id"`
	Timestamp int64            `json:"timestamp"`

	VisibleLimit udecimal.Decimal `json:"visible_limit,omitzero"`
	HiddenSize   udecimal.Decimal `json:"hidden_size,omitzero"`

	next *Order
	prev *Order
}

type DepthChange struct {
	Side     Side
	Price    udecimal.Decimal
	SizeDiff udecimal.Decimal
}

type InputEvent struct {
	Request any

	Query any
	Resp  chan any
}
