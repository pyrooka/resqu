package db

import (
	"database/sql"
	"encoding/json"
)

// processSQLRows takes a SQL query result and reads it to the correct format.
func processSQLRows(rows *sql.Rows) (result []byte, err error) {
	columns, err := rows.Columns()
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows in the result so return an empty array without an error.
			return []byte{}, nil
		}

		return
	}
	count := len(columns)

	data := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		entry := make(map[string]interface{})

		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		data = append(data, entry)
	}

	result, err = json.Marshal(data)

	return
}
