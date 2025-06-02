package query

import (
	_ "embed"
)

var (
	// Embed SQL files

	//go:embed scripts/ReadUser.sql
	ReadUser string

	//go:embed scripts/ReadNewsletter.sql
	ListUser string

	//go:embed scripts/ReadNewsletter.sql
	ReadNewsletter string

	//go:embed scripts/ListNewsletters.sql
	ListNewsletters string

	//go:embed scripts/InsertNewsletter.sql
	InsertNewsletter string

	//go:embed scripts/UpdateNewsletter.sql
	UpdateNewsletter string
)
