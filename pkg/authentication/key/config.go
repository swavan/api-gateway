package key

type Config struct {
	GenerateKey struct {
		UsrFor []string `mapstructure:"use_for"`
	} `mapstructure:"generate_key"`
	Migration struct {
		Run     bool     `mapstructure:"run"`
		Scripts []string `mapstructure:"scripts"`
	} `mapstructure:"migration"`
	Scripts struct {
		DeleteByID     string `mapstructure:"delete_by_id"`
		DeleteByUseFor string `mapstructure:"delete_by_use_for"`
		GetByID        string `mapstructure:"get_by_id"`
		GetByUseFor    string `mapstructure:"get_by_use_for"`
		FetchAll       string `mapstructure:"fetch_all"`
		Save           string `mapstructure:"save"`
	} `mapstructure:"scripts"`
}

func NewConfig() *Config {
	return &Config{}
}

func (kc *Config) SetDefaultIfEmpty() *Config {
	if len(kc.Migration.Scripts) == 0 {
		kc.Migration.Scripts = []string{
			`CREATE TABLE IF NOT EXISTS keys_store (
					id SERIAL PRIMARY KEY,
					use_for VARCHAR(255) NOT NULL unique,
					private Text NOT NULL,
					public Text NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				);`}
	}
	if kc.Scripts.DeleteByID == "" {
		kc.Scripts.DeleteByID = "DELETE FROM keys_store WHERE id = $1"
	}
	if kc.Scripts.DeleteByUseFor == "" {
		kc.Scripts.DeleteByUseFor = "DELETE FROM keys_store WHERE use_for = $1"
	}
	if kc.Scripts.GetByID == "" {
		kc.Scripts.GetByID = "SELECT id, use_for, private, public, created_at FROM keys_store WHERE id = $1 limit 1"
	}
	if kc.Scripts.GetByUseFor == "" {
		kc.Scripts.GetByUseFor = "SELECT id, use_for, private, public, created_at FROM keys_store WHERE use_for = $1 limit 1"
	}
	if kc.Scripts.Save == "" {
		kc.Scripts.Save = `
		INSERT INTO keys_store
			(use_for, private, public)
		VALUES
			($1, $2, $3)
		ON CONFLICT (use_for) DO NOTHING`
	}
	if kc.Scripts.FetchAll == "" {
		kc.Scripts.FetchAll = `
		SELECT 
			id,
			use_for,
			private,
			public,
			created_at
		FROM keys_store
		`
	}
	return kc
}
