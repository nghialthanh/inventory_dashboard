# API cURL examples

Client-side stub for the backend assessment: explore the REST API without the React UI. Base URL assumes Compose on port 8080. Full contract: [architecture/api-outline.md](architecture/api-outline.md).

```bash
BASE=http://localhost:8080
CIVIC=4
```

## Health

```bash
curl -s "$BASE/health"
```

## Catalog

```bash
curl -s "$BASE/api/v1/catalog/makes" | jq .
curl -s "$BASE/api/v1/catalog/models?make_id=2" | jq .
```

## Record a vehicle

```bash
curl -s -X POST "$BASE/api/v1/vehicles" \
  -H 'Content-Type: application/json' \
  -d "{\"vin\":\"1HGCM82633A004352\",\"model_id\":$CIVIC,\"entered_on\":\"2026-10-26\"}" | jq .
```

Save the returned `data.id` as `ID` for later calls.

## List and filter

```bash
curl -s "$BASE/api/v1/vehicles?page=1&page_size=20" | jq .
curl -s "$BASE/api/v1/vehicles?aging=true&page=1&page_size=20" | jq .
curl -s "$BASE/api/v1/vehicles?stock_state=in_stock&make_id=2" | jq .
```

## Read one

```bash
curl -s "$BASE/api/v1/vehicles/$ID" | jq .
```

## Correct while in stock

```bash
curl -s -X PATCH "$BASE/api/v1/vehicles/$ID" \
  -H 'Content-Type: application/json' \
  -d '{"entered_on":"2026-07-01"}' | jq .
```

## Proposed action (aging stock only)

```bash
curl -s -X PUT "$BASE/api/v1/vehicles/$ID/proposed-action" \
  -H 'Content-Type: application/json' \
  -d '{"proposed_action":"wholesale_planned"}' | jq .
```

## Stock out

```bash
curl -s -X POST "$BASE/api/v1/vehicles/$ID/stock-out" \
  -H 'Content-Type: application/json' \
  -d '{"stocked_out_on":"2026-10-26"}' | jq .
```
