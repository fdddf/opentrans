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
CREATE INDEX IF NOT EXISTS idx_sync_history_user_id ON sync_history(user_id);
CREATE INDEX IF NOT EXISTS idx_sync_history_app_id ON sync_history(app_id);
CREATE INDEX IF NOT EXISTS idx_sync_history_config_id ON sync_history(config_id);
CREATE INDEX IF NOT EXISTS idx_sync_history_created_at ON sync_history(created_at DESC);