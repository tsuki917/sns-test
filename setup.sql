-- DROP table posts cascade if exists;
-- drop table comments if exists;
SET TIME ZONE  'Asia/Tokyo';
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS comments;


CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    content TEXT,
    author VARCHAR(255),
    day_post  TEXT
);


CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    content TEXT,
    author VARCHAR(255),
    post_id INTEGER REFERENCES posts(id),
    day_comment TEXT
);
