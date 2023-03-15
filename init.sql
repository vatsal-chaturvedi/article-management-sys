CREATE DATABASE IF NOT EXISTS articleDb;
USE articleDb;

CREATE TABLE articleTable (
                         id VARCHAR(255) NOT NULL PRIMARY KEY,
                         title VARCHAR(255) NOT NULL,
                         author VARCHAR(255) NOT NULL,
                         content TEXT NOT NULL
);
