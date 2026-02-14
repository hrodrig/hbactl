package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Client connects to PostgreSQL and runs HBA-related queries.
type Client struct {
	pool *pgxpool.Pool
}

// NewClient creates a client using connStr (e.g. from DATABASE_URL or --conn).
func NewClient(ctx context.Context, connStr string) (*Client, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &Client{pool: pool}, nil
}

// Close releases the connection pool.
func (c *Client) Close() {
	if c.pool != nil {
		c.pool.Close()
	}
}

// HBAFilePath returns the path to pg_hba.conf by running SHOW hba_file.
func (c *Client) HBAFilePath(ctx context.Context) (string, error) {
	var path string
	err := c.pool.QueryRow(ctx, "SHOW hba_file").Scan(&path)
	return path, err
}

// HBAFileError describes a syntax error in pg_hba.conf reported by pg_hba_file_rules.
type HBAFileError struct {
	LineNumber int
	Error      string
}

// ReloadConf makes PostgreSQL reload its configuration (including pg_hba.conf). No restart required.
func (c *Client) ReloadConf(ctx context.Context) error {
	_, err := c.pool.Exec(ctx, "SELECT pg_reload_conf()")
	return err
}

// HBAFileErrors returns any syntax errors in the current pg_hba.conf (from pg_hba_file_rules where error IS NOT NULL).
func (c *Client) HBAFileErrors(ctx context.Context) ([]HBAFileError, error) {
	rows, err := c.pool.Query(ctx, "SELECT line_number, error FROM pg_hba_file_rules WHERE error IS NOT NULL ORDER BY line_number")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var errs []HBAFileError
	for rows.Next() {
		var e HBAFileError
		if err := rows.Scan(&e.LineNumber, &e.Error); err != nil {
			return nil, err
		}
		errs = append(errs, e)
	}
	return errs, rows.Err()
}
