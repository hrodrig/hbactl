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

// WriteRulesTableGroupedByUser prints rules in a table with "=== user: X ===" separators.
// Rules should be sorted by user so that consecutive rules with the same user form a group.
func WriteRulesTableGroupedByUser(w io.Writer, rules []hba.Rule) {
	if len(rules) == 0 {
		return
	}
	var group []hba.Rule
	curUser := ""
	for _, r := range rules {
		if r.User != curUser {
			if len(group) > 0 {
				fmt.Fprintf(w, "\n=== user: %s ===\n\n", curUser)
				writeRulesTableTo(w, group)
			}
			curUser = r.User
			group = group[:0]
		}
		group = append(group, r)
	}
	if len(group) > 0 {
		fmt.Fprintf(w, "\n=== user: %s ===\n\n", curUser)
		writeRulesTableTo(w, group)
	}
}

func writeRulesTableTo(w io.Writer, rules []hba.Rule) {
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
