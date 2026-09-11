package main

import "fmt"

type PaymentEvent struct {
	ID          string
	AccountID   string
	AmountCents int64
	Currency    string
	RiskScore   int
}
type AuditEntry struct {
	PaymentID string
	Decision  string
	MessageID string
}

func notifyPayment(e PaymentEvent) (AuditEntry, error) {
	decision := "approved"
	if e.RiskScore >= 80 {
		decision = "held_for_review"
	}
	entry := AuditEntry{PaymentID: e.ID, Decision: decision}
	if decision != "approved" {
		return entry, nil
	}
	message := fmt.Sprintf("Payment %s approved: %s %.2f", e.ID, e.Currency, float64(e.AmountCents)/100)
	id, err := sendSMS(e.AccountID, message, "payment-"+e.ID)
	if err != nil {
		return entry, err
	}
	entry.MessageID = id
	return entry, nil
}
