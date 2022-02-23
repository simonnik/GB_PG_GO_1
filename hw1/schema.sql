-- Возможность добавлять пост с фотографиями в свой профиль. Основные действия
-- Авторизация/Регистрация
-- CRUD поста
-- Активность поста
-- Лайк поста
-- CRUD сомментария к посту

CREATE USER instabank WITH PASSWORD 'iniT11';
CREATE DATABASE instabank
    TEMPLATE = 'template0'
    ENCODING = 'utf-8'
    LC_COLLATE = 'C.UTF-8'
    LC_CTYPE = 'C.UTF-8';


\connect instabank

CREATE TABLE users
(
    id         INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    first_name VARCHAR(200),
    last_name  VARCHAR(200),
    username   VARCHAR(200),
    email      INT
);

CREATE TABLE posts
(
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id     INT CONSTRAINT users_posts_id_fk REFERENCES users,
    description TEXT,
    is_active   BOOLEAN
);

CREATE TABLE posts_images
(
    id         INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    post_id    INT CONSTRAINT posts_posts_images_fk REFERENCES posts,
    path       VARCHAR(255)
);

CREATE TABLE posts_favorites
(
    user_id INT CONSTRAINT users_posts_id_fk REFERENCES users,
    post_id INT CONSTRAINT posts_posts_images_fk REFERENCES posts,
    PRIMARY KEY (user_id, post_id)
);

CREATE TABLE posts_comments
(
    id      INT GENERATED ALWAYS AS IDENTITY,
    user_id INT CONSTRAINT users_posts_id_fk REFERENCES users,
    post_id INT CONSTRAINT posts_posts_images_fk REFERENCES posts,
    text    TEXT
);