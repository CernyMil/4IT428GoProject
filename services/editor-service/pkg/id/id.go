package id

import (
	"fmt"

	"github.com/google/uuid"
)

type Editor uuid.UUID

func (u *Editor) FromString(s string) error {
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*u = Editor(id)
	return nil
}

func (u Editor) String() string {
	return uuid.UUID(u).String()
}

func (u *Editor) Scan(data any) error {
	return scanUUID((*uuid.UUID)(u), "Editor", data)
}

func (u Editor) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(u).String()), nil
}

func (u *Editor) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(u), "Editor", data)
}

func scanUUID(u *uuid.UUID, idTypeName string, data any) error {
	if err := u.Scan(data); err != nil {
		return fmt.Errorf("scanning %q id value: %w", idTypeName, err)
	}
	return nil
}

func unmarshalUUID(u *uuid.UUID, idTypeName string, data []byte) error {
	if err := u.UnmarshalText(data); err != nil {
		return fmt.Errorf("parsing %q id value: %w", idTypeName, err)
	}
	return nil
}
