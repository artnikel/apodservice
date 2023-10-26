CREATE TABLE apod (
    id              SERIAL PRIMARY KEY,
    copyright       VARCHAR(200),
    date            DATE NOT NULL,
    explanation     VARCHAR(2000) NOT NULL,
    media_type      VARCHAR(20) NOT NULL,
    service_version VARCHAR(10) NOT NULL,
    title           VARCHAR(200) NOT NULL,
    url             VARCHAR(200) NOT NULL
);
