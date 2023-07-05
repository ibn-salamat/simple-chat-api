CREATE TABLE users (
    email character varying(30) NOT NULL UNIQUE,
    password text NOT NULL
);

CREATE TABLE users_confirmation (
    email character varying(30) NOT NULL UNIQUE,
    left_tries_count numeric DEFAULT 5 CHECK (left_tries_count <= 5 AND left_tries_count >= 0),
    confirmed boolean DEFAULT false,
    confirmation_code character(6)
);