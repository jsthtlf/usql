package metacmd

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"

	"github.com/jsthtlf/usql/env"
	"github.com/jsthtlf/usql/text"
)

// Help is a Help meta command (\?). Writes a help message to the output.
//
// Descs:
//
//	?	[commands]	show help on {{CommandName}}'s meta (backslash) commands
//	?	options	show help on {{CommandName}} command-line options
//	?	variables	show help on special {{CommandName}} variables
func Help(p *Params) error {
	stdout, stderr := p.Handler.IO().Stdout(), p.Handler.IO().Stderr()
	var err error
	var cmd *exec.Cmd
	var wc io.WriteCloser
	if pager := env.Get("PAGER"); p.Handler.IO().Interactive() && pager != "" {
		if wc, cmd, err = env.Pipe(stdout, stderr, pager); err != nil {
			return err
		}
		stdout = wc
	}
	_ = Dump(stdout, false)
	if cmd != nil {
		if err := wc.Close(); err != nil {
			return err
		}
		return cmd.Wait()
	}
	return nil
}

// Transact is a Transaction meta command (\begin, \commit, \rollback). Begins,
// commits, or aborts (rollback) the current database transaction on the open
// database connection.
//
// Descs:
//
//	begin	[-read-only [ISOLATION]]	begin transaction, with optional isolation level
//	commit	commit current transaction
//	rollback	rollback (abort) current transaction
//	abort:rollback
func Transact(p *Params) error {
	switch p.Name {
	case "commit":
		return p.Handler.Commit()
	case "rollback", "abort":
		return p.Handler.Rollback()
	}
	// read begin params
	readOnly := false
	n, ok, err := p.NextOpt(true)
	if ok {
		if n != "read-only" {
			return fmt.Errorf(text.InvalidOption, n)
		}
		readOnly = true
		if n, err = p.Next(true); err != nil {
			return err
		}
	}
	// build tx options
	var txOpts *sql.TxOptions
	if readOnly || n != "" {
		isolation := sql.LevelDefault
		switch strings.ToLower(n) {
		case "default", "":
		case "read-uncommitted":
			isolation = sql.LevelReadUncommitted
		case "read-committed":
			isolation = sql.LevelReadCommitted
		case "write-committed":
			isolation = sql.LevelWriteCommitted
		case "repeatable-read":
			isolation = sql.LevelRepeatableRead
		case "snapshot":
			isolation = sql.LevelSnapshot
		case "serializable":
			isolation = sql.LevelSerializable
		case "linearizable":
			isolation = sql.LevelLinearizable
		default:
			return text.ErrInvalidIsolationLevel
		}
		txOpts = &sql.TxOptions{
			Isolation: isolation,
			ReadOnly:  readOnly,
		}
	}
	// begin
	return p.Handler.Begin(txOpts)
}

// Describe is a Informational meta command (\d and variants). Queries the open
// database connection for information about the database schema and writes the
// information to the output.
//
// Descs:
//
//	d[S+]	[NAME]	list tables, views, and sequences or describe table, view, sequence, or index
//	da[S+]	[PATTERN]	list aggregates
//	df[S+]	[PATTERN]	list functions
//	di[S+]	[PATTERN]	list indexes
//	dm[S+]	[PATTERN]	list materialized views
//	dn[S+]	[PATTERN]	list schemas
//	dp[S]	[PATTERN]	list table, view, and sequence access privileges
//	ds[S+]	[PATTERN]	list sequences
//	dt[S+]	[PATTERN]	list tables
//	dv[S+]	[PATTERN]	list views
//	l[+]	list databases
func Describe(p *Params) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	m, err := p.Handler.MetadataWriter(ctx)
	if err != nil {
		return err
	}
	verbose := strings.ContainsRune(p.Name, '+')
	showSystem := strings.ContainsRune(p.Name, 'S')
	name := strings.TrimRight(p.Name, "S+")
	pattern, err := p.Next(true)
	if err != nil {
		return err
	}
	switch name {
	case "d":
		if pattern != "" {
			return m.DescribeTableDetails(p.Handler.URL(), pattern, verbose, showSystem)
		}
		return m.ListTables(p.Handler.URL(), "tvmsE", pattern, verbose, showSystem)
	case "df", "da":
		return m.DescribeFunctions(p.Handler.URL(), name, pattern, verbose, showSystem)
	case "dt", "dtv", "dtm", "dts", "dv", "dm", "ds":
		return m.ListTables(p.Handler.URL(), name, pattern, verbose, showSystem)
	case "dn":
		return m.ListSchemas(p.Handler.URL(), pattern, verbose, showSystem)
	case "di":
		return m.ListIndexes(p.Handler.URL(), pattern, verbose, showSystem)
	case "l":
		return m.ListAllDbs(p.Handler.URL(), pattern, verbose)
	case "dp":
		return m.ListPrivilegeSummaries(p.Handler.URL(), pattern, showSystem)
	}
	return nil
}

// Stats is a Informational meta command (\ss and variants). Queries the open
// database connection for stats and writes it to the output.
//
// Descs:
//
//	ss[+]	[TABLE|QUERY] [k]	show stats for a table or a query
func Stats(p *Params) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	m, err := p.Handler.MetadataWriter(ctx)
	if err != nil {
		return err
	}
	verbose := strings.ContainsRune(p.Name, '+')
	name := strings.TrimRight(p.Name, "+")
	pattern, err := p.Next(true)
	if err != nil {
		return err
	}
	k := 0
	if verbose {
		k = 3
	}
	if name == "ss" {
		name = "sswnulhmkf"
	}
	val, ok, err := p.NextOK(true)
	switch {
	case err != nil:
		return err
	case ok:
		verbose = true
		if k, err = strconv.Atoi(val); err != nil {
			return err
		}
	}
	return m.ShowStats(p.Handler.URL(), name, pattern, verbose, k)
}
