package migrations

import "embed"

// FS holds our embedded SQL migration files.
//go:embed *.sql
var FS embed.FS