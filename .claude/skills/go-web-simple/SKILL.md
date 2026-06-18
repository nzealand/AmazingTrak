---
name: go-web-simple
description: Use when building or modifying simple Go web apps that use SQLite, Caddy, and systemd. Prefer standard library patterns, direct SQL, environment-variable config, tests, and a /healthz endpoint. Avoid Kubernetes, microservices, Redis, and premature optimization.
---

# Go web app rules

Build simple Go web apps.

## Stack
- Go
- SQLite
- Caddy
- systemd

## Rules
- Standard library first.
- Minimize dependencies.
- Direct SQL.
- Env vars for config.
- Add tests.
- Add /healthz.
- Keep architecture simple.

## Avoid
- Kubernetes
- Microservices
- Redis
- Premature optimization

Choose the simplest solution that works.
