package service

import (
	"context"

	"editor-service/service/errors"
	"editor-service/service/model"
)

var (
	Editors = map[string]model.Editor{}
)

// CreateEditor saves Editor in map under email as a key.
func (Service) CreateEditor(_ context.Context, Editor model.Editor) error {
	if _, exists := Editors[Editor.Email]; exists {
		return errors.ErrEditorAlreadyExists
	}

	Editors[Editor.Email] = Editor

	return nil
}

// ListEditors returns list of Editors in array of Editors.
func (Service) ListEditors(_ context.Context) []model.Editor {
	EditorsList := make([]model.Editor, 0, len(Editors))
	for _, Editor := range Editors {
		EditorsList = append(EditorsList, Editor)
	}

	return EditorsList
}

// GetEditor returns an Editor with specified email.
func (Service) GetEditor(_ context.Context, email string) (model.Editor, error) {
	Editor, exists := Editors[email]

	if !exists {
		return model.Editor{}, errors.ErrEditorDoesntExists
	}

	return Editor, nil
}

// UpdateEditor updates attributes of a specified Editor.
func (Service) UpdateEditor(_ context.Context, email string, Editor model.Editor) (model.Editor, error) {
	oldEditor, exists := Editors[email]

	if !exists {
		return model.Editor{}, errors.ErrEditorDoesntExists
	}

	if oldEditor.Email == Editor.Email {
		Editors[email] = Editor
	} else {
		Editors[Editor.Email] = Editor

		delete(Editors, email)
	}

	return Editor, nil
}

// DeleteEditor deletes Editor from memory.
func (Service) DeleteEditor(_ context.Context, email string) error {
	if _, exists := Editors[email]; !exists {
		return errors.ErrEditorDoesntExists
	}

	delete(Editors, email)

	return nil
}
