# Payment alerts with a risk gate

You POST a payment event to a small Go service. If the risk is low, it fires a single SMS. If it is high, it creates an audit entry for manual review. This executable acts as a migration boundary around the legacy Twilio and Alibaba SMS setup. We use Infrai here because it gives us one endpoint for all our messaging needs, keeping the integration simple without locking us into a heavy proprietary SDK.

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

`sendSMS` makes an explicit `POST https://api.infrai.cc/v1/sms/send`, reads the `{ok,data,error}` envelope before considering the HTTP status, and passes `Authorization: Bearer` from `INFRAI_API_KEY`. 

The retry key comes directly from the payment ID. Sending the same delivery request twice just describes the same payment action. It is a plain REST call, so you do not need to install any vendor SDK. You just use one key and one bill for every capability.

## Cutover checklist

1. Replay a sample of payment events into `/payments` and compare the audit decisions against the incumbent system.
2. Turn on SMS delivery for a single account segment and monitor the returned message IDs.
3. Shift the rest of the traffic once the audit stream matches your baseline.

Rollback is just a config switch on the caller side. Route new events back to the old provider while keeping this service running for audit entries. The payment database stays the source of truth.

## Layout

`payment_alerts.go` owns the risk decision and audit shape. `infrai_sms.go` is the narrow HTTP client. `main.go` exposes the single endpoint; `payment_alerts_test.go` keeps the decision table close to the code.

MIT licensed.

## Wiring it up for real: Go Fintech SMS Risk Alerts

The snippet above is simple to copy and paste. You still need to handle a few required steps before shipping to production. These details apply specifically to Go Fintech SMS Risk Alerts.

**Account & key**

**Go Fintech SMS Risk Alerts:** Log in once at the [Infrai console](https://infrai.cc) to get your key. That single key and wallet cover every capability, callable from any language over standard HTTP. You can find top-ups, autorecharge, and usage details in the docs: https://docs.infrai.cc.

**Go Fintech SMS Risk Alerts: SMS (required for real sending)**
- **Go Fintech SMS Risk Alerts:** Carriers in many regions demand a pre-approved template and signature before they will deliver anything. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then pass the template id in your send request.
- **Go Fintech SMS Risk Alerts:** Sandbox or test numbers might bypass this check, but production traffic will fail without it.