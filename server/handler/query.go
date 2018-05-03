package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dpordomingo/gitbase-playground/server/serializer"
	"gopkg.in/bblfsh/sdk.v1/uast"

	"github.com/go-sql-driver/mysql"
)

type queryRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

// genericVals returns a slice of interface{}, each one a pointer to the proper
// type for each column
func genericVals(colTypes []*sql.ColumnType) []interface{} {
	columnValsPtr := make([]interface{}, len(colTypes))

	for i, colType := range colTypes {
		switch colType.DatabaseTypeName() {
		case "BIT":
			columnValsPtr[i] = new(sql.NullBool)
		case "TIMESTAMP":
			columnValsPtr[i] = new(mysql.NullTime)
		case "INT":
			columnValsPtr[i] = new(sql.NullInt64)
		case "JSON":
			columnValsPtr[i] = new([]byte)
		default: // "TEXT" and any others
			columnValsPtr[i] = new(sql.NullString)
		}
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

		columnNames, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		colTypes, err := rows.ColumnTypes()
		if err != nil {
			return nil, err
		}
		columnValsPtr := genericVals(colTypes)

		tableData := make([]map[string]interface{}, 0)

		for rows.Next() {
			if err := rows.Scan(columnValsPtr...); err != nil {
				return nil, err
			}

			colData := make(map[string]interface{})

			for i, val := range columnValsPtr {
				colData[columnNames[i]] = nil

				switch val.(type) {
				case *sql.NullBool:
					sqlVal, _ := val.(*sql.NullBool)
					if sqlVal.Valid {
						colData[columnNames[i]] = sqlVal.Bool
					}
				case *mysql.NullTime:
					sqlVal, _ := val.(*mysql.NullTime)
					if sqlVal.Valid {
						colData[columnNames[i]] = sqlVal.Time
					}
				case *sql.NullInt64:
					sqlVal, _ := val.(*sql.NullInt64)
					if sqlVal.Valid {
						colData[columnNames[i]] = sqlVal.Int64
					}
				case *sql.NullString:
					sqlVal, _ := val.(*sql.NullString)
					if sqlVal.Valid {
						colData[columnNames[i]] = sqlVal.String
					}
				case *[]byte:
					// TODO (carlosms) this may not be an array always
					var protobufs [][]byte
					if err := json.Unmarshal(*val.(*[]byte), &protobufs); err != nil {
						return nil, err
					}

					nodes := make([]*uast.Node, len(protobufs))

					for i, v := range protobufs {
						node := uast.NewNode()
						if err = node.Unmarshal(v); err != nil {
							return nil, err
						}
						nodes[i] = node
					}

					colData[columnNames[i]] = nodes
				}
			}

			tableData = append(tableData, colData)
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

		return serializer.NewQueryResponse(tableData, columnNames), nil
	}
}
