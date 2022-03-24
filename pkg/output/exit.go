package output

import (
	"fmt"
	"os"

	"github.com/rafaylabs/rcloud-cli/pkg/context"
	"github.com/rafaylabs/rcloud-cli/pkg/exit"
	"github.com/rafaylabs/rcloud-cli/pkg/log"
)

/* This function exits the rctl program. In case 'exit' is not set,
   nothing will be printed to the console. The return code of
   program is set to zero.

   When 'exit' is set, the exit message will be printed to
   the console before the program exits with the return code
   set in the 'exit' structure.
*/
func Exit() {
	e := exit.Get()
	if e == nil {
		log.GetLogger().Debugf("Exit 0")
		os.Exit(0)
	}

	log.GetLogger().Debugf("Exit with Error")
	if context.GetContext().UseStructuredOutput() == true {
		PrintJson(e, true)
	} else {
		fmt.Println(e.String())
	}

	os.Exit(e.ReturnCode)
}
