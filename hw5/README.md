# Миграции

Запускаем Postgres:

```bash
docker run \
    -d \
    -p 54320:5432 \
    --name instabank \
    -e POSTGRES_PASSWORD=iniT11 \
    -e PGDATA=/var/lib/postgresql/data \
    -v "$(pwd)/mntdata":/var/lib/postgresql/data \
    -v "$(pwd)/init-db":/docker-entrypoint-initdb.d \
    postgres:14.1
```

Запускаем миграцию схемы БД:

```bash
docker run \
    -v "$(pwd)/migrations":/migrations \
    migrate/migrate \
    -path=/migrations/ \
    -database "postgresql://instabank:iniT11@172.17.0.2:5432/instabank?sslmode=disable" \
    -verbose \
    up
```

Проверим:

```bash
psql -h 127.0.0.1 -p 5432 -U gopher -d gopher_corp
```

Попробуем откатиться на первую версию:

```
docker run \
    -v "$(pwd)/gopher-corp-backend/migrations":/migrations \
    migrate/migrate \
    -path=/migrations/ \
    -database postgresql://gopher:P@ssw0rd@172.17.0.2:5432/gopher_corp?sslmode=disable \
    -verbose \
    down 1
```

И вернемся на вторую:

```
docker run \
    -v $(pwd)/gopher-corp-backend/migrations:/migrations \
    migrate/migrate \
    -path=/migrations/ \
    -database postgresql://gopher:P@ssw0rd@172.17.0.2:5432/gopher_corp?sslmode=disable \
    -verbose \
    up 1
```

## Наполнение таблиц

Вручную создадим записи в таблицах `Departments` и `Positions`:

```sql
BEGIN;

INSERT INTO positions (title)
VALUES
    ('CTO'),
    ('CEO'),
    ('CSO'),
    ('Backend Dev'),
    ('Frontend Dev'),
    ('Fullstack Dev'),
    ('QA'),
    ('Technical writer')
ON CONFLICT(title) DO NOTHING;

INSERT INTO departments (id, parent_id, name)
OVERRIDING SYSTEM VALUE
VALUES
    (0, 0, 'root') ON CONFLICT(id) DO NOTHING;

INSERT INTO departments (parent_id, name)
VALUES
    (
        (
            SELECT id
            FROM departments
            WHERE name = 'root'
        ),
        'executives'
    ),
    (
        (
            SELECT id
            FROM departments
            WHERE name = 'root'
        ),
        'R&D'
    ),
    (
        (
            SELECT id
            FROM departments
            WHERE name = 'root'
        ),
        'Accounting'
    ),
    (
        (
            SELECT id
            FROM departments
            WHERE name = 'root'
        ),
        'Sales'
    ) ON CONFLICT(name) DO NOTHING;

COMMIT;
```

И также добавим менеджеров в `Employees`:

```sql
BEGIN DEFERRABLE;
    INSERT INTO employees (first_name, last_name, phone, email, salary, manager_id, department, position)
    VALUES
        (
            'Bob',
            'Morane',
            '+79231234567',
            'bmorane@gopher_corp.com',
            500000,
            42,
            (SELECT id FROM departments WHERE name = 'executives'),
            (SELECT id FROM positions WHERE title = 'CSO')
        );
    UPDATE employees
    SET manager_id = (
        SELECT id
        FROM employees
        WHERE
            first_name = 'Bob'
            AND last_name = 'Morane'
    )
    WHERE
        first_name = 'Bob'
        AND last_name = 'Morane';

    INSERT INTO employees (first_name, last_name, phone, email, salary, manager_id, department, position)
    VALUES
        (
            'Charley',
            'Bucket',
            '+79159876543',
            'cbucket@gopher_corp.com',
            1000000,
            42,
            (SELECT id FROM departments WHERE name = 'executives'),
            (SELECT id FROM positions WHERE title = 'CEO')
        );
    UPDATE employees
    SET manager_id = (
        SELECT id
        FROM employees
        WHERE
            first_name = 'Charley'
            AND last_name = 'Bucket'
    )
    WHERE
        first_name = 'Charley'
        AND last_name = 'Bucket';

    INSERT INTO employees (first_name, last_name, phone, email, salary, manager_id, department, position)
    VALUES
        (
            'Alice',
            'Liddell',
            '+79169008070',
            'aliddell@gopher_corp.com',
            500000,
            42,
            (SELECT id FROM departments WHERE name = 'executives'),
            (SELECT id FROM positions WHERE title = 'CTO')
        );
    UPDATE employees
    SET manager_id = (
        SELECT id
        FROM employees
        WHERE
            first_name = 'Alice'
            AND last_name = 'Liddell'
    )
    WHERE
        first_name = 'Alice'
        AND last_name = 'Liddell';
COMMIT;
```

Как вы думаете, что произойдет, если мы откатимся к первой миграции?

# Интеграционные тесты

В репозитории надо перейти на коммит, посвященный интеграционным тестам.

Для их запуска используйте команду:

```bash
go test --tags=integration_tests ./... -v -count=1
```

# Хранимые функции

Скопируем новые миграции из `dept_budget_migrations` в `gopher-corp-backend/migrations`.

Запустим БД и накатим миграции:

```bash
docker run \
    -v "$(pwd)/gopher-corp-backend/migrations":/migrations \
    migrate/migrate \
    -path=/migrations/ \
    -database postgresql://gopher:P@ssw0rd@172.17.0.2:5432/gopher_corp?sslmode=disable \
    -verbose \
    up
```

Реализуем функцию:

```sql
CREATE FUNCTION budget_change() RETURNS trigger AS $budget_change$
    BEGIN
        INSERT INTO departments_budget(id, budget)
        VALUES(NEW.department, 0)
        ON CONFLICT DO NOTHING;

        UPDATE departments_budget db
        SET budget = budget + NEW.salary
        WHERE NEW.department = db.id;

        RETURN NEW;
    END;
$budget_change$ LANGUAGE plpgsql;

CREATE TRIGGER budget_change BEFORE INSERT ON employees
    FOR EACH ROW EXECUTE FUNCTION budget_change();
```

Вставим первых сотрудников (можно использовать `prepopulate_db.sql`) и посмотрим на результат.

