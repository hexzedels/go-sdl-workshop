# GoSDLWorkshop — curl cheatsheet

Every endpoint as a copy-pasteable curl. Default baseURL: `http://localhost:8080`.

## Setup

```bash
BASE=http://localhost:8080
TOKEN=$(curl -s -X POST $BASE/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"user","password":"password123"}' | jq -r .token)
echo "$TOKEN"
```

Swap `user`/`password123` for `admin`/`admin123` for the admin account.

---

## Auth

### Login

```bash
curl -s -X POST $BASE/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"user","password":"password123"}'
```

Response includes `token`, `session_id`, `user`.

---

## Health

```bash
curl -s $BASE/health
```

---

## Documents (JWT required)

### List own documents

```bash
curl -s -H "Authorization: Bearer $TOKEN" $BASE/api/documents
```

### Search by title

```bash
curl -sG -H "Authorization: Bearer $TOKEN" \
  --data-urlencode 'q=project' $BASE/api/documents/search
```

### Get by id

```bash
curl -s -H "Authorization: Bearer $TOKEN" $BASE/api/documents/1
```

### Create

```bash
curl -s -X POST -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"Hello","content":"world","locale":"en"}' \
  $BASE/api/documents
```

---

## Files (JWT required)

### Download

```bash
curl -sG -H "Authorization: Bearer $TOKEN" \
  --data-urlencode 'name=example.txt' $BASE/api/files
```

### Upload

```bash
curl -s -X POST -H "Authorization: Bearer $TOKEN" \
  -F "file=@./localfile.txt" $BASE/api/files
```

---

## Webhooks (API key required)

```bash
API_KEY=workshop-api-key-f8c2d1a0
curl -s -X POST -H "X-API-Key: $API_KEY" \
  -H 'Content-Type: application/json' \
  -d '{"urls":["https://example.com/hook"]}' \
  $BASE/api/webhooks/notify
```

---

## Admin (JWT with `role=admin` required)

```bash
ADMIN_TOKEN=$(curl -s -X POST $BASE/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123"}' | jq -r .token)
```

### Sessions

```bash
curl -s -H "Authorization: Bearer $ADMIN_TOKEN" $BASE/api/admin/sessions
```

### Audit

```bash
curl -s -H "Authorization: Bearer $ADMIN_TOKEN" $BASE/api/admin/audit
```

---

## Debug (default build)

Available unless built with `-tags nodebug`.

```bash
curl -s $BASE/debug/pprof/
curl -s "$BASE/debug/pprof/goroutine?debug=1" | head -50
curl -s "$BASE/debug/pprof/profile?seconds=5" -o cpu.pprof
```
