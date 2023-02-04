package api

import (
	"github.com/cyhe50/simple_bank/util"
	"github.com/go-playground/validator"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}

	return false
}
