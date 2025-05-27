package api

import (
	v1 "editor-service/transport/api/v1"
	"net/http"
)

type Controller struct {
	handler *v1.EditorHandler
}

func NewController(handler *v1.EditorHandler) *Controller {
	return &Controller{handler: handler}
}

func (c *Controller) SignUp(w http.ResponseWriter, r *http.Request) {
	c.handler.SignUp(w, r)
}

func (c *Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	c.handler.SignIn(w, r)
}

func (c *Controller) ChangePassword(w http.ResponseWriter, r *http.Request) {
	c.handler.ChangePassword(w, r)
}
