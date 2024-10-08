-- +goose Up
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    ticket_id UUID NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now(),  
    created_by VARCHAR(100) NULL,             
    updated_at TIMESTAMP NULL,  
    updated_by VARCHAR(100) NULL,             
    deleted_at TIMESTAMP NULL,                    
    deleted_by VARCHAR(100) NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (ticket_id) REFERENCES tickets(id)
);

-- +goose Down
DROP TABLE IF EXISTS payments;
