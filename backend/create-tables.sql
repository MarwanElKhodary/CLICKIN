DROP TABLE IF EXISTS `count_table`;
CREATE TABLE `count_table` (
    `id` INT(11) NOT NULL AUTO_INCREMENT,
    `slot` INT(11) NOT NULL DEFAULT 0,
    `count` INT UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`)
);

INSERT INTO `count_table`
    (slot, count)
VALUES
    (RAND() * 100, 1),
    (RAND() * 100, 1),
    (RAND() * 100, 1);