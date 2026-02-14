package cli

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/hrodrig/hbactl/internal/hba"
)

// WriteRulesTable prints rules in a formatted table to w.
func WriteRulesTable(w io.Writer, rules []hba.Rule) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TYPE\tDATABASE\tUSER\tADDRESS\tMETHOD")
	fmt.Fprintln(tw, "----\t--------\t----\t-------\t------")
	for _, r := range rules {
		addr := r.Address
		if r.Netmask != "" {
			addr = r.Address + " / " + r.Netmask
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", r.Type, r.Database, r.User, addr, r.Method)
	}
	tw.Flush()
}
