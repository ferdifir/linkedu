# Linkedu — Student Attendance Management SaaS: Design Spec

**Date:** 2026-05-13  
**Stack:** Go + Fiber + GORM + MySQL (backend), Next.js App Router (frontend)  
**Scope:** MVP

---

## 1. Overview

Multi-tenant SaaS for student attendance management. Each tenant is one school. Attendance is taken by teachers using Web NFC API on Android Chrome — students tap RFID cards on the teacher's phone. Parents can submit leave permits with supporting documents. Admin approves or rejects permits.

---

## 2. Multi-Tenancy

**Strategy:** Shared database, shared schema. Every table has a `tenant_id` column.

- GORM middleware auto-injects `WHERE tenant_id = ?` on every query using the JWT claim.
- `super_admin` role has `tenant_id = NULL` and can access all tenants.
- Tenant onboarding: self-registration via landing page → creates `tenant` record + initial `school_admin` account.

---

## 3. Roles

| Role | Scope |
|------|-------|
| `super_admin` | All tenants — platform management |
| `school_admin` | One tenant — master data, schedules, events, permit approval |
| `teacher` | One tenant — open/close attendance sessions, view permits |
| `student` | One tenant — view own schedule and attendance history |
| `parent` | One tenant — view child's attendance, submit and track permits |

---

## 4. Database Schema

```sql
tenants
  id, name, school_name, email, subdomain, status, created_at

users
  id, tenant_id, name, email, password_hash, role, created_at

academic_years
  id, tenant_id, name, start_date, end_date, is_active

classrooms
  id, tenant_id, name, grade_level

subjects
  id, tenant_id, name, code

teachers
  id, tenant_id, user_id, nip, phone

students
  id, tenant_id, user_id, nis, classroom_id, rfid_uid

parents
  id, tenant_id, user_id, student_id

schedules
  id, tenant_id, academic_year_id, classroom_id, subject_id,
  teacher_id, day_of_week, start_time, end_time

events
  id, tenant_id, name, event_type, date, start_time, end_time, description

attendance_sessions
  id, tenant_id, schedule_id (nullable), event_id (nullable),
  teacher_id, opened_at, closed_at, status ENUM(open, closed)

attendance_records
  id, tenant_id, session_id, student_id,
  status ENUM(present, absent, permit), tapped_at, method ENUM(nfc, manual)

permits
  id, tenant_id, student_id, parent_id, session_id (nullable),
  date, reason, attachment_path, status ENUM(pending, approved, rejected),
  reviewed_by, reviewed_at, notes
```

**Key constraints:**
- `attendance_sessions` links to either `schedules` (regular) or `events`, not both.
- `attendance_records` are idempotent per `(session_id, student_id)` — duplicate tap is a no-op.
- Absent records auto-inserted on session close have `tapped_at = NULL`, `method = NULL`.
- When a permit is approved and an absent record exists for that date, the record status updates to `permit`.

---

## 5. Attendance Flow (NFC)

1. Teacher opens attendance page on Chrome Android.
2. Teacher selects schedule/event → clicks "Mulai Presensi".
   - `POST /api/v1/sessions` → creates `attendance_session` (status: open).
3. Teacher taps student RFID card on phone.
   - Web NFC API reads `rfid_uid`.
   - `POST /api/v1/sessions/:id/tap` `{rfid_uid}`.
   - Backend: find student by `rfid_uid + tenant_id`.
   - Unknown card → return error "Kartu tidak dikenal".
   - Known card → upsert `attendance_record` (status: present, method: nfc).
   - Frontend: show student name + realtime checklist.
4. Teacher clicks "Tutup Presensi".
   - `PATCH /api/v1/sessions/:id/close`.
   - Backend: insert absent record for every student in classroom without a record.
5. Students and parents see results on their dashboards.

**Fallback:** Teacher can manually input attendance via `POST /api/v1/sessions/:id/manual`.

---

## 6. Permit Flow

1. Parent opens dashboard → "Ajukan Izin".
2. Fills form: date, reason, uploads supporting document (saved to server local storage).
   - `POST /api/v1/permits`.
3. Admin sees permit badge notification in sidebar.
4. Admin reviews → Approve or Reject with optional notes.
   - `PATCH /api/v1/permits/:id/review` `{status, notes}`.
5. If approved and absent record exists for that date → update `attendance_record.status = permit`.
6. Teacher can view permits for students in their class (read-only).
7. Parent sees updated permit status on dashboard.

**Note:** Permits can be submitted without a specific session (date-based). Decision authority is with admin only.

---

## 7. API Structure

Base prefix: `/api/v1`

```
Auth
  POST /auth/register-tenant
  POST /auth/login
  POST /auth/refresh

Master Data (school_admin)
  CRUD /academic-years
  CRUD /classrooms
  CRUD /subjects
  CRUD /teachers
  CRUD /students
  CRUD /parents
  CRUD /schedules
  CRUD /events

Attendance (teacher)
  POST   /sessions
  PATCH  /sessions/:id/close
  POST   /sessions/:id/tap
  POST   /sessions/:id/manual
  GET    /sessions/:id/records

Permits
  POST   /permits
  GET    /permits
  PATCH  /permits/:id/review

Reports
  GET    /reports/attendance
  GET    /reports/attendance/export

Dashboards
  GET    /dashboard/student/:id
  GET    /dashboard/parent/:id
```

---

## 8. Frontend Pages (Next.js App Router)

```
/ (public)
  /                    — landing + school registration form
  /login

/(admin)
  /admin/dashboard
  /admin/master/classrooms
  /admin/master/subjects
  /admin/master/teachers
  /admin/master/students
  /admin/master/parents
  /admin/schedules
  /admin/events
  /admin/permits
  /admin/reports

/(teacher)
  /teacher/dashboard
  /teacher/attendance  — open session, NFC tap, close session
  /teacher/permits     — read-only

/(student)
  /student/dashboard   — schedule + attendance history

/(parent)
  /parent/dashboard
  /parent/permits      — submit + track
```

---

## 9. Export & Reports

- **PDF:** `gofpdf` — attendance table per classroom/subject/period.
- **Excel:** `excelize` — sheet per classroom, column per date.
- **Filters:** `classroom_id`, `subject_id`, `date_from`, `date_to`.
- **Delivery:** Direct file download — not stored on server.

---

## 10. Project Structure

```
linkedu/
├── backend/
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── middleware/       — JWT auth, tenant injector
│   │   ├── domain/
│   │   │   ├── tenant/
│   │   │   ├── user/
│   │   │   ├── master/       — classroom, subject, teacher, student, parent
│   │   │   ├── schedule/
│   │   │   ├── attendance/   — session, record, NFC handler
│   │   │   ├── permit/
│   │   │   └── report/
│   │   └── shared/           — response helpers, pagination, errors
│   ├── storage/              — uploaded permit attachments
│   └── go.mod
│
└── frontend/
    ├── app/
    │   ├── (public)/
    │   ├── (admin)/
    │   ├── (teacher)/
    │   ├── (student)/
    │   └── (parent)/
    ├── components/
    ├── lib/                  — API client, auth helpers
    └── package.json
```

Each backend domain follows:
```
domain/attendance/
  model.go       — GORM structs
  repository.go  — DB queries
  service.go     — business logic
  handler.go     — Fiber route handlers
```

---

## 11. Auth

- JWT with claims: `user_id`, `tenant_id`, `role`
- Access token: short-lived (15 min)
- Refresh token: long-lived (7 days), stored in httpOnly cookie
- Tenant middleware: extracts `tenant_id` from JWT, sets on request context, GORM scopes use it automatically

---

## 12. Out of Scope (MVP)

- Email / WhatsApp notifications
- Object storage (S3/R2) — local storage for now
- Mobile native app
- Multiple parents per student
- Detailed analytics / charts
- Payment / subscription management
