CREATE USER instabank WITH PASSWORD 'iniT11';

CREATE DATABASE instabank
    WITH OWNER instabank
    TEMPLATE = 'template0'
    ENCODING = 'utf-8'
    LC_COLLATE = 'C.UTF-8'
    LC_CTYPE = 'C.UTF-8';