-- Create translation_queues table for batch translation processing
CREATE TABLE IF NOT EXISTS translation_queues (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP WITH TIME ZONE,

	-- Foreign keys
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
	app_id INTEGER REFERENCES apps(id) ON DELETE CASCADE,

	-- Queue job details
	type VARCHAR(50) NOT NULL,
	status VARCHAR(50) NOT NULL,
	priority INTEGER NOT NULL DEFAULT 0,
	progress INTEGER DEFAULT 0,
	total INTEGER DEFAULT 0,
	done INTEGER DEFAULT 0,
	error TEXT,

	-- Configuration
	provider_type VARCHAR(50) NOT NULL,
	source_language VARCHAR(50) NOT NULL,
	target_languages JSONB NOT NULL,
	config_data JSONB NOT NULL DEFAULT '{}',

	-- Result data
	result_data JSONB DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_translation_queues_user_id ON translation_queues(user_id);
CREATE INDEX IF NOT EXISTS idx_translation_queues_project_id ON translation_queues(project_id);
CREATE INDEX IF NOT EXISTS idx_translation_queues_app_id ON translation_queues(app_id);
CREATE INDEX IF NOT EXISTS idx_translation_queues_status ON translation_queues(status);
CREATE INDEX IF NOT EXISTS idx_translation_queues_priority ON translation_queues(priority);
CREATE INDEX IF NOT EXISTS idx_translation_queues_created_at ON translation_queues(created_at);