package role

type Config struct {
	Migration struct {
		Run     bool     `mapstructure:"run"`
		Scripts []string `mapstructure:"scripts"`
	}
	Scripts struct {
		FetchAll   string `mapstructure:"fetch_all"`
		FetchByID  string `mapstructure:"fetch_by_id"`
		Save       string `mapstructure:"save"`
		DeleteByID string `mapstructure:"delete_by_id"`
		UpdateByID string `mapstructure:"update_by_id"`
	} `mapstructure:"scripts"`
}

func (c *Config) SetDefaultIfEmpty() *Config {
	if c.Migration.Run {
		if len(c.Migration.Scripts) == 0 {
			c.Migration.Scripts = []string{
				`
					CREATE TABLE IF NOT EXISTS roles_store (
					id SERIAL PRIMARY KEY,
					name VARCHAR(255) NOT NULL UNIQUE,
					description TEXT,
					modifier VARCHAR(255),
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
			modifier,
			created_at,
			updated_at
		FROM
			roles_store`
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
			roles_store
		WHERE
			id = $1		
		`
	}
	if c.Scripts.Save == "" {
		c.Scripts.Save = `
		INSERT INTO roles_store (
			name,
			description,
			modifier)
		VALUES (
			$1,
			$2,
			$3
		)
		ON CONFLICT DO NOTHING
		`
	}
	if c.Scripts.DeleteByID == "" {
		c.Scripts.DeleteByID = `
		DELETE
			FROM
		roles_store
			WHERE 
		id = $1`
	}

	if c.Scripts.UpdateByID == "" {
		c.Scripts.UpdateByID = `
		UPDATE roles_store
		SET
			name = $2,
			description = $3,
			modifier = $4
		WHERE id = $1`
	}
	return c
}
