package service

import (
	"embed"
)

//go:embed templates/confirmation_request.html
var templateFS_ConfReq embed.FS

//go:embed templates/confirmation.html
var templateFS_Conf embed.FS

//go:embed templates/newsletter_post.html
var templateFS_NewsPost embed.FS
