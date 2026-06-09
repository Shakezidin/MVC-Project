package utils

// MaskAccountNumber masks all but the last 4 digits of an account number.
func MaskAccountNumber(accountNumber string) string {
	if len(accountNumber) <= 4 {
		return accountNumber
	}
	masked := ""
	for i := 0; i < len(accountNumber)-4; i++ {
		masked += "X"
	}
	return masked + accountNumber[len(accountNumber)-4:]
}
