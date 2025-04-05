package hooks

import (
	"context"
	"database/sql/driver"
	"errors"
	"time"

	"github.com/artifact-space/ArtiSpace/log"
)

type HookedSQLConn struct {
	conn driver.Conn
}

// implement driver.Conn interface

func (c *HookedSQLConn) Prepare(query string) (driver.Stmt, error) {
	start := time.Now()
	stmt, err := c.conn.Prepare(query)
	elapsed := time.Since(start)

	log.Logger().Debug().Msgf("Prepared query: %s , elapsed: %v", query, elapsed)
	return stmt, err
}

func (c *HookedSQLConn) Close() error {
	return c.conn.Close()
}

func (c *HookedSQLConn) Begin() (driver.Tx, error) {
	return c.conn.Begin()
}

// implement ConnBeginTx
func (c *HookedSQLConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	start := time.Now()
	var elapsed time.Duration
	var tx driver.Tx
	var err error
	if connBeginTx, ok := c.conn.(driver.ConnBeginTx); ok {
		tx, err = connBeginTx.BeginTx(ctx, opts)
		elapsed = time.Since(start)
	} else {
		tx, err = c.conn.Begin()
		elapsed = time.Since(start)
	}

	log.Logger().Debug().Msgf("Starting database transaction, elapsed: %v", elapsed)
	return tx, err
}

// implement ConnPrepareContext
func (c *HookedSQLConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	start := time.Now()
	var elapsed time.Duration
	var stmt driver.Stmt
	var err error
	if pc, ok := c.conn.(driver.ConnPrepareContext); ok {
		stmt, err = pc.PrepareContext(ctx, query)
		elapsed = time.Since(start)
	} else {
		return stmt, errors.ErrUnsupported
	}
	log.Logger().Debug().Msgf("PrepareContext, query: %s, elapsed: %v", query, elapsed)
	return stmt, err
}

// implement ExecerContext
func (c *HookedSQLConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	start := time.Now()
	var result driver.Result
	var err error
	var elapsed time.Duration

	if ec, ok := c.conn.(driver.ExecerContext); ok {
		result, err = ec.ExecContext(ctx, query, args)
		elapsed = time.Since(start)
	} else {
		return result, errors.ErrUnsupported
	}
	log.Logger().Debug().Msgf("ExecContext, query: %s, args: %v, elapsed: %v", query, args, elapsed)
	return result, err
}

// implement QueryerContext
func (c *HookedSQLConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	start := time.Now()
	var rows driver.Rows
	var err error
	var elapsed time.Duration
	if qc, ok := c.conn.(driver.QueryerContext); ok {
		rows, err = qc.QueryContext(ctx, query, args)
		elapsed = time.Since(start)
	} else {
		return rows, errors.ErrUnsupported
	}
	log.Logger().Debug().Msgf("QueryContext, query: %s, args: %v, elapsed: %v", query, args, elapsed)

	return rows, err
}

// implement Queryer
func (c *HookedSQLConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	start := time.Now()
	result, err := c.conn.(driver.Queryer).Query(query, args)
	elapsed := time.Since(start)
	log.Logger().Debug().Msgf("Query: %s, args: %s, elapsed: %v", query, args, elapsed)
	return result, err
}

// implement Execer
func (c *HookedSQLConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	start := time.Now()
	result, err := c.conn.(driver.Execer).Exec(query, args)
	elapsed := time.Since(start)
	log.Logger().Debug().Msgf("Exec query: %s, args: %v, elapsed: %v", query, args, elapsed)
	return result, err
}

type HookedDriver struct {
	driver.Driver
}

func (d *HookedDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		return nil, err
	}
	return &HookedSQLConn{conn: conn}, nil
}
