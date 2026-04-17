-- migrations/000002_create_tasks.up.sql

BEGIN;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'task_status') THEN
        CREATE TYPE task_status AS ENUM ('New', 'InProgress', 'Done', 'Trashed');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS tasks (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    namespace_id BIGINT      NOT NULL REFERENCES namespaces(id) ON DELETE CASCADE,
    title        TEXT        NOT NULL,
    code         TEXT        NOT NULL DEFAULT ' ',
    status       task_status NOT NULL DEFAULT 'New',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ
    );

CREATE UNIQUE INDEX IF NOT EXISTS idx_tasks_code  ON tasks (code);
CREATE INDEX IF NOT EXISTS idx_tasks_namespace_id ON tasks (namespace_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status       ON tasks (status);
CREATE INDEX IF NOT EXISTS idx_tasks_deleted_at   ON tasks (deleted_at);

CREATE OR REPLACE FUNCTION generate_task_code()
RETURNS TRIGGER AS $$
DECLARE
    v_namespace_name TEXT;
    v_seq            BIGINT;
BEGIN
    UPDATE namespaces
    SET    task_seq = task_seq + 1
    WHERE  id = NEW.namespace_id
    RETURNING name, task_seq INTO v_namespace_name, v_seq;

    NEW.code := v_namespace_name || '-' || v_seq;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_generate_task_code
    BEFORE INSERT ON tasks
    FOR EACH ROW EXECUTE FUNCTION generate_task_code();
       
COMMIT;