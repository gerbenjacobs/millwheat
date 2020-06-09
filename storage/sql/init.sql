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
