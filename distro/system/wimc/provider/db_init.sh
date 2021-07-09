#!/bin/sh

#Script to initialize the sqlite3 db.

sqlite3 ./file.db <<EOF

DROP TABLE IF EXISTS Containers;
CREATE TABLE Containers (Name TEXT, Tag TEXT, Type TEXT, Method TEXT, URL TEXT); 
INSERT INTO Containers VALUES('Name1', 'latest', '1', 'CHUNK' ,'/Containers/Name1');
INSERT INTO Containers VALUES('Name2', 'latest', '1', 'CA-SYNC' ,'/Containers/Name2');
INSERT INTO Containers VALUES('Name3', 'latest', '1', 'CA-SYNC' ,'/Containers/Name3');
INSERT INTO Containers VALUES('Name4', 'latest', '1', 'CA-SYNC' ,'/Containers/Name4');
INSERT INTO Containers VALUES('Name5', 'latest', '1', 'CA-SYNC' ,'/Containers/Name5');

SELECT * FROM Containers;
EOF

echo "All done."
