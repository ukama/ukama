#!/usr/bin/env python3

import struct, sys, json, os
from pprint import pprint

# our module
import bitfieldsparser

def pmecc_get_ecc_bytes(cap, sector_size):
	m = 12 + sector_size // 512
	return (m * cap + 7) // 8

def gen_pmecc_header(page_size, oob_size, ecc_bits, sector_size):
	'''
	generate a pmecc header value, which is 32bits.
	'''
	# print 'page_size: %d, oob_size: %d, ecc_bits: %d, sector_size: %d' \
	#	% (page_size, oob_size, ecc_bits, sector_size)
	nbSectorPerPage = page_size // sector_size
	eccOffset = oob_size - nbSectorPerPage * pmecc_get_ecc_bytes(ecc_bits, sector_size)
	#print 'eccOffset = %d, nbSector = %d, nbSectorPerPage = %d' % (eccOffset, nbSectorPerPage, nbSectorPerPage)

	######## name : value (the excceed field will be ignored)
	#
	# name should be same as .json struct[] array's name.
	maps = {
		"key" : 0xc,
		"usePmecc" : (ecc_bits > 1),
		"nbSectorPerPage" : nbSectorPerPage,
		"spareSize" : oob_size,
		"eccBitReq" : 2 if ecc_bits < 2 else ecc_bits,
		"sectorSize" : sector_size,
		"eccOffset" : eccOffset,
		}

	path_list = [
		os.getcwd(),
		os.path.join(os.getcwd(), 'scripts'),
		]
	for path in path_list:
		fname = os.path.join(path, 'pmecc_head.json')
		if os.path.isfile(fname):
			json_file=open(fname)
			data = json.load(json_file)
			val = bitfieldsparser.convert_bitfield(data["pmecc_header"]["struct"], maps)
			#print 'val = 0x%x' % val
			json_file.close()

			return val

	sys.exit('Cannot find pmecc_head.json!')

# given
#   1. .json file name, which descripts every bit field of a structure
#   2. structure name: which is for the val.
#   3. val: values that need to be parsered by .json
#
# output
#   human readable string description for each bit fields
def display_pmecc_header(file_name, val):
	name = "pmecc_header"
	json_file=open(file_name)
	data = json.load(json_file)
	#pprint(data["pmecc_header"]["usePmecc"])
	if name in data:
		bitfieldsparser.parse_bitfield(data[name]["struct"], val)
	else:
		print('Error: Cannot find the key name: %s in the json file: %s' % (name, file_name))
	json_file.close()

# Test code start here...
# val = gen_pmecc_header(8192, 448, 8, 1024)
#val = gen_pmecc_header(2048, 64, 4, 512)
#print('val = 0x%x' % val)
#print('')
#val = 0xc1304805
# display_pmecc_header("pmecc_head.json", val)
