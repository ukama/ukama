#!/usr/bin/env python3

import struct, sys

import pmecc_head

# open bootstrap file
fd = open(sys.argv[1], "rb")
line = fd.read()
fd.close()

if sys.argv[3] == "sama5d2_ptc_ek":
	pmecc_word = pmecc_head.gen_pmecc_header(2048, 64, 4, 512)
elif sys.argv[3].startswith("sama5d3"):
	pmecc_word = pmecc_head.gen_pmecc_header(2048, 64, 4, 512)
elif sys.argv[3] == "sama5d4ek" or sys.argv[3] == "sama5d4_xplained":
	pmecc_word = pmecc_head.gen_pmecc_header(4096, 224, 8, 512)
elif sys.argv[3] == "at91sam9x5ek" or sys.argv[3] == "at91sam9n12ek":
	pmecc_word = pmecc_head.gen_pmecc_header(2048, 64, 2, 512)
elif sys.argv[3] == "sam9x60ek":
	pmecc_word = pmecc_head.gen_pmecc_header(4096, 224, 8, 512)
elif sys.argv[3] == "sam9x60_ddr2_sip_eb":
	pmecc_word = pmecc_head.gen_pmecc_header(4096, 224, 8, 512)
elif sys.argv[3] == "sam9x60_sdr_sip_eb":
	pmecc_word = pmecc_head.gen_pmecc_header(4096, 224, 8, 512)
elif sys.argv[3] == "core9g25":
	pmecc_word = pmecc_head.gen_pmecc_header(2048, 64, 2, 512)
else:
	sys.exit("Not support board!")

vec = struct.pack("<I", pmecc_word)

# generate a new file with pmecc header
fd = open(sys.argv[2], "wb")

for i in range(0, 52):
	fd.write(vec)

fd.write(line[0:])

fd.close()
