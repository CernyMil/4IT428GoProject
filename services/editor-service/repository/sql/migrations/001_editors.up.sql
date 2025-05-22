CREATE TABLE editors (
    id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
    Email     string NOT NULL,
	FirstName string NOT NULL,
	LastName  string NOT NULL,
    Password  string NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL,
    updated_at  TIMESTAMPTZ  NOT NULL
);

INSERT INTO editors (
    id,
    created_at,
    updated_at
) VALUES (
    '00000000-0000-0000-0000-000000000000',
    NOW(),
    NOW()
);
