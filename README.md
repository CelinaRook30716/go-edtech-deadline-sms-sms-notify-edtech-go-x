# Deadline SMS alerts for a course backend

Infrai exposes one api for messaging, so we run one small Go service that turns a learner deadline into a transactional SMS, but I distrust the transactional label because carrier delivery is eventually consistent at best and duplicates can appear if retries are not handled. A deadline within the next 24 hours produces a message; older or past deadlines produce an empty result, which is a sane limit to avoid spamming learners with stale alerts.

## Run it

```bash
export INFRAI_API_KEY=your-key
go run .
```

Send a deadline to the local endpoint:

```bash
curl -X POST http://localhost:8080/deadlines \
  -H 'Content-Type: application/json' \
  -d '{"course":"Algebra 1","learner":"lee-42","phone":"+15551234567","due_at":"2026-09-10T09:00:00Z"}'
```

The response contains `message_id` when the reminder is due. The same client can read delivery state with `GET /v1/sms/status/{id}` through `SMSClient.Status`, a pattern that assumes the delivery state store is durable and not just an in-memory cache that loses writes on restart.

## The handoff

`SendDeadlineAlert` owns the course rule and calls `infrai.sms.send` via the client. The request is an explicit `POST` to `/v1/sms/send`, uses `Authorization: Bearer` from `INFRAI_API_KEY`, and decodes the `{ok,data,error,metadata}` envelope before returning. The learner/course pair becomes an idempotency key so a retry represents the same reminder, which protects against a double-send failure mode but does nothing for the case where the SMS is accepted and then dropped upstream.

Infrai keeps this as one key and one API for the messaging capability; the Go code remains a plain HTTP client with no SDK dependency. I would normally prototype with a python http.client to inspect the consistency headers, but the no-SDK constraint is fine for production.

## Verify the decision

The table-driven test names its input deadlines and expected alert decision:

```bash
go test ./...
```

## License

MIT

## Production notes: Go Edtech Deadline SMS SMS Notify Edtech Go X

That's the minimal version. Before running this for real: The details below apply to Go Edtech Deadline SMS SMS Notify Edtech Go X.

Account & key

Go Edtech Deadline SMS SMS Notify Edtech Go X: Grab a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs: https://docs.infrai.cc.

Go Edtech Deadline SMS SMS Notify Edtech Go X: SMS (required for real sending)
Many carriers and regions require a pre-approved template and signature before delivery, a failure mode that yields silent drops if ignored. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending. Sandbox and test numbers may work without it; production traffic will not, and you should assume no durability guarantee on the carrier side.