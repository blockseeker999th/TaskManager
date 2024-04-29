CREATE TABLE IF NOT EXISTS users
(
    id        SERIAL PRIMARY KEY,
    email     VARCHAR(255) NOT NULL,
    firstname VARCHAR(255) NOT NULL,
    lastname  VARCHAR(255) NOT NULL,
    password  VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (email)
);


CREATE TABLE IF NOT EXISTS projects
(
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    assignedToID INT          NOT NULL,
    createdAt    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

    /*UNIQUE (name, assignedToID)*/
    FOREIGN KEY (assignedToID) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS tasks
(
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    status       VARCHAR(20)  NOT NULL DEFAULT 'TODO',
    projectId    INT          NOT NULL,
    assignedToID INT          NOT NULL,
    createdAt    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT valid_status CHECK (status IN ('TODO', 'IN_PROGRESS', 'IN_TESTING', 'DONE')),

    FOREIGN KEY (assignedToID) REFERENCES users (id),
    FOREIGN KEY (projectId) REFERENCES projects (id)
);