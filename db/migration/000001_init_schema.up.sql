CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "fullname" varchar NOT NULL,
  "username" varchar NOT NULL,
  "email" varchar NOT NULL,
  "role" varchar NOT NULL,
  "gender" varchar NOT NULL,
  "phone_number" varchar NOT NULL,
  "password" varchar NOT NULL,
  "address" varchar NOT NULL,
  "avatar" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "homes" (
  "id" uuid PRIMARY KEY,
  "owner_id" uuid NOT NULL,
  "title" varchar NOT NULL,
  "featured_image" varchar NOT NULL,
  "bedrooms" int NOT NULL,
  "bathrooms" int NOT NULL,
  "type_rent" varchar NOT NULL,
  "price" bigint NOT NULL,
  "province_id" int NOT NULL,
  "city_id" int NOT NULL,
  "description" varchar NOT NULL,
  "amenities" varchar NOT NULL,
  "area" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "images" (
  "id" uuid PRIMARY KEY,
  "house_id" uuid NOT NULL,
  "url" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "transactions" (
  "id" uuid PRIMARY KEY,
  "tenant_id" uuid NOT NULL,
  "owner_id" uuid NOT NULL,
  "house_id" uuid NOT NULL,
  "payment_status" varchar NOT NULL,
  "payment_proof" varchar NOT NULL,
  "total_payment" bigint NOT NULL,
  "check_in" timestamp NOT NULL,
  "check_out" timestamp NOT NULL,
  "time_rent" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "homes" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id");

ALTER TABLE "images" ADD FOREIGN KEY ("house_id") REFERENCES "homes" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("house_id") REFERENCES "homes" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("tenant_id") REFERENCES "users" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id");
