#!/bin/bash
#
# Handy script to rebuild the markdown version of the man pages.
# This uses https://github.com/mle86/man-to-md if it is installed.

if [[ -z "$(which man-to-md)" ]]; then
    echo "man-to-md not found - skipping conversion"
    exit 0
fi

outdir="$1"
if [[ -z "${outdir}" ]]; then
    echo "usage $0 <outdir>"
    exit 1
fi

mkdir -p "${outdir}"
if [[ $? -ne 0 ]]; then
    echo "failed to make output directory: ${outdir}"
    exit 1
fi

index="${outdir}/index.md"

function do_page () {
    m="$1"
    base="${m%.*}"
    sect="${m#*.}"
    output="${base}-${sect}.md"

    redir="$(grep '^.so man' "${m}")"
    if [[ $? -eq 0 ]]; then
	r="${redir#*/}"
	rbase="${r%.*}"
	rsect="${r#*.}"
	echo "* [${base}(${sect})](${rbase}-${rsect}.md)" >> "${index}"
	return
    fi

    man-to-md -f < "${m}" | sed -e 's/^\*\*\([^\*]\+\)\*\*(\([138]\+\))/[\1(\2)](\1-\2.md)/' > "${outdir}/${base}-${sect}.md"
    echo "* [${base}(${sect})](${base}-${sect}.md)" >> "${index}"
}

cat > "${index}" <<EOF
# Manpages for libcap and libpsx

## Individual reference pages
EOF

# Assumes the m's are listed alphabetically.
for n in 1 3 8 ; do
	cat >> "${index}" <<EOF

### Section ${n}

EOF
    for m in *.${n}; do
	do_page "${m}"
    done
done

cat >> "${index}" <<EOF

## More information

For further information, see the
[FullyCapable](https://sites.google.com/site/fullycapable/) homepage
for libcap.

## MD page generation

These official man pages for libcap were converted to markdown using
[man-to-md](https://github.com/mle86/man-to-md).

EOF
