package sqlhandler

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"test_ina_bank/config"
	"test_ina_bank/pkg/baselogger"
)

const (
	optionSingleStatement = "?parseTime=true&loc=UTC&multiStatements=false"
	optionMultiStatements = "?parseTime=true&loc=UTC&multiStatements=true"
)

type SqlHandler interface {
	Exec(string, ...interface{}) (Result, error)
	Query(string, ...interface{}) (Row, error)
	Transaction(func() (interface{}, error)) (interface{}, error)
	MultiExec(string) error
	Close() error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type Row interface {
	Scan(...interface{}) error
	Next() bool
	Close() error
}

// sqlHandler has a connection with DB
type sqlHandler struct {
	log     *baselogger.Logger
	DB      *sql.DB
	connect string
}

func newDB(connect, option string) (*sql.DB, error) {
	dbms := "mysql"
	connect = strings.Join([]string{connect, option}, "")
	return sql.Open(dbms, connect)
}

func newConnect(host, database, user, password string) string {
	return strings.Join([]string{user, ":", password, "@", "tcp(", host, ":3306)/", database}, "")
}

func NewSqlHandler(log *baselogger.Logger, config *config.Cfg) SqlHandler {
	host := config.Database.Host
	database := config.Database.DbName
	user := config.Database.User
	password := config.Database.Password
	log.Debug("SqlHandler created variables from Config")

	connect := newConnect(host, database, user, password)
	db, err := newDB(connect, optionSingleStatement)
	if err != nil {
		log.Panic(err)
	}
	log.Debug("SqlHandler prepared connection to database in single statement mode")

	db.SetMaxIdleConns(300)
	db.SetMaxOpenConns(300)

	return &sqlHandler{
		log:     log,
		DB:      db,
		connect: connect}
}

// Use it for database initialization only
func (handler *sqlHandler) MultiExec(multiStatements string) (err error) {
	handler.log.Debug("Connect to MySQL Database in multi statement mode")
	db, err := newDB(handler.connect, optionMultiStatements)
	if err != nil {
		handler.log.Error(err)
		return
	}
	defer db.Close()

	handler.log.Debug("Exec multi statements SQL")
	_, err = db.Exec(multiStatements)
	if err != nil {
		handler.log.Error(err)
	}
	return
}

// Exec executes the SQL that manipulates the data of the table
func (handler *sqlHandler) Exec(statement string, args ...interface{}) (Result, error) {
	handler.log.Debug("Prepare SQL statement for execution")
	stmt, err := handler.DB.Prepare(statement)
	if err != nil {
		handler.log.Error(err)
		return nil, err
	}
	defer stmt.Close()

	handler.log.Debug("Execute prepared SQL statement")
	res, err := stmt.Exec(args...)
	if err != nil {
		handler.log.Error(err)
		return nil, err
	}

	return &SqlResult{Result: res}, nil
}

// Query gets data from the database
func (handler *sqlHandler) Query(statement string, args ...interface{}) (Row, error) {
	handler.log.Debug("Prepare SQL statement for query")
	stmt, err := handler.DB.Prepare(statement)
	if err != nil {
		handler.log.Error(err)
		return nil, err
	}
	defer stmt.Close()

	handler.log.Debug("Query prepared SQL statement")
	rows, err := stmt.Query(args...)
	if err != nil {
		handler.log.Error(err)
		return nil, err
	}

	return &SqlRow{Rows: rows}, nil
}

// Transaction ...
func (handler *sqlHandler) Transaction(f func() (interface{}, error)) (interface{}, error) {
	handler.log.Debug("Begin SQL transaction")
	tx, err := handler.DB.Begin()
	if err != nil {
		handler.log.Error(err)
		return nil, err
	}

	v, err := f()
	if err != nil {
		handler.log.Error(err)
		handler.log.Warn("Rollback transaction")
		eRollback := tx.Rollback()
		if eRollback != nil {
			err = eRollback
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		handler.log.Error(err)
		handler.log.Warn("Rollback transaction")
		tx.Rollback()
		return nil, err
	}

	return v, nil
}

func (handler *sqlHandler) Close() error {
	if handler.DB != nil {
		handler.DB.Close()
	}
	return nil
}

// SqlResult ...
type SqlResult struct {
	Result sql.Result
}

// LastInsertId ...
func (r *SqlResult) LastInsertId() (int64, error) {
	res, err := r.Result.LastInsertId()
	if err != nil {
		return res, err
	}
	return res, nil
}

// RowsAffected ...
func (r *SqlResult) RowsAffected() (int64, error) {
	res, err := r.Result.LastInsertId()
	if err != nil {
		return res, err
	}
	return res, nil
}

// SqlRow ...
type SqlRow struct {
	Rows *sql.Rows
}

// Scan ...
func (r *SqlRow) Scan(dest ...interface{}) error {
	if err := r.Rows.Scan(dest...); err != nil {
		return err
	}
	return nil
}

// Next ...
func (r *SqlRow) Next() bool {
	return r.Rows.Next()
}

// Close ...
func (r *SqlRow) Close() error {
	if err := r.Rows.Close(); err != nil {
		return err
	}
	return nil
}
