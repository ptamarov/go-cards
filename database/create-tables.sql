CREATE TABLE IF NOT EXISTS history (
        user_id   VARCHAR(36) NOT NULL, 
        deck_id   TEXT NOT NULL, 
        card_id   TEXT NOT NULL, 
        guess     TEXT NOT NULL, 
        duration  REAL NOT NULL, 
        time_stamp TEXT NOT NULL
    );
    
CREATE TABLE IF NOT EXISTS cards (
    card_id         VARCHAR(36) NOT NULL, 
    lang_learn      TEXT NOT NULL, 
    lang_user       TEXT NOT NULL, 
    prompt          TEXT NOT NULL,
    hint            TEXT NOT NULL,
    translation     TEXT NOT NULL,
    grammar         TEXT NOT NULL
    PRIMARY KEY(card_id)
);