package secret

type Config struct {
	Migration struct {
		Run     bool     `mapstructure:"run"`
		Scripts []string `mapstructure:"scripts"`
	}
	Scripts struct {
		FetchByIDs    string `mapstructure:"fetch_by_ids"`
		Save          string `mapstructure:"save"`
		DeleteByID    string `mapstructure:"delete_by_id"`
		ArchiveByID   string `mapstructure:"archive_by_id"`
		FetchByDomain string `mapstructure:"fetch_by_domain"`
		FetchByUser   string `mapstructure:"fetch_by_user"`
	} `mapstructure:"scripts"`
}

func (c *Config) SetDefaultIfEmpty() *Config {
	if c.Migration.Run {
		if len(c.Migration.Scripts) == 0 {
			c.Migration.Scripts = []string{
				`
					CREATE TABLE IF NOT EXISTS secret_store (
					id VARCHAR(255) PRIMARY KEY,
					description VARCHAR(255) NOT NULL,
					type VARCHAR(255) NOT NULL,
					domain VARCHAR(255) NOT NULL,
					issue_at TIMESTAMP NOT NULL,
					expires_at TIMESTAMP NOT NULL,
					alert_to VARCHAR(255),
					modifier VARCHAR(255) NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					delete_at TIMESTAMP);
				`}
		}
	}

	if c.Scripts.FetchByIDs == "" {
		c.Scripts.FetchByIDs = `
		SELECT
			id,
			type,
			description,
			domain,
			issue_at,
			expires_at,
			alert_to,
			modifier,
			created_at
		FROM
			secret_store
		WHERE
			id = $1 and delete_at is null
		LIMIT 1`
	}

	if c.Scripts.Save == "" {
		c.Scripts.Save = `
		INSERT INTO secret_store (
			id,
			description,
			type,
			domain,
			issue_at,
			expires_at,
			alert_to,
			modifier
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8
		)`
	}

	if c.Scripts.DeleteByID == "" {
		c.Scripts.DeleteByID = `
		DELETE FROM secret_store
		WHERE
			id = $1`
	}

	if c.Scripts.FetchByDomain == "" {
		c.Scripts.FetchByDomain = `
		SELECT
			id,
			description,
			type,
			domain,
			issue_at,
			expires_at,
			alert_to,
			modifier,
			created_at
		FROM
			secret_store
		WHERE
			domain = $1 and delete_at is null`
	}

	if c.Scripts.FetchByUser == "" {
		c.Scripts.FetchByUser = `
		SELECT
			id,
			description,
			type,
			domain,
			issue_at,
			expires_at,
			alert_to,
			modifier,
			created_at
		FROM
			secret_store
		WHERE
			modifier = $1 and delete_at is null`
	}

	if c.Scripts.ArchiveByID == "" {
		c.Scripts.ArchiveByID = `
		UPDATE secret_store
		SET delete_at = CURRENT_TIMESTAMP
		WHERE
			id = $1`

	}

	return c
}
