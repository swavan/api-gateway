package domain

type Config struct {
	Migration struct {
		Run     bool     `mapstructure:"run"`
		Scripts []string `mapstructure:"scripts"`
	}
	Scripts struct {
		FetchAll    string `mapstructure:"fetch_all"`
		FetchByIDs  string `mapstructure:"fetch_by_ids"`
		FetchByID   string `mapstructure:"fetch_by_id"`
		FetchByName string `mapstructure:"fetch_by_name"`
		Save        string `mapstructure:"save"`
		DeleteByID  string `mapstructure:"delete_by_id"`
		UpdateByID  string `mapstructure:"update_by_id"`
	} `mapstructure:"scripts"`
}

func (c *Config) SetDefaultIfEmpty() *Config {
	if c.Migration.Run {
		if len(c.Migration.Scripts) == 0 {
			c.Migration.Scripts = []string{
				`
					CREATE TABLE IF NOT EXISTS domains_store (
					id VARCHAR(255) PRIMARY KEY,
					name VARCHAR(255) NOT NULL UNIQUE,
					description TEXT,
					modifier VARCHAR(255),
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
				`}
		}
	}

	if c.Scripts.FetchByName == "" {
		c.Scripts.FetchByName = `
		SELECT
			id,
			name,
			description,
			modifier,
			created_at,
			updated_at
		FROM
			domains_store
		WHERE
			name = $1`
	}

	if c.Scripts.FetchAll == "" {
		c.Scripts.FetchAll = `
		SELECT
			id,
			name,
			description,
			modifier,
			created_at,
			updated_at
		FROM domains_store`
	}
	if c.Scripts.FetchByID == "" {
		c.Scripts.FetchByID = `
		SELECT
			id,
			name,
			description,
			modifier,
			created_at,
			updated_at
		FROM
			domains_store
		WHERE
			id = $1		
		`
	}
	if c.Scripts.FetchByIDs == "" {
		c.Scripts.FetchByIDs = `
		SELECT
			id,
			name,
			description,
			modifier,
			created_at,
			updated_at
		FROM
			domains_store
		WHERE
			id IN (?)		
		`
	}
	if c.Scripts.Save == "" {
		c.Scripts.Save = `
		INSERT 
			INTO 
		domains_store
			(id, name, description, modifier)
		VALUES
			($1, $2, $3, $4)
		ON CONFLICT DO NOTHING`
	}
	if c.Scripts.DeleteByID == "" {
		c.Scripts.DeleteByID = `
		DELETE
			FROM
		domains_store
			WHERE 
		id = $1`
	}
	if c.Scripts.UpdateByID == "" {
		c.Scripts.UpdateByID = `
		UPDATE domains_store
		SET
			name = $1,
			description = $2,
			modifier = $3
		WHERE id = $4`
	}
	return c
}
