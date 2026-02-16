package hba

import "sort"

// SortColumns are the allowed column names for sorting.
var SortColumns = []string{"type", "database", "user", "address", "method"}

// SortRules sorts rules by the given column (type, database, user, address, method).
// Order in pg_hba.conf matters for matching; this is for display only.
func SortRules(rules []Rule, by string) {
	switch by {
	case "type":
		sort.Slice(rules, func(i, j int) bool { return rules[i].Type < rules[j].Type })
	case "database":
		sort.Slice(rules, func(i, j int) bool { return rules[i].Database < rules[j].Database })
	case "user":
		sort.Slice(rules, func(i, j int) bool { return rules[i].User < rules[j].User })
	case "address":
		sort.Slice(rules, func(i, j int) bool { return rules[i].Address < rules[j].Address })
	case "method":
		sort.Slice(rules, func(i, j int) bool { return rules[i].Method < rules[j].Method })
	}
}

// SortRulesWithLine sorts RuleWithLine slice by the Rule field (same columns as SortRules).
func SortRulesWithLine(rwl []RuleWithLine, by string) {
	switch by {
	case "type":
		sort.Slice(rwl, func(i, j int) bool { return rwl[i].Rule.Type < rwl[j].Rule.Type })
	case "database":
		sort.Slice(rwl, func(i, j int) bool { return rwl[i].Rule.Database < rwl[j].Rule.Database })
	case "user":
		sort.Slice(rwl, func(i, j int) bool { return rwl[i].Rule.User < rwl[j].Rule.User })
	case "address":
		sort.Slice(rwl, func(i, j int) bool { return rwl[i].Rule.Address < rwl[j].Rule.Address })
	case "method":
		sort.Slice(rwl, func(i, j int) bool { return rwl[i].Rule.Method < rwl[j].Rule.Method })
	}
}

// ValidSortColumn returns true if col is one of the sortable columns.
func ValidSortColumn(col string) bool {
	for _, c := range SortColumns {
		if col == c {
			return true
		}
	}
	return false
}
