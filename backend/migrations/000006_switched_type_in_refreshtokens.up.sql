ALTER TABLE refresh_tokens
ALTER COLUMN token_hash TYPE BYTEA USING token_hash::BYTEA;