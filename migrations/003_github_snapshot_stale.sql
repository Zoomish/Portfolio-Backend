-- +goose Up
UPDATE content_snapshot
SET fetched_at = TIMESTAMPTZ '1970-01-01 00:00:00+00'
WHERE source = 'github_repos';

-- +goose Down
SELECT 1;
