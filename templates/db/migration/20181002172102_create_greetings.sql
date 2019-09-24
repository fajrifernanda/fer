-- +migrate Up notransaction
CREATE SEQUENCE IF NOT EXISTS "greetings_seq";

CREATE TABLE IF NOT EXISTS "greetings" (
    "id" BIGINT NOT NULL DEFAULT nextval('greetings_seq'),
    "name" VARCHAR(100) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NULL,
    CONSTRAINT "greetings_pkey" PRIMARY KEY ("id")
);

-- +migrate Down
DROP TABLE IF EXISTS "greetings";
DROP SEQUENCE IF EXISTS "greetings_seq";
