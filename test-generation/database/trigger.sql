USE learned_vocabulary;

DELIMITER //

CREATE TRIGGER before_schedule_insert
BEFORE INSERT ON schedule
FOR EACH ROW
BEGIN
    SET NEW.id = CAST(
        CONCAT(
            FLOOR(RAND() * 9 + 1), 
            LPAD(
                FLOOR(RAND() * 999999999), 
                9, 
                '0'
            )
        ) AS SIGNED
    );
END//

DELIMITER ;

