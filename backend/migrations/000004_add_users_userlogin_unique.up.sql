ALTER TABLE users
ADD CONSTRAINT unique_user_login
UNIQUE(user_login);