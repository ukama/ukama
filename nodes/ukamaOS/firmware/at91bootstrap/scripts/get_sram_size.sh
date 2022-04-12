#!/bin/sh

# find the _sramsize line.
# get the hex number, and remove the '0x' as bc cannot recognize it.
# also we need convert all lower case letter to upper case.
size_hex=`awk '$2 == "_sramsize" { print $1}' $1 | sed 's/^0x//' | tr '[a-z]' '[A-Z]'`

# quit if not sucessful
[ $? -ne 0 ] && exit 1;
# quit if $size_hex is empty
[ "x"$size_hex = "x" ] && exit 1;

size_int=`echo "ibase=16;"$size_hex | bc`
# quit if not sucessful
[ $? -ne 0 ] && exit 1;
echo $size_int
