package sqlbuilder

import (
	"fmt"
	"strings"
)

// Sangat disederhanakan: builder menerima mapping column->value placeholder,
// dan menghasilkan SQL + args untuk 4 mode (insert|upsert|update|delete).

type Stmt struct {
	SQL  string
	Args []interface{}
}

func Insert(table string, cols []string, vals []interface{}) Stmt {
	ph := make([]string, len(cols))
	for i := range cols {
		ph[i] = fmt.Sprintf("$%d", i+1)
	}
	sql := fmt.Sprintf(`insert into %s (%s) values (%s)`,
		table, strings.Join(cols, ","), strings.Join(ph, ","))
	return Stmt{SQL: sql, Args: vals}
}

func Upsert(table string, cols []string, vals []interface{}, conflictCols []string, updates []string) Stmt {
	base := Insert(table, cols, vals)
	sql := fmt.Sprintf(`%s on conflict (%s) do update set %s`,
		base.SQL,
		strings.Join(conflictCols, ","),
		strings.Join(updates, ","))
	return Stmt{SQL: sql, Args: base.Args}
}

func Update(table string, sets []string, where []string, args []interface{}) Stmt {
	sql := fmt.Sprintf(`update %s set %s where %s`,
		table, strings.Join(sets, ","), strings.Join(where, " and "))
	return Stmt{SQL: sql, Args: args}
}

func Delete(table string, where []string, args []interface{}) Stmt {
	sql := fmt.Sprintf(`delete from %s where %s`, table,
		strings.Join(where, " and "))
	return Stmt{SQL: sql, Args: args}
}
