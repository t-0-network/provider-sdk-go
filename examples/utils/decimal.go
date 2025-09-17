package utils

import (
	"github.com/shopspring/decimal"
	"github.com/t-0-network/provider-sdk-go/api/tzero/v1/common"
)

func DecimalFromProto(amount *common.Decimal) decimal.Decimal {
	return decimal.New(amount.Unscaled, amount.Exponent)
}

func DecimalToProto(amount decimal.Decimal) *common.Decimal {
	return &common.Decimal{
		Unscaled: amount.CoefficientInt64(),
		Exponent: amount.Exponent(),
	}
}
