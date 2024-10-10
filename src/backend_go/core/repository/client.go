package repository

import (
	"context"
	"github.com/go-pg/pg/v10"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// Config represents the configuration for the application.
type Config struct {
	URL       string `env:"COCKROACH_CLIENT_DATABASE_URL_GOLANG"`
	DebugMode bool   `env:"COCKROACH_CLIENT_DB_DEBUG" envDefault:"false"`
}

// NewDatabase creates a new database connection using the provided configuration.
// It takes a pointer to a Config struct as input. The Config struct represents the
// required configuration parameters for the database connection, including the URL
// and the debug mode. The URL is used to parse the connection options, and if any
// parse error occurs, an error is returned. Then, a new database connection is established
// using the parsed options by calling the Connect function from the pg package.
// If the debug mode is enabled, a query hook is added to the database connection, which
// logs the executed queries. After the connection is established, a ping request is sent
// to the database to verify its availability. If the ping request fails, an error is returned.
// Finally, the function returns the established database connection and nil error if no errors occurred.
func NewDatabase(cfg *Config) (*pg.DB, error) {
	opt, err := pg.ParseURL(cfg.URL)
	if err != nil {
		return nil, err
	}
	db := pg.Connect(opt)
	if cfg.DebugMode {
		db.AddQueryHook(dbLogger{})
	}
	if err = db.Ping(context.Background()); err != nil {
		return nil, err
	}
	return db, nil
}

// dbLogger is a type that implements the pg.QueryHook interface.
// It allows logging of SQL queries before and after their execution in a PostgreSQL database.
// Use the BeforeQuery function to perform any operations before the query is executed.
// Use the AfterQuery function to log the query and any errors that occurred during execution.
type dbLogger struct{}

// BeforeQuery is a method of dbLogger that is called before a query is executed.
// It takes a context.Context object and a pointer to a pg.QueryEvent object as parameters.
// It returns the context.Context object and an error.
// This method does not perform any actions and always returns the input context and a nil error.
func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

// AfterQuery is a method of the dbLogger struct that is called after a database query is executed.
// It logs the executed query using zap.S().Infof().
// If an error occurs while getting the formatted query, it logs the error instead.
// This method returns nil error as it is not expected to fail.
func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	if query, err := q.FormattedQuery(); err != nil {
		zap.S().Infof("[SQL]: %s", err.Error())
	} else {
		zap.S().Infof("[SQL]: %s", string(query))
	}

	return nil
}
