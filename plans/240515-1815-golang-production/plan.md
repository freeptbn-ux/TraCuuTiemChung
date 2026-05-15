# Plan: Golang Production Backend Upgrade
Created: 2026-05-15 18:15
Status: 🟡 In Progress
Project: TraCuuTiemChung (Golang Backend)

## Overview
Nâng cấp bản Golang hiện tại lên chuẩn Production, tập trung vào hiệu năng Serverless (Vercel) và tính ổn định. Giải quyết vấn đề session persistence bằng Redis và gia cố bảo mật.

## Tech Stack
- **Language**: Go 1.25+
- **Framework**: Gin-gonic
- **HTTP Client**: Resty (v2)
- **Database/Cache**: Redis (Upstash recommended for Vercel)
- **Logging**: Structured Logging (slog)

## Phases

| Phase | Name | Status | Progress |
|-------|------|--------|----------|
| 01 | [Foundation & Observability](./phase-01-foundation.md) | ⬜ Pending | 0% |
| 02 | [Redis Session Persistence](./phase-02-redis-session.md) | ⬜ Pending | 0% |
| 03 | [Concurrency & Locking](./phase-03-concurrency.md) | ⬜ Pending | 0% |
| 04 | [API Hardening & Middleware](./phase-04-hardening.md) | ⬜ Pending | 0% |
| 05 | [Full Integration & Testing](./phase-05-integration.md) | ⬜ Pending | 0% |

## Quick Commands
- Bắt đầu Phase 1: `/code phase-01`
- Kiểm tra tiến độ: `/next`
- Lưu context: `/save-brain`
