
CREATE TABLE general_chat_messages (
    id UUID NOT NULL,
    email character varying(30) NOT NULL,
    message_type character varying(30) NOT NULL,
    message_content text,
  	created_at timestamp
);
