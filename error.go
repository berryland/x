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
	OK                     = ApiCode(1000)
	GeneralError           = ApiCode(1001)
	InternalError          = ApiCode(1002)
	AuthenticationFailed   = ApiCode(1003)
	FundPasswordLocked     = ApiCode(1004)
	IncorrectFundPassword  = ApiCode(1005)
	AuthenticationAuditing = ApiCode(1006)
	EmptyChannel           = ApiCode(1007)
	EmptyEvent             = ApiCode(1008)
	Maintained             = ApiCode(1009)
	InsufficientFund       = ApiCode(2009)
	OrderNotFound          = ApiCode(3001)
	InvalidPrice           = ApiCode(3002)
	InvalidAmount          = ApiCode(3003)
	UserNotFound           = ApiCode(3004)
	InvalidArgument        = ApiCode(3005)
	InvalidIpAddress       = ApiCode(3006)
	RequestTimeExpired     = ApiCode(3007)
	TradeRecordNotFound    = ApiCode(3008)
	Unavailable            = ApiCode(4001)
	TooFrequent            = ApiCode(4002)
)
