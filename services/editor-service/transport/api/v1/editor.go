package v1

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"

	"editor-service/service/model"
	"editor-service/transport/util"
)

var validate = validator.New()

func getEmailFromURL(r *http.Request) string {
	email := chi.URLParam(r, "email")
	return email
}

func (h *Handler) CreateEditor(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}
	var Editor model.Editor
	if err := json.Unmarshal(b, &Editor); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}
	if err := validate.Struct(Editor); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return

	}

	if err := h.service.CreateEditor(r.Context(), Editor); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}
	util.WriteResponse(w, http.StatusCreated, Editor)
}

func (h *Handler) GetEditor(w http.ResponseWriter, r *http.Request) {
	Editor, err := h.service.GetEditor(r.Context(), getEmailFromURL(r))
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}
	util.WriteResponse(w, http.StatusOK, Editor)
}

func (h *Handler) ListEditors(w http.ResponseWriter, r *http.Request) {
	Editors := h.service.ListEditors(r.Context())
	util.WriteResponse(w, http.StatusOK, Editors)
}

func (h *Handler) UpdateEditor(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}
	var Editor model.Editor
	if err := json.Unmarshal(b, &Editor); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}
	email := Editor.Email
	EditorUpdated, err := h.service.UpdateEditor(r.Context(), email, Editor)
	if err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}
	util.WriteResponse(w, http.StatusOK, EditorUpdated)
}

func (h *Handler) DeleteEditor(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteEditor(r.Context(), getEmailFromURL(r)); err != nil {
		util.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}
}
