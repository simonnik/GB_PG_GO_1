CREATE USER gopher WITH PASSWORD 'iniT11';
CREATE DATABASE gopher_corp
    TEMPLATE = 'template0'
    ENCODING = 'utf-8'
    LC_COLLATE = 'C.UTF-8'
    LC_CTYPE = 'C.UTF-8';


\connect gopher_corp

CREATE TABLE departments
(
    id     INT GENERATED ALWAYS AS IDENTITY,
    parent INT,
    name   VARCHAR(200)
);

CREATE TABLE employees
(
    id         INT GENERATED ALWAYS AS IDENTITY,
    first_name VARCHAR(200),
    last_name  VARCHAR(200),
    salary     MONEY,
    manager    INT,
    department INT,
    position   INT
);

CREATE TABLE positions
(
    id    INT GENERATED ALWAYS AS IDENTITY,
    title VARCHAR(200)
);

INSERT INTO departments(parent, name)
VALUES
    (1, 'root'),
    (1, 'R&D'),
    (1, 'HOME');

INSERT INTO employees (first_name, last_name, salary, manager, department, position)
VALUES
    ('Jane', 'Doe', 75000.00, NULL, 2, 2),
    ('John', 'Doe', 50000.00, NULL, 2, 1),
    ('Piter', 'Doe', 800000.00, NULL, 3, 3);

INSERT INTO positions(title)
VALUES
    ('Software Engineer I'),
    ('Software Engineer II'),
    ('Software Engineer III');

ALTER DATABASE gopher_corp OWNER TO gopher;
GRANT ALL PRIVILEGES ON DATABASE gopher_corp to gopher;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO gopher;