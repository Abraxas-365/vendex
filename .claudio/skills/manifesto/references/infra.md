# Infrastructure Packages

## Logging (logx)

Import: `github.com/Abraxas-365/manifesto/internal/logx`

Always use `logx` — never `fmt.Println` or `log.Printf`.

```go
log := logx.New()
log.Info("user created", logx.Field("user_id", id), logx.Field("tenant", tenantID))
log.Error("operation failed", logx.Err(err), logx.Field("op", "CreateUser"))
log.Warn("retrying", logx.Field("attempt", n))
log.WithContext(ctx).Debug("processing request")
```

---

## Job Queue (jobx)

Import: `github.com/Abraxas-365/manifesto/internal/jobx`

Backed by Redis.

```go
// Register handler (at startup)
client.Register("send-email", func(ctx context.Context, job jobx.Job) error {
    var payload EmailPayload
    json.Unmarshal(job.Payload, &payload)
    return sendEmail(payload)
})

// Enqueue
err := client.Enqueue(ctx, jobx.Job{
    Type:    "send-email",
    Payload: payloadBytes,
})

// Start worker loop
client.Start(ctx)
```

---

## File Storage (fsx)

Import: `github.com/Abraxas-365/manifesto/internal/fsx`

`fsx.FileSystem` is an interface. Implementations: local disk and S3. Swap via config/constructor injection.

```go
// Write
err := fs.Write(ctx, "path/to/file.pdf", data)

// Read
data, err := fs.Read(ctx, "path/to/file.pdf")

// Delete
err := fs.Delete(ctx, "path/to/file.pdf")

// Presigned URL (S3 implementation)
url, err := fs.PresignedURL(ctx, "path/to/file.pdf", 15*time.Minute)
```

---

## Email (notifx)

Import: `github.com/Abraxas-365/manifesto/internal/notifx`

```go
err := sender.Send(ctx, notifx.Email{
    To:      []string{"user@example.com"},
    Subject: "Welcome",
    Body:    renderedTemplate,
})
```

Supports HTML templating. Inject `notifx.EmailSender` interface; swap implementations without changing callers.

---

## Config

Import: `github.com/Abraxas-365/manifesto/internal/config`

Load app configuration at startup. All infrastructure clients (DB, Redis, S3, LLM keys) are wired via config into the container.
