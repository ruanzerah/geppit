CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
   id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
   username VARCHAR(80) NOT NULL UNIQUE,
   email VARCHAR(90) NOT NULL UNIQUE,
   hash VARCHAR(255) NOT NULL UNIQUE,
   balance int8 NOT NULL DEFAULT 0,
   created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
   updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);