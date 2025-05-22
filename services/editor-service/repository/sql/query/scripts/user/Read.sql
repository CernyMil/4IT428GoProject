SELECT
	e.id,
	e.created_at,
	e.updated_at
FROM
	editors as e
WHERE
	id = @id
