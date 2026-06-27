ALTER TABLE refresh_tokens
ADD CONSTRAINT unique_token_hash
UNIQUE(token_hash);