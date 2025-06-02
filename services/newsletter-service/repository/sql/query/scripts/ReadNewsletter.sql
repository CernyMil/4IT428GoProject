SELECT id, subject, body, created_at
        FROM newsletters
        WHERE id = $1