## 1. Documentation Structure & Header Setup

- [x] 1.1 Add status badges (Go 1.26+, Chi v5, PostgreSQL, Prisma Go Client, Kafka, Redis, Cloudflare R2) and introductory microservice summary
- [x] 1.2 Add navigable Table of Contents with emoji headers linking to all sections

## 2. Architecture & Domain Workflow Documentation

- [x] 2.1 Add Mermaid architecture flowchart illustrating Gateway offloading, Chi router, Redis, Kafka consumer, Prisma PostgreSQL, and Cloudflare R2
- [x] 2.2 Add Mermaid inventory synchronization lifecycle diagram detailing order event handling, stock adjustments, and cache invalidation
- [x] 2.3 Document Key Features with bulleted summaries

## 3. Configuration & Developer Setup Guides

- [x] 3.1 Update Technology Stack list and ASCII repository structure tree
- [x] 3.2 Update Prerequisites and complete `.env` environment variables table
- [x] 3.3 Detail Database Setup with Prisma ORM (`db push` and client generation) and Local Development Quickstart

## 4. API Catalog & Integration Mechanics

- [x] 4.1 Document interactive OpenAPI 3.1 documentation links (Scalar UI and Swagger UI) and complete API Endpoint Catalog table
- [x] 4.2 Document Kafka Pipeline (`order.events` schemas, consumer group, idempotency ledger in `processed_events`)
- [x] 4.3 Document Redis Caching strategies and Sliding-Window Rate Limiting policies with fail-open circuit breaker behavior
- [x] 4.4 Document Cloudflare R2 direct presigned image upload workflow
- [x] 4.5 Document Testing execution commands and Docker container deployment instructions
