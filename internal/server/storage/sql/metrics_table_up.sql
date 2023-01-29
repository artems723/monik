CREATE TABLE IF NOT EXISTS "metrics" (
    id serial PRIMARY KEY,
    name text NOT NULL,
    type text NOT NULL,
    delta bigint,
    value double precision,
    CONSTRAINT unique_name_index UNIQUE (name)
);