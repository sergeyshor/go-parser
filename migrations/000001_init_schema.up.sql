CREATE TABLE "books" (
  "id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "author" varchar NOT NULL,
  "price" integer NOT NULL
);