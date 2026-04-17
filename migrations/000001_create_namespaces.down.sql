-- migrations/000001_create_namespaces.down.sql

BEGIN;

DROP TRIGGER  IF EXISTS trg_protect_default_namespace ON namespaces;
DROP FUNCTION IF EXISTS protect_default_namespace;
DROP TABLE    IF EXISTS namespaces CASCADE;

COMMIT;