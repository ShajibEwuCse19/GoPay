package migrations

import "embed"

// FS embeds all .sql migration files directly into the compiled Go binary.
//
//go:embed *.sql
var FS embed.FS
