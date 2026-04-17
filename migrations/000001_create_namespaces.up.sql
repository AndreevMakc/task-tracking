-- migrations/000001_create_namespaces.up.sql

BEGIN;

CREATE TABLE IF NOT EXISTS namespaces (
    id         BIGSERIAL   PRIMARY KEY,
    name       TEXT        NOT NULL UNIQUE,
    task_seq   BIGINT      NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

INSERT INTO namespaces (name) VALUES ('default')
    ON CONFLICT (name) DO NOTHING;

CREATE OR REPLACE FUNCTION protect_default_namespace()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.name = 'default' THEN
        RAISE EXCEPTION 'default namespace cannot be modified or deleted';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_protect_default_namespace
    BEFORE UPDATE OR DELETE ON namespaces
    FOR EACH ROW EXECUTE FUNCTION protect_default_namespace();

COMMIT;