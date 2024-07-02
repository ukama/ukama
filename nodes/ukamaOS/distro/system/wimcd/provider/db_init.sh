#!/bin/sh

export PATH=../../../vendor/build/bin/:$PATH

sqlite3 $1 <<EOF

DROP TABLE IF EXISTS Containers;
CREATE TABLE Containers (
  Name      TEXT,
  Version   TEXT,
  Type      TEXT,
  URL       TEXT,
  CreatedAt TEXT,
  SizeBytes INTEGER,
  ChunkURL  TEXT,
  ChunksURL TEXT
);

INSERT INTO Containers VALUES('test-capp', '0.0.3', 'tar.gz', '/v1/capps/test-capp/0.0.3.tar.gz', '2024-06-11T18:24:09.128Z', 25864577, '', '');
INSERT INTO Containers VALUES('test-capp', '0.0.3', 'chunk', '/v1/capps/test-capp/0.0.3.caibx', '0001-01-01T00:00:00Z', NULL, '/v1/capps/test-capp/0.0.3.caibx', '/v1/chunks/');
INSERT INTO Containers VALUES('test-capp', 'latest', 'chunk', '/v1/capps/test-capp/latest.caibx', '0001-01-01T00:00:00Z', NULL, '/v1/capps/test-capp/latest.caibx', '/v1/chunks/');
INSERT INTO Containers VALUES('another-capp', '1.2.0', 'tar.gz', '/v1/capps/another-capp/1.2.0.tar.gz', '2024-05-01T12:00:00.000Z', 12345678, '', '');
INSERT INTO Containers VALUES('another-capp', '1.2.0', 'chunk', '/v1/capps/another-capp/1.2.0.caibx', '0001-01-01T00:00:00Z', NULL, '/v1/capps/another-capp/1.2.0.caibx', '/v1/chunks/');
INSERT INTO Containers VALUES('my-app', '2.3.4', 'tar.gz', '/v1/capps/my-app/2.3.4.tar.gz', '2024-04-20T08:30:45.123Z', 98765432, '', '');
INSERT INTO Containers VALUES('my-app', '2.3.4', 'chunk', '/v1/capps/my-app/2.3.4.caibx', '0001-01-01T00:00:00Z', NULL, '/v1/capps/my-app/2.3.4.caibx', '/v1/chunks/');
INSERT INTO Containers VALUES('example-app', '3.0.1', 'tar.gz', '/v1/capps/example-app/3.0.1.tar.gz', '2024-03-15T14:15:22.789Z', 65432198, '', '');
INSERT INTO Containers VALUES('example-app', '3.0.1', 'chunk', '/v1/capps/example-app/3.0.1.caibx', '0001-01-01T00:00:00Z', NULL, '/v1/capps/example-app/3.0.1.caibx', '/v1/chunks/');
INSERT INTO Containers VALUES('sample-app', '4.5.6', 'tar.gz', '/v1/capps/sample-app/4.5.6.tar.gz', '2024-02-10T10:10:10.100Z', 11223344, '', '');
INSERT INTO Containers VALUES('sample-app', '4.5.6', 'chunk', '/v1/capps/sample-app/4.5.6.caibx', '0001-01-01T00:00:00Z', NULL, '/v1/capps/sample-app/4.5.6.caibx', '/v1/chunks/');
INSERT INTO Containers VALUES('demo-app', '5.0.0', 'tar.gz', '/v1/capps/demo-app/5.0.0.tar.gz', '2024-01-01T00:00:00.000Z', 77889900, '', '');
INSERT INTO Containers VALUES('demo-app', '5.0.0', 'chunk', '/v1/capps/demo-app/5.0.0.caibx', '0001-01-01T00:00:00Z', NULL, '/v1/capps/demo-app/5.0.0.caibx', '/v1/chunks/');

SELECT * FROM Containers;
EOF

echo "All done."
