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

// WriteRulesTableWithIndex prints rules with a 1-based index column (#). Use for list when remove is available.
func WriteRulesTableWithIndex(w io.Writer, rwl []hba.RuleWithLine) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "#\tTYPE\tDATABASE\tUSER\tADDRESS\tMETHOD")
	fmt.Fprintln(tw, "-\t----\t--------\t----\t-------\t------")
	for _, x := range rwl {
		r := x.Rule
		addr := r.Address
		if r.Netmask != "" {
			addr = r.Address + " / " + r.Netmask
		}
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\n", x.Index, r.Type, r.Database, r.User, addr, r.Method)
	}
	tw.Flush()
}

// WriteRulesTableNoIndex prints rules in a table without the index column, for copy-paste friendly output.
// Header and separator lines are prefixed with "# " so they become pg_hba.conf comments when pasted.
// No tab after "# " so TYPE aligns with the first data column.
func WriteRulesTableNoIndex(w io.Writer, rwl []hba.RuleWithLine) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "# TYPE\tDATABASE\tUSER\tADDRESS\tMETHOD")
	fmt.Fprintln(tw, "# ----\t--------\t----\t-------\t------")
	for _, x := range rwl {
		r := x.Rule
		addr := r.Address
		if r.Netmask != "" {
			addr = r.Address + " / " + r.Netmask
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", r.Type, r.Database, r.User, addr, r.Method)
	}
	tw.Flush()
}

// WriteRulesTableGroupedByUser prints rules in a table with "=== user: X ===" separators and index column.
// Rules should be sorted by user so that consecutive rules with the same user form a group.
func WriteRulesTableGroupedByUser(w io.Writer, rwl []hba.RuleWithLine) {
	if len(rwl) == 0 {
		return
	}
	var group []hba.RuleWithLine
	curUser := ""
	for _, x := range rwl {
		if x.Rule.User != curUser {
			if len(group) > 0 {
				fmt.Fprintf(w, "\n=== user: %s ===\n\n", curUser)
				writeRulesTableWithIndexTo(w, group)
			}
			curUser = x.Rule.User
			group = group[:0]
		}
		group = append(group, x)
	}
	if len(group) > 0 {
		fmt.Fprintf(w, "\n=== user: %s ===\n\n", curUser)
		writeRulesTableWithIndexTo(w, group)
	}
}

// WriteRulesTableGroupedByUserNoIndex prints rules grouped by user without the index column (copy-paste friendly).
// Group separators and table headers are prefixed with "# " so they are pg_hba.conf comments when pasted.
func WriteRulesTableGroupedByUserNoIndex(w io.Writer, rwl []hba.RuleWithLine) {
	if len(rwl) == 0 {
		return
	}
	var group []hba.RuleWithLine
	curUser := ""
	for _, x := range rwl {
		if x.Rule.User != curUser {
			if len(group) > 0 {
				fmt.Fprintf(w, "\n# === user: %s ===\n\n", curUser)
				writeRulesTableNoIndexTo(w, group)
			}
			curUser = x.Rule.User
			group = group[:0]
		}
		group = append(group, x)
	}
	if len(group) > 0 {
		fmt.Fprintf(w, "\n# === user: %s ===\n\n", curUser)
		writeRulesTableNoIndexTo(w, group)
	}
}

func writeRulesTableNoIndexTo(w io.Writer, rwl []hba.RuleWithLine) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "# TYPE\tDATABASE\tUSER\tADDRESS\tMETHOD")
	fmt.Fprintln(tw, "# ----\t--------\t----\t-------\t------")
	for _, x := range rwl {
		r := x.Rule
		addr := r.Address
		if r.Netmask != "" {
			addr = r.Address + " / " + r.Netmask
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", r.Type, r.Database, r.User, addr, r.Method)
	}
	tw.Flush()
}

func writeRulesTableWithIndexTo(w io.Writer, rwl []hba.RuleWithLine) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "#\tTYPE\tDATABASE\tUSER\tADDRESS\tMETHOD")
	fmt.Fprintln(tw, "-\t----\t--------\t----\t-------\t------")
	for _, x := range rwl {
		r := x.Rule
		addr := r.Address
		if r.Netmask != "" {
			addr = r.Address + " / " + r.Netmask
		}
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\n", x.Index, r.Type, r.Database, r.User, addr, r.Method)
	}
	tw.Flush()
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
