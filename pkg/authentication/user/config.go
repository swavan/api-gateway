package user

type Config struct {
	Migration struct {
		Run     bool     `mapstructure:"run"`
		Scripts []string `mapstructure:"scripts"`
	}
	Scripts struct {
		FetchAll              string `mapstructure:"fetch_all"`
		FetchByUsername       string `mapstructure:"fetch_by_username"`
		FetchDomainByUsername string `mapstructure:"fetch_domain_by_username"`
		FetchByID             string `mapstructure:"fetch_by_id"`
		Save                  string `mapstructure:"save"`
		DeleteByID            string `mapstructure:"delete_by_id"`
		UpdateByUsername      string `mapstructure:"update_by_username"`
		ChangePassword        string `mapstructure:"change_password"`
		CheckCredentials      string `mapstructure:"check_credentials"`
	} `mapstructure:"scripts"`
}

func (c *Config) SetDefaultIfEmpty() *Config {
	if c.Migration.Run {
		if len(c.Migration.Scripts) == 0 {
			c.Migration.Scripts = []string{
				`
					CREATE TABLE IF NOT EXISTS users_store (
						id VARCHAR(255) PRIMARY KEY,
						user_name VARCHAR(255) NOT NULL unique,
						secret TEXT NOT NULL DEFAULT '',
						name VARCHAR(255) NOT NULL DEFAULT '',
						preferred_username VARCHAR(255) NOT NULL,
						given_name VARCHAR(255) NOT NULL DEFAULT '',
						family_name VARCHAR(255) NOT NULL DEFAULT '',
						email VARCHAR(255) NOT NULL DEFAULT '',
						email_verified BOOLEAN NOT NULL DEFAULT FALSE,
						avatar VARCHAR(255) NOT NULL DEFAULT '',
						domains TEXT NOT NULL DEFAULT '',
						none_user BOOLEAN NOT NULL DEFAULT FALSE,
						created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
					);
				`}
		}
	}

	sqlSelect := `
	SELECT 
		id,
		user_name,
		name,
		preferred_username,
		given_name,
		family_name,
		email,
		email_verified,
		avatar,
		domains,
		none_user,
		created_at
	FROM
		users_store`

	if c.Scripts.FetchAll == "" {
		c.Scripts.FetchAll = sqlSelect
	}
	if c.Scripts.FetchByUsername == "" {
		c.Scripts.FetchByUsername = sqlSelect + `
		WHERE
			user_name=$1`
	}
	if c.Scripts.Save == "" {
		c.Scripts.Save = `
		INSERT INTO users_store (
			id,
			user_name,
			preferred_username,
			name,
			given_name,
			family_name,
			email,
			email_verified,
			avatar,
			domains,
			none_user)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11)
		ON CONFLICT (user_name) DO UPDATE
		SET
			preferred_username=$3,
			name=$4,
			given_name=$5,
			family_name=$6,
			email=$7,
			email_verified=$8,
			avatar=$9,
			domains=$10,
			none_user=$11
		`
	}
	if c.Scripts.DeleteByID == "" {
		c.Scripts.DeleteByID = `
		DELETE
			FROM
		users_store
			WHERE
		user_name=$1
		`
	}
	if c.Scripts.FetchByID == "" {
		c.Scripts.FetchByID = sqlSelect + `
			WHERE
				id=$1
		`

	}
	if c.Scripts.UpdateByUsername == "" {
		c.Scripts.UpdateByUsername = "UPDATE users_store SET domains=$1 WHERE user_name=$2"
	}
	if c.Scripts.FetchDomainByUsername == "" {
		c.Scripts.FetchDomainByUsername = "SELECT domains FROM users_store WHERE user_name=$1"
	}
	if c.Scripts.ChangePassword == "" {
		c.Scripts.ChangePassword = "UPDATE users_store SET secret=$1 WHERE user_name=$2"
	}
	if c.Scripts.CheckCredentials == "" {
		c.Scripts.CheckCredentials = `
		SELECT
			id,
			user_name,
			name,
			preferred_username,
			given_name,
			family_name,
			email,
			email_verified,
			avatar,
			domains,
			none_user,
		FROM
			users_store
		WHERE
			user_name=$1`
	}
	return c
}
