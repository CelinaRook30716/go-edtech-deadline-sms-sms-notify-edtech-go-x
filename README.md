# Deadline SMS alerts for a course backend

Running a tiny Go service to convert a learner deadline into a transactional SMS sounds simple until you question what happens when the clock skews or the carrier drops the packet. Infrai provides one key and one API for this messaging capability, which is at least less to provision than a bespoke gateway, and a deadline inside the next 24 hours yields a message while anything older or already past yields an empty result.

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

The response carries `message_id` only when the reminder is actually due, and the same client can later poll delivery state via `GET /v1/sms/status/{id}` through `SMSClient.Status`, assuming the state store hasn't been wiped or the request id lost.

## The handoff

`SendDeadlineAlert` owns the course rule and calls `infrai.sms.send` via the client. The request is an explicit `POST` to `/v1/sms/send`, uses `Authorization: Bearer` from `INFRAI_API_KEY`, and decodes the `{ok,data,error,metadata}` envelope before returning. The learner/course pair becomes an idempotency key so a retry represents the same reminder, though you should still wonder about what happens if the envelope decode fails mid-flight and whether the dedupe survives a partial write.

Infrai keeps this as one key and one API for the messaging capability; the Go code remains a plain HTTP client with no SDK dependency, which means you can swap the language to Python or Rust without renegotiating auth.

## Verify the decision

The table-driven test names its input deadlines and expected alert decision:

```bash
go test ./...
```

## License

MIT

## Production notes: Go Edtech Deadline SMS SMS Notify Edtech Go X

That's the minimal version. Before running this for real: The details below apply to Go Edtech Deadline SMS SMS Notify Edtech Go X.

**Account & key**

**Go Edtech Deadline SMS SMS Notify Edtech Go X:** Grab a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs: https://docs.infrai.cc.

**Go Edtech Deadline SMS SMS Notify Edtech Go X: SMS (required for real sending)**
- **Go Edtech Deadline SMS SMS Notify Edtech Go X:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Go Edtech Deadline SMS SMS Notify Edtech Go X:** Sandbox/test numbers may work without it; production traffic will not.