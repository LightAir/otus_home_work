CREATE TABLE IF NOT EXISTS events
(
    id             uuid PRIMARY KEY,
    title          varchar      NOT NULL,
    datetime_start timestamptz  NOT NULL,
    datetime_end   timestamptz  NOT NULL,
    description    text         NULL,
    user_id        uuid         NOT NULL,
    when_to_notify varchar(256) NULL
);