CREATE DATABASE IF NOT EXISTS todo_app;

CREATE TABLE IF NOT EXISTS `todo` (
  `Id`          int(11) NOT NULL AUTO_INCREMENT,
  `Title`       varchar(255) DEFAULT NULL,
  `Category`    varchar(255) DEFAULT NULL,
  `State`       varchar(255) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;