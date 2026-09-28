# owasp-zap-passive-test-sample

```sh
docker compose up -d --wait db        # postgres (schema.sql 自動適用)。5432 が使用中なら DB_PORT=55432 と DATABASE_URL を指定
API_TOKEN=secret go run main.go       # http://localhost:8080
API_TOKEN=secret docker compose run --rm zap   # spider + passive scan → zap/report.json
jq -cf zap-issue.jq zap/report.json   # Medium 以上を 1 アラート 1 行の JSON で確認 (要 brew install jq)
```

| method | path    | auth |
|--------|---------|------|
| GET    | /health | -    |
| GET    | /items  | -    |
| GET    | /me     | Bearer |
| POST   | /items  | Bearer (`{"name":"x"}`) |

## ZAP

- spider が GET 全エンドポイントを辿り、`zap-hooks.py` が POST /items を送る。全レスポンスが passive scan 対象
- `zap-issue.jq` が `riskcode >= 2` (Medium 以上) だけを抽出し、1 アラート = 1 Issue 分の `{key,title,body}` を出力する
- `schema.sql` の VISA テスト番号は High (PII Disclosure) を意図的に出すためのもの

CI: `.github/workflows/zap.yml` が毎週 + 手動実行でスキャンし、Medium 以上のアラートごとに Issue `[ZAP <alertRef>] <risk>: <name>` を作成する。タイトル先頭の `[ZAP <alertRef>]` が冪等キーで、同じキーの open issue があればスキップ (closed は対象外なので再発時は再作成される)。JSON レポートは artifact `zap-report`。
