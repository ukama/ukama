#!/bin/bash
# IBM_PROLOG_BEGIN_TAG
# This is an automatically generated prolog.
#
# $Source: ffs/test/ffs_tool_test.sh $
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
# ffs_tool_test.sh
#
#  Test case to perform tests for supported options in ffs tool
#
# Author: Shekar Babu <shekbabu@in.ibm.com>
#

FFS_TOOL="ffs"
NOR_IMAGE="/tmp/pnor"
OFFSET="0x3F0000"
SIZE="8MiB"
BLOCK="64KiB"
LOGICAL="logical"
DATA="data"
PAD="0xFF"

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

	# Check if nor image already exist
	if [ -f $1 ];then
		rm $1
	fi
	echo Creating nor image $1
	$FFS_TOOL --create $1 -p $2 -s $3 -b $4
	RC=$?
	if [ $RC -ne 0 ]; then
		echo Error, creating $1 image
		exit $RC
	fi
	echo Success, creating $1 image
}

add_logical_partition() {
	echo Adding logical partition $3
	echo $FFS_TOOL --add $1 -O $2 --flags 0x0 --pad $PAD -n $3 -t $4
	$FFS_TOOL --add $1 -O $2 --flags 0x0 --pad $PAD -n $3 -t $4
	RC=$?
	if [ $RC -ne 0 ]; then
		echo Error, adding $4 partition $3
		exit $RC
	fi
	echo Success, adding $4 partition $3
}

add_data_partition() {
	echo Adding data partition $3
	echo $FFS_TOOL --add $1 -O $2 --flags 0x0 --pad $PAD -s $5 -o $6 -n $3 -t $4
	$FFS_TOOL --add $1 -O $2 --flags 0x0 --pad $PAD -s $5 -o $6 -n $3 -t $4
	RC=$?
	if [ $RC -ne 0 ]; then
		echo Error, adding $4 partition $3
		exit $RC
	fi
	echo Success, adding $4 partition $3
}

read_partition_entry() {
	echo Reading partition entry $3
        echo $FFS_TOOL --read $1 -O $2 --name $3 -d $4 --force
        $FFS_TOOL --read $1 -O $2 --name $3 -d $4 --force
        RC=$?
        if [ $RC -ne 0 ]; then
                echo Error, reading partition entry $3
                exit $RC
        fi
        echo Success, reading partition entry $3
}

write_partition_entry() {
        echo Writing to partition entry $3
        echo $FFS_TOOL --write $1 -O $2 --name $3 -d $4 --force
        $FFS_TOOL --write $1 -O $2 --name $3 -d $4 --force
        RC=$?
        if [ $RC -ne 0 ]; then
                echo Error, writing to partition entry $3
                exit $RC
        fi
        echo Success, writing to partition entry $3
}

list_partition_table_entries() {
	echo Listing partition table entries in $1
	echo $FFS_TOOL --list $1 -O $2
	$FFS_TOOL --list $1 -O $2
        RC=$?
        if [ $RC -ne 0 ]; then
                echo Error, Listing partition table entries in $1
                exit $RC
        fi
        echo Success, Listing partition table entries in $1
}
hexdump_partition_entry() {
	echo Hexdump partition entry $3 into $4
	echo "$FFS_TOOL --hexdump $1 -O $2 --name $3 > $4"
	$FFS_TOOL --hexdump $1 -O $2 --name $3 > $4
        RC=$?
        if [ $RC -ne 0 ]; then
                echo Error, hexdump partition entry $3 into $4
                exit $RC
        fi
        echo Success, hexdump partition entry $3 into $4
}

delete_partition_entry() {
	echo Delete partition entry $3
	echo $FFS_TOOL --delete $1 -O $2 --name $3
	$FFS_TOOL --delete $1 -O $2 --name $3
        RC=$?
        if [ $RC -ne 0 ]; then
                echo Error, deleting partition entry $3
                exit $RC
        fi
        echo Success, deleting partition entry $3
}

get_partition_entry_user_word() {
	echo Get user word from partition entry $3
	echo $FFS_TOOL --modify $1 -O $2 --name $3 -u $4
	$FFS_TOOL --modify $1 -O $2 --name $3 -u $4 > /tmp/GETUW
	sed 's/^\(.\)\{7\}//g' /tmp/GETUW > /tmp/chop_GETUW
        RC=$?
        if [ $RC -ne 0 ]; then
                echo Error, Getting user word from partition entry $3
                exit $RC
        fi
        echo Success, Getting user word from partition entry $3
}

put_partition_entry_user_word() {
	echo Put user word to partition entry $3
	echo $5 > /tmp/PUTUW
	echo $FFS_TOOL --modify $1 -O $2 --name $3 -u $4 --value $5
	$FFS_TOOL --modify $1 -O $2 --name $3 -u $4 --value $5
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
        cmp /tmp/chop_GETUW /tmp/PUTUW > /dev/null
        RC=$?
        if [ $RC -ne 0 ]; then
                echo FAIL, user word get/put mismatch -- entry $3
                exit $RC
        fi
        echo PASS, user word get/put matches -- entry $3
	rm /tmp/GETUW /tmp/PUTUW /tmp/chop_GETUW
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
add_data_partition $NOR_IMAGE $OFFSET boot0/bootenv $DATA 1MiB 0M
add_data_partition $NOR_IMAGE $OFFSET boot0/ipl $DATA 1MiB 2M
add_data_partition $NOR_IMAGE $OFFSET boot0/spl $DATA 960K 3M
# Add logical partition
add_logical_partition $NOR_IMAGE $OFFSET boot1 $LOGICAL
# Creating data partition
add_data_partition $NOR_IMAGE $OFFSET boot1/uboot $DATA 1MiB 4M
add_data_partition $NOR_IMAGE $OFFSET boot1/fsp $DATA 1MiB 6M
add_data_partition $NOR_IMAGE $OFFSET boot1/bootfsp $DATA 960K 7M
# Listing all created partition entries (logical+data)
list_partition_table_entries $NOR_IMAGE $OFFSET

# Perform read and write operations on all partition entries
read_write_part_entry $NOR_IMAGE $OFFSET boot0/bootenv /tmp/in_file /tmp/out_file
read_write_part_entry $NOR_IMAGE $OFFSET boot0/ipl /tmp/in_file /tmp/out_file
read_write_part_entry $NOR_IMAGE $OFFSET boot1/uboot /tmp/in_file /tmp/out_file
read_write_part_entry $NOR_IMAGE $OFFSET boot1/fsp /tmp/in_file /tmp/out_file

# Perform get and put user words on all partition entries
get_put_user_word $NOR_IMAGE $OFFSET boot0/bootenv 0 0x0000000a
get_put_user_word $NOR_IMAGE $OFFSET boot0/ipl 1 0x0000000b
get_put_user_word $NOR_IMAGE $OFFSET boot0/spl 2 0x0000000c
get_put_user_word $NOR_IMAGE $OFFSET boot1/uboot 3 0x0000000d
get_put_user_word $NOR_IMAGE $OFFSET boot1/fsp 4 0x0000000f

# Hexdump partition entry
compare_hexdump $NOR_IMAGE $OFFSET boot0/bootenv /tmp/hexdump

# Delete a partition entry
delete_partition_entry $NOR_IMAGE $OFFSET boot0/bootenv

# Listing all created partition entries (logical+data)
list_partition_table_entries $NOR_IMAGE $OFFSET

# Clean/remove all temporary files
clean_data
