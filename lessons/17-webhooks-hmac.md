# 17 — Webhooks & HMAC

## Goal

Build the Clerk webhook handler that keeps our local `users` table in
sync, and understand *why* webhook signature verification works the way
it does.

## Concepts

- HMAC (`crypto/hmac`, `crypto/sha256`): proving a request really came
  from Clerk (who shares a secret with us) without needing a shared
  network channel — the receiver recomputes the signature and compares.
- **Constant-time comparison** (`crypto/subtle.ConstantTimeCompare`) —
  why `sig1 == sig2` is a security bug here (timing attacks) and what a
  constant-time comparison actually changes about the code.
- Replay protection: checking the signed timestamp is recent, not just
  that the signature is valid.
- Why we hand-roll this instead of a dependency — recap of the
  svix-webhooks monorepo-weight decision from the reference
  implementation.

*(Task and checkpoint written when we start this lesson.)*
