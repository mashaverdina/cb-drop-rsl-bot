CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    alias  VARCHAR(128),
    text VARCHAR(4096),
    fire_at TIMESTAMP NOT NULL,
    last_fire_time TIMESTAMP NOT NULL
);

CREATE TABLE disabled_notifications (
    user_id         INT,
    notification_id INT,
    PRIMARY KEY (user_id, notification_id),
    CONSTRAINT notification_id_fk
        FOREIGN KEY(notification_id)
            REFERENCES notifications(id)
            ON UPDATE CASCADE
            ON DELETE CASCADE,
    CONSTRAINT user_id_fk
        FOREIGN KEY(user_id)
            REFERENCES users(user_id)
            ON UPDATE CASCADE
            ON DELETE CASCADE
);

CREATE INDEX notifications_alias_idx ON notifications(alias);

insert into notifications (alias, text, fire_at, last_fire_time) values ('fill_drop', 'Не забудь заполнить дроп с КБ!', '2021-01-01 13:30', '2021-01-01 13:30');