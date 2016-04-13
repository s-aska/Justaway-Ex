BEGIN;

DROP TABLE IF EXISTS crawler;
CREATE TABLE crawler (
    id                  BIGINT UNSIGNED      NOT NULL AUTO_INCREMENT,
    url                 VARCHAR(255)
                        CHARACTER SET ascii
                        COLLATE ascii_bin    NOT NULL,
    status              ENUM(
                            'STANDBY',
                            'ACTIVE',
                            'INACTIVE'
                        )                    NOT NULL DEFAULT 'STANDBY',
    created_on          INTEGER UNSIGNED     NOT NULL,
    updated_on          INTEGER UNSIGNED     NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (url)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS account;
CREATE TABLE account (
    id                  BIGINT UNSIGNED      NOT NULL AUTO_INCREMENT,
    crawler_id          BIGINT UNSIGNED      NOT NULL,
    user_id             BIGINT UNSIGNED      NOT NULL,
    name                VARCHAR(64)          NOT NULL,
    screen_name         VARCHAR(64)
                        CHARACTER SET ascii
                        COLLATE ascii_bin    NOT NULL,
    api_token           VARCHAR(64)
                        CHARACTER SET ascii
                        COLLATE ascii_bin    NOT NULL,
    access_token        VARCHAR(64)
                        CHARACTER SET ascii
                        COLLATE ascii_bin    NOT NULL,
    access_token_secret VARCHAR(64)
                        CHARACTER SET ascii
                        COLLATE ascii_bin    NOT NULL,
    status              ENUM(
                            'ACTIVE',
                            'REVOKE',
                            'DELETE'
                        )                    NOT NULL DEFAULT 'ACTIVE',
    created_on          INTEGER UNSIGNED     NOT NULL,
    updated_on          INTEGER UNSIGNED     NOT NULL,
    revoked_on          INTEGER UNSIGNED     NOT NULL DEFAULT 0,
    deleted_on          INTEGER UNSIGNED     NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS notification_settings;
CREATE TABLE notification_settings (
    id                  BIGINT UNSIGNED      NOT NULL AUTO_INCREMENT,
    user_id             BIGINT UNSIGNED      NOT NULL,
    data                MEDIUMBLOB           NOT NULL,
    created_on          INTEGER UNSIGNED     NOT NULL,
    updated_on          INTEGER UNSIGNED     NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS notification_device;
CREATE TABLE notification_device (
    id                  BIGINT UNSIGNED      NOT NULL AUTO_INCREMENT,
    user_id             BIGINT UNSIGNED      NOT NULL,
    name                TEXT                 NOT NULL,
    token               VARCHAR(255)
                        CHARACTER SET ascii
                        COLLATE ascii_bin    NOT NULL,
    platform            ENUM('APNS', 'GCM')  NOT NULL,
    created_on          INTEGER UNSIGNED     NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (platform, token),
    INDEX (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS activity;
CREATE TABLE activity (
    id                  BIGINT UNSIGNED      NOT NULL AUTO_INCREMENT,
    event               VARCHAR(32)
                        CHARACTER SET ascii
                        COLLATE ascii_bin    NOT NULL,
    target_id           BIGINT UNSIGNED      NOT NULL,
    source_id           BIGINT UNSIGNED      NOT NULL,
    target_object_id    BIGINT UNSIGNED      NOT NULL,
    retweeted_status_id BIGINT UNSIGNED      NULL,
    created_on          INTEGER UNSIGNED     NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (target_object_id, event, source_id),
    KEY (retweeted_status_id),
    KEY (target_id)
) ENGINE=InnoDB DEFAULT CHARSET=ascii;

DROP TABLE IF EXISTS favorite;
CREATE TABLE favorite (
    id                  BIGINT UNSIGNED      NOT NULL AUTO_INCREMENT,
    target_id           BIGINT UNSIGNED      NOT NULL,
    source_user_id      BIGINT UNSIGNED      NOT NULL,
    created_on          INTEGER UNSIGNED     NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (target_id, source_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

COMMIT;