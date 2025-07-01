package dto

import "time"

type BankSubmission struct {
	BankName    string               `json:"bankName"`
	Status      BankSubmissionStatus `json:"status"`
	BankID      string               `json:"bankId,omitempty"`
	SubmittedAt time.Time            `json:"submittedAt"`
	CompletedAt *time.Time           `json:"completedAt,omitempty"`
	Error       string               `json:"error,omitempty"`
}

type BankSubmissionStatus string

const (
	BankStatusPending   BankSubmissionStatus = "PENDING"
	BankStatusSubmitted BankSubmissionStatus = "SUBMITTED"
	BankStatusSuccess   BankSubmissionStatus = "SUCCESS"
	BankStatusRejected  BankSubmissionStatus = "REJECTED"
	BankStatusError     BankSubmissionStatus = "ERROR"
)
