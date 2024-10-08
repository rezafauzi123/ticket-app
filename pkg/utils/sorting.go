package utils

import "fmt"

func ApplySorting(query string, sortBy, sortOrder string) string {
	if sortBy != "" {
		query += fmt.Sprintf(" ORDER BY %s", sortBy)
		if sortOrder == "desc" {
			query += " DESC"
		} else {
			query += " ASC"
		}
	}
	return query
}
