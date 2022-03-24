CREATE TABLE players (
    id INT NOT NULL AUTO_INCREMENT,
    name CHAR(30) NOT NULL,
    in_game BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (id)
);

CREATE TABLE games (
    id INT NOT NULL AUTO_INCREMENT,
    title CHAR(30) NOT NULL,
    PRIMARY KEY (id)
);

INSERT INTO players (name) VALUES ('Dutch'), ('Bill'), ('Mika'), ('Javier'), ('Arthur'), ('Hosea');
INSERT INTO games (title) VALUES ('Poker'), ('Blackjack'), ('Five finger fillet'), ('Domino');