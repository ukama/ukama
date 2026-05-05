# casync vendor recipe

This directory contains the pinned build recipe for `casync`.

Ukama uses `casync` as a separate runtime executable for `wimc.d`.
It is not linked into Ukama binaries.

## Layout

```text
distro/vendor/casync/
  build.sh
  casync.version
  casync.sha256
  patches/
  README.md
  NOTICE
