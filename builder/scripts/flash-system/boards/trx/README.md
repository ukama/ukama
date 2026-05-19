# TRX Board (Tower Node Radio)

## Hardware setup

```
[Linux host] ── ethernet ──┐
                           ├─→ [common subnet 192.168.53.x]
[BDI 192.168.53.72] ───────┘
       │
       └── JTAG cable ──→ [TRX board] ── UART (USB) ──→ [host /dev/ttyUSB0]
                                 │
                                 └── PoE or DC power
```

## Prerequisites

| Item | Where it lives | Set in |
|---|---|---|
| `oct-remote-boot` | Octeon TIP SDK | `oct_remote_boot.path` in board.yaml |
| `CNF71XX.cfg` | BDI config | `bdi.config_file` in board.yaml |
| `lsm_os_trx.gz`, `lsm_rd_trx.gz`, `u-boot-octeon_trx.bin` | Phase 1 artifacts | `phase1.artifacts.*.path` |
| 8 `flash_*.img` files | Phase 2 artifacts | `phase2.images.*.src` |
| Band config `<band>.cfg` | Per-band tuning | `band.configs_dir/<band>.cfg` |
| TRX root password | per-device firmware default | env var `TRX_ROOT_PASSWORD` |

Adjust paths in `board.yaml` to match where the files actually live on your build host.

## Usage

```bash
export TRX_ROOT_PASSWORD='...'
./flash trx
```

To select a different band:

```bash
BAND=b7 ./flash trx
```

## Flow

**Phase 1 — JTAG bringup** (~5 min, automated)
1. Validates artifacts, BDI reachability, `oct-remote-boot` presence.
2. Starts a local TFTP server with the Phase 1 artifacts.
3. Telnets to BDI, sends `go 0x400000`.
4. Runs `oct-remote-boot`, watches serial for `Octeon zen(ram)=>` prompt.
5. Pushes u-boot env vars + `saveenv`.
6. Enables ethernet and pings.
7. TFTPs OS, RD, uboot into DDR and flashes each at the correct address.

**Manual pause**
- Operator powers OFF the TRX
- Disconnects the BDI/JTAG cable
- Powers ON the TRX (boots from newly-flashed u-boot)
- Presses ENTER to continue

**Phase 2 — Image flash via SSH** (~12 min, automated)
8. SSHs into TRX as root.
9. For each of 8 `.img` files: scp to `/mnt/tmp`, `dd` to its `/dev/flash_*`, delete from `/mnt/tmp`.
10. Copies the selected band config to `/etc/trx/band.cfg`.

**Verify**
11. SSHs back in, confirms the band config is applied.
