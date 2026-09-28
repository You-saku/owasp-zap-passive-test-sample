# owasp-zap-passive-test-sample

```sh
docker compose up -d --wait db        # postgres (schema.sql 自動適用)。5432 が使用中なら DB_PORT=55432 と DATABASE_URL を指定
API_TOKEN=secret go run .             # http://localhost:8080
API_TOKEN=secret docker compose run --rm zap   # spider + passive scan → zap/report.html
```

| method | path    | auth |
|--------|---------|------|
| GET    | /health | -    |
| GET    | /items  | -    |
| GET    | /me     | Bearer |
| POST   | /items  | Bearer (`{"name":"x"}`) |

CI: `.github/workflows/zap.yml` が毎週 + 手動実行で ZAP baseline を走らせ、結果を Issue「ZAP Scan Baseline Report」に登録/更新する。
