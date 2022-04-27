#!/usr/bin/env python3

from __future__ import print_function
import struct, sys, json
from pprint import pprint

# it can parse a json structure like following:
#   "some_name_you_defined": {
#       "struct": [
#       {
#           "name": "xxx",
#           "type": "bit_field",    /* now only support bit_field */
#           "start": 0,             /* bit start */
#           "len": 2,               /* length of this bit field */
#           "meaning": {
#               /* bit val: [ "meaningful string", values that convert to bit val] */
#               "0": ["512 bytes", 512],
#           }
#       },
#       {
#       ...
#       },
#
#       ...
#       ]
#   }
#

def find_meaning(maps, bits):
	'''Conver value to correcpoding bits'''
	if bits["name"] in maps:
		start = bits["start"]
		len = bits["len"]
		mask = (2 ** len - 1) << start

		if "meaning" in bits:
			# find the corresponding name, value in the maps,
			# then the .json's meaning part, convert it to bit field
			for name, array in bits["meaning"].items():
				# array[1] present the convert sitation. if equal then convert to the name, which is number in bits
				if maps[bits["name"]] == array[1]:
					return (int(name) << start) & mask

			print("Error: not found")
		else:
			# find the corresponding name, value in the maps, and shift as bit field.
			return (maps[bits["name"]] << start) & mask

	return 0

# given
#   1. a .json data that descripts the bitfield structure.
#          - for example: pmecc_header: {struct : [, , , ...]}
#   2. a maps that include the bit field value to set for each bit field.
#          - that we can generate a whole word according the every value in the bit fields.
# output
#   a value of 32bit number, so far.
def convert_bitfield(fields, maps):
	'''
	according to .json file's struct[] array, convert given maps into 32bit value
	'''
	val = 0
	for bits in fields:
		if bits["type"] == "bit_field":
			# convert the each maps value to corresponing bits field
			val = val | find_meaning(maps, bits)

	return val

# given a value, parse to readable string format.
def parse_bitfield(fields, whole_word):
	for bits in fields:
		if bits["type"] == "bit_field":
			start = bits["start"]
			len = bits["len"]
			# print 'start = %d, len = %d' % (start, len)
			mask = (2 ** len - 1) << start
			# print 'mask: %x' % mask
			value = (whole_word & mask) >> start
			meaning_str = ""
			if "meaning" in bits:
				if repr(value) in bits["meaning"]:
					meaning_str = bits["meaning"][repr(value)][0]

			if meaning_str == "":
				print(bits["name"], ':', value)
			else:
				print(bits["name"], ':', value, '->', meaning_str)
		else:
			print('wrong type')
