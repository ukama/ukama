#!/bin/bash
# IBM_PROLOG_BEGIN_TAG
# This is an automatically generated prolog.
#
# $Source: ffs/test/libffs_test.sh $
#
# OpenPOWER FFS Project
#
# Contributors Listed Below - COPYRIGHT 2014,2015
# [+] International Business Machines Corp.
#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
# implied. See the License for the specific language governing
# permissions and limitations under the License.
#
# IBM_PROLOG_END_TAG
# libffs_test.sh
#
#  Test case to perform unit tests for all api's in libffs.so
#
# Author: Shekar Babu <shekbabu@in.ibm.com>
#


NOR_IMAGE=/tmp/sunray2.nor
OFFSET=0
#OFFSET=4128768
SIZE=8388608
#SIZE=67108864	#For 64MB nor
BLOCK=65536
LOGICAL=logical
DATA=data

create_dummy_file() {
	echo Creating a dummy file $1 of size $2 with sample data $3
	yes $3 | head -$2 > $1
        RC=$?
        if [ $RC -ne 0 ]; then
                echo Error, creating dummy file $1
                exit $RC
        fi
	echo Success, creating $1
}

create_nor_image() {
	echo Creating nor image $1
	echo test_libffs -c $1 -O $2 -s $3 -b $4
	./test_libffs -c $1 -O $2 -s $3 -b $4
	RC=$?
	if [ $RC -ne 0 ]; then
		echo Error, creating $1 image
		exit $RC
	fi
	echo Success, creating $1 image
}

add_logical_partition() {
	echo Adding logical partition $3
	echo test_libffs -a $1 -O $2 -n $3 -t $4
	./test_libffs -a $1 -O $2 -n $3 -t $4
	RC=$?
	if [ $RC -ne 0 ]; then
		echo Error, adding $4 partition $3
		exit $RC
	fi
	echo Success, adding $4 partition $3
}

add_data_partition() {
	echo Adding data partition $3
	echo test_libffs -a $1 -O $2 -n $3 -t $4 -s $5 -o $6
	./test_libffs -a $1 -O $2 -n $3 -t $4 -s $5 -o $6
	RC=$?
	if [ $RC -ne 0 ]; then
		echo Error, adding $4 partition $3
		exit $RC
	fi
	echo Success, adding $4 partition $3
}

read_partition_entry() {
	echo Reading partition entry $3
	echo test_libffs -r $1 -O $2 -n $3 -o $4
	./test_libffs -r $1 -O $2 -n $3 -o $4
        RC=$?
        if [ $RC -ne 0 ]; then
                echo Error, reading partition entry $3
                exit $RC
        fi
        echo Success, reading partition entry $3
}

write_partition_entry() {
        echo Writing to partition entry $3
        echo test_libffs -w $1 -O $2 -n $3 -i $4
        ./test_libffs -w $1 -O $2 -n $3 -i $4 > /dev/null
        RC=$?
        if [ $RC -ne 0 ]; then
                echo Error, writing to partition entry $3
                exit $RC
        fi
        echo Success, writing to partition entry $3
}

list_partition_table_entries() {
	echo Listing partition table entries in $1
	echo test_libffs -l $1 -O $2
	./test_libffs -l $1 -O $2
        RC=$?
        if [ $RC -ne 0 ]; then
                echo Error, Listing partition table entries in $1
                exit $RC
        fi
        echo Success, Listing partition table entries in $1
}
hexdump_partition_entry() {
	echo Hexdump partition entry $3 into $4
	echo "test_libffs -h $1 -O $2 -n $3 > $4"
	./test_libffs -h $1 -O $2 -n $3 > $4
        RC=$?
        if [ $RC -ne 0 ]; then
                echo Error, hexdump partition entry $3 into $4
                exit $RC
        fi
        echo Success, hexdump partition entry $3 into $4
}

delete_partition_entry() {
	echo Delete partition entry $3
	echo test_libffs -d $1 -O $2 -n $3
	./test_libffs -d $1 -O $2 -n $3
        RC=$?
        if [ $RC -ne 0 ]; then
                echo Error, deleting partition entry $3
                exit $RC
        fi
        echo Success, deleting partition entry $3
}

get_partition_entry_user_word() {
	echo Get user word from partition entry $3
	echo test_libffs -m $1 -O $2 -n $3 -u $4 -g
	./test_libffs -m $1 -O $2 -n $3 -u $4 -g > /tmp/GETUW
        RC=$?
        if [ $RC -ne 0 ]; then
                echo Error, Getting user word from partition entry $3
                exit $RC
        fi
        echo Success, Getting user word from partition entry $3
}

put_partition_entry_user_word() {
	echo Put user word to partition entry $3
	echo test_libffs -m $1 -O $2 -n $3 -u $4 -p -v $5
	./test_libffs -m $1 -O $2 -n $3 -u $4 -p -v $5 > /tmp/PUTUW
        RC=$?
        if [ $RC -ne 0 ]; then
                echo Error, Putting user word to partition entry $3
                exit $RC
        fi
        echo Success, Putting user word to partition entry $3
}

read_write_part_entry() {
	write_partition_entry $1 $2 $3 $4
	read_partition_entry $1 $2 $3 $5
	cmp $4 $5 > /dev/null
	RC=$?
        if [ $RC -ne 0 ]; then
                echo FAIL, data read/write mismatch -- entry $3
                exit $RC
        fi
        echo PASS, data read/write matches -- entry $3
}

get_put_user_word() {
        put_partition_entry_user_word $1 $2 $3 $4 $5
        get_partition_entry_user_word $1 $2 $3 $4
        cmp /tmp/GETUW /tmp/PUTUW > /dev/null
        RC=$?
        if [ $RC -ne 0 ]; then
                echo FAIL, user word get/put mismatch -- entry $3
                exit $RC
        fi
        echo PASS, user word get/put matches -- entry $3
	rm /tmp/GETUW /tmp/PUTUW
}

compare_hexdump() {
	hexdump_partition_entry $1 $2 $3 $4
	HEXFILE=/tmp/hex_sz0
	stat -c %s $4 > $HEXFILE
	if [[ -s $HEXFILE ]] ; then
		echo PASS, hexdump test on entry $3
	else
		echo FAIL, hexdump test on entry $3
	        exit $RC
        fi
	rm $4 $HEXFILE
}

clean_data() {
	rm $NOR_IMAGE /tmp/in_file /tmp/out_file
	exit 0
}

# Main program starts

# Create a dummy file as 'filename size data'
create_dummy_file /tmp/in_file 131072 WELCOME

# Create nor image
create_nor_image $NOR_IMAGE $OFFSET $SIZE $BLOCK

# Add logical partition
add_logical_partition $NOR_IMAGE $OFFSET boot0 $LOGICAL
# Creating data partition
add_data_partition $NOR_IMAGE $OFFSET boot0/bootenv $DATA 1048576 65536
add_data_partition $NOR_IMAGE $OFFSET boot0/ipl $DATA 1048576 2097152
add_data_partition $NOR_IMAGE $OFFSET boot0/spl $DATA 1048576 3145728
# Add logical partition
add_logical_partition $NOR_IMAGE $OFFSET boot1 $LOGICAL
# Creating data partition
add_data_partition $NOR_IMAGE $OFFSET boot1/uboot $DATA 1048576 4194304
add_data_partition $NOR_IMAGE $OFFSET boot1/fsp $DATA 1048576 5242880
# Add logical partition
add_logical_partition $NOR_IMAGE $OFFSET linux0 $LOGICAL
# Creating data partition
add_data_partition $NOR_IMAGE $OFFSET linux0/vpd $DATA 1048576 6291456
add_data_partition $NOR_IMAGE $OFFSET linux0/hostboot $DATA 1048576 7340032

# Listing all created partition entries (logical+data)
list_partition_table_entries $NOR_IMAGE $OFFSET

# Perform read and write operations on all partition entries
read_write_part_entry $NOR_IMAGE $OFFSET boot0/bootenv /tmp/in_file /tmp/out_file
read_write_part_entry $NOR_IMAGE $OFFSET boot0/ipl /tmp/in_file /tmp/out_file
read_write_part_entry $NOR_IMAGE $OFFSET boot0/spl /tmp/in_file /tmp/out_file
read_write_part_entry $NOR_IMAGE $OFFSET boot1/uboot /tmp/in_file /tmp/out_file
read_write_part_entry $NOR_IMAGE $OFFSET boot1/fsp /tmp/in_file /tmp/out_file
read_write_part_entry $NOR_IMAGE $OFFSET linux0/vpd /tmp/in_file /tmp/out_file
read_write_part_entry $NOR_IMAGE $OFFSET linux0/hostboot /tmp/in_file /tmp/out_file

# Perform get and put user words on all partition entries
get_put_user_word $NOR_IMAGE $OFFSET boot0/bootenv 0 28
get_put_user_word $NOR_IMAGE $OFFSET boot0/ipl 1 56
get_put_user_word $NOR_IMAGE $OFFSET boot0/spl 2 96
get_put_user_word $NOR_IMAGE $OFFSET boot1/uboot 3 16
get_put_user_word $NOR_IMAGE $OFFSET boot1/fsp 4 84
get_put_user_word $NOR_IMAGE $OFFSET linux0/vpd 8 64
get_put_user_word $NOR_IMAGE $OFFSET linux0/hostboot 15 42

# Hexdump partition entry
compare_hexdump $NOR_IMAGE $OFFSET boot0/bootenv /tmp/hexdump

# Delete a partition entry
delete_partition_entry $NOR_IMAGE $OFFSET boot0/bootenv

# Listing all created partition entries (logical+data)
list_partition_table_entries $NOR_IMAGE $OFFSET

# Clean/remove all temporary files
clean_data
