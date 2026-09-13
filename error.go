package match

import "errors"

var (
	// ErrInsufficientLiquidity 在没有足够的深度来填充市价/限价单/立即成交或取消单时，被返回
	ErrInsufficientLiquidity = errors.New("there is not enough depth to fill the order")
	// ErrInvalidParam 当参数无效时，被返回
	ErrInvalidParam = errors.New("the param is invalid")
	// ErrInternal 当内部服务器错误发生时，被返回
	ErrInternal = errors.New("internal server error")
	// ErrTimeout 当操作超时时，被返回
	ErrTimeout = errors.New("timeout")
	// ErrShutdown 当引擎关闭时，被返回
	ErrShutdown = errors.New("order book is shutting down")
	// ErrNotFound 当资源未找到时，被返回
	ErrNotFound = errors.New("not found")
	// ErrInvalidPrice 当价格无效时，被返回
	ErrInvalidPrice = errors.New("invalid price")
	// ErrInvalidSize 当大小无效时，被返回
	ErrInvalidSize = errors.New("invalid size")
	// ErrOrderBookClosed 当订单簿关闭时，被返回
	ErrOrderBookClosed = errors.New("order book is closed")
	// ErrUnknownCommand 当命令类型未知时，被返回
	ErrUnknownCommand = errors.New("unknown command")
	// ErrUnknownQuery 当查询类型未知时，被返回
	ErrUnknownQuery = errors.New("unknown query")
)
