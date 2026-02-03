-- ============================================
-- xcstrings-translator Database Schema
-- Combined migration - all tables and fields
-- ============================================

-- Create users table
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	username VARCHAR(255) UNIQUE NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	password VARCHAR(255) NOT NULL,
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
	is_activated BOOLEAN NOT NULL DEFAULT FALSE,
	activation_code VARCHAR(255),
	role VARCHAR(20) NOT NULL DEFAULT 'user',
	is_subscribed BOOLEAN NOT NULL DEFAULT FALSE,
	subscription_type VARCHAR(50) NOT NULL DEFAULT 'free',
	subscription_end TIMESTAMP WITH TIME ZONE,
	max_apps INTEGER NOT NULL DEFAULT 1,
	max_translations INTEGER NOT NULL DEFAULT 1000,
	current_usage INTEGER NOT NULL DEFAULT 0,
	current_app_count INTEGER NOT NULL DEFAULT 0
);

-- Create projects table
CREATE TABLE IF NOT EXISTS projects (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	name VARCHAR(255) NOT NULL,
	description TEXT,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	file_content TEXT,
	file_name VARCHAR(255),
	source_language VARCHAR(50),
	content_structure JSONB,
	settings JSONB
);

-- Create translations table
CREATE TABLE IF NOT EXISTS translations (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	key VARCHAR(255) NOT NULL,
	source_text TEXT NOT NULL,
	target_text TEXT,
	target_language VARCHAR(50) NOT NULL,
	state VARCHAR(50),
	translation_provider VARCHAR(50)
);

-- Create provider_configs table
CREATE TABLE IF NOT EXISTS provider_configs (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	provider_type VARCHAR(50) NOT NULL,
	config_data JSONB NOT NULL,
	is_default BOOLEAN NOT NULL DEFAULT FALSE
);

-- Create user_activities table
CREATE TABLE IF NOT EXISTS user_activities (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	action VARCHAR(255) NOT NULL,
	details TEXT,
	ip_address VARCHAR(255),
	user_agent TEXT
);

-- Create apps table
CREATE TABLE IF NOT EXISTS apps (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	name VARCHAR(255) NOT NULL,
	description TEXT,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	bundle_id VARCHAR(255) UNIQUE NOT NULL,
	apple_id VARCHAR(255),
	primary_locale VARCHAR(50),
	apple_connect_token TEXT,
	origin VARCHAR(20) NOT NULL DEFAULT 'manual',
	short_description TEXT,
	long_description TEXT,
	keywords TEXT,
	support_url VARCHAR(255),
	marketing_url VARCHAR(255),
	privacy_url VARCHAR(255),
	version VARCHAR(50),
	app_category VARCHAR(100),
	is_ready_for_review BOOLEAN NOT NULL DEFAULT FALSE,
	subtitle TEXT
);

-- Create app_localizations table
CREATE TABLE IF NOT EXISTS app_localizations (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	app_id INTEGER NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
	language_code VARCHAR(50) NOT NULL,
	name VARCHAR(255),
	subtitle VARCHAR(255),
	privacy_url VARCHAR(255),
	marketing_url VARCHAR(255),
	support_url VARCHAR(255),
	description TEXT,
	keywords TEXT,
	whats_new TEXT,
	promotional_text TEXT,
	what_to_test TEXT,
	locale VARCHAR(50),
	synced_at TIMESTAMP WITHOUT TIME ZONE,
	source VARCHAR(20) NOT NULL DEFAULT 'local',
	sync_status VARCHAR(20) NOT NULL DEFAULT 'pending',
	version VARCHAR(50),
	version_state VARCHAR(50),
	UNIQUE(app_id, language_code)
);

-- Create subscriptions table
CREATE TABLE IF NOT EXISTS subscriptions (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	stripe_subscription_id VARCHAR(255) UNIQUE,
	stripe_customer_id VARCHAR(255),
	subscription_type VARCHAR(50) NOT NULL,
	subscription_status VARCHAR(50) NOT NULL,
	current_period_start TIMESTAMP WITH TIME ZONE NOT NULL,
	current_period_end TIMESTAMP WITH TIME ZONE NOT NULL,
	trial_end TIMESTAMP WITH TIME ZONE,
	cancel_at_period_end BOOLEAN NOT NULL DEFAULT FALSE
);

-- Create app_users table
CREATE TABLE IF NOT EXISTS app_users (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	app_id INTEGER NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	role VARCHAR(50) NOT NULL DEFAULT 'viewer',
	UNIQUE(app_id, user_id)
);

-- Create translation_queue table (legacy, kept for compatibility)
CREATE TABLE IF NOT EXISTS translation_queue (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
	app_id INTEGER REFERENCES apps(id) ON DELETE CASCADE,
	type VARCHAR(50) NOT NULL,
	status VARCHAR(50) NOT NULL DEFAULT 'pending',
	priority INTEGER NOT NULL DEFAULT 0,
	progress INTEGER NOT NULL DEFAULT 0,
	total INTEGER NOT NULL DEFAULT 0,
	done INTEGER NOT NULL DEFAULT 0,
	error TEXT,
	provider_type VARCHAR(50),
	source_language VARCHAR(50),
	target_languages JSONB,
	config_data JSONB,
	result_data JSONB
);

-- Create translation_queues table (new plural form)
CREATE TABLE IF NOT EXISTS translation_queues (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
	app_id INTEGER REFERENCES apps(id) ON DELETE CASCADE,
	type VARCHAR(50) NOT NULL,
	status VARCHAR(50) NOT NULL,
	priority INTEGER NOT NULL DEFAULT 0,
	progress INTEGER DEFAULT 0,
	total INTEGER DEFAULT 0,
	done INTEGER DEFAULT 0,
	error TEXT,
	provider_type VARCHAR(50) NOT NULL,
	source_language VARCHAR(50) NOT NULL,
	target_languages JSONB NOT NULL,
	config_data JSONB NOT NULL DEFAULT '{}',
	result_data JSONB DEFAULT '{}'
);

-- Create app_provider_configs table
CREATE TABLE IF NOT EXISTS app_provider_configs (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	app_id INTEGER NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
	provider_config_id INTEGER NOT NULL REFERENCES provider_configs(id) ON DELETE CASCADE,
	provider_type VARCHAR(50) NOT NULL,
	is_default BOOLEAN NOT NULL DEFAULT FALSE,
	UNIQUE(app_id, provider_config_id)
);

-- Create sync_history table
CREATE TABLE IF NOT EXISTS sync_history (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	app_id INTEGER NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
	config_id INTEGER NOT NULL REFERENCES provider_configs(id) ON DELETE CASCADE,
	direction VARCHAR(20) NOT NULL,
	strategy VARCHAR(20) NOT NULL DEFAULT 'apple_first',
	status VARCHAR(20) NOT NULL,
	total_synced INTEGER NOT NULL DEFAULT 0,
	total_failed INTEGER NOT NULL DEFAULT 0,
	total_conflicts INTEGER NOT NULL DEFAULT 0,
	language_codes JSONB,
	error TEXT,
	snapshot_data JSONB,
	ip_address VARCHAR(45),
	user_agent TEXT
);

-- ============================================
-- Create Indexes
-- ============================================

-- users indexes (none needed beyond primary key)

-- projects indexes
CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id);

-- translations indexes
CREATE INDEX IF NOT EXISTS idx_translations_project_id ON translations(project_id);

-- provider_configs indexes
CREATE INDEX IF NOT EXISTS idx_provider_configs_user_id ON provider_configs(user_id);

-- user_activities indexes
CREATE INDEX IF NOT EXISTS idx_user_activities_user_id ON user_activities(user_id);

-- apps indexes
CREATE INDEX IF NOT EXISTS idx_apps_user_id ON apps(user_id);

-- app_localizations indexes
CREATE INDEX IF NOT EXISTS idx_app_localizations_app_id ON app_localizations(app_id);

-- subscriptions indexes
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);

-- app_users indexes
CREATE INDEX IF NOT EXISTS idx_app_users_app_id ON app_users(app_id);
CREATE INDEX IF NOT EXISTS idx_app_users_user_id ON app_users(user_id);

-- translation_queue indexes
CREATE INDEX IF NOT EXISTS idx_translation_queue_user_id ON translation_queue(user_id);
CREATE INDEX IF NOT EXISTS idx_translation_queue_project_id ON translation_queue(project_id);
CREATE INDEX IF NOT EXISTS idx_translation_queue_app_id ON translation_queue(app_id);

-- translation_queues indexes
CREATE INDEX IF NOT EXISTS idx_translation_queues_user_id ON translation_queues(user_id);
CREATE INDEX IF NOT EXISTS idx_translation_queues_project_id ON translation_queues(project_id);
CREATE INDEX IF NOT EXISTS idx_translation_queues_app_id ON translation_queues(app_id);
CREATE INDEX IF NOT EXISTS idx_translation_queues_status ON translation_queues(status);
CREATE INDEX IF NOT EXISTS idx_translation_queues_priority ON translation_queues(priority);
CREATE INDEX IF NOT EXISTS idx_translation_queues_created_at ON translation_queues(created_at);

-- app_provider_configs indexes
CREATE INDEX IF NOT EXISTS idx_app_provider_configs_app_id ON app_provider_configs(app_id);
CREATE INDEX IF NOT EXISTS idx_app_provider_configs_provider_config_id ON app_provider_configs(provider_config_id);

-- sync_history indexes
CREATE INDEX IF NOT EXISTS idx_sync_history_user_id ON sync_history(user_id);
CREATE INDEX IF NOT EXISTS idx_sync_history_app_id ON sync_history(app_id);
CREATE INDEX IF NOT EXISTS idx_sync_history_config_id ON sync_history(config_id);
CREATE INDEX IF NOT EXISTS idx_sync_history_created_at ON sync_history(created_at DESC);