CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),        
    username STRING NOT NULL,      
    email STRING UNIQUE NOT NULL,  
    age INT,                       
    created_at TIMESTAMP DEFAULT now() 
);