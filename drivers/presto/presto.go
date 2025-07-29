// Package presto defines and registers usql's Presto driver.
//
// See: https://github.com/prestodb/presto-go-client
package presto

import (
	"context"

	"github.com/jsthtlf/usql/drivers"
	_ "github.com/prestodb/presto-go-client/presto" // DRIVER
)

func init() {
	drivers.Register("presto", drivers.Driver{
		AllowMultilineComments: true,
		Process:                drivers.StripTrailingSemicolon,
		Version: func(ctx context.Context, db drivers.DB) (string, error) {
			var ver string
			err := db.QueryRowContext(
				ctx,
				`SELECT node_version FROM system.runtime.nodes LIMIT 1`,
			).Scan(&ver)
			if err != nil {
				return "", err
			}
			return "Presto " + ver, nil
		},
	})
}
