package query

import (
	_ "embed"
)

var (
	//go:embed query/CreateEditor.sql
	CreateEditor string
	//go:embed query/GetEditorByEmail.sql
	GetEditor string
)
