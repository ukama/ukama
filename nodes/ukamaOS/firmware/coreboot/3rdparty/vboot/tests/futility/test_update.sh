#!/bin/bash -eux
# Copyright 2018 The Chromium OS Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

me=${0##*/}
TMP="$me.tmp"

# Test --sys_props (primitive test needed for future updating tests).
test_sys_props() {
	! "${FUTILITY}" --debug update --sys_props "$*" 2>&1 |
		sed -n 's/.*property\[\(.*\)].value = \(.*\)/\1,\2,/p' |
		tr '\n' ' '
}

test "$(test_sys_props "1,2,3")" = "0,1, 1,2, 2,3, "
test "$(test_sys_props "1 2 3")" = "0,1, 1,2, 2,3, "
test "$(test_sys_props "1, 2,3 ")" = "0,1, 1,2, 2,3, "
test "$(test_sys_props "   1,, 2")" = "0,1, 2,2, "
test "$(test_sys_props " , 4,")" = "1,4, "

test_quirks() {
	! "${FUTILITY}" --debug update --quirks "$*" 2>&1 |
		sed -n 's/.*Set quirk \(.*\) to \(.*\)./\1,\2/p' |
		tr '\n' ' '
}

test "$(test_quirks "enlarge_image")" = "enlarge_image,1 "
test "$(test_quirks "enlarge_image=2")" = "enlarge_image,2 "
test "$(test_quirks " enlarge_image, enlarge_image=2")" = \
	"enlarge_image,1 enlarge_image,2 "

# Test data files
LINK_BIOS="${SCRIPTDIR}/data/bios_link_mp.bin"
PEPPY_BIOS="${SCRIPTDIR}/data/bios_peppy_mp.bin"
RO_VPD_BLOB="${SCRIPTDIR}/data/ro_vpd.bin"

# Work in scratch directory
cd "$OUTDIR"
set -o pipefail

# In all the test scenario, we want to test "updating from PEPPY to LINK".
TO_IMAGE=${TMP}.src.link
FROM_IMAGE=${TMP}.src.peppy
TO_HWID="X86 LINK TEST 6638"
FROM_HWID="X86 PEPPY TEST 4211"
cp -f ${LINK_BIOS} ${TO_IMAGE}
cp -f ${PEPPY_BIOS} ${FROM_IMAGE}
"${FUTILITY}" load_fmap "${FROM_IMAGE}" \
	RO_VPD:"${RO_VPD_BLOB}" RW_VPD:"${RO_VPD_BLOB}"

patch_file() {
	local file="$1"
	local section="$2"
	local section_offset="$3"
	local data="$4"

	# NAME OFFSET SIZE
	local fmap_info="$(${FUTILITY} dump_fmap -p ${file} ${section})"
	local base="$(echo "${fmap_info}" | sed 's/^[^ ]* //; s/ [^ ]*$//')"
	local offset=$((base + section_offset))
	echo "offset: ${offset}"
	printf "${data}" | dd of="${file}" bs=1 seek="${offset}" conv=notrunc
}

# PEPPY and LINK have different platform element ("Google_Link" and
# "Google_Peppy") in firmware ID so we want to hack them by changing
# "Google_" to "Google.".
patch_file ${TO_IMAGE} RW_FWID_A 0 Google.
patch_file ${TO_IMAGE} RW_FWID_B 0 Google.
patch_file ${TO_IMAGE} RO_FRID 0 Google.
patch_file ${FROM_IMAGE} RW_FWID_A 0 Google.
patch_file ${FROM_IMAGE} RW_FWID_B 0 Google.
patch_file ${FROM_IMAGE} RO_FRID 0 Google.

unpack_image() {
	local folder="${TMP}.$1"
	local image="$2"
	mkdir -p "${folder}"
	(cd "${folder}" && ${FUTILITY} dump_fmap -x "../${image}")
	${FUTILITY} gbb -g --rootkey="${folder}/rootkey" "${image}"
}

# Unpack images so we can prepare expected results by individual sections.
unpack_image "to" "${TO_IMAGE}"
unpack_image "from" "${FROM_IMAGE}"

# Hack FROM_IMAGE so it has same root key as TO_IMAGE (for RW update).
FROM_DIFFERENT_ROOTKEY_IMAGE="${FROM_IMAGE}2"
cp -f "${FROM_IMAGE}" "${FROM_DIFFERENT_ROOTKEY_IMAGE}"
"${FUTILITY}" gbb -s --rootkey="${TMP}.to/rootkey" "${FROM_IMAGE}"

# Hack for quirks
cp -f "${FROM_IMAGE}" "${FROM_IMAGE}.large"
truncate -s $((8388608 * 2)) "${FROM_IMAGE}.large"

# Create GBB v1.2 images (for checking digest)
GBB_OUTPUT="$("${FUTILITY}" gbb --digest "${TO_IMAGE}")"
[ "${GBB_OUTPUT}" = "digest: <none>" ]
TO_IMAGE_GBB12="${TO_IMAGE}.gbb12"
HWID_DIGEST="adf64d2a434b610506153da42440b0b498d7369c0e98b629ede65eb59f4784fa"
cp -f "${TO_IMAGE}" "${TO_IMAGE_GBB12}"
patch_file "${TO_IMAGE_GBB12}" GBB 6 "\x02"
"${FUTILITY}" gbb -s --hwid="${TO_HWID}" "${TO_IMAGE_GBB12}"
GBB_OUTPUT="$("${FUTILITY}" gbb --digest "${TO_IMAGE_GBB12}")"
[ "${GBB_OUTPUT}" = "digest: ${HWID_DIGEST}   valid" ]

# Generate expected results.
cp -f "${TO_IMAGE}" "${TMP}.expected.full"
cp -f "${FROM_IMAGE}" "${TMP}.expected.rw"
cp -f "${FROM_IMAGE}" "${TMP}.expected.a"
cp -f "${FROM_IMAGE}" "${TMP}.expected.b"
cp -f "${FROM_IMAGE}" "${TMP}.expected.legacy"
"${FUTILITY}" gbb -s --hwid="${FROM_HWID}" "${TMP}.expected.full"
"${FUTILITY}" load_fmap "${TMP}.expected.full" \
	RW_VPD:${TMP}.from/RW_VPD \
	RO_VPD:${TMP}.from/RO_VPD
"${FUTILITY}" load_fmap "${TMP}.expected.rw" \
	RW_SECTION_A:${TMP}.to/RW_SECTION_A \
	RW_SECTION_B:${TMP}.to/RW_SECTION_B \
	RW_SHARED:${TMP}.to/RW_SHARED \
	RW_LEGACY:${TMP}.to/RW_LEGACY
"${FUTILITY}" load_fmap "${TMP}.expected.a" \
	RW_SECTION_A:${TMP}.to/RW_SECTION_A
"${FUTILITY}" load_fmap "${TMP}.expected.b" \
	RW_SECTION_B:${TMP}.to/RW_SECTION_B
"${FUTILITY}" load_fmap "${TMP}.expected.legacy" \
	RW_LEGACY:${TMP}.to/RW_LEGACY
cp -f "${TMP}.expected.full" "${TMP}.expected.full.gbb12"
patch_file "${TMP}.expected.full.gbb12" GBB 6 "\x02"
"${FUTILITY}" gbb -s --hwid="${FROM_HWID}" "${TMP}.expected.full.gbb12"
cp -f "${TMP}.expected.full" "${TMP}.expected.full.gbb0"
"${FUTILITY}" gbb -s --flags=0 "${TMP}.expected.full.gbb0"
cp -f "${FROM_IMAGE}" "${FROM_IMAGE}.gbb0"
"${FUTILITY}" gbb -s --flags=0 "${FROM_IMAGE}.gbb0"
cp -f "${TMP}.expected.full" "${TMP}.expected.large"
dd if=/dev/zero bs=8388608 count=1 | tr '\000' '\377' >>"${TMP}.expected.large"
cp -f "${TMP}.expected.full" "${TMP}.expected.me_unlocked"
patch_file "${TMP}.expected.me_unlocked" SI_DESC 128 \
	"\x00\xff\xff\xff\x00\xff\xff\xff\x00\xff\xff\xff"

# A special set of images that only RO_VPD is preserved (RW_VPD is wiped) using
# FMAP_AREA_PRESERVE (\010=0x08).
TO_IMAGE_WIPE_RW_VPD="${TO_IMAGE}.wipe_rw_vpd"
cp -f "${TO_IMAGE}" "${TO_IMAGE_WIPE_RW_VPD}"
patch_file ${TO_IMAGE_WIPE_RW_VPD} FMAP 0x3fc "$(printf '\010')"
cp -f "${TMP}.expected.full" "${TMP}.expected.full.empty_rw_vpd"
"${FUTILITY}" load_fmap "${TMP}.expected.full.empty_rw_vpd" \
	RW_VPD:"${TMP}.to/RW_VPD"
patch_file "${TMP}.expected.full.empty_rw_vpd" FMAP 0x3fc "$(printf '\010')"

test_update() {
	local test_name="$1"
	local emu_src="$2"
	local expected="$3"
	local error_msg="${expected#!}"
	local msg

	shift 3
	cp -f "${emu_src}" "${TMP}.emu"
	echo "*** Test Item: ${test_name}"
	if [ "${error_msg}" != "${expected}" ] && [ -n "${error_msg}" ]; then
		msg="$(! "${FUTILITY}" update --emulate "${TMP}.emu" "$@" 2>&1)"
		echo "${msg}" | grep -qF -- "${error_msg}"
	else
		"${FUTILITY}" update --emulate "${TMP}.emu" "$@"
		cmp "${TMP}.emu" "${expected}"
	fi
}

# --sys_props: mainfw_act, tpm_fwver, is_vboot2, platform_ver, [wp_hw, wp_sw]
# tpm_fwver = <data key version:16><firmware version:16>.
# TO_IMAGE is signed with data key version = 1, firmware version = 4 => 0x10004.

# Test Full update.
test_update "Full update" \
	"${FROM_IMAGE}" "${TMP}.expected.full" \
	-i "${TO_IMAGE}" --wp=0 --sys_props 0,0x10001,1

test_update "Full update (incompatible platform)" \
	"${FROM_IMAGE}" "!platform is not compatible" \
	-i "${LINK_BIOS}" --wp=0 --sys_props 0,0x10001,1

test_update "Full update (TPM Anti-rollback: data key)" \
	"${FROM_IMAGE}" "!Data key version rollback detected (2->1)" \
	-i "${TO_IMAGE}" --wp=0 --sys_props 1,0x20001,1

test_update "Full update (TPM Anti-rollback: kernel key)" \
	"${FROM_IMAGE}" "!Firmware version rollback detected (5->4)" \
	-i "${TO_IMAGE}" --wp=0 --sys_props 1,0x10005,1

test_update "Full update (TPM Anti-rollback: 0 as tpm_fwver)" \
	"${FROM_IMAGE}" "${TMP}.expected.full" \
	-i "${TO_IMAGE}" --wp=0 --sys_props 0,0x0,1

test_update "Full update (TPM check failure due to invalid tpm_fwver)" \
	"${FROM_IMAGE}" "!Invalid tpm_fwver: -1" \
	-i "${TO_IMAGE}" --wp=0 --sys_props 0,-1,1

test_update "Full update (Skip TPM check with --force)" \
	"${FROM_IMAGE}" "${TMP}.expected.full" \
	-i "${TO_IMAGE}" --wp=0 --sys_props 0,-1,1 --force

test_update "Full update (from stdin)" \
	"${FROM_IMAGE}" "${TMP}.expected.full" \
	-i - --wp=0 --sys_props 0,-1,1 --force <"${TO_IMAGE}"

test_update "Full update (GBB=0 -> 0)" \
	"${FROM_IMAGE}.gbb0" "${TMP}.expected.full.gbb0" \
	-i "${TO_IMAGE}" --wp=0 --sys_props 0,0x10001,1

test_update "Full update (--host_only)" \
	"${FROM_IMAGE}" "${TMP}.expected.full" \
	-i "${TO_IMAGE}" --wp=0 --sys_props 0,0x10001,1 \
	--host_only --ec_image non-exist.bin --pd_image non_exist.bin

test_update "Full update (GBB1.2 hwid digest)" \
	"${FROM_IMAGE}" "${TMP}.expected.full.gbb12" \
	-i "${TO_IMAGE_GBB12}" --wp=0 --sys_props 0,0x10001,1

test_update "Full update (Preserve VPD using FMAP_AREA_PRESERVE)" \
	"${FROM_IMAGE}" "${TMP}.expected.full.empty_rw_vpd" \
	-i "${TO_IMAGE_WIPE_RW_VPD}" --wp=0 --sys_props 0,0x10001,1


# Test RW-only update.
test_update "RW update" \
	"${FROM_IMAGE}" "${TMP}.expected.rw" \
	-i "${TO_IMAGE}" --wp=1 --sys_props 0,0x10001,1

test_update "RW update (incompatible platform)" \
	"${FROM_IMAGE}" "!platform is not compatible" \
	-i "${LINK_BIOS}" --wp=1 --sys_props 0,0x10001,1

test_update "RW update (incompatible rootkey)" \
	"${FROM_DIFFERENT_ROOTKEY_IMAGE}" "!RW signed by incompatible root key" \
	-i "${TO_IMAGE}" --wp=1 --sys_props 0,0x10001,1

test_update "RW update (TPM Anti-rollback: data key)" \
	"${FROM_IMAGE}" "!Data key version rollback detected (2->1)" \
	-i "${TO_IMAGE}" --wp=1 --sys_props 1,0x20001,1

test_update "RW update (TPM Anti-rollback: kernel key)" \
	"${FROM_IMAGE}" "!Firmware version rollback detected (5->4)" \
	-i "${TO_IMAGE}" --wp=1 --sys_props 1,0x10005,1

# Test Try-RW update (vboot2).
test_update "RW update (A->B)" \
	"${FROM_IMAGE}" "${TMP}.expected.b" \
	-i "${TO_IMAGE}" -t --wp=1 --sys_props 0,0x10001,1

test_update "RW update (B->A)" \
	"${FROM_IMAGE}" "${TMP}.expected.a" \
	-i "${TO_IMAGE}" -t --wp=1 --sys_props 1,0x10001,1

test_update "RW update -> fallback to RO+RW Full update" \
	"${FROM_IMAGE}" "${TMP}.expected.full" \
	-i "${TO_IMAGE}" -t --wp=0 --sys_props 1,0x10002,1
test_update "RW update (incompatible platform)" \
	"${FROM_IMAGE}" "!platform is not compatible" \
	-i "${LINK_BIOS}" -t --wp=1 --sys_props 0x10001,1

test_update "RW update (incompatible rootkey)" \
	"${FROM_DIFFERENT_ROOTKEY_IMAGE}" "!RW signed by incompatible root key" \
	-i "${TO_IMAGE}" -t --wp=1 --sys_props 0,0x10001,1

test_update "RW update (TPM Anti-rollback: data key)" \
	"${FROM_IMAGE}" "!Data key version rollback detected (2->1)" \
	-i "${TO_IMAGE}" -t --wp=1 --sys_props 1,0x20001,1

test_update "RW update (TPM Anti-rollback: kernel key)" \
	"${FROM_IMAGE}" "!Firmware version rollback detected (5->4)" \
	-i "${TO_IMAGE}" -t --wp=1 --sys_props 1,0x10005,1

test_update "RW update -> fallback to RO+RW Full update (TPM Anti-rollback)" \
	"${FROM_IMAGE}" "!Firmware version rollback detected (6->4)" \
	-i "${TO_IMAGE}" -t --wp=0 --sys_props 1,0x10006,1

# Test Try-RW update (vboot1).
test_update "RW update (vboot1, A->B)" \
	"${FROM_IMAGE}" "${TMP}.expected.b" \
	-i "${TO_IMAGE}" -t --wp=1 --sys_props 0,0 --sys_props 0,0x10001,0

test_update "RW update (vboot1, B->B)" \
	"${FROM_IMAGE}" "${TMP}.expected.b" \
	-i "${TO_IMAGE}" -t --wp=1 --sys_props 1,0 --sys_props 0,0x10001,0

# Test 'factory mode'
test_update "Factory mode update (WP=0)" \
	"${FROM_IMAGE}" "${TMP}.expected.full" \
	-i "${TO_IMAGE}" --wp=0 --sys_props 0,0x10001,1 --mode=factory

test_update "Factory mode update (WP=0)" \
	"${FROM_IMAGE}" "${TMP}.expected.full" \
	--factory -i "${TO_IMAGE}" --wp=0 --sys_props 0,0x10001,1

test_update "Factory mode update (WP=1)" \
	"${FROM_IMAGE}" "!remove write protection for factory mode" \
	-i "${TO_IMAGE}" --wp=1 --sys_props 0,0x10001,1 --mode=factory

test_update "Factory mode update (WP=1)" \
	"${FROM_IMAGE}" "!remove write protection for factory mode" \
	--factory -i "${TO_IMAGE}" --wp=1 --sys_props 0,0x10001,1

test_update "Factory mode update (GBB=0 -> 39)" \
	"${FROM_IMAGE}.gbb0" "${TMP}.expected.full" \
	--factory -i "${TO_IMAGE}" --wp=0 --sys_props 0,0x10001,1

# Test legacy update
test_update "Legacy update" \
	"${FROM_IMAGE}" "${TMP}.expected.legacy" \
	-i "${TO_IMAGE}" --mode=legacy

# Test quirks
test_update "Full update (wrong size)" \
	"${FROM_IMAGE}.large" "!Image size is different" \
	-i "${TO_IMAGE}" --wp=0 --sys_props 0,0x10001,1

test_update "Full update (--quirks enlarge_image)" \
	"${FROM_IMAGE}.large" "${TMP}.expected.large" --quirks enlarge_image \
	-i "${TO_IMAGE}" --wp=0 --sys_props 0,0x10001,1

test_update "Full update (--quirks unlock_me_for_update)" \
	"${FROM_IMAGE}" "${TMP}.expected.me_unlocked" \
	--quirks unlock_me_for_update \
	-i "${TO_IMAGE}" --wp=0 --sys_props 0,0x10001,1

test_update "Full update (failure by --quirks min_platform_version)" \
	"${FROM_IMAGE}" "!Need platform version >= 3 (current is 2)" \
	--quirks min_platform_version=3 \
	-i "${TO_IMAGE}" --wp=0 --sys_props 0,0x10001,1,2

test_update "Full update (--quirks min_platform_version)" \
	"${FROM_IMAGE}" "${TMP}.expected.full" \
	--quirks min_platform_version=3 \
	-i "${TO_IMAGE}" --wp=0 --sys_props 0,0x10001,1,3

# Test archive and manifest.
A="${TMP}.archive"
mkdir -p "${A}/bin"
echo 'echo "${WL_TAG}"' >"${A}/bin/vpd"
chmod +x "${A}/bin/vpd"

cp -f "${LINK_BIOS}" "${A}/bios.bin"
echo "TEST: Manifest (--manifest, bios.bin)"
${FUTILITY} update -a "${A}" --manifest >"${TMP}.json.out"
cmp "${TMP}.json.out" "${SCRIPTDIR}/link_bios.manifest.json"

mv -f "${A}/bios.bin" "${A}/image.bin"
echo "TEST: Manifest (--manifest, image.bin)"
${FUTILITY} update -a "${A}" --manifest >"${TMP}.json.out"
cmp "${TMP}.json.out" "${SCRIPTDIR}/link_image.manifest.json"


cp -f "${TO_IMAGE}" "${A}/image.bin"
test_update "Full update (--archive, single package)" \
	"${FROM_IMAGE}" "${TMP}.expected.full" \
	-a "${A}" --wp=0 --sys_props 0,0x10001,1,3

echo "TEST: Output (--mode=output)"
mkdir -p "${TMP}.output"
${FUTILITY} update -i "${LINK_BIOS}" --mode=output --output_dir="${TMP}.output"
cmp "${LINK_BIOS}" "${TMP}.output/image.bin"

mkdir -p "${A}/keyset"
cp -f "${LINK_BIOS}" "${A}/image.bin"
cp -f "${TMP}.to/rootkey" "${A}/keyset/rootkey.WL"
cp -f "${TMP}.to/VBLOCK_A" "${A}/keyset/vblock_A.WL"
cp -f "${TMP}.to/VBLOCK_B" "${A}/keyset/vblock_B.WL"
${FUTILITY} gbb -s --rootkey="${TMP}.from/rootkey" "${A}/image.bin"
${FUTILITY} load_fmap "${A}/image.bin" VBLOCK_A:"${TMP}.from/VBLOCK_A"
${FUTILITY} load_fmap "${A}/image.bin" VBLOCK_B:"${TMP}.from/VBLOCK_B"

test_update "Full update (--archive, whitelabel, no VPD)" \
	"${A}/image.bin" "!Need VPD set for white" \
	-a "${A}" --wp=0 --sys_props 0,0x10001,1,3

test_update "Full update (--archive, whitelabel, no VPD - factory mode)" \
	"${LINK_BIOS}" "${A}/image.bin" \
	-a "${A}" --wp=0 --sys_props 0,0x10001,1,3 --mode=factory

test_update "Full update (--archive, whitelabel, no VPD - quirk mode)" \
	"${LINK_BIOS}" "${A}/image.bin" \
	-a "${A}" --wp=0 --sys_props 0,0x10001,1,3 --quirks=allow_empty_wltag

test_update "Full update (--archive, WL, single package)" \
	"${A}/image.bin" "${LINK_BIOS}" \
	-a "${A}" --wp=0 --sys_props 0,0x10001,1,3 --signature_id=WL

WL_TAG="WL" PATH="${A}/bin:${PATH}" \
	test_update "Full update (--archive, WL, fake vpd)" \
	"${A}/image.bin" "${LINK_BIOS}" \
	-a "${A}" --wp=0 --sys_props 0,0x10001,1,3

echo "TEST: Output (-a, --mode=output)"
mkdir -p "${TMP}.outa"
cp -f "${A}/image.bin" "${TMP}.emu"
WL_TAG="WL" PATH="${A}/bin:${PATH}" \
	${FUTILITY} update -a "${A}" --mode=output --emu="${TMP}.emu" \
	--output_dir="${TMP}.outa"
cmp "${LINK_BIOS}" "${TMP}.outa/image.bin"

# Test archive with Unified Build contents.
cp -r "${SCRIPTDIR}/models" "${A}/"
mkdir -p "${A}/images"
mv "${A}/image.bin" "${A}/images/bios_coral.bin"
cp -f "${PEPPY_BIOS}" "${A}/images/bios_peppy.bin"
cp -f "${LINK_BIOS}" "${A}/images/bios_link.bin"
cp -f "${TMP}.to/rootkey" "${A}/keyset/rootkey.whitetip-wl"
cp -f "${TMP}.to/VBLOCK_A" "${A}/keyset/vblock_A.whitetip-wl"
cp -f "${TMP}.to/VBLOCK_B" "${A}/keyset/vblock_B.whitetip-wl"
cp -f "${PEPPY_BIOS}" "${FROM_IMAGE}.ap"
cp -f "${LINK_BIOS}" "${FROM_IMAGE}.al"
patch_file ${FROM_IMAGE}.ap FW_MAIN_A 0 "corrupted"
patch_file ${FROM_IMAGE}.al FW_MAIN_A 0 "corrupted"
test_update "Full update (--archive, model=link)" \
	"${FROM_IMAGE}.al" "${LINK_BIOS}" \
	-a "${A}" --wp=0 --sys_props 0,0x10001,1,3 --model=link
test_update "Full update (--archive, model=peppy)" \
	"${FROM_IMAGE}.ap" "${PEPPY_BIOS}" \
	-a "${A}" --wp=0 --sys_props 0,0x10001,1,3 --model=peppy
test_update "Full update (--archive, model=unknown)" \
	"${FROM_IMAGE}.ap" "!Unsupported model: 'unknown'" \
	-a "${A}" --wp=0 --sys_props 0,0x10001,1,3 --model=unknown
test_update "Full update (--archive, model=whitetip, signature_id=WL)" \
	"${FROM_IMAGE}.al" "${LINK_BIOS}" \
	-a "${A}" --wp=0 --sys_props 0,0x10001,1,3 --model=whitetip \
	--signature_id=whitetip-wl

WL_TAG="wl" PATH="${A}/bin:${PATH}" \
	test_update "Full update (-a, model=WL, fake VPD)" \
	"${FROM_IMAGE}.al" "${LINK_BIOS}" \
	-a "${A}" --wp=0 --sys_props 0,0x10001,1,3 --model=whitetip

# WL-Unibuild without default keys
test_update "Full update (--a, model=WL, no VPD, no default keys)" \
	"${FROM_IMAGE}.al" "!Need VPD set for white" \
	-a "${A}" --wp=0 --sys_props 0,0x10001,1,3 --model=whitetip

# WL-Unibuild with default keys as model name
cp -f "${TMP}.to/rootkey" "${A}/keyset/rootkey.whitetip"
cp -f "${TMP}.to/VBLOCK_A" "${A}/keyset/vblock_A.whitetip"
cp -f "${TMP}.to/VBLOCK_B" "${A}/keyset/vblock_B.whitetip"
test_update "Full update (-a, model=WL, no VPD, default keys)" \
	"${FROM_IMAGE}.al" "${LINK_BIOS}" \
	-a "${A}" --wp=0 --sys_props 0,0x10001,1,3 --model=whitetip

# Test special programmer
if type flashrom >/dev/null 2>&1; then
	echo "TEST: Full update (dummy programmer)"
	cp -f "${FROM_IMAGE}" "${TMP}.emu"
	sudo "${FUTILITY}" update --programmer \
		dummy:emulate=VARIABLE_SIZE,image=${TMP}.emu,size=8388608 \
		-i "${TO_IMAGE}" --wp=0 --sys_props 0,0x10001,1,3 >&2
	cmp "${TMP}.emu" "${TMP}.expected.full"
fi

if type cbfstool >/dev/null 2>&1; then
	echo "SMM STORE" >"${TMP}.smm"
	truncate -s 262144 "${TMP}.smm"
	cp -f "${FROM_IMAGE}" "${TMP}.from.smm"
	cp -f "${TMP}.expected.full" "${TMP}.expected.full_smm"
	cbfstool "${TMP}.from.smm" add -r RW_LEGACY -n "smm_store" \
		-f "${TMP}.smm" -t raw
	cbfstool "${TMP}.expected.full_smm" add -r RW_LEGACY -n "smm_store" \
		-f "${TMP}.smm" -t raw -b 0x1bf000
	test_update "Legacy update (--quirks eve_smm_store)" \
		"${TMP}.from.smm" "${TMP}.expected.full_smm" \
		-i "${TO_IMAGE}" --wp=0 --sys_props 0,0x10001,1 \
		--quirks eve_smm_store
fi
