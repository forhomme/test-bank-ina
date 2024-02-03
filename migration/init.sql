CREATE TABLE IF NOT EXISTS `users` (
    id BIGINT AUTO_INCREMENT NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(64) NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS `tasks` (
    id BIGINT AUTO_INCREMENT NOT NULL,
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

DROP PROCEDURE IF EXISTS `sp_DeleteUser`;
CREATE PROCEDURE `sp_DeleteUser`
(
    IN userId int
)
BEGIN
    START TRANSACTION;
    SET SESSION FOREIGN_KEY_CHECKS=OFF;
    UPDATE `tasks` set user_id = "" WHERE user_id = userId;
    DELETE FROM `users` WHERE id = userId;
    SET SESSION FOREIGN_KEY_CHECKS=ON;
    COMMIT;
END;

/*https://dba.stackexchange.com/questions/24531/mysql-create-index-if-not-exists*/
DROP PROCEDURE IF EXISTS `sp_CreateIndex`;
CREATE PROCEDURE `sp_CreateIndex`
(
    given_database VARCHAR(64),
    given_table    VARCHAR(64),
    given_index    VARCHAR(64),
    given_columns  VARCHAR(64)
)
BEGIN

    DECLARE IndexIsThere INTEGER;

SELECT COUNT(1) INTO IndexIsThere
FROM INFORMATION_SCHEMA.STATISTICS
WHERE table_schema = given_database
  AND   table_name   = given_table
  AND   index_name   = given_index;

IF IndexIsThere = 0 THEN
        SET @sqlstmt = CONCAT('CREATE INDEX ',given_index,' ON ',
        given_database,'.',given_table,' (',given_columns,')');
PREPARE st FROM @sqlstmt;
EXECUTE st;
DEALLOCATE PREPARE st;
ELSE
SELECT CONCAT('Index ',given_index,' already exists on Table ',
              given_database,'.',given_table) CreateindexErrorMessage;
END IF;

END;

DROP PROCEDURE IF EXISTS `sp_CreateUniqueIndex`;
CREATE PROCEDURE `sp_CreateUniqueIndex`
(
    given_database VARCHAR(64),
    given_table    VARCHAR(64),
    given_index    VARCHAR(64),
    given_columns  VARCHAR(64)
)
BEGIN

    DECLARE IndexIsThere INTEGER;

SELECT COUNT(1) INTO IndexIsThere
FROM INFORMATION_SCHEMA.STATISTICS
WHERE table_schema = given_database
  AND   table_name   = given_table
  AND   index_name   = given_index;

IF IndexIsThere = 0 THEN
        SET @sqlstmt = CONCAT('CREATE UNIQUE INDEX ',given_index,' ON ',
        given_database,'.',given_table,' (',given_columns,')');
PREPARE st FROM @sqlstmt;
EXECUTE st;
DEALLOCATE PREPARE st;
ELSE
SELECT CONCAT('Index ',given_index,' already exists on Table ',
              given_database,'.',given_table) CreateindexErrorMessage;
END IF;

END;

DROP PROCEDURE IF EXISTS `sp_AlterTable`;
CREATE PROCEDURE sp_AlterTable()
BEGIN
END;
CALL sp_AlterTable();

CALL sp_CreateUniqueIndex('bank_ina', 'users', 'unique_users_email_idx', 'email');