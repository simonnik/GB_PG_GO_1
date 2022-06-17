BEGIN;

CREATE TABLE IF NOT EXISTS users
(
    id         INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    first_name VARCHAR(200),
    last_name  VARCHAR(200),
    username   VARCHAR(200) NOT NULL,
    email      VARCHAR(200) NOT NULL,
    password   VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS posts
(
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id     INT NOT NULL,
    description TEXT NOT NULL,
    is_active   BOOLEAN,
    CONSTRAINT posts_fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS posts_images
(
    id         INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    post_id    INT NOT NULL,
    path       VARCHAR(255) NOT NULL,
    CONSTRAINT posts_images_fk_post_id FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS posts_favorites
(
    user_id INT NOT NULL,
    post_id INT NOT NULL,
    PRIMARY KEY (user_id, post_id),
    CONSTRAINT posts_favorites_fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT posts_favorites_fk_post_id FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS posts_comments
(
    id      INT GENERATED ALWAYS AS IDENTITY,
    user_id INT NOT NULL,
    post_id INT NOT NULL,
    text    TEXT NOT NULL,
    CONSTRAINT posts_comments_fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT posts_comments_fk_post_id FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS users_email_idx ON users (email);
CREATE INDEX IF NOT EXISTS users_username_idx ON users (username);
CREATE INDEX IF NOT EXISTS posts_user_id_idx ON posts (user_id);

COMMIT;