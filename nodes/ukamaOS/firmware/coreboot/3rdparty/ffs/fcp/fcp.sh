#!/bin/bash
# IBM_PROLOG_BEGIN_TAG
# This is an automatically generated prolog.
#
# $Source: fcp/fcp.sh $
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
FFS=ffs
FCP=fcp
MKDIR=mkdir
GREP=grep
HEAD=head
HEX=hexdump
CRC=crc32
DD=dd
DIFF=diff
CAT=cat
CP=cp
TRUNC=truncate

TARGET=test.nor
COPY=copy.nor
TMP=/tmp/fcp.$$
URANDOM=/dev/urandom
POFFSET="0x00000,0x10000"

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
	local crc=$(${CRC} ${2})
	if [[ ${1} == ${crc} ]]; then
		echo "[PASSED] crc: '${2}' ===> expect=${1}, actual=${crc}"
	else
		echo "[FAILED] crc: '${2}' ===> expect=${1}, actual=${crc}"
		exit 1
	fi
}

function setup()
{
	local target=${TMP}/${TARGET}
	local output=${TMP}/${TARGET}.txt

	pass ${RM} -rf ${TMP}
	pass ${MKDIR} -p ${TMP}
	pass ${RM} -f ${target}

	pass ${FFS} -t ${target} -s 64M -b 64K -p 0x3F0000 -C
	pass ${FFS} -t ${target} -s 64M -b 64K -p 0x7F0000 -C

	for ((j=0; j<2; j++)); do
		local name="logical${j}"
		local base=$((${j}*$MB*4))

		pass ${FFS} -t ${target} -n ${name} -g 0 -l -A
		pass ${FFS} -t ${target} -n ${name} -L > ${output}

		pass ${GREP} ${name} ${output} > /dev/null
		pass ${GREP} "l-----" ${output} > /dev/null

		for ((i=0; i<4; i++)); do
			local full="${name}/entry${i}"
			local offset=$((${base}+${i}*$MB))

			if [[ ${i} -eq 3 ]]; then
				local size=$(($MB-64*${KB}))
			else
				local size=$MB
			fi

			pass ${FFS} -t ${target} -o ${offset} -s ${size} \
			     -g 0 -n ${full} -a ${i} -A

			pass ${FFS} -t ${target} -n ${full} -L > ${output}
			pass ${GREP} ${full} ${output} > /dev/null
			pass ${GREP} "d-----" ${output} > /dev/null

			local range=$(printf "%.8x-%.8x" ${offset} \
				    $((${offset}+${size}-1)))

			pass ${GREP} ${range} ${output} > /dev/null
			pass ${GREP} $(printf "%x" ${size}) ${output} > \
			     /dev/null
			pass ${RM} -f ${output}
		done
	done
}

function cleanup()
{
	pass ${RM} -rf ${TMP}
}

function erase()
{
	local target=${TMP}/${TARGET}
	local name="logical"

	for ((i=0; i<4; i++)); do
		local output=${TMP}/entry${i}.hex
		local full=${name}0/entry${i}

		pass ${FCP} ${target}:${full} -E ${i}
		pass ${FCP} -T ${target}:${full}
		pass ${FCP} -R ${target}:${full} - | ${HEX} -C > ${output}

		local p=$(printf "%2.2x %2.2x %2.2x %2.2x" $i $i $i $i)
		pass ${GREP} \"${p} ${p}\" ${output}

		pass ${RM} -f ${output}
	done

	pass ${FCP} ${target}:${name}0 -E 0x00
	pass ${FCP} ${target}:${name}1 -E 0x00
}

function write()
{
	local target=${TMP}/${TARGET}
	local output=${TMP}/${TARGET}.txt

	local name="logical"
	local input=${TMP}/write.in
	local output=${TMP}/write.out

	if [[ -z "${1}" ]]; then
		local block=$((64*KB))
	else
		local block=${1}
	fi

	for ((i=0; i<4; i++)); do
		if [[ ${i} -eq 3 ]]; then
			local size=$(($MB-64*${KB}))
		else
			local size=$MB
		fi

		local count=$((${size}/$block))

		for ((c=1; c<=${count}; c++)); do
			pass ${DD} if=${URANDOM} of=${input} bs=${block} \
			     count=${c} 2> /dev/null

			pass ${CRC} ${input} > ${output}
			local crc=$(${CAT} ${output})
			local full=${name}0/entry${i}

			pass ${FCP} ${target}:${full} -E 0x00
			pass ${FCP} -W ${input} ${target}:${full}
			pass ${FCP} ${target}:${full} -U 0=${crc}
			pass ${FCP} ${target}:${full} -R - -f | cat > ${output}

			size $((${c}*${block})) ${output}
			crc  ${crc} ${output}

			local crc=$(printf "%x" ${crc})
			pass "${FCP} ${target}:${full} -o 0x3F0000 -U 0 | \
			     ${GREP} ${crc}" > /dev/null
			pass "${FCP} ${target}:${full} -o 0x7F0000 -U 0 | \
			     ${GREP} ${crc}" > /dev/null

			pass ${RM} -f ${input} ${output}
		done
	done
}

function copy()
{
	local src=${TMP}/${TARGET}
	local dst=${TMP}/${TARGET}.copy

	local input=${TMP}/copy.in
	local output=${TMP}/copy.out

	if [[ -z "${1}" ]]; then
		local block=$((64*KB))
	else
		local block=${1}
	fi

	pass ${TRUNC} -s $(stat -L -c %s ${src}) ${dst}

	for ((i=0; i<4; i++)); do
		if [[ ${i} -eq 3 ]]; then
			local size=$(($MB-64*${KB}))
		else
			local size=$MB
		fi

		local count=$((${size}/$block))

		for ((c=1; c<=${count}; c++)); do
			pass ${DD} if=${URANDOM} of=${input} bs=${block} count=${c} 2> /dev/null
			local crc=$(${CRC} ${input})

			local name="logical0/entry${i}"

			pass ${FCP} ${src}:${name} -E 0x00
			pass ${FCP} ${input} ${src}:${name} -W
			pass ${FCP} ${src}:${name} ${output}.src -f -R
			pass ${FCP} ${src}:${name} -U 0 0=${crc}

			pass ${FCP} ${src}:${name} ${dst}:${name} -C -v -f
			pass ${FCP} ${src}:${name} ${dst}:${name} -M -v

			pass ${FCP} ${dst}:${name} ${output}.dst -R -f
			size $((${c}*${block})) ${output}.src
			size $((${c}*${block})) ${output}.dst
			crc  ${crc} ${output}.src
			crc  ${crc} ${output}.dst
			pass ${DIFF} ${output}.src ${output}.dst

			local crc=$(printf "%x" ${crc})
			pass "${FCP} ${dst}:${name} -o 0x3F0000 -U 0 | \
				${GREP} ${crc}" > /dev/null
			pass "${FCP} ${dst}:${name} -o 0x7F0000 -U 0 | \
				${GREP} ${crc}" > /dev/null

			pass ${RM} -f ${input}* ${output}*
		done
	done

	pass ${FCP} ${src}":logical0" ${src}":logical1" -C -v # logical mirror
	pass ${FCP} ${src}":logical1" ${dst}":logical1" -C -v # logical copy

	expect 2 ${DIFF} ${src} ${dst}
}

function main()
{
	erase
	write $((21*$KB))
	write $((64*$KB))
	copy $((21*$KB))
	copy $((64*$KB))
}

setup
if [[ -z "${1:-}" ]]; then
	main
else
	case "$1" in
		erase	) erase 				;;
		write	) write 				;;
		copy	) copy	 				;;
		*	) echo "$1 not implemented"; exit 1	;;
	esac
	exit 0;
fi
cleanup
