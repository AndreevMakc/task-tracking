-- migrations/000002_create_tasks.down.sql

BEGIN;

DROP TRIGGER  IF EXISTS trg_generate_task_code ON tasks;
DROP FUNCTION IF EXISTS generate_task_code;
DROP TABLE    IF EXISTS tasks CASCADE;
DROP TYPE     IF EXISTS task_status;

COMMIT;