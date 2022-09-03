CREATE TABLE tasks (
    id          TEXT PRIMARY KEY,
    title       TEXT NOT NULL,
    done        INTEGER NOT NULL,
    created_ts  DATETIME NOT NULL
);
