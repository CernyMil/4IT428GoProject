 UPDATE newsletters
        SET subject = $1, body = $2
        WHERE id = $3
        RETURNING id, subject, body, created_at