#!/bin/bash
# IBM_PROLOG_BEGIN_TAG
# This is an automatically generated prolog.
#
# $Source: fpart/fpart.sh $
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
#
#     File: ffs cmdline test script
#   Author: Shaun wetzstein <shaun@us.ibm.com>
#     Date: 02/06/13
#
#     Todo:
#	1) function to check syntax of each command

CP=cp
RM=rm
FPART=fpart
MKDIR=mkdir
GREP=grep
HEAD=head
HEX=hexdump
CRC=sha1sum
DD=dd
DIFF=diff

IMG=sunray2_nor64M_flash.mif
TMP=/tmp/ffs.$$
ORIG=${TMP}/nor.orig
URANDOM=/dev/urandom

FAIL=1
PASS=0

KB=$((1*1024))
MB=$(($KB*1024))

function expect()
{
	local expect=${1}
	shift
	local command=${*}
	eval $command
	local actual=${?}

	if [[ ${expect} -eq ${actual} ]]; then
		echo "[PASSED] rc: '${command}' ===> expect=${expect}, actual=${actual}" >&2
	else
		echo "[FAILED] rc: '${command}' ===> expect=${expect}, actual=${actual}" >&2
		exit 1
	fi
}

function pass()
{
	expect ${PASS} ${*}
}

function fail()
{
	expect ${FAIL} ${*}
}

function size()
{
	local size=$(stat -L -c %s ${2})
	if [[ ${1} == ${size} ]]; then
		echo "[PASSED] size: '${2}' ===> expect=${1}, actual=${size}"
	else
		echo "[FAILED] size: '${2}' ===> expect=${1}, actual=${size}"
		exit 1
	fi
}

function crc()
{
	local crc=$(${CRC} ${2}|cut -f 1 -d ' ')
	if [[ ${1} == ${crc} ]]; then
		echo "[PASSED] crc: '${2}' ===> expect=${1}, actual=${crc}"
	else
		echo "[FAILED] crc: '${2}' ===> expect=${1}, actual=${crc}"
		exit 1
	fi
}

function setup()
{
	pass ${RM} -rf ${TMP}
	pass ${MKDIR} -p ${TMP}
}

function cleanup()
{
	pass ${RM} -rf ${TMP}
}

function create()
{
	local target=${TMP}/create.nor
	pass ${RM} -f ${target}

	pass ${FPART} -t ${target} -s 64M -b 64K -p 0x3f0000 -C
	pass ${FPART} --target ${target} --size 64MiB --block 64kb \
		   	      --partition-offset 0x7f0000 --create

	local output=${TMP}/create.txt
	pass ${FPART} -t ${target} -L > ${output}
	pass ${GREP} \'0x3f0000\' ${output} > /dev/null
	pass ${GREP} \'0x7f0000\' ${output} > /dev/null
	pass ${GREP} \'blk:010000\' ${output} > /dev/null
	pass ${GREP} \'blk\(s\):000400\' ${output} > /dev/null
	pass ${GREP} \'entsz:000080\' ${output} > /dev/null
	pass ${GREP} \'ent\(s\):000001\' ${output} > /dev/null
	pass ${RM} -f ${output}

	pass ${RM} -f ${target}
}

function add()
{
	local target=${TMP}/add.nor
	pass ${RM} -f ${target}

	pass ${FPART} -t ${target} -s 64M -b 64K --partition-offset 0x3f0000 -C
	pass ${FPART} --target ${target} --size 64MiB --block 64kb -p 0x7f0000 --create

	local output=${TMP}/add.txt
	local name="logical"
	pass ${FPART} -t ${target} -l -n ${name} -g 0 -A
	pass ${FPART} -t ${target}    -n ${name} -L > ${output}
	pass ${GREP} ${name} ${output} > /dev/null
	pass ${GREP} "l-----" ${output} > /dev/null

	for ((i=0; i<9; i++)); do
		local full="${name}/test${i}"
		local offset=$((${i}*$MB))
		local output=${TMP}/add.txt

		# avoid clobbering 'part'
		if [[ ${i} -eq 4 ]]; then
			local size=$MB
		else
			local size=$(($MB-64*KB))
		fi

		pass ${FPART} -t ${target} -o ${offset} -s ${size} -g 0 -n ${full} -a ${i} -A
		pass ${FPART} -t ${target} -n ${full} -T
		pass ${FPART} -t ${target} -n ${full} -L > ${output}
		pass ${GREP} ${full} ${output} > /dev/null
		pass ${GREP} "d-----" ${output} > /dev/null
		local range=$(printf "%.8x-%.8x" ${offset} $((${offset}+${size}-1)))
		pass ${GREP} ${range} ${output} > /dev/null
		pass ${GREP} $(printf "%x" ${size}) ${output} > /dev/null
		pass ${RM} -f ${output}
	done

	pass ${RM} -f ${target}
}

function delete()
{
	local target=${TMP}/delete.nor
	pass ${RM} -f ${target}

	pass ${FPART} -t ${target} -s 64M -b 64K --partition-offset 0x3f0000 -C
	pass ${FPART} --target ${target} --size 64MiB --block 64kb -p 0x7f0000 --create

	local output=${TMP}/delete.txt

	for ((i=0; i<9; i++)); do
		local full="test${i}"
		local offset=$((${i}*$MB))
		local output=${TMP}/add.txt

		# avoid clobbering 'part'
		if [[ ${i} -eq 4 ]]; then
			local size=$MB
		else
			local size=$(($MB-64*KB))
		fi

		pass ${FPART} -t ${target} -o ${offset} -s ${size} -g 0 -n ${full} -a ${i} -A
		pass ${FPART} -t ${target} -n ${full} -T
		pass ${FPART} -t ${target} -n ${full} -L > ${output}
		pass ${GREP} ${full} ${output} > /dev/null
		pass ${FPART} -t ${target} -n ${full} -D
		pass ${FPART} -t ${target} -n ${full} -L > ${output}
		fail ${GREP} ${full} ${output} > /dev/null

		pass ${RM} -f ${output}
	done

	pass ${RM} -f ${target}
}

function hex()
{
	local target=${TMP}/hexdump.nor
	pass ${RM} -f ${target}

	pass ${FPART} -t ${target} -s 64M -b 64K --partition-offset 0x3f0000 -C
	pass ${FPART} --target ${target} --size 64MiB --block 64kb -p 0x7f0000 --create

	local name="logical"
	pass ${FPART} -t ${target} -l -n ${name} -g 0 -A

	for ((i=0; i<9; i++)); do
		local full="${name}/test${i}"
		local offset=$((${i}*$MB))
		local output=${TMP}/hexdump.txt

		# avoid clobbering 'part'
		if [[ ${i} -eq 4 ]]; then
			local size=$MB
		else
			local size=$(($MB-64*KB))
		fi

		pass ${FPART} -t ${target} -o ${offset} -s ${size} -g 0 -n ${full} -a ${i} -A
		pass ${FPART} -t ${target} -n ${full} -T
		pass ${FPART} -t ${target} -n ${full} -H > ${output}
		local p=$(printf "%2.2x%2.2x%2.2x%2.2x" $i $i $i $i)
		pass ${GREP} \'$p $p $p $p\' ${output} > /dev/null
		pass ${RM} -f ${output}
	done

	pass ${RM} -f ${target}
}

function read()
{
	local target=${TMP}/read.nor
	pass ${RM} -f ${target}

	pass ${FPART} -t ${target} -s 64M -b 64K --partition-offset 0x3f0000 -C
	pass ${FPART} --target ${target} --size 64MiB --block 64kb -p 0x7f0000 --create

	local name="logical"
	pass ${FPART} -t ${target} -l -n ${name} -g 0 -A

	for ((i=0; i<9; i++)); do
		local full="${name}/test${i}"
		local offset=$((${i}*$MB))
		local output=${TMP}/${i}.bin

		# avoid clobbering 'part'
		if [[ ${i} -eq 4 ]]; then
			local size=$MB
		else
			local size=$(($MB-64*KB))
		fi

		pass ${FPART} -t ${target} -o ${offset} -s ${size} -g 0 -n ${full} -a ${i} -A
		pass ${FPART} -t ${target} -n ${full} -T
		pass ${FPART} -t ${target} -n ${full} -R ${output}
		pass ${HEX} -C ${output} > ${output}.hex
		local p=$(printf "%2.2x %2.2x %2.2x %2.2x" $i $i $i $i)
		pass ${GREP} \"${p} ${p}\" ${output}.hex
		pass ${RM} -f ${output} ${output}.hex
	done

	pass ${RM} -f ${target}
}

function write()
{
	local target=${TMP}/write.nor
	pass ${RM} -f ${target}

	pass ${FPART} -t ${target} -s 64M -b 64K --partition-offset 0x3f0000 -C
	pass ${FPART} --target ${target} --size 64MiB --block 64kb -p 0x7f0000 --create

	local name="write"
	local input=${TMP}/write.in
	local output=${TMP}/write.out
	local block=${1}
	local count=$(($MB/$block))

	pass ${FPART} -t ${target} -o 1M -s 1M -g 0 -n ${name} -a 0xFF -A
	for ((i=1; i<${count}; i++)); do
		pass ${DD} if=${URANDOM} of=${input} bs=${block} count=${i} 2> /dev/null
		local crc=$(${CRC} ${input})

		pass ${FPART} -t ${target} -n ${name} -E -a 0xFF
		pass ${FPART} -t ${target} -n ${name} -W ${input}
		pass ${FPART} -t ${target} -n ${name} -U 0 -u ${crc}
		pass ${FPART} -t ${target} -n ${name} -R ${output}

		size $((${i}*${block})) ${output}
		crc  ${crc} ${output}

		local crc=$(printf "%x" ${crc})
		pass "${FPART} -t ${target} -n ${name} -U 0 | ${GREP} ${crc}" > /dev/null

		pass ${RM} -f ${input} ${output}
	done
	pass ${FPART} -t ${target} -n ${name} -D

	pass ${RM} -f ${target}
}

function copy()
{
	local src=${TMP}/copy.src
	local dst=${TMP}/copy.dst
	pass ${RM} -f ${src} ${dst}

	local part="0x0,0x20000"
	pass ${FPART} -t ${src} -s 64M -b 64K -p ${part} -C
	pass ${FPART} -t ${dst} -s 64M -b 64K -p ${part} -C

	local name="copy"
	pass ${FPART} -t ${src} -o 1M -s 1M -g 0 -n ${name} -a 0xFF -p ${part} -A
	pass ${FPART} -t ${dst} -o 1M -s 1M -g 0 -n ${name} -a 0x00 -p ${part} -A

	local input=${TMP}/copy.in
	local output=${TMP}/copy.out
	if [[ -z "${1}" ]]; then
		local block=64
	else
		local block=${1}
	fi
	local count=$(($MB/$block))

	for ((i=1; i<${count}; i++)); do
		pass ${DD} if=${URANDOM} of=${input} bs=${block} count=${i} 2> /dev/null
		local crc=$(${CRC} ${input})

		pass ${FPART} -t ${src} -n ${name} -p ${part} -E -a 0xFF
		pass ${FPART} -t ${src} -n ${name} -p ${part} -W ${input}
		pass ${FPART} -t ${src} -n ${name} -p ${part} -R ${output}.src
		pass ${FPART} -t ${src} -n ${name} -p ${part} -U 0 -u ${crc}
		pass ${FPART} -t ${dst} -n ${name} -p ${part} -O ${src}
		pass ${FPART} -t ${dst} -n ${name} -p ${part} -M ${src}
		pass ${FPART} -t ${src} -n ${name} -p ${part} -R ${output}.dst

		size $((${i}*${block})) ${output}.src
		size $((${i}*${block})) ${output}.dst
		crc  ${crc} ${output}.src
		crc  ${crc} ${output}.dst
		pass ${DIFF} ${output}.src ${output}.dst

		local crc=$(printf "%x" ${crc})
		pass "${FPART} -t ${src} -n ${name} -p ${part} -U 0 | ${GREP} ${crc}" > /dev/null
		pass "${FPART} -t ${dst} -n ${name} -p ${part} -U 0 | ${GREP} ${crc}" > /dev/null

		pass ${RM} -f ${input} ${output}.*

	done

	expect 2 ${DIFF} ${src} ${dst}

	pass ${RM} -f ${src} ${dst}
}

function main()
{
	create
	add
	delete
#	hex
#	read
#	copy  $((15*$KB))
#	copy  $((21*$KB))
#	copy  $((64*$KB))
#	write $((15*$KB))
#	write $((21*$KB))
#	write $((64*$KB))
}

setup
if [[ -z "${1:-}" ]]; then
	main
else
	case "$1" in
		create	) create		;;
		add	) add			;;
		delete	) delete		;;
#		hex	) hex			;;
#		read	) read			;;
#		copy	) copy	$((${2}*$KB))	;;
#		write	) write $((${2}*$KB))	;;
		*	) echo "$1 not implemented"; exit 1	;;
	esac

	exit 0;
fi
cleanup
