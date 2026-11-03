---
name: sast-audit
description: Static application security testing (SAST) checklist for Go SQL injection and HTTP middleware hardening.
---

# Go SAST & OWASP Hardening Checklist

1. **No Raw SQL Interpolation (G201/G202):**
   - Never pass `fmt.Sprintf` strings into `QueryContext` or `ExecContext`.
   - Always use positional placeholders (`$1`, `$2`, ...) and pass untrusted values in `args ...any`.
2. **Rate Limiting Middleware:**
   - Guard per-client state with a `sync.Mutex`.
   - Return `http.StatusTooManyRequests` (`429`) when a client's token bucket is exhausted.
