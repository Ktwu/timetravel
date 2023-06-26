package data

const TIMETRAVEL_DB = "timetravel.db"

const RECORDS_TABLE = "records"
const CREATE_RECORDS_TABLE = `CREATE TABLE IF NOT EXISTS ` +
	RECORDS_TABLE + `(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		version INTEGER NOT NULL,
		jsonData TEXT NOT NULL
	);`
const INSERT_RECORD = `INSERT INTO ` + RECORDS_TABLE +
	` (id, version, jsonData) VALUES (?, 1, ?)`
const UPDATE_RECORD = `UPDATE ` + RECORDS_TABLE +
	` SET version = ?, jsonData = ? WHERE id = ?`
const QUERY_RECORD = `SELECT * FROM ` + RECORDS_TABLE +
	` WHERE id = ?`

const RECORD_DELTAS_TABLE = "record_deltas"
const CREATE_RECORD_DELTAS_TABLE = `CREATE TABLE IF NOT EXISTS ` +
	RECORD_DELTAS_TABLE + `(
		id INTEGER NOT NULL,
		versionBeforeDelta INTEGER NOT NULL,
		inverseDelta TEXT NOT NULL,
		PRIMARY KEY (id, versionBeforeDelta)
	);`
const INSERT_RECORD_DELTA = `INSERT INTO ` + RECORD_DELTAS_TABLE +
	` (id, versionBeforeDelta, inverseDelta) VALUES (?, ?, ?)`

// When calculating record versions, we apply inverse updates on the current
// version. Make sure we sort the results of this query so that we iterate
// through the most recent updates (that should be applied first).
const QUERY_RECORD_DELTAS = `SELECT * FROM ` + RECORD_DELTAS_TABLE +
	` WHERE versionBeforeDelta >= ? AND id = ?
	  ORDER BY versionBeforeDelta DESC`
