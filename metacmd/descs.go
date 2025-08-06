package metacmd

import (
	"github.com/jsthtlf/usql/text"
)

// sections are the command description sections.
var sections = []string{
	"Help",
	"Query Buffer",
	"Informational",
	"Transaction",
}

// descs are the command descriptions.
var descs [][]desc

// cmds are the command lookup map.
var cmds map[string]func(*Params) error

func init() {
	descs = [][]desc{
		// Help
		{
			{Help, `?`, `[commands]`, `show help on ` + text.CommandName + `'s meta (backslash) commands`, true, false},
		},
		// Query Buffer
		{
			{Reset, `r`, ``, `reset (clear) the query buffer`, false, false},
			{Reset, `reset`, ``, `alias for \r`, true, false},
		},
		// Informational
		{
			{Describe, `d[S+]`, `[NAME]`, `list tables, views, and sequences or describe table, view, sequence, or index`, false, false},
			{Describe, `da[S+]`, `[PATTERN]`, `list aggregates`, false, false},
			{Describe, `df[S+]`, `[PATTERN]`, `list functions`, false, false},
			{Describe, `di[S+]`, `[PATTERN]`, `list indexes`, false, false},
			{Describe, `dm[S+]`, `[PATTERN]`, `list materialized views`, false, false},
			{Describe, `dn[S+]`, `[PATTERN]`, `list schemas`, false, false},
			{Describe, `dp[S]`, `[PATTERN]`, `list table, view, and sequence access privileges`, false, false},
			{Describe, `ds[S+]`, `[PATTERN]`, `list sequences`, false, false},
			{Describe, `dt[S+]`, `[PATTERN]`, `list tables`, false, false},
			{Describe, `dv[S+]`, `[PATTERN]`, `list views`, false, false},
			{Describe, `l[+]`, ``, `list databases`, false, false},
			{Stats, `ss[+]`, `[TABLE|QUERY] [k]`, `show stats for a table or a query`, false, false},
		},
		// Transaction
		{
			{Transact, `begin`, `[-read-only [ISOLATION]]`, `begin transaction, with optional isolation level`, false, false},
			{Transact, `commit`, ``, `commit current transaction`, false, false},
			{Transact, `rollback`, ``, `rollback (abort) current transaction`, false, false},
			{Transact, `abort`, ``, `alias for \rollback`, false, false},
		},
	}
	cmds = make(map[string]func(*Params) error)
	for i := range sections {
		for _, desc := range descs[i] {
			for _, n := range desc.Names() {
				cmds[n] = desc.Func
			}
		}
	}
}
