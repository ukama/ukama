# Ukama UE E2E Test Stack

This directory contains the UE-side end-to-end test stack used to drive
traffic through a virtual Ukama tower node. The containers are Alpine-based
and the scripts use Podman.

The stack lives under:

```text
testing/ue
```

It is intentionally small:

```text
csv/        SIM/UE inventory used by the scripts
ue-agent/   C dataplane process that bridges UE tun0 to epcemu UDP
ue/         UE container image wrapper around ue-agent
media/      Simple Alpine traffic target container, HTTP + iperf3
scripts/    Podman build, run, traffic, status, and cleanup helpers
```

## What ue-agent does

`ue-agent` is the UE dataplane shim. It is not EPC, PCRF, or a media
gateway.

For each simulated UE it:

1. creates a `tun0` interface inside the UE container,
2. assigns the UE IP from the CSV,
3. attaches the IMSI to `epcemu` over the control API,
4. reads raw IP packets from `tun0`,
5. sends those packets over UDP to `epcemu-data`,
6. receives downlink packets back over UDP, and
7. writes them back into `tun0`.

The scripts keep the name `ue-agent` because it describes the process role.
The container image is named `ukama/ue:dev` because, operationally, each
container represents one UE.

## Traffic path

```text
UE app: ping/curl/iperf
  |
  v
UE tun0
  |
  v
ue-agent
  |
  | UDP tunnel
  v
epcemu-data on virtual tower
  |
  v
epcemu tun3
  |
  v
OVS / PCRF datapath
  |
  v
media target
```

Return traffic follows the same path in reverse.

## Control path

`run-ue.sh` starts one UE container. During startup, `ue-agent` sends:

```text
POST http://$TOWER_IP:18092/v1/ue/attach
```

with the IMSI, ICCID, UE IP, APN, and UE UDP endpoint. That tells `epcemu`
where to send downlink user-plane packets for that IMSI.

Policy creation is not handled here. The backend-to-node PCRF flow should
already install the UE policy on the virtual node.

## Ports

Default tower-side ports:

```text
init-network  18091/tcp
pcrf          18090/tcp
epcemu ctrl   18092/tcp
epcemu data   18110/udp
```

Default UE-side UDP ports start at `41000`. The scripts derive a per-UE port
from the last three digits of the IMSI:

```text
UE_DATA_PORT = UE_BASE_PORT + last_three_digits(IMSI)
```

The UE container publishes that UDP port to the host so the virtual tower can
send downlink packets back to the UE agent.

## CSV

The scripts expect a CSV with at least these columns:

```text
IMSI,ICCID,UE_IP,APN,Enabled
```

Only rows with `Enabled=TRUE` are used by `run-many-ues.sh`.

## Typical files

```text
testing/ue/csv/SimPool.with-imsi.csv
testing/ue/scripts/build-images.sh
testing/ue/scripts/run-media.sh
testing/ue/scripts/run-ue.sh
testing/ue/scripts/run-many-ues.sh
testing/ue/scripts/traffic-ue.sh
testing/ue/scripts/detach-ue.sh
testing/ue/scripts/cleanup-ues.sh
testing/ue/scripts/status.sh
```

See `RUN.md` for the command sequence.
