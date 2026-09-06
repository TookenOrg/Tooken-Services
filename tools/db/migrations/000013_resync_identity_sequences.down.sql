-- Reverse of 000013: deliberately empty.
--
-- The up migration only repaired sequences that were behind their data.
-- Putting them back where they were would immediately hand out ids that
-- already exist, which is the very failure it fixed. There is nothing to
-- undo here, and undoing it would be the bug.

SELECT 1;
