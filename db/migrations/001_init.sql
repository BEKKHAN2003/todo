CREATE TABLE IF NOT EXISTS users
(
    id       serial PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255)        NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks
(
    id         serial PRIMARY KEY,
    user_id    int references users (id) on delete cascade not null,
    title      TEXT                                        NOT NULL,
    completed  BOOLEAN DEFAULT FALSE,
    deadline   TIMESTAMPTZ,
    is_overdue BOOLEAN DEFAULT FALSE
);


