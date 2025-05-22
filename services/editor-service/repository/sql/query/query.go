package query

import (
	_ "embed"
)

var (
	//go:embed scripts/editor/Read.sql
	ReadEditor string
	//go:embed scripts/editor/List.sql
	ListEditor string
)
