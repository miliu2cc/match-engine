package protocol

// DepthItem 表示订单簿深度中的单个价格层级。
type DepthItem struct {
	Price string `json:"price"`
	Size  string `json:"size"`
	Count string `json:"count"`
}

// GetDepthResponse 表示获取订单簿深度的响应。
type GetDepthResponse struct {
	UpdateID uint64      `json:"update_id"`
	Asks     []DepthItem `json:"asks"`
	Bids     []DepthItem `json:"bids"`
}

// Side 表示订单的方向。
type Side int8

const (
	// SideBuy 表示买入方向。
	SideBuy Side = 1
	// SideSell 表示卖出方向。
	SideSell Side = 2
)

// OrderType 表示订单的类型。
type OrderType string

const (
	// OrderTypeMarket 表示市价单。
	OrderTypeMarket OrderType = "market"
	// OrderTypeLimit 表示限价单。
	OrderTypeLimit OrderType = "limit"
	// OrderTypeIOC 表示立即成交或取消订单。
	OrderTypeIOC OrderType = "ioc"
	// OrderTypeFOK 表示全部成交或取消订单。
	OrderTypeFOK OrderType = "fok"
	// OrderTypePostOnly 表示仅在有足够流动性时才成交。
	OrderTypePostOnly OrderType = "post_only"
	// OrderTypeCancel 表示取消订单。
	OrderTypeCancel OrderType = "cancel"
)

// OrderBookState 表示订单簿的状态。
type OrderBookState uint8

const (
	// OrderBookStateRunning 表示订单簿正在运行。
	OrderBookStateRunning OrderBookState = 1
	// OrderBookStateSuspended 表示订单簿已暂停。
	OrderBookStateSuspended OrderBookState = 2 // Suspended 允许撤单
	// OrderBookStateHalted 表示订单簿已停止。
	OrderBookStateHalted OrderBookState = 3
)

// LogType 表示日志类型。
type LogType string

const (
	// LogTypeOpen 表示订单打开日志。
	LogTypeOpen LogType = "open"
	// LogTypeMatch 表示订单匹配日志。
	LogTypeMatch LogType = "match"
	// LogTypeCancel 表示订单取消日志。
	LogTypeCancel LogType = "cancel"
	// LogTypeAmend 表示订单修改日志。
	LogTypeAmend LogType = "amend"
	// LogTypeReject 表示订单拒绝日志。
	LogTypeReject LogType = "reject"
	// LogTypeUser 表示用户事件日志。
	LogTypeUser LogType = "user_event"
	// LogTypeAdmin 表示管理命令成功日志。
	LogTypeAdmin LogType = "admin_event"
)

// RejectReason 表示订单拒绝的原因。
type RejectReason string

const (
	// RejectReasonNone 表示订单拒绝原因为空。
	RejectReasonNone RejectReason = ""
	// RejectReasonNoLiquidity 表示订单拒绝原因是无流动性。
	RejectReasonNoLiquidity RejectReason = "no_liquidity"
	// RejectReasonPriceMismatch 表示订单拒绝原因是价格不匹配。
	RejectReasonPriceMismatch RejectReason = "price_mismatch"
	// RejectReasonInsufficientSize 表示订单拒绝原因是数量不足。
	RejectReasonInsufficientSize RejectReason = "insufficient_size"
	// RejectReasonPostOnlyMatch 表示订单拒绝原因是仅在有足够流动性时才成交。
	RejectReasonPostOnlyMatch RejectReason = "post_only_match"
	// RejectReasonDuplicateID 表示订单拒绝原因是订单 ID 重复。
	RejectReasonDuplicateID RejectReason = "duplicate_order_id"
	// RejectReasonOrderNotFound 表示订单拒绝原因是订单未找到。
	RejectReasonOrderNotFound RejectReason = "order_not_found"
	// RejectReasonInvalidPayload 表示订单拒绝原因是无效的负载。
	RejectReasonInvalidPayload RejectReason = "invalid_payload"
	// RejectReasonUnknownCommand 表示订单拒绝原因是未知的命令。
	RejectReasonUnknownCommand RejectReason = "unknown_command"
	// RejectReasonMarketNotFound 表示订单拒绝原因是市场未找到。
	RejectReasonMarketNotFound RejectReason = "market_not_found"
	// RejectReasonMarketAlreadyExists 表示订单拒绝原因是市场已存在。
	RejectReasonMarketAlreadyExists RejectReason = "market_already_exists"
	// RejectReasonMarketSuspended 表示订单拒绝原因是市场已暂停。
	RejectReasonMarketSuspended RejectReason = "market_suspended"
	// RejectReasonMarketHalted 表示订单拒绝原因是市场已停止。
	RejectReasonMarketHalted RejectReason = "market_halted"
	// RejectReasonUnauthorized 表示订单拒绝原因是未经授权。
	RejectReasonUnauthorized RejectReason = "unauthorized"
)

const (
	// OrderTypeUnknownUint8 表示订单类型未知。
	OrderTypeUnknownUint8 uint8 = 0
	// OrderTypeMarketUint8 表示订单类型是市价订单。
	OrderTypeMarketUint8 uint8 = 1
	// OrderTypeLimitUint8 表示订单类型是限价订单。
	OrderTypeLimitUint8 uint8 = 2
	// OrderTypeFOKUint8 表示订单类型是立即成交或取消订单。
	OrderTypeFOKUint8 uint8 = 3
	// OrderTypeIOCUint8 表示订单类型是立即成交或取消订单。
	OrderTypeIOCUint8 uint8 = 4
	// OrderTypePostUint8 表示订单类型是后成交订单。
	OrderTypePostUint8 uint8 = 5
	// OrderTypeCancelUint8 表示订单类型是取消订单。
	OrderTypeCancelUint8 uint8 = 6
)

// ToUint8 将 OrderType 转换为对应的 uint8 值。
func (ot OrderType) ToUint8() uint8 {
	switch ot {
	case OrderTypeMarket:
		return OrderTypeMarketUint8
	case OrderTypeLimit:
		return OrderTypeLimitUint8
	case OrderTypeFOK:
		return OrderTypeFOKUint8
	case OrderTypeIOC:
		return OrderTypeIOCUint8
	case OrderTypePostOnly:
		return OrderTypePostUint8
	case OrderTypeCancel:
		return OrderTypeCancelUint8
	default:
		return OrderTypeUnknownUint8
	}
}

// OrderTypeFromUint8 将 uint8 转换为对应的 OrderType。
func OrderTypeFromUint8(v uint8) OrderType {
	switch v {
	case OrderTypeMarketUint8:
		return OrderTypeMarket
	case OrderTypeFOKUint8:
		return OrderTypeFOK
	case OrderTypeIOCUint8:
		return OrderTypeIOC
	case OrderTypePostUint8:
		return OrderTypePostOnly
	case OrderTypeCancelUint8:
		return OrderTypeCancel
	case OrderTypeLimitUint8:
		fallthrough
	default:
		return OrderTypeLimit // Default to limit
	}
}

// GetStatsResponse 表示获取统计信息的响应。
type GetStatsResponse struct {
	AskDepthCount int64 `json:"ask_depth_count"`
	AskOrderCount int64 `json:"ask_order_count"`
	BidDepthCount int64 `json:"bid_depth_count"`
	BidOrderCount int64 `json:"bid_order_count"`
}
