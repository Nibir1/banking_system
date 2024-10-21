package api

import (
	"github.com/go-playground/validator/v10" // Import for using validator library
	"github.com/nibir1/banking_system/util"  // Import for util functions (assumed to have IsSupportedCurrency function)
)

// Define a custom validation function named validCurrency
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	// Check if the field being validated is a string type
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// If it's a string, attempt to validate the currency using the IsSupportedCurrency function (assumed to exist in util)
		return util.IsSupportedCurrency(currency)
	}

	// If the field is not a string or the validation fails, return false (invalid)
	return false
}
