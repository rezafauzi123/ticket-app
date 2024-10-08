-- +goose Up
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    message TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now(),  
    created_by VARCHAR(100) NULL,             
    updated_at TIMESTAMP NULL,  
    updated_by VARCHAR(100) NULL,             
    deleted_at TIMESTAMP NULL,                    
    deleted_by VARCHAR(100) NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE IF EXISTS notifications;
