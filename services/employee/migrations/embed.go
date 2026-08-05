package migrations

import "embed"

// Files содержит SQL-миграции employee.
//
//go:embed *.sql
var Files embed.FS
