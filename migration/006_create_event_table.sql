-- +goose Up
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL,
    location VARCHAR(255) NOT NULL,
    available_tickets INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),  
    created_by VARCHAR(100) NULL,             
    updated_at TIMESTAMP NULL,  
    updated_by VARCHAR(100) NULL,             
    deleted_at TIMESTAMP NULL,                    
    deleted_by VARCHAR(100) NULL
);

-- +goose Down
DROP TABLE IF EXISTS events;
