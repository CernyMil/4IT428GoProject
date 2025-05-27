package query

const (
	// Query to read a user by ID
	ReadUser = `
        SELECT id, email, first_name, last_name, created_at
        FROM users
        WHERE id = :id
    `

	// Query to list all users
	ListUser = `
        SELECT id, email, first_name, last_name, created_at
        FROM users
        ORDER BY created_at DESC
    `

	// Query to read a newsletter by ID
	ReadNewsletter = `
        SELECT id, subject, body, created_at
        FROM newsletters
        WHERE id = $1
    `

	// Query to list all newsletters
	ListNewsletters = `
        SELECT id, subject, body, created_at
        FROM newsletters
        ORDER BY created_at DESC
    `

	// Query to insert a new newsletter
	InsertNewsletter = `
        INSERT INTO newsletters (id, subject, body, created_at)
        VALUES ($1, $2, $3, $4)
    `

	// Query to update a newsletter
	UpdateNewsletter = `
        UPDATE newsletters
        SET subject = $1, body = $2
        WHERE id = $3
        RETURNING id, subject, body, created_at
    `
)
