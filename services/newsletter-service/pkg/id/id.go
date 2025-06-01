package id

import (
	"fmt"

	"github.com/google/uuid"
)

// Specific ID types
type (
	Newsletter uuid.UUID
	Post       uuid.UUID
)

// Newsletter methods
func (n *Newsletter) FromString(str string) error {
	return fromString((*uuid.UUID)(n), str)
}

func (n Newsletter) String() string {
	return uuid.UUID(n).String()
}

func (n *Newsletter) Scan(data any) error {
	return scanUUID((*uuid.UUID)(n), "Newsletter", data)
}

func (n Newsletter) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(n).String()), nil
}

func (n *Newsletter) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(n), "Newsletter", data)
}

func NewNewsletter() Newsletter {
	return Newsletter(uuid.New())
}

// Post methods
func (p *Post) FromString(str string) error {
	return fromString((*uuid.UUID)(p), str)
}

func (p Post) String() string {
	return uuid.UUID(p).String()
}

func (p *Post) Scan(data any) error {
	return scanUUID((*uuid.UUID)(p), "Post", data)
}

func (p Post) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(p).String()), nil
}

func (p *Post) UnmarshalText(data []byte) error {
	return unmarshalUUID((*uuid.UUID)(p), "Post", data)
}

func NewPost() Post {
	return Post(uuid.New())
}

// Common methods for all ID types
func fromString(u *uuid.UUID, s string) error {
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}
	*u = uuid.UUID(id)
	return nil
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
