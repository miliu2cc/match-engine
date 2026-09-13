package protocol

import (
	"github.com/quagmt/udecimal"
)

// CommandType 定义命令类型
type CommandType uint8

const (
	// CmdUnknown 未知命令
	CmdUnknown CommandType = 0
	// CmdPlaceOrder 下单命令
	CmdPlaceOrder CommandType = 1
	// CmdCancelOrder 撤单命令
	CmdCancelOrder CommandType = 2
	// CmdAmendOrder 改单命令
	CmdAmendOrder CommandType = 3
	// CmdCreateMarket 创建市场命令
	CmdCreateMarket CommandType = 11
	// CmdSuspendMarket 挂起市场命令
	CmdSuspendMarket CommandType = 12
	// CmdResumeMarket 恢复市场命令
	CmdResumeMarket CommandType = 13
	// CmdUpdateConfig 更新配置命令
	CmdUpdateConfig CommandType = 14
	// CmdUserStatusEvent 用户状态事件命令
	CmdUserStatusEvent CommandType = 21
)

// BaseCommand 是所有命令的基础结构体
type BaseCommand struct {
	SeqID     uint64 // 上游分配的单调序列，用于保持逻辑命令的顺序。
	CommandID string
	UserID    string
	MarketID  string // 交易对名
	Timestamp int64
}

// CommandRequest 是所有命令请求的接口
type CommandRequest interface {
	Base() BaseCommand
}

// GetRequestBase 返回一个类型化请求的共享元数据和一个表示成功与否的布尔值
func GetRequestBase(req any) (BaseCommand, bool) {
	if cr, ok := req.(CommandRequest); ok {
		return cr.Base(), true
	}
	// used in tests
	switch r := req.(type) {
	case *BaseCommand:
		return *r, true
	case BaseCommand:
		return r, true
	}
	return BaseCommand{}, false
}

// Base 方法返回基础命令
func (r *BaseCommand) Base() BaseCommand {
	return *r
}

// PlaceOrderRequest 是下单请求的结构体
type PlaceOrderRequest struct {
	BaseCommand

	OrderID   string           `json:"order_id"`
	Side      Side             `json:"side"`
	OrderType OrderType        `json:"order_type"`
	Price     udecimal.Decimal `json:"price"`
	Size      udecimal.Decimal `json:"size"`
	Visible   udecimal.Decimal `json:"visible_size"`
	QuoteSize udecimal.Decimal `json:"quote_size"`
}

// Base 方法返回嵌入的 BaseCommand
func (r *PlaceOrderRequest) Base() BaseCommand {
	return r.BaseCommand
}

type CancelOrderRequest struct {
	BaseCommand

	OrderID string `json:"order_id"`
}

func (r *CancelOrderRequest) Base() BaseCommand {
	return r.BaseCommand
}

type AmendOrderRequest struct {
	BaseCommand

	OrderID  string           `json:"order_id"`
	NewPrice udecimal.Decimal `json:"new_price"`
	NewSize  udecimal.Decimal `json:"new_size"`
}

func (r *AmendOrderRequest) Base() BaseCommand {
	return r.BaseCommand
}

type CreateMarketRequest struct {
	BaseCommand

	MinLotSize udecimal.Decimal `json:"min_lot_size"`
}

func (r *CreateMarketRequest) Base() BaseCommand {
	return r.BaseCommand
}

type SuspendMarketRequest struct {
	BaseCommand

	Reason string `json:"reason"`
}

func (r *SuspendMarketRequest) Base() BaseCommand {
	return r.BaseCommand
}

type ResumeMarketRequest struct {
	BaseCommand
}

func (r *ResumeMarketRequest) Base() BaseCommand {
	return r.BaseCommand
}

type UpdateConfigRequest struct {
	BaseCommand

	MinLotSize udecimal.Decimal `json:"min_lot_size"`
}

func (r *UpdateConfigRequest) Base() BaseCommand {
	return r.BaseCommand
}

type UserEventRequest struct {
	BaseCommand

	EventType string `json:"event_type"`
	Key       string `json:"key"`
	Data      []byte `json:"data"`
}

func (r *UserEventRequest) Base() BaseCommand {
	return r.BaseCommand
}
