package x

import "fmt"

type ApiError struct {
	Code    ApiCode
	Message string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("Fail to invoke api (%d, %v)", e.Code, e.Message)
}

type ApiCode uint16

const (
	OK                     ApiCode = iota
	GeneralError
	InvalidArgument
	InternalError
	Maintained
	Unavailable
	RequestTimeExpired
	TooFrequent
	Unknown
	AuthenticationFailed
	FundPasswordLocked
	IncorrectFundPassword
	AuthenticationAuditing
	EmptyChannel
	EmptyEvent
	InsufficientFund
	OrderNotFound
	InvalidPrice
	InvalidAmount
	UserNotFound
	InvalidIpAddress
	TradeRecordNotFound
)
