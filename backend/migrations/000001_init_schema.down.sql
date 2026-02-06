-- ============================================
-- Drop all tables and indexes
-- ============================================

-- Drop indexes (explicitly drop all indexes first)
DROP INDEX IF EXISTS idx_sync_history_created_at;
DROP INDEX IF EXISTS idx_sync_history_config_id;
DROP INDEX IF EXISTS idx_sync_history_app_id;
DROP INDEX IF EXISTS idx_sync_history_user_id;

DROP INDEX IF EXISTS idx_app_provider_configs_provider_config_id;
DROP INDEX IF EXISTS idx_app_provider_configs_app_id;

DROP INDEX IF EXISTS idx_translation_queues_created_at;
DROP INDEX IF EXISTS idx_translation_queues_priority;
DROP INDEX IF EXISTS idx_translation_queues_status;
DROP INDEX IF EXISTS idx_translation_queues_app_id;
DROP INDEX IF EXISTS idx_translation_queues_project_id;
DROP INDEX IF EXISTS idx_translation_queues_user_id;

DROP INDEX IF EXISTS idx_translation_queue_app_id;
DROP INDEX IF EXISTS idx_translation_queue_project_id;
DROP INDEX IF EXISTS idx_translation_queue_user_id;

DROP INDEX IF EXISTS idx_app_users_user_id;
DROP INDEX IF EXISTS idx_app_users_app_id;

DROP INDEX IF EXISTS idx_subscriptions_user_id;

DROP INDEX IF EXISTS idx_app_localizations_app_id;

DROP INDEX IF EXISTS idx_apps_user_id;

DROP INDEX IF EXISTS idx_user_activities_user_id;

DROP INDEX IF EXISTS idx_provider_configs_user_id;

DROP INDEX IF EXISTS idx_translations_project_id;

DROP INDEX IF EXISTS idx_apps_project_id;
DROP INDEX IF EXISTS idx_projects_user_id;

-- Drop tables (in correct order - respecting foreign key dependencies)
DROP TABLE IF EXISTS sync_history;

DROP TABLE IF EXISTS app_provider_configs;

DROP TABLE IF EXISTS translation_queues;

DROP TABLE IF EXISTS translation_queue;

DROP TABLE IF EXISTS app_users;

DROP TABLE IF EXISTS subscriptions;

DROP TABLE IF EXISTS app_localizations;

DROP TABLE IF EXISTS apps;

DROP TABLE IF EXISTS user_activities;

DROP TABLE IF EXISTS provider_configs;

DROP TABLE IF EXISTS translations;

DROP TABLE IF EXISTS projects;

DROP TABLE IF EXISTS users;
