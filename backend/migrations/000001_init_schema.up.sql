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
CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id);

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
CREATE INDEX IF NOT EXISTS idx_translations_project_id ON translations(project_id);

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
CREATE INDEX IF NOT EXISTS idx_provider_configs_user_id ON provider_configs(user_id);

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
CREATE INDEX IF NOT EXISTS idx_user_activities_user_id ON user_activities(user_id);

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
	is_ready_for_review BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX IF NOT EXISTS idx_apps_user_id ON apps(user_id);

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
	download_description TEXT,
	short_description TEXT,
	long_description TEXT,
	keywords TEXT,
	release_notes TEXT,
	synced_at TIMESTAMP WITH TIME ZONE,
	source VARCHAR(20) NOT NULL DEFAULT 'local',
	sync_status VARCHAR(20) NOT NULL DEFAULT 'pending',
	UNIQUE(app_id, language_code)
);
CREATE INDEX IF NOT EXISTS idx_app_localizations_app_id ON app_localizations(app_id);

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
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);

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
CREATE INDEX IF NOT EXISTS idx_app_users_app_id ON app_users(app_id);
CREATE INDEX IF NOT EXISTS idx_app_users_user_id ON app_users(user_id);

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
CREATE INDEX IF NOT EXISTS idx_translation_queue_user_id ON translation_queue(user_id);
CREATE INDEX IF NOT EXISTS idx_translation_queue_project_id ON translation_queue(project_id);
CREATE INDEX IF NOT EXISTS idx_translation_queue_app_id ON translation_queue(app_id);

CREATE TABLE IF NOT EXISTS app_provider_configs (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH time ZONE,
	app_id INTEGER NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
	provider_config_id INTEGER NOT NULL REFERENCES provider_configs(id) ON DELETE CASCADE,
	provider_type VARCHAR(50) NOT NULL,
	is_default BOOLEAN NOT NULL DEFAULT FALSE,
	UNIQUE(app_id, provider_config_id)
);
CREATE INDEX IF NOT EXISTS idx_app_provider_configs_app_id ON app_provider_configs(app_id);
CREATE INDEX IF NOT EXISTS idx_app_provider_configs_provider_config_id ON app_provider_configs(provider_config_id);