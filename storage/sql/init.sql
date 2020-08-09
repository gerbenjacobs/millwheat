CREATE TABLE `users`
(
    `id`          binary(16)   NOT NULL,
    `email`       varchar(100) NOT NULL,
    `password`    varchar(100) NOT NULL,
    `token`       varchar(255) NOT NULL,
    `currentTown` binary(16)   NOT NULL,
    `createdAt`   datetime     NOT NULL,
    `updatedAt`   datetime     NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

ALTER TABLE `users`
    ADD PRIMARY KEY (`id`),
    ADD UNIQUE KEY `unique_email` (`email`);
COMMIT;

CREATE TABLE `towns`
(
    `id`        binary(16)   NOT NULL,
    `owner`     binary(16)   NOT NULL,
    `name`      varchar(100) NOT NULL,
    `warehouse` json         NOT NULL,
    `createdAt` datetime     NOT NULL,
    `updatedAt` datetime     NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

ALTER TABLE `towns`
    ADD PRIMARY KEY (`id`),
    ADD INDEX (`owner`);
COMMIT;

CREATE TABLE `buildings`
(
    `id`             binary(16)   NOT NULL,
    `townId`         binary(16)   NOT NULL,
    `type`           int unsigned NOT NULL,
    `level`          int unsigned NOT NULL,
    `lastCollection` datetime     NOT NULL,
    `createdAt`      datetime     NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

ALTER TABLE `buildings`
    ADD PRIMARY KEY (`id`),
    ADD INDEX (`townId`);
COMMIT;

CREATE TABLE `jobs`
(
    `id`        binary(16)   NOT NULL,
    `townId`    binary(16)   NOT NULL,
    `type`      int unsigned NOT NULL,
    `jobData`   json         NOT NULL,
    `queued`    datetime     NOT NULL,
    `started`   datetime     NOT NULL,
    `completed` datetime     NOT NULL,
    `status`    int unsigned NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

ALTER TABLE `jobs`
    ADD PRIMARY KEY (`id`),
    ADD INDEX (`townId`);
COMMIT;

CREATE TABLE `warriors`
(
    `battleId`    binary(16)   NOT NULL,
    `armyId`      binary(16)   NOT NULL,
    `townId`      binary(16)   NOT NULL,
    `warriorType` int unsigned NOT NULL,
    `quantity`    int          NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

ALTER TABLE `warriors`
    ADD PRIMARY KEY (battleId, armyId, townId, warriorType);
COMMIT;