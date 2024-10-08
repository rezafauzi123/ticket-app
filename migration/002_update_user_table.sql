-- +goose Up
ALTER TABLE users
    ADD COLUMN address TEXT,                        
    ADD COLUMN gender VARCHAR(10),                  
    ADD COLUMN marital_status VARCHAR(15),          

    ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT now(),  
    ADD COLUMN created_by VARCHAR(100) NULL,             
    ADD COLUMN updated_at TIMESTAMP NULL,  
    ADD COLUMN updated_by VARCHAR(100) NULL,             
    ADD COLUMN deleted_at TIMESTAMP NULL,                    
    ADD COLUMN deleted_by VARCHAR(100) NULL;                 


-- +goose Down
DROP TABLE IF EXISTS users;
