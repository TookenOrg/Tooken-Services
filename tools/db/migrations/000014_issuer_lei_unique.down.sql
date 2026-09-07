-- Rollback 000014 — l'unicité du LEI redevient une simple convention.

DROP INDEX IF EXISTS ass.issuer_lei_uk;
