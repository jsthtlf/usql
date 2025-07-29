// Package chai defines and registers usql's ChaiSQL driver.
//
// See: https://github.com/chaisql/chai
package chai

import (
	_ "github.com/chaisql/chai/driver" // DRIVER
	"github.com/jsthtlf/usql/drivers"
)

func init() {
	drivers.Register("chai", drivers.Driver{})
}
