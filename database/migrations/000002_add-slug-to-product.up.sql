CREATE EXTENSION IF NOT EXISTS unaccent WITH SCHEMA public;

CREATE OR REPLACE FUNCTION slugify (v text)
    RETURNS text
    AS $$
BEGIN
    -- 1. trim trailing and leading whitespaces from text
    -- 2. remove accents (diacritic signs) from a given text
    -- 3. lowercase unaccented text
    -- 4. replace non-alphanumeric (excluding hyphen, underscore) with a hyphen
    -- 5. trim leading and trailing hyphens
    RETURN trim(BOTH '-' FROM regexp_replace(lower(public.unaccent (trim(v))), '[^a-z0-9\\-_]+', '-', 'gi'));
END;
$$
LANGUAGE PLPGSQL
STRICT IMMUTABLE;

ALTER TABLE products
    ADD slug text NOT NULL GENERATED ALWAYS AS (slugify (name)) STORED;

