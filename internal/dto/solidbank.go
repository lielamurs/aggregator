package dto

type SolidBankApplicationRequest struct {
	Phone           string  `json:"phone"`
	Email           string  `json:"email"`
	MonthlyIncome   float64 `json:"monthlyIncome"`
	MonthlyExpenses float64 `json:"monthlyExpenses"`
	MaritalStatus   string  `json:"maritalStatus"`
	AgreeToBeScored bool    `json:"agreeToBeScored"`
	Amount          float64 `json:"amount"`
}

type SolidBankApplication struct {
	ID     string          `json:"id"`
	Status string          `json:"status"`
	Offer  *SolidBankOffer `json:"offer,omitempty"`
}

type SolidBankOffer struct {
	MonthlyPaymentAmount float64 `json:"monthlyPaymentAmount"`
	TotalRepaymentAmount float64 `json:"totalRepaymentAmount"`
	NumberOfPayments     int     `json:"numberOfPayments"`
	AnnualPercentageRate float64 `json:"annualPercentageRate"`
	FirstRepaymentDate   string  `json:"firstRepaymentDate"`
}
