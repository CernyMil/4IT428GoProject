package errors

import "errors"

var (
	ErrEditorAlreadyExists = errors.New("Editor already exists")
	ErrEditorDoesntExists  = errors.New("Editor does not exist")
)
