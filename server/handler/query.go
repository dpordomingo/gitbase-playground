package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/src-d/gitbase-playground/server/serializer"
)

type queryRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

// genericVals returns a slice of interface{}, each one a pointer to a NullString
func genericVals(nColumns int) []interface{} {
	columnVals := make([]sql.NullString, nColumns)
	columnValsPtr := make([]interface{}, nColumns)

	for i := range columnVals {
		columnValsPtr[i] = &columnVals[i]
	}

	return columnValsPtr
}

// Query returns a function that forwards an SQL query to gitbase and returns
// the rows as JSON
func Query(db *sql.DB) RequestProcessFunc {
	return func(r *http.Request) (*serializer.Response, error) {
		var queryRequest queryRequest
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err == nil {
			err = json.Unmarshal(body, &queryRequest)
		}

		if err != nil {
			return nil, err
		}

		// TODO (carlosms) this only works if the query does not end in ;
		// and does not have a limit. It will also fail for queries like
		// DESCRIBE TABLE
		query := fmt.Sprintf("%s LIMIT %d", queryRequest.Query, queryRequest.Limit)
		rows, err := db.Query(query)
		if err != nil {
			return nil, serializer.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		defer rows.Close()

		columnNames, _ := rows.Columns()
		nColumns := len(columnNames)
		columnValsPtr := genericVals(nColumns)

		tableData := make([]map[string]string, 0)

		for rows.Next() {
			if err := rows.Scan(columnValsPtr...); err != nil {
				return nil, err
			}

			colData := make(map[string]string)

			for i, val := range columnValsPtr {
				var st string
				sqlSt, _ := val.(*sql.NullString)

				if sqlSt.Valid {
					st = sqlSt.String
				}

				colData[columnNames[i]] = st
			}

			tableData = append(tableData, colData)
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

		return serializer.NewQueryResponse(tableData, columnNames), nil
	}
}
