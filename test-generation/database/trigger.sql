USE learned_vocabulary;

DELIMITER //

CREATE TRIGGER before_schedule_insert
BEFORE INSERT ON schedule
FOR EACH ROW
BEGIN
    SET NEW.id = FLOOR(
        UNIX_TIMESTAMP(NOW(6)) * 10000 + 
        FLOOR(RAND() * 9999)
    );
END//

DELIMITER ;
