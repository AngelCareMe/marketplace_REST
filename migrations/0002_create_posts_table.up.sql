CREATE TABLE posts (
    id UUID PRIMARY KEY,
    header VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    image VARCHAR(255),
    price FLOAT8 NOT NULL,
    author_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);