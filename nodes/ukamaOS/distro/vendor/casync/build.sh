#!/bin/sh
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

set -eu

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
VERSION_FILE="${SCRIPT_DIR}/casync.version"
SHA_FILE="${SCRIPT_DIR}/casync.sha256"
PATCH_DIR="${SCRIPT_DIR}/patches"

if [ ! -f "${VERSION_FILE}" ]; then
    echo "missing ${VERSION_FILE}" >&2
    exit 1
fi

# shellcheck disable=SC1090
. "${VERSION_FILE}"

: "${CASYNC_REPO:?CASYNC_REPO is not set}"
: "${CASYNC_REF:?CASYNC_REF is not set}"

DEP_DIR="${DEP_DIR:-${SCRIPT_DIR}/../build}"
BUILD_DIR="${BUILD_DIR:-${SCRIPT_DIR}/tmpbuild}"
DIST_DIR="${DIST_DIR:-${SCRIPT_DIR}/distfiles}"
SRC_DIR="${BUILD_DIR}/src"
OBJ_DIR="${BUILD_DIR}/obj"
STAGE_DIR="${BUILD_DIR}/stage"

CC="${CC:-gcc}"
NPROCS="${NPROCS:-1}"

ARCHIVE="${DIST_DIR}/casync-${CASYNC_REF}.tar.gz"
EXPECTED_SHA=""

mkdir -p "${DEP_DIR}/bin"
mkdir -p "${DEP_DIR}/share/licenses/casync"
mkdir -p "${DIST_DIR}"
mkdir -p "${BUILD_DIR}"

if [ -f "${SHA_FILE}" ]; then
    EXPECTED_SHA="$(cat "${SHA_FILE}" | tr -d '[:space:]')"
fi

download_source() {
    if [ -f "${ARCHIVE}" ]; then
        return 0
    fi

    tmp_git="${BUILD_DIR}/git"

    rm -rf "${tmp_git}"
    git clone --no-checkout "${CASYNC_REPO}" "${tmp_git}"

    (
        cd "${tmp_git}"
        git fetch --depth 1 origin "${CASYNC_REF}"
        git checkout FETCH_HEAD
        git archive --format=tar.gz \
            --prefix="casync-${CASYNC_REF}/" \
            -o "${ARCHIVE}" HEAD
    )

    rm -rf "${tmp_git}"
}

verify_source() {
    actual_sha="$(sha256sum "${ARCHIVE}" | awk '{print $1}')"

    if [ "${EXPECTED_SHA}" = "UPDATE_ME" ] || [ -z "${EXPECTED_SHA}" ]; then
        echo "casync source archive SHA256:"
        echo "${actual_sha}"
        echo
        echo "Put this value into:"
        echo "${SHA_FILE}"
        exit 2
    fi

    if [ "${actual_sha}" != "${EXPECTED_SHA}" ]; then
        echo "casync SHA256 mismatch" >&2
        echo "expected: ${EXPECTED_SHA}" >&2
        echo "actual:   ${actual_sha}" >&2
        exit 1
    fi
}

extract_source() {
    rm -rf "${SRC_DIR}" "${OBJ_DIR}" "${STAGE_DIR}"
    mkdir -p "${SRC_DIR}" "${OBJ_DIR}" "${STAGE_DIR}"

    tar -xzf "${ARCHIVE}" -C "${SRC_DIR}" --strip-components=1
}

apply_patches() {
    if [ ! -d "${PATCH_DIR}" ]; then
        return 0
    fi

    for patch in "${PATCH_DIR}"/*.patch; do
        if [ -f "${patch}" ]; then
            echo "Applying patch: ${patch}"
            patch -d "${SRC_DIR}" -p1 < "${patch}"
        fi
    done
}

build_casync() {
    (
        cd "${SRC_DIR}"

        CASYNC_CFLAGS="${CFLAGS:-}"
        CASYNC_CFLAGS="-D_GNU_SOURCE ${CASYNC_CFLAGS}"
        CASYNC_CFLAGS="${CASYNC_CFLAGS} -Wno-error=nonnull"
        CASYNC_CFLAGS="${CASYNC_CFLAGS} -Wno-error=maybe-uninitialized"
        CASYNC_CFLAGS="${CASYNC_CFLAGS} -Wno-error=stringop-truncation"
        CASYNC_CFLAGS="${CASYNC_CFLAGS} -Wno-error=deprecated-declarations"

        CC="${CC}" CFLAGS="${CASYNC_CFLAGS}" meson setup "${OBJ_DIR}" \
            --prefix=/usr \
            --buildtype=release \
            --default-library=shared

        ninja -C "${OBJ_DIR}" -j "${NPROCS}"
        DESTDIR="${STAGE_DIR}" ninja -C "${OBJ_DIR}" install
    )
}

install_runtime() {
    if [ ! -x "${STAGE_DIR}/usr/bin/casync" ]; then
        echo "casync binary not found after build" >&2
        exit 1
    fi

    cp "${STAGE_DIR}/usr/bin/casync" "${DEP_DIR}/bin/casync"
    chmod 0755 "${DEP_DIR}/bin/casync"

    if [ -f "${SRC_DIR}/LICENSE.LGPL2.1" ]; then
        cp "${SRC_DIR}/LICENSE.LGPL2.1" \
            "${DEP_DIR}/share/licenses/casync/LICENSE.LGPL2.1"
    fi

    if [ -f "${SRC_DIR}/LICENSE" ]; then
        cp "${SRC_DIR}/LICENSE" \
            "${DEP_DIR}/share/licenses/casync/LICENSE"
    fi

    cat > "${DEP_DIR}/share/licenses/casync/NOTICE" <<EOF
casync is built from upstream source and installed as a separate runtime
executable used by wimc.d.

Repository: ${CASYNC_REPO}
Ref:        ${CASYNC_REF}
SHA256:     $(sha256sum "${ARCHIVE}" | awk '{print $1}')

Ukama does not link casync into wimc.d. wimc.d executes /usr/bin/casync
as an external program.
EOF
}

download_source
verify_source
extract_source
apply_patches
build_casync
install_runtime

echo "casync installed to ${DEP_DIR}/bin/casync"

