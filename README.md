# events-email-service

Simple HTTP service to send emails via SMTP.

- Language: Go
- Web framework: gin
- Port: 8082

## Build

Prerequisites: Go 1.20+ (or the version in the Dockerfile).

Build locally:
```
go build -o events-email-service .
./events-email-service
```

Run with environment variables (example, Windows PowerShell):
```
$env:SMTP_HOST="smtp.example.com"; $env:SMTP_PORT="587"; $env:SMTP_USER="you@example.com"; $env:SMTP_PASS="yourpassword"; ./events-email-service
```

## Docker

Build the image:
```
docker build -t events-email-service .
```

Run with explicit env vars:
```
docker run --rm -p 8082:8082 \
  -e SMTP_HOST=smtp.example.com \
  -e SMTP_PORT=587 \
  -e SMTP_USER=you@example.com \
  -e SMTP_PASS=yourpassword \
  events-email-service
```

Or use an env file:
```
docker run --rm -p 8082:8082 --env-file .env events-email-service
```

## Example environment variables (.env)
```
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=you@example.com
SMTP_PASS=supersecret
```
Do not commit real credentials.

## API

POST /send-email
- Content-Type: application/json
- Body:
```json
{
  "to": "recipient@example.com",
  "subject": "Test",
  "body": "Hello!"
}
```
- Success: 200 {"message":"Email enviado com sucesso!"}
- Bad request: 400 (validation error)
- Server error: 500 (send failure)

Health check
- The repository does not include a /health endpoint by default. Add one if needed for readiness/liveness probes.

## Notes
- Uses gomail for SMTP delivery and godotenv for local env file loading.
- For testing, prefer a sandbox SMTP service (Mailtrap, Ethereal) instead of real inbox credentials.