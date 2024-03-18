package access

const (
	// maxParameterCount is the maximum number of parameters that a rule can have.
	maxParameterCount  = 7
	defaultPlaceholder = "?"
)

const (
	CREATE_QL = `
	CREATE TABLE IF NOT EXISTS %[1]s(
		p_type VARCHAR(32)  DEFAULT '' NOT NULL,
		v0     VARCHAR(255) DEFAULT '' NOT NULL,
		v1     VARCHAR(255) DEFAULT '' NOT NULL,
		v2     VARCHAR(255) DEFAULT '' NOT NULL,
		v3     VARCHAR(255) DEFAULT '' NOT NULL,
		v4     VARCHAR(255) DEFAULT '' NOT NULL,
		v5     VARCHAR(255) DEFAULT '' NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_%[1]s ON %[1]s (p_type,v0,v1);`
	TRUNCATE_QL     = "TRUNCATE TABLE %s"
	INSERT_QL       = "INSERT INTO %s (p_type,v0,v1,v2,v3,v4,v5) VALUES ($1,$2,$3,$4,$5,$6,$7)"
	UPDATE_QL       = "UPDATE %s SET p_type=$1,v0=$2,v1=$3,v2=$4,v3=$5,v4=$6,v5=$7 WHERE p_type=$8 AND v0=$9 AND v1=$10 AND v2=$11 AND v3=$12 AND v4=$13 AND v5=$14"
	DELETE_ALL_QL   = "DELETE FROM %s"
	DELETE_QL       = "DELETE FROM %s WHERE p_type=$1 AND v0=$2 AND v1=$3 AND v2=$4 AND v3=$5 AND v4=$6 AND v5=$7"
	DELETE_BY_ARGS  = "DELETE FROM %s WHERE p_type=$1"
	SELECT_ALL_QL   = "SELECT p_type,v0,v1,v2,v3,v4,v5 FROM %s"
	SELECT_WHERE_QL = "SELECT p_type,v0,v1,v2,v3,v4,v5 FROM %s WHERE"
	TABLE_EXIST_QL  = "SELECT 1 FROM %s"
)
