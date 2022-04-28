find -type f \( -name 'GNUmakefile' -o -name 'makefile' -o -name 'Makefile' \) \
-exec bash -c 'cd "$(dirname "{}")" && go mod tidy -compat=1.17' \;