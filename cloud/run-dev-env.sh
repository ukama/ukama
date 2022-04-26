find -type f \( -name 'GNUmakefile' -o -name 'makefile' -o -name 'Makefile' \) \
-exec bash -c 'cd "$(dirname "{}")" && make build' \;
docker-compose up --build