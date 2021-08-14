CREATE DATABSE IF NOT EXISTS findmypet;

CREATE TABLE IF NOT EXISTS UserPet(
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(60) UNIQUE NOT NULL,
    password_hash VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS Post(
    id UUID NOT NULL PRIMARY KEY,
    title VARCHAR(120) NOT NULL, 
    description TEXT
);

CREATE TABLE IF NOT EXISTS UserPost(
    id uuid NOT NULL primary key,
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    CONSTRAINT user_fk FOREIGN KEY(user_id) REFERENCES UserPet(id),
    CONSTRAINT post_fk FOREIGN KEY(post_id) REFERENCES Post(id)
);
