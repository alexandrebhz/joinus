package utils

import "strings"

// SanitizeOrder allowlists ORDER BY column and direction to prevent SQL injection.
// allowed maps request keys (lowercase) to safe SQL identifiers.
func SanitizeOrder(orderBy, orderDir string, allowed map[string]string, defaultCol string) (col, dir string) {
	col = defaultCol
	if allowed != nil {
		if c, ok := allowed[strings.ToLower(strings.TrimSpace(orderBy))]; ok {
			col = c
		}
	}
	dir = "DESC"
	if strings.EqualFold(strings.TrimSpace(orderDir), "ASC") {
		dir = "ASC"
	}
	return col, dir
}

// JobOrderColumns are safe sort keys for the jobs list.
var JobOrderColumns = map[string]string{
	"created_at": "created_at",
	"updated_at": "updated_at",
	"title":      "title",
	"status":     "status",
	"city":       "city",
	"country":    "country",
}

// StartupOrderColumns are safe sort keys for the startups list.
var StartupOrderColumns = map[string]string{
	"created_at": "created_at",
	"updated_at": "updated_at",
	"name":       "name",
	"slug":       "slug",
	"status":     "status",
	"industry":   "industry",
}
