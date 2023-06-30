CREATE TABLE users (
    id SERIAL UNIQUE,
    email character varying(30) NOT NULL UNIQUE
);
