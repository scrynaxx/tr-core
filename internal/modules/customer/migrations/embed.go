package migrations

import "embed"

// Files содержит SQL-миграции customer.
//
//go:embed *.sql
var Files embed.FS
