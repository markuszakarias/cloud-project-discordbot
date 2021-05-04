CREATE TABLE todo (
  Id          INT IDENTITY PRIMARY KEY,
  UserId      NVARCHAR(255) NOT NULL,
  Title       NVARCHAR(255) NOT NULL,
  Category    NVARCHAR(255) NOT NULL,
  State       NVARCHAR(255) NOT NULL,
)

INSERT INTO table1(UserId, Title, Category, State)
VALUES ('abcdefgh12345678', 'FirstTask', 'Work', 'active')