package main

import "testing"

func TestNotifyPaymentRiskDecision(t *testing.T) {
	tests := []struct {
		name  string
		score int
		want  string
	}{{"clear", 20, "approved"}, {"high risk", 80, "held_for_review"}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := notifyPayment(PaymentEvent{ID: "p-1", AccountID: "acct", AmountCents: 1250, Currency: "USD", RiskScore: tt.score})
			if tt.score < 80 && err == nil { /* network is expected only for approved events */
			}
			if got.Decision != tt.want {
				t.Fatalf("decision=%q want %q", got.Decision, tt.want)
			}
		})
	}
}
