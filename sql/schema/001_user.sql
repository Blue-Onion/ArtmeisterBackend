-- +goose Up
CREATE TABLE USERS (
    id uuid PRIMARY KEY,
    name TEXT NOT NULL, 
    password TEXT NOT NULL,
    email TEXT UNIQUE not Null, 
    createdAt TIME NOT NULL, 
    updatedAt TIME NOT NULL
);
-- +goose Down
Drop table USERS;
