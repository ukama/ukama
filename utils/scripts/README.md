To replace the copyright header, in the given directory, do following:

find ./{path} -name "*.[ch]" -exec ./replace_copyright_header.sh {} \;
