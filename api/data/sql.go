package data

const TIMETRAVEL_DB = "timetravel.db"

const RECORDS_TABLE = "records"
const CREATE_RECORDS_TABLE =
	`CREATE TABLE IF NOT EXISTS ` + RECORDS_TABLE + `(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		version INTEGER NOT NULL,
		jsonData TEXT NOT NULL
	);`
const INSERT_RECORD =
	`INSERT INTO ` + RECORDS_TABLE +
	` (id, version, jsonData) VALUES (?, 1, ?)`
const UPDATE_RECORD =
	`UPDATE ` + RECORDS_TABLE +
	` SET version = ?, jsonData = ? WHERE id = ?`
const QUERY_RECORD =
	`SELECT * FROM ` + RECORDS_TABLE +
	` WHERE id = ?`

const RECORD_DELTAS_TABLE = "record_deltas"
const CREATE_RECORD_DELTAS_TABLE = 
	`CREATE TABLE IF NOT EXISTS ` + RECORD_DELTAS_TABLE + `(
		id INTEGER NOT NULL,
		version INTEGER NOT NULL,
		jsonDelta TEXT NOT NULL,
		PRIMARY KEY (id, version)
	);`
const INSERT_RECORD_DELTA =
	`INSERT INTO ` + RECORD_DELTAS_TABLE +
	` (id, version, jsonDelta) VALUES (?, ?, ?)`
const QUERY_RECORD_DELTAS = 
	`SELECT * FROM ` + RECORD_DELTAS_TABLE +
	` WHERE version >= ? AND id = ?`

