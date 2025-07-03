package dto

type FastBankApplicationRequest struct {
	PhoneNumber              string  `json:"phoneNumber"`
	Email                    string  `json:"email"`
	MonthlyIncomeAmount      float64 `json:"monthlyIncomeAmount"`
	MonthlyCreditLiabilities float64 `json:"monthlyCreditLiabilities"`
	Dependents               int     `json:"dependents"`
	AgreeToDataSharing       bool    `json:"agreeToDataSharing"`
	Amount                   float64 `json:"amount"`
}

type FastBankApplication struct {
	ID     string         `json:"id"`
	Status string         `json:"status"`
	Offer  *FastBankOffer `json:"offer,omitempty"`
}

type FastBankOffer struct {
	MonthlyPaymentAmount float64 `json:"monthlyPaymentAmount"`
	TotalRepaymentAmount float64 `json:"totalRepaymentAmount"`
	NumberOfPayments     int     `json:"numberOfPayments"`
	AnnualPercentageRate float64 `json:"annualPercentageRate"`
	FirstRepaymentDate   string  `json:"firstRepaymentDate"`
}
