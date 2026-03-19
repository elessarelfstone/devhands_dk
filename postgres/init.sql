CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT
);

INSERT INTO documents (title, content) VALUES
('Docker', 'Containerization platform'),
('Compose', 'Orchestration tool for Docker');
