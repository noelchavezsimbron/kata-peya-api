CREATE TABLE IF NOT EXISTS pets
(
    id         INT PRIMARY KEY,
    name       VARCHAR(50) NOT NULL,
    vaccines   VARCHAR(100),
    age_months INT         NOT NULL
);


INSERT INTO pets.pets (id, name, vaccines, age_months) VALUES (1, 'Rino', 'distemper,parvovirus', 36);
INSERT INTO pets.pets (id, name, vaccines, age_months) VALUES (2, 'Braco', 'rabies', 3);
INSERT INTO pets.pets (id, name, vaccines, age_months) VALUES (3, 'Duke', null, 12);
INSERT INTO pets.pets (id, name, vaccines, age_months) VALUES (4, 'Bolt', 'rabies,distemper,parvovirus', 3);
