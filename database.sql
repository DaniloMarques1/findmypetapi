CREATE TABLE IF NOT EXISTS UserPet(
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(60) UNIQUE NOT NULL,
    password_hash VARCHAR(100) NOT NULL
);

DO $$
BEGIN
    IF NOT EXISTS (select * from pg_type where pg_type.typname = 'petstatus') THEN
        CREATE TYPE PetStatus AS ENUM('missing', 'found');
    END IF ;
END$$;

CREATE TABLE IF NOT EXISTS Post(
    id UUID NOT NULL PRIMARY KEY,
    author_id UUID Not null,
    title VARCHAR(120) NOT NULL, 
    description TEXT,
    image_url VARCHAR(100),
    status PetStatus default 'missing',
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT post_user_fk FOREIGN KEY(author_id) REFERENCES UserPet(id)
);

CREATE TABLE IF NOT EXISTS Comment(
    id UUID NOT NULL PRIMARY KEY,
    comment_text TEXT NOT NULL,
    author_id UUID NOT NULL,
    post_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT user_fk FOREIGN KEY(author_id) REFERENCES UserPet(id),
    CONSTRAINT post_fk FOREIGN KEY(post_id) REFERENCES Post(id)
);
