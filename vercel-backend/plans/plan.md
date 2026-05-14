# Plan: Python to Go Migration (Vercel Backend)
Created: 2026-05-14
Status: 🟡 In Progress
Project: TraCuuTiemChung (Backend)

## 🎯 Overview
Di cư toàn bộ hệ thống từ Python sang Golang. Tập trung vào tính module hóa cao, dễ test và hiệu suất tối đa trên Vercel.

## 🛠 Project Structure (Standard 2026)
```text
.
├── api/                # Vercel entry points (Handlers)
│   └── index.go        # Main entry
├── internal/           # Private business logic
│   ├── portal/         # Portal Client logic
│   ├── analyzer/       # Vaccine rules engine
│   ├── models/         # Shared structs
│   └── config/         # App configuration
├── assets/             # vaccine_rules.json
├── go.mod
└── vercel.json
```

## 📊 Phases Progress

| Phase | Name | Test Status | Progress |
|-------|------|-------------|----------|
| 01 | [Setup Environment](./phase-01-setup.md) | ✅ Pass | 100% |
| 02 | [Portal Service Migration](./phase-02-portal-service.md) | ✅ Pass | 100% |
| 03 | [Analyzer Logic Migration](./phase-03-analyzer-logic.md) | ✅ Pass | 100% |
| 03a| [Base Group Checkers](./phase-03a-group-checkers.md) | ⬜ Pending | 0% |
| 03b| [Specialized Checkers](./phase-03b-specialized-checkers.md) | ⬜ Pending | 0% |
| 03c| [Pneumo Group & Final Parity](./phase-03c-pneumo-final.md) | ⬜ Pending | 0% |
| 04 | [API Handlers & Routing](./phase-04-api-handlers.md) | ⬜ Pending | 0% |
| 05 | [Vercel Config & Cleanup](./phase-05-vercel-cleanup.md) | ⬜ Pending | 0% |

## 🚀 Quick Commands
- Start Migration: `/code phase-01`
- Run Tests: `go test ./...`
- Local Dev: `vercel dev`
