-- Drop translation_queues table
DROP INDEX IF EXISTS idx_translation_queues_created_at;
DROP INDEX IF EXISTS idx_translation_queues_priority;
DROP INDEX IF EXISTS idx_translation_queues_status;
DROP INDEX IF EXISTS idx_translation_queues_app_id;
DROP INDEX IF EXISTS idx_translation_queues_project_id;
DROP INDEX IF EXISTS idx_translation_queues_user_id;
DROP TABLE IF EXISTS translation_queues;