# Release Tracker — REST API Contract

**Version:** 0.1 (Draft)  
**Base URL:** `http://localhost:8080`  
**API Prefix:** `/v1`  
**Resource:** `releases`  
**Backend:** Go  
**Database:** PostgreSQL

## 1. Overview

Release Tracker is a REST API for managing movie releases.

For the MVP, a release represents a movie with basic metadata and a release date.

The API supports five operations:

- Retrieve all movies.
- Retrieve a movie by its ID.
- Create a new movie.
- Update an existing movie.
- Delete an existing movie.

### 1.1 Scope

**Included:**
- CRUD operations for movies
- Request validation
- PostgreSQL persistence
- JSON responses
- HTTP error handling

**Excluded from MVP:**
- Authentication and authorization
- External movie APIs
- Notifications
- Scheduled refresh jobs
- Pagination, filtering, and sorting
- Frontend-specific business logic

## 2. Resource Model

### Movie

| Field | Type | Required | Description |
|---|---|---|---|
| `id` | integer | Server-generated | Unique movie identifier |
| `title` | string | Yes | Movie title |
| `description` | string | No | Optional movie description |
| `release_date` | string | Yes | Movie release date, `YYYY-MM-DD` |
| `created_at` | string | Server-generated | Creation timestamp in UTC |
| `updated_at` | string | Server-generated | Last update timestamp in UTC |

**Example movie:**

```json
{
  "id": 42,
  "title": "Interstellar",
  "description": "A science fiction film",
  "release_date": "2014-11-07",
  "created_at": "2026-11-17T18:00:00Z",
  "updated_at": "2026-11-17T18:00:00Z"
}
```

The `id`, `created_at`, and `updated_at` fields are managed by the server and must not be supplied when creating or replacing movies.

### 2.1 Validation Rules

| Field | Validation |
|---|---|
| `id` | Positive integer |
| `title` | Required, nonempty after trimming whitespace, maximum 255 characters |
| `description` | Optional, maximum 2,000 characters |
| `release_date` | Required, valid calendar date in `YYYY-MM-DD` format |

Additional rules:

- Requests with bodies must contain valid JSON.
- Unknown JSON fields are rejected.
- Missing required fields are rejected.
- Client-supplied server-managed fields are rejected.
- Frontend validation does not replace backend validation.
- No uniqueness constraint is imposed on movie titles.

An omitted description is represented as an empty string in the API response.

## 3. Response Conventions

### 3.1 Successful Responses

Successful requests that return a representation use a `data` envelope.

```json
{
  "data": {}
}
```

Collection responses contain an array:

```json
{
  "data": []
}
```

A successful DELETE returns `204 No Content`, without a response body.

### 3.2 Error Responses

All errors use the following structure:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message"
  }
}
```

### Standard Error Codes

| HTTP Status | Error Code | Meaning |
|---|---|---|
| 400 | `INVALID_REQUEST` | Malformed JSON, invalid ID, or invalid fields |
| 404 | `MOVIE_NOT_FOUND` | Requested movie does not exist |
| 500 | `INTERNAL_SERVER_ERROR` | Unexpected backend or database failure |

Internal database details, SQL statements, and stack traces must not be exposed to clients.

### 3.3 Headers

All JSON requests should use:

`Content-Type: application/json`

JSON responses must include:

`Content-Type: application/json`

## 4. API Endpoints

### 4.1 GET /v1/releases

**Description:** Retrieve all movies stored in the database.

**Request body:** None.

**Success: 200 OK**

```json
{
  "data": [
    {
      "id": 42,
      "title": "Interstellar",
      "description": "A science fiction film",
      "release_date": "2014-11-07",
      "created_at": "2026-11-17T18:00:00Z",
      "updated_at": "2026-11-17T18:00:00Z"
    }
  ]
}
```

If no movies exist:

```json
{
  "data": []
}
```

**Failure: 500 Internal Server Error**

Triggered when an unexpected backend or database failure prevents the list from being retrieved.

**Behavior:**
1. Frontend sends GET request.
2. Backend queries PostgreSQL for movies.
3. Database returns records or an error.
4. Backend responds with the appropriate status and JSON body.

Movie ordering is not guaranteed in v0.1.

---

### 4.2 GET /v1/releases/{id}

**Description:** Retrieve a specific movie using its ID.

**Path parameters:**

| Parameter | Type | Required |
|---|---|---|
| `id` | Positive integer | Yes |

**Example request:**

`GET /v1/releases/42`

**Success: 200 OK**

```json
{
  "data": {
    "id": 42,
    "title": "Interstellar",
    "description": "A science fiction film",
    "release_date": "2014-11-07",
    "created_at": "2026-11-17T18:00:00Z",
    "updated_at": "2026-11-17T18:00:00Z"
  }
}
```

**Failure cases:**

- `400 Bad Request`: ID format is invalid.
- `404 Not Found`: No movie matches the requested ID.
- `500 Internal Server Error`: Unexpected database or server failure.

**Example 404 response:**

```json
{
  "error": {
    "code": "MOVIE_NOT_FOUND",
    "message": "Movie with the specified ID was not found."
  }
}
```

---

### 4.3 POST /v1/releases

**Description:** Create a new movie.

**Request body:**

```json
{
  "title": "Interstellar",
  "description": "A science fiction film",
  "release_date": "2014-11-07"
}
```

**Success: 201 Created**

```json
{
  "data": {
    "id": 42,
    "title": "Interstellar",
    "description": "A science fiction film",
    "release_date": "2014-11-07",
    "created_at": "2026-11-17T18:00:00Z",
    "updated_at": "2026-11-17T18:00:00Z"
  }
}
```

**Response header:**

`Location: /v1/releases/42`

**Failure cases:**

- `400 Bad Request`: Malformed JSON or failed validation.
- `500 Internal Server Error`: Unexpected database insertion failure.

**Example 400 response:**

```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "The title field is required."
  }
}
```

**Behavior:**
1. Receive JSON request.
2. Decode and validate fields.
3. Reject invalid requests with 400.
4. Insert movie into PostgreSQL.
5. Retrieve the generated ID and timestamps.
6. Return 201 with the created movie.

---

### 4.4 PUT /v1/releases/{id}

**Description:** Replace the editable fields of an existing movie.

PUT uses full replacement semantics for client-editable fields. All required fields must be supplied, and omitted optional fields reset to their defaults.

**Example request:**

`PUT /v1/releases/42`

**Request body:**

```json
{
  "title": "Interstellar (Updated)",
  "description": "Updated description",
  "release_date": "2014-11-07"
}
```

**Success: 200 OK**

```json
{
  "data": {
    "id": 42,
    "title": "Interstellar (Updated)",
    "description": "Updated description",
    "release_date": "2014-11-07",
    "created_at": "2026-11-17T18:00:00Z",
    "updated_at": "2026-11-17T19:00:00Z"
  }
}
```

**Failure cases:**

- `400 Bad Request`: Invalid ID or request body.
- `404 Not Found`: Movie does not exist.
- `500 Internal Server Error`: Unexpected database update failure.

**Behavior:**
1. Validate ID and JSON payload.
2. Execute an UPDATE operation.
3. Determine whether a matching movie was updated.
4. Return 404 when no matching movie exists.
5. Return 200 with the updated movie on success.

**Constraint:** PUT does not create a new movie if the requested ID does not exist.

---

### 4.5 DELETE /v1/releases/{id}

**Description:** Delete an existing movie.

**Example request:**

`DELETE /v1/releases/42`

**Request body:** None.

**Success: 204 No Content**

No JSON response body is returned.

**Failure cases:**

- `400 Bad Request`: Invalid ID format.
- `404 Not Found`: Movie does not exist.
- `500 Internal Server Error`: Unexpected database deletion failure.

**Behavior:**
1. Validate the movie ID.
2. Execute the DELETE operation.
3. Check whether a matching row was deleted.
4. Return 204 on success.
5. Return 404 if no matching row exists.

## 5. HTTP Status Code Summary

| Endpoint | 200 | 201 | 204 | 400 | 404 | 500 |
|---|---|---|---|---|---|---|
| GET /releases | Yes | | | | | Yes |
| GET /releases/{id} | Yes | | | Yes | Yes | Yes |
| POST /releases | | Yes | | Yes | | Yes |
| PUT /releases/{id} | Yes | | | Yes | Yes | Yes |
| DELETE /releases/{id} | | | Yes | Yes | Yes | Yes |

All paths in this table use the `/v1` prefix.

## 6. Database Expectations

PostgreSQL serves as the persistent source of truth.

The backend is responsible for translating database outcomes into HTTP responses.

| Database outcome | API behavior |
|---|---|
| Query succeeds and returns records | Return requested data |
| List query succeeds with zero records | `200 OK` with empty array |
| Single-record query returns no row | `404 Not Found` |
| UPDATE affects zero rows | `404 Not Found` |
| DELETE affects zero rows | `404 Not Found` |
| Unexpected database error | `500 Internal Server Error` |

The API must not report a successful write before PostgreSQL confirms that the write has succeeded.

## 7. Implementation Expectations

- All endpoints must support JSON response serialization.
- POST and PUT must validate request bodies before writing to the database.
- All SQL queries must be parameterized.
- Database failures must be logged internally.
- User-facing errors must not expose sensitive implementation details.
- Database operations should use request-scoped contexts for cancellation.
- Each endpoint must have tests covering success, validation errors where applicable, not-found outcomes, and database failures.
- Unexpected server failures return a standardized 500 response.

## 8. Open Decisions

The following decisions can be revisited before implementation:

- **ID strategy:** Currently positive integer IDs, backed by PostgreSQL-generated BIGINT values.
- **Date representation:** `release_date` uses a date-only format, not a timestamp.
- **Future pagination:** List endpoint returns all records in MVP; pagination may be added later.
- **Conflict errors:** No explicit `409 Conflict` scenario is defined because the current movie model has no unique business fields.
- **Rate limits and request size limits:** To be defined when operational requirements are established.

## 9. Acceptance Criteria

The API MVP is complete when:

- All five routes are implemented.
- Movies are persisted in PostgreSQL.
- Request validation follows this contract.
- JSON responses match the documented schemas.
- HTTP responses distinguish client errors from server failures.
- CRUD operations have automated success and failure tests.
- The API can be exercised independently of the frontend.

**Contract status:** Draft v0.1  
**Next step:** Define the database schema and map each endpoint to its Go handler and SQL operation.
