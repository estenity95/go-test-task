-- +goose Up
-- +goose StatementBegin
INSERT INTO subscriptions (service_name, price, user_id, start_month, end_month)
VALUES
    ('Yandex Plus', 400, '60601fee-2bf1-4721-ae6f-7636e79a0cba', DATE '2025-07-01', NULL),
    ('Netflix',     999, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', DATE '2025-01-01', DATE '2025-06-01'),
    ('Spotify',     299, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', DATE '2025-03-01', NULL),
    ('YouTube',     249, '60601fee-2bf1-4721-ae6f-7636e79a0cba', DATE '2025-05-01', DATE '2025-09-01'),
    ('Apple Music', 169, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', DATE '2025-04-01', NULL),
    ('Okko',        399, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', DATE '2025-02-01', DATE '2025-05-01'),
    ('Kinopoisk',   299, '60601fee-2bf1-4721-ae6f-7636e79a0cba', DATE '2025-08-01', NULL),
    ('Xbox Game Pass', 599, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', DATE '2025-07-01', NULL),
    ('PS Plus',     749, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', DATE '2025-01-01', NULL),
    ('iCloud',      149, '60601fee-2bf1-4721-ae6f-7636e79a0cba', DATE '2025-06-01', NULL);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM subscriptions
WHERE (service_name, price, user_id, start_month, end_month) IN (
     ('Yandex Plus', 400, '60601fee-2bf1-4721-ae6f-7636e79a0cba', DATE '2025-07-01', NULL),
     ('Netflix',     999, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', DATE '2025-01-01', DATE '2025-06-01'),
     ('Spotify',     299, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', DATE '2025-03-01', NULL),
     ('YouTube',     249, '60601fee-2bf1-4721-ae6f-7636e79a0cba', DATE '2025-05-01', DATE '2025-09-01'),
     ('Apple Music', 169, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', DATE '2025-04-01', NULL),
     ('Okko',        399, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', DATE '2025-02-01', DATE '2025-05-01'),
     ('Kinopoisk',   299, '60601fee-2bf1-4721-ae6f-7636e79a0cba', DATE '2025-08-01', NULL),
     ('Xbox Game Pass', 599, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', DATE '2025-07-01', NULL),
     ('PS Plus',     749, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', DATE '2025-01-01', NULL),
     ('iCloud',      149, '60601fee-2bf1-4721-ae6f-7636e79a0cba', DATE '2025-06-01', NULL)
);
-- +goose StatementEnd
