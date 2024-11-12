#!/bin/bash

# Generate static.h

echo "#ifndef STATIC"        > static.h
echo "#ifdef  UNIT_TEST"     >> static.h
echo "#define STATIC"        >> static.h
echo "#else"                 >> static.h 
echo "#define STATIC static" >> static.h
echo "#endif"                >> static.h
echo "#endif"                >> static.h
