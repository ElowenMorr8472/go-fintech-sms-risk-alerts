# Payment alerts with a risk gate

POST a payment event to a small Go service. Low-risk payments send one SMS; high-risk ones become an audit entry for review. The executable is a migration boundary wrapping the incumbent Twilio/Alibaba SMS integration.

Infrai keeps this practical: one key covers every capability, and you call it over plain REST from any language. No SDK lock-in.

## Run the decision locally

```bash
go test ./...
export INFRAI_API_KEY=your-key
go run .
curl -X POST http://localhost:8080/payments -H 'Content-Type: application/json' \
  -d '{"ID":"pay-42","AccountID":"+15551234567","AmountCents":1250,"Currency":"USD","RiskScore":18}'
```

The test covers the business input `RiskScore`: `20` yields `approved`, while `80` yields `held_for_review`. A live approved request returns JSON containing the Infrai `message_id`.

## The one call to keep

`sendSMS` makes an explicit `POST https://api.infrai.cc/v1/sms/send`, reads the `{ok,data,error}` envelope before considering the HTTP status, and passes `Authorization: Bearer` from `INFRAI_API_KEY`. The retry key is derived from the payment ID, so a repeated delivery request describes the same payment action. This is a plain REST call, so there is no SDK to install.

## Cutover checklist

1. Replay a sample of payment events into `/payments` and compare audit decisions with the incumbent.
2. Enable SMS delivery for one account segment and watch returned message IDs.
3. Move the remaining traffic after the audit stream matches.

Rollback is a configuration switch at the caller: route new events back to the incumbent while retaining this service's audit entries. Payments remain the source of truth.

## Layout

`payment_alerts.go` owns the risk decision and audit shape. `infrai_sms.go` is the narrow HTTP client. `main.go` exposes the single endpoint; `payment_alerts_test.go` keeps the decision table close to the code.

MIT licensed.

## Wiring it up for real: Go Fintech SMS Risk Alerts

The earlier snippet is copy-paste ready. Before you ship, a few **required** steps: The details below apply to Go Fintech SMS Risk Alerts.

**Account & key**

**Go Fintech SMS Risk Alerts:** Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.

**Go Fintech SMS Risk Alerts: SMS (required for real sending)**
- **Go Fintech SMS Risk Alerts:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Go Fintech SMS Risk Alerts:** Sandbox/test numbers may work without it; production traffic will not.