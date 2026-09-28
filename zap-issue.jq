# ZAP JSON レポートから riskcode >= 2 (Medium 以上) を抽出し、1 アラート = 1 行の JSON {key,title,body} を出力する
# key = "[ZAP <alertRef>]" を Issue タイトルの先頭に付け、これを冪等キーとして重複登録を防ぐ
# usage: jq -cf zap-issue.jq zap/report.json
.site[].alerts[]
| select((.riskcode | tonumber) >= 2)
| "[ZAP \(.alertRef)]" as $key
| {
    key: $key,
    title: "\($key) \(.riskdesc | split(" ")[0]): \(.name)",
    body: (
      "**Risk**: \(.riskdesc)\n\n"
      + (.instances | map("- `\(.method)` \(.uri)") | join("\n"))
      + "\n\n<details><summary>solution</summary>\n\n\(.solution | gsub("<[^>]*>"; ""))\n\n</details>\n"
    )
  }
