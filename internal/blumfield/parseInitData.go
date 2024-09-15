package blumfield

import (
	"blumfield/internal/models"
	"encoding/json"
	"net/url"
)

func ParseQueryToStruct(data string) (models.QueryData, error) {
	parsedData, _ := url.ParseQuery(data)
	queryData := models.QueryData{}

	// Manually assign fields
	queryData.QueryID = parsedData.Get("query_id")
	queryData.AuthDate = int64(parsedData.Get("auth_date")[0]) // you can use strconv.ParseInt to convert string to int
	queryData.Hash = parsedData.Get("hash")

	// Parse the 'user' field, which is a JSON string
	userData := parsedData.Get("user")
	err := json.Unmarshal([]byte(userData), &queryData.User)
	if err != nil {
		return queryData, err
	}

	return queryData, nil
}
