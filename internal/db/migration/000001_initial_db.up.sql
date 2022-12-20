CREATE TABLE IF NOT EXISTS Album (
    "albumName" VARCHAR(100) PRIMARY KEY
    );

CREATE TABLE IF NOT EXISTS Image (
    "imageName" TEXT NOT NULL UNIQUE,
    "albumName" VARCHAR(100) NOT NULL,
    "image" TEXT
);

ALTER TABLE Image ADD FOREIGN KEY ("albumName") REFERENCES Album ("albumName") ;