-- SQL schema for TranslateHub
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    creator TEXT NOT NULL,
    filename TEXT NOT NULL,
    status TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS dialogues (
    id SERIAL PRIMARY KEY,
    task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    trans TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS changes (
    id SERIAL PRIMARY KEY,
    task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    dialogue_id INTEGER REFERENCES dialogues(id) ON DELETE CASCADE,
    "user" TEXT NOT NULL,
    new_trans TEXT NOT NULL,
    status TEXT NOT NULL
);
