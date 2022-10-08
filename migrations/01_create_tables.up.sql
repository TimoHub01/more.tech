CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS news (
    id  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    link TEXT NOT NULL,
    topic TEXT NOT NULL,
    text TEXT NOT NULL,
    date TEXT NOT NULL
);