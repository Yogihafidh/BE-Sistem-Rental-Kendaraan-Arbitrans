# ARBITRANS — Product Requirements Document (PRD)

This document defines the functional scope, actors, business rules, and acceptance criteria for the **ARBITRANS** vehicle rental management backend system.

---

## 1. Product Overview

### 1.1 Product Vision
ARBITRANS is a backend system designed to manage the end-to-end operational lifecycle of commercial vehicle rentals. It centralizes fleet tracking, customer identity verification, reservation scheduling, pricing calculations, status history, and security auditing into a unified platform.

### 1.2 Business Problem
Vehicle rental operations often suffer from:
- **Double bookings** caused by race conditions during simultaneous reservations.
- **Fraud and identity risk** from unverified customer identities or forged documents.
- **Fragmented scheduling** where vehicle availability, maintenance downtime, and rental handoffs are not tracked in real time.
- **Lack of auditability** when booking statuses, pricing, or vehicle assignments are modified.

### 1.3 Main Objectives
1. Provide accurate, real-time vehicle availability tracking with zero double-booking tolerance.
2. Automate KYC (Know Your Customer) identity document collection and verification workflows.
3. Manage complete reservation lifecycles from initial request to vehicle return.
4. Maintain an immutable audit log of sensitive system mutations.

### 1.4 Scope and Boundaries
- **In-Scope (Target MVP)**:
  - User and role management (RBAC).
  - Customer profile management and KYC document storage via MinIO.
  - Vehicle catalog, category management, and vehicle media.
  - Booking reservations, scheduling conflict detection, itemization, and status histories.
  - Structured audit logging.
- **Out-of-Scope (Current Phase)**:
  - Automated payment gateway integration (Midtrans/Xendit/Stripe) — planned for future phase.
  - Real-time GPS vehicle telematics and tracking.
  - Multi-tenant SaaS capability (the system operates as a single-tenant deployment for PT Arbitrans).

---

## 2. Implementation Status: Actual vs. Planned

To maintain transparency between what exists in the codebase today and what is planned for feature development:

| Component / Layer | Current Code State | Actual Status |
| :--- | :--- | :--- |
| **Database Migrations** | Complete DDL in `database/migrations/` (Extensions, Schema, Booking Source Enum) | **Implemented & Verified** |
| **SQL Queries (sqlc)** | Query files in `database/queries/*.sql` & compiled code in `database/sqlc/` | **Implemented & Verified** |
| **Shared Infrastructure** | PostgreSQL pool, Redis client, MinIO client, JWT, logging, middlewares, validator | **Implemented & Verified** |
| **Users & Auth Module** | No handlers, DTOs, or services in `internal/modules/auth` | **Planned (Not Yet Implemented)** |
| **Customers Module** | No handlers, DTOs, or services in `internal/modules/customer` | **Planned (Not Yet Implemented)** |
| **Vehicles Module** | No handlers, DTOs, or services in `internal/modules/vehicle` | **Planned (Not Yet Implemented)** |
| **Bookings Module** | No handlers, DTOs, or services in `internal/modules/booking` | **Planned (Not Yet Implemented)** |
| **Audit Log Module** | Low-level helper in `internal/shared/audit`, but no business service or handler | **Planned (Not Yet Implemented)** |
| **Payments Integration** | Schema placeholder comment in `database/queries/payments.sql` | **Planned (Future Scope)** |

> *Rule: Having a table in PostgreSQL or a query in sqlc does NOT mean a business capability is implemented.*

---

## 3. Actors and Responsibilities

| Actor | Role Identifier | Responsibilities | Permissions & Capabilities |
| :--- | :--- | :--- | :--- |
| **Administrator** | `admin` | Overall system oversight, staff management, fleet configuration. | Full read/write across all modules; user account creation; role assignment; system audit inspection. |
| **Staff / Operator** | `staff`, `operator` | Daily operational desk: processing bookings, verifying customer documents, updating vehicle status. | Review and approve/reject customer documents; create and update bookings; transition vehicle status (`available`, `maintenance`). |
| **Customer** | `customer` | Vehicle renter renting personal or corporate vehicles. | Register profile; upload KYC documents (KTP, SIM); view vehicle catalog; request bookings. |

> [!NOTE]
> **Needs clarification — Customer Auth Model**: The database schema defines roles via the `roles` table (seeded with `admin`, `staff`, `operator`). Whether customers authenticate through the same `users` table or a separate customer portal authentication flow remains **to be determined**. Currently, `customers` and `users` are separate database tables without a direct login credential link on `customers`.

---

## 4. Functional Requirements

### 4.1 User Management & Access Control
- **User Accounts**: Manage internal operators with `full_name`, `email` (unique), `phone`, and `password_hash` (bcrypt).
- **Activation Status**: Accounts can be activated (`is_active = true`) or deactivated.
- **Soft Deletion**: Users are soft-deleted via `deleted_at`.
- **Role Assignment**: Every user belongs to exactly one role (`role_id REFERENCES roles(id)`).
- **Authentication**: Stateless JWT token issued upon valid email and password credentials.

### 4.2 Customer Management
- **Customer Identity**: Collect `full_name`, `email`, `phone`, `nik` (16-digit national ID), `address`, and `date_of_birth`.
- **Unique Identifiers**: `nik` must be strictly 16 numeric digits and unique across the customer base. `email` must be valid and unique.
- **Verification Status**: Each customer possesses an overall `verification_status`:
  - `pending`: Registered, documents under review.
  - `verified`: Approved by staff; authorized to rent vehicles.
  - `rejected`: Incomplete or fraudulent identification.

### 4.3 Customer Document Management (KYC)
- **Supported Document Types** (`customer_document_type_enum`):
  - `ktp-penyewa` (Renter National ID)
  - `ktp-penjamin` (Guarantor National ID)
  - `sim` (Driver's License)
  - `passport` (Passport for foreign renters)
  - `tiket` (Travel ticket confirmation)
  - `npwp` (Tax ID)
- **Document Storage**: Document files are stored securely in MinIO/S3 object storage, referenced by `file_url`.
- **Verification Audit**: Each document tracks `verified_by` (FK to `users`), `verified_at`, and `uploaded_at`.
- **Constraint**: Exactly one document per `document_type` per customer (`UNIQUE(customer_id, document_type)`). Re-uploading replaces the previous document file.

### 4.4 Vehicle Category Management
- **Categories**: Organize fleet by category (e.g., MPV, SUV, Sedan, Commercial Van) with unique names and descriptions.

### 4.5 Vehicle Fleet Management & Images
- **Vehicle Specifications**: `name`, `brand`, `model`, `year` (1900 to current year + 1), `transmission` (`manual`, `automatic`), `fuel_type` (`gasoline`, `diesel`, `electric`, `hybrid`), `seat_capacity`, `luggage_capacity`, and `plate_number` (unique).
- **Operational Status** (`vehicle_status_enum`):
  - `available`: Ready for reservation.
  - `booked`: Currently on rent.
  - `maintenance`: In workshop for repairs or regular service.
  - `inactive`: Decommissioned from fleet.
- **Media**: Multiple photos stored in MinIO; supports setting a primary photo (`is_thumbnail = true`).

### 4.6 Booking Management
- **Identifiers**: Unique human-readable `booking_code` (e.g., `ARB-20261026-001`).
- **Booking Source** (`booking_source_enum`): `website`, `mobile_app`, `call_center`, `walk_in`, `online`.
- **Logistics**: Pickup and return locations (text), dates (`start_date`, `end_date`), and optional times (`pickup_time`, `return_time`).
- **Constraint**: `end_date >= start_date`.
- **Pricing**: `subtotal` and `total_amount` (numeric $\ge 0$).

### 4.7 Booking Items & Status History
- **Booking Items**: Represents reserved vehicles associated with a booking. Tracks `price_per_day`, `total_days`, and calculated `subtotal`.
- **Status Lifecycle** (`booking_status_enum`): `pending`, `confirmed`, `ongoing`, `completed`, `cancelled`.
- **Audit History**: Every status modification records a `booking_status_histories` entry capturing `changed_by`, `status`, `notes`, and timestamp.

### 4.8 Audit Logging
- **Immutable Log**: Captures state changes across all modules. Records `user_id`, `action` (e.g., `CREATE_BOOKING`, `VERIFY_CUSTOMER`), `entity_type`, `entity_id`, and `old_value` / `new_value` as JSONB payloads.

---

## 5. Core Business Rules & Clarifications

### 5.1 Vehicle Availability & Scheduling Conflicts
A vehicle is considered **unavailable** for a requested date range $[S_{req}, E_{req}]$ if there exists any active booking item for that vehicle where the booking status is `pending`, `confirmed`, or `ongoing`, and the date ranges overlap:

$$\text{Overlap Condition: } \quad S_{existing} \le E_{req} \quad \text{AND} \quad E_{existing} \ge S_{req}$$

- **Conflict Rule**: A booking cannot include a vehicle that fails this overlap check.
- **Concurrency Locking Rule**: Booking creation must lock the vehicle row (`SELECT ... FOR UPDATE`) inside an ACID transaction to prevent two simultaneous requests from passing the check at the exact same millisecond.
- > [!NOTE]
  > **Needs clarification — Turnaround Buffer Time**: Does PT Arbitrans require a mandatory inspection/cleaning buffer between rentals (e.g., 2 hours or 1 calendar day gap), or can Vehicle B be picked up immediately on the same date another rental ends?

### 5.2 Booking Status State Machine
The permitted status transitions for a booking are:

```
              ┌───────────────┐
              │    pending    │
              └──┬─────────┬──┘
                 │         │
      [Confirm]  │         │  [Cancel]
                 ▼         ▼
  ┌───────────────┐       ┌───────────────┐
  │   confirmed   │       │   cancelled   │
  └──┬─────────┬──┘       └───────────────┘
     │         │
[Dispatch]     │ [Cancel]
     │         ▼
     │   ┌───────────────┐
     │   │   cancelled   │
     ▼   └───────────────┘
┌───────────────┐
│    ongoing    │
└──┬────────────┘
   │ [Return]
   ▼
┌───────────────┐
│   completed   │
└───────────────┘
```

- **Transitions**:
  - `pending` $\rightarrow$ `confirmed`: Customer identity is verified and booking is locked.
  - `confirmed` $\rightarrow$ `ongoing`: Vehicle is handed over (dispatched) to the customer; vehicle status updates to `'booked'`.
  - `ongoing` $\rightarrow$ `completed`: Vehicle is returned and inspected; vehicle status returns to `'available'`.
  - `pending` $\rightarrow$ `cancelled` or `confirmed` $\rightarrow$ `cancelled`: Rental cancelled prior to dispatch.
- **Terminal States**: `completed` and `cancelled` cannot be transitioned to any other status.
- > [!NOTE]
  > **Needs clarification — Late Return Handling**: If a customer returns a vehicle after the scheduled `end_date` / `return_time`, does the status remain `ongoing` until formal return, and what penalty rate applies?

### 5.3 Rental Duration and Pricing Calculations
1. **Total Days Calculation**:
   $$\text{total\_days} = (\text{end\_date} - \text{start\_date}) + 1 \quad \text{(or exact 24h intervals)}$$
   > [!NOTE]
   > **Needs clarification — Pricing Duration Unit**: Does PT Arbitrans charge on a calendar-day basis (e.g., Oct 26 to Oct 27 = 2 days) or a 24-hour block basis? The database currently stores `start_date`, `end_date` as `DATE` and optional `pickup_time`, `return_time` as `TIME`.
2. **Item Subtotal**:
   $$\text{subtotal} = \text{price\_per\_day} \times \text{total\_days}$$
3. **Total Amount**:
   $$\text{total\_amount} = \sum \text{subtotals} + \text{additional fees} - \text{discounts}$$
4. > [!NOTE]
   > **Needs clarification — Extra Charges & Deposits**: Does the total amount include driver fees, insurance add-ons, or a refundable security deposit? Currently, only `subtotal` and `total_amount` columns exist.

### 5.4 Cancellation Policy
- > [!NOTE]
  > **Needs clarification — Cancellation & Refund Rules**:
  > - Can customers cancel directly via the app, or must cancellation be processed by staff?
  > - Is there a cutoff window (e.g., 24 hours prior to pickup) for cancellations?
  > - Is there a cancellation fee or deposit forfeiture policy?

### 5.5 Payment Handling
- The current database schema does not yet include a `payments` table (only a placeholder comment exists in `database/queries/payments.sql`).
- > [!NOTE]
  > **Needs clarification — Payment Workflow**:
  > - Is a down payment (DP) required to transition a booking from `pending` to `confirmed`?
  > - What payment methods will be supported initially (Cash on Delivery/Pickup, Manual Bank Transfer, or Payment Gateway like Midtrans)?
  > - Will payments be modeled as a separate `payments` table referencing `bookings(id)`?

### 5.6 Customer Document Verification Rules
- Required documents: At least a valid national ID (`ktp-penyewa` or `passport`) and a driver's license (`sim`).
- > [!NOTE]
  > **Needs clarification — Verification Prerequisites**: Can an unverified customer create a `pending` reservation, or must documents be approved by staff *before* a customer is allowed to book?

---

## 6. Non-Functional Requirements

### 6.1 Security
- **Authentication**: JWT signed with HMAC-SHA256. Minimum 32-character secret in production. Tokens must have bounded expiration (`JWT_EXPIRATION`).
- **Passwords**: Hashed with bcrypt (cost $\ge 10$). Plaintext passwords must never be stored, logged, or returned in API responses.
- **Input Validation**: All incoming requests validated before reaching the service layer.
- **SQL Injection Prevention**: All SQL executed via parameterized sqlc queries. Raw SQL concatenation is forbidden.
- **CORS Protection**: Origin whitelist enforced; wildcard `*` prohibited when `AllowCredentials = true` in production.

### 6.2 Data Integrity
- Referential integrity guaranteed via foreign keys with explicit delete rules (`ON DELETE RESTRICT` for primary entities, `CASCADE` for owned children like `booking_items`).
- ACID transactions enforced for all multi-step operations via `sharedDatabase.WithTransaction`.

### 6.3 Reliability and Performance
- HTTP server configured with bounded timeouts (`ReadTimeout: 10s`, `WriteTimeout: 15s`, `IdleTimeout: 60s`).
- Graceful shutdown handles in-flight requests within a bounded 10-second window.
- Database connection pooling managed via `pgxpool` with tuned connection limits and automatic idle pruning.

### 6.4 Observability
- Structured logging using `log/slog` (Human-readable in development, JSON in production).
- Standard health check probe at `/health` reporting multi-component status (`UP`, `DEGRADED`, or `DOWN`).
