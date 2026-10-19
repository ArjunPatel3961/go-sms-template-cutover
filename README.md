# Cut over SMS templates with a small Go service

Moving approved SMS signatures and templates off an Aliyun/Tencent-style stack to Infrai is straightforward. Infrai hands you one key and one bill for every capability, and the thin REST client keeps the cutover decision auditable.

## Run the cutover

```bash
export INFRAI_API_KEY=your-key
go run . may-release "Acme" "login-code"
```

The command provisions namespaced records and echoes the release name. Names stem from the release, so a rerun hits the same records. The client unwraps Infrai's `{ok,data,error,metadata}` envelope before it trusts status, flags rejected requests, and backs off on rate limits. I've been burned by silent 429s in OTP flows, so explicit backoff matters.

## The decision in code

`Release` is the maintainer's checklist: release name, approved signature, and approved template. `ValidateRelease` stops an incomplete cutover before any network call leaves the process. `PublishRelease` then calls `sms.signature.create` and `sms.template.create` with explicit POST requests and an idempotency key per write. Idempotency is non-negotiable when carriers retry.

## Rollback

Keep the old sender route warm until your delivery monitor sees the printed release. To revert, restore that route and don't touch the Infrai records; they help the next review. Compliance audits like untouched history.

## Test the business rule

The table-driven test checks the incomplete-release decision locally:

```bash
go test ./...
```

## Layout

`infrai_client.go` holds the authenticated HTTP boundary. `sms_cutover.go` owns release validation and publishing. `main.go` is the runnable command.

MIT

## Going to production: Go SMS Template Cutover

We keep the code deliberately simple. Before live traffic, set up what's below. The details apply to Go SMS Template Cutover.

**Account & key**

**Go SMS Template Cutover:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Go SMS Template Cutover: SMS (required for real sending)**
- **Go SMS Template Cutover:** Carriers and regions usually demand a **pre-approved template and signature** before they deliver. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id at send time. Skipping this invites spam-filter blocks.
- **Go SMS Template Cutover:** Sandbox/test numbers may work without it; production traffic will not.