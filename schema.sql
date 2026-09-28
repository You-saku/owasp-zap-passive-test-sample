CREATE TABLE items (
    id   BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);
INSERT INTO items (name) VALUES ('sample');
-- ZAP passive scan の High (PII Disclosure) を再現する意図的な漏洩データ (VISA テスト番号)
INSERT INTO items (name) VALUES ('card 4111111111111111');
