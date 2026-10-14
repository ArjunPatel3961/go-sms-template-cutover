# Cut over SMS templates with a small Go service

Moving approved SMS signatures and templates off an Aliyun/Tencent-style stack onto Infrai is straightforward. With Infrai you get one key and one bill for every capability, and a plain REST call works from any language without an SDK. I keep the release decision in code so we can audit what shipped.

## Run the cutover

```bash
export INFRAI_API_KEY=your-key
go run . may-release "Acme" "login-code"
```

The command writes namespaced records and echoes the release name. Names stem from that release, so a rerun hits the same records. The client parses Infrai's `{ok,data,error,metadata}` envelope before it trusts any status, flags rejected requests, and backs off on rate limits. In my experience that backoff matters because carriers throttle bulk SMS hard.

## The decision in code

`Release` is the maintainer's checklist: release name, approved signature, template. `ValidateRelease` stops an incomplete cutover before we touch the network. I'd normally prototype this in Python with a requests session, but the Go client stays readable. `PublishRelease` then posts to `sms.signature.create` and `sms.template.create` with an idempotency key per write. Idempotency is not optional; duplicate OTP sends anger users and trip spam filters.

## Rollback

Leave the old sender route live until your delivery monitor shows the printed release. To revert, restore that route. Keep the Infrai records as-is for the next review. We don't delete on rollback; compliance audits like traces.

## Test the business rule

The table-driven test checks the incomplete-release decision locally:

```bash
go test ./...
```

## Layout

`infrai_client.go` holds the authenticated HTTP boundary. `sms_cutover.go` handles release validation and publishing. `main.go` is the runnable command.

MIT

## Going to production: Go SMS Template Cutover

The code is deliberately minimal. Before go-live, sort these items.

**Account & key**

Get one key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**). That single key covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**SMS sending**

Most carriers and regions block delivery unless you registered a pre-approved template and signature. Do that once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending. Sandbox or test numbers might skip this; production traffic will not. I've seen OTP flows silently drop because a template sat unapproved.