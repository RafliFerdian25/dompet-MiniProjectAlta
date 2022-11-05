package constantError

import "gorm.io/gorm"

const (
	// Error
	// ErrorAccountNotFound is error message when account not found
	ErrorAccountNotFound = "account not found"
	// ErrorAccountNotFound is error message when account not found
	ErrorOldAccountNotFound = "old account not found"
	// ErrorAccountNotFound is error message when account not found
	ErrorNewAccountNotFound = "new account not found"
	// ErrorTransactionNotFound is error message when transaction not found
	ErrorTransactionNotFound = "transaction not found"
	// ErrorDebtNotFound is error message when debt not found
	ErrorDebtNotFound = "debt not found"
	// ErrorSubCategoryNotFound is error message when subcategory not found
	ErrorSubCategoryNotFound = "subcategory not found"
	// ErrorAccountNotEnoughBalance is error message when account not enough balance
	ErrorAccountNotEnoughBalance = "account not enough balance"
	// ErrorNotAuthorized is error message when user not authorized
	ErrorNotAuthorized = "you are not authorized"
	// ErrorBalanceLessThanZero is error message when balance less than zero
	ErrorBalanceLessThanZero = "balance less than zero"
	// ErrorCannotChangeSubCategory is error message when user cannot change sub category
	ErrorCannotChangeSubCategory = "you cannot change sub category"
	// ErrorDestinationAccountNotEnoughBalance is error message when destination account not enough balance
	ErrorRecipientAccountNotEnoughBalance = "not enough recipient account balance"
	// ErrorSenderAndRecipientAccountIsSame is error message when sender and recipient account is same
	ErrorSenderAndRecipientAccountIsSame = "sender account and recipient account is same"
	// ErrorOldAccountBalanceNotEnough is error message when old account balance not enough
	ErrorOldAccountBalanceNotEnough = "old account balance is not enough"
	// ErrorNewAccountBalanceNotEnough is error message when new account balance not enough
	ErrorNewAccountBalanceNotEnough = "new account balance is not enough"
	// ErrorEmailOrPasswordNotMatch is error message when email or password not match
	ErrorEmailOrPasswordNotMatch = "email or password not match"
)

var ErrorCode = map[string]int{
	gorm.ErrRecordNotFound.Error():   404,
	"account not found":              404,
	"old account not found":          404,
	"new account not found":          404,
	"transaction not found":          404,
	"debt not found":                 404,
	"you are not authorized":         401,
	"account not enough balance":     400,
	"balance less than zero":         400,
	"you cannot change sub category": 400,
	"not enough recipient account balance": 400,
	"sender account and recipient account is same": 400,
	"old account balance is not enough": 400,
	"new account balance is not enough": 400,
	"email or password not match": 400,
}
