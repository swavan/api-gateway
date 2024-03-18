package resource

type Config struct {
	Migration struct {
		Run     bool     `mapstructure:"run"`
		Scripts []string `mapstructure:"scripts"`
	}
	Scripts struct {
		FetchAll             string `mapstructure:"fetch_all"`
		FetchByID            string `mapstructure:"fetch_by_id"`
		FetchActionsByID     string `mapstructure:"fetch_actions_by_id"`
		FetchActionsBySource string `mapstructure:"fetch_actions_by_source"`
		Create               string `mapstructure:"save"`
		DeleteByID           string `mapstructure:"delete_by_id"`
		UpdateByID           string `mapstructure:"updated_by_id"`
	} `mapstructure:"scripts"`
}

func (c *Config) SetDefaultIfEmpty() *Config {
	if c.Migration.Run {
		if len(c.Migration.Scripts) == 0 {
			c.Migration.Scripts = []string{
				`
					CREATE TABLE IF NOT EXISTS resource_store (
					id VARCHAR(255) PRIMARY KEY,
					name VARCHAR(255) NOT NULL UNIQUE,
					description TEXT NOT NULL DEFAULT '',
					source VARCHAR(255) NOT NULL,
					actions VARCHAR(255) NOT NULL DEFAULT 'ANY',
					modifier VARCHAR(255) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
				`}
		}
	}
	if c.Scripts.FetchAll == "" {
		c.Scripts.FetchAll = `
		SELECT
			id,
			name,
			description,
			source,
			actions,
			modifier,
			created_at,
			updated_at
		FROM
			resource_store`
	}
	if c.Scripts.FetchByID == "" {
		c.Scripts.FetchByID = `
		SELECT
			id,
			name,
			description,
			source,
			actions,
			modifier,
			created_at,
			updated_at
		FROM
			resource_store
		WHERE
			id = $1
		LIMIT 1
		`
	}
	if c.Scripts.Create == "" {
		c.Scripts.Create = `
		INSERT INTO resource_store (
			id,
			name,
			description,
			source,
			actions,
			modifier)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6)
		ON CONFLICT DO NOTHING
		`
	}
	if c.Scripts.DeleteByID == "" {
		c.Scripts.DeleteByID = `
		DELETE
			FROM
		resource_store
			WHERE 
		id = $1`
	}

	if c.Scripts.UpdateByID == "" {
		c.Scripts.UpdateByID = `
		UPDATE resource_store
		SET
			name = $2,
			description = $3,
			source = $4,
			actions = $5,
			modifier = $6
		WHERE id = $1`
	}
	if c.Scripts.FetchActionsByID == "" {
		c.Scripts.FetchActionsByID = `
		SELECT
			actions
		FROM
			resource_store
		WHERE
			id = $1
		LIMIT 1
		`
	}

	if c.Scripts.FetchActionsBySource == "" {
		c.Scripts.FetchActionsBySource = `
		SELECT
			actions
		FROM
			resource_store
		WHERE
			source = $1
		`
	}

	return c
}
