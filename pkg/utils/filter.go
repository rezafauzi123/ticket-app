package utils

import "fmt"

func ApplyFilters(query string, args []interface{}, filters map[string]interface{}, idx *int) (string, []interface{}) {
	for key, value := range filters {
		if value != nil {
			switch v := value.(type) {
			case string:
				if v != "" {
					query += fmt.Sprintf(" AND %s ILIKE $%d", key, *idx)
					args = append(args, "%"+v+"%")
					*idx++
				}
			case *int:
				if v != nil {
					query += fmt.Sprintf(" AND %s = $%d", key, *idx)
					args = append(args, *v)
					*idx++
				}
			}
		}
	}
	return query, args
}
