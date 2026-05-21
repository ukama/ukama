# Running the UE E2E Stack

Run these commands from the repository root unless noted otherwise. The stack is
Alpine-based and uses Podman only.

## 1. Build images

```bash
testing/ue/scripts/build-images.sh
```

This builds Alpine-based images:

```text
ukama/ue:dev
ukama/media:dev
```

Override names if needed:

```bash
UE_IMAGE=my/ue:dev MEDIA_IMAGE=my/media:dev     testing/ue/scripts/build-images.sh
```

## 2. Export environment

On the machine running the UE containers:

```bash
export TOWER_IP=<virtual-node-or-tower-ip>
export MEDIA_IP=<media-host-ip>
export UE_DATA_HOST=<this-host-ip-reachable-from-tower>
```

Optional defaults:

```bash
export EPCEMU_PORT=18092
export EPCEMU_DATA_PORT=18110
export PCRF_PORT=18090
export INITNET_PORT=18091
export UE_BASE_PORT=41000
```

`UE_DATA_HOST` matters. It must be an address that the virtual tower can reach.
The UE containers publish their UDP data ports on this host address.

## 3. Start media target

On the media host:

```bash
testing/ue/scripts/run-media.sh
```

This starts:

```text
HTTP   :8080
iperf3 :5201
```

If media is on a different machine from the UE host, build or copy the media
image there first using Podman and run only the media script on that machine.

## 4. Check tower-side services

From the UE host:

```bash
testing/ue/scripts/status.sh
```

Expected services:

```text
init-network /v1/status
pcrf         /v1/status
epcemu       /v1/status
```

## 5. Start one UE

```bash
testing/ue/scripts/run-ue.sh     --csv testing/ue/csv/SimPool.with-imsi.csv     --imsi 001010000000001
```

This starts a container named:

```text
ue-001010000000001
```

Inside that container, `ue-agent` creates `tun0`, attaches the IMSI to
`epcemu`, and waits for traffic.

## 6. Send traffic

Ping:

```bash
testing/ue/scripts/traffic-ue.sh     --imsi 001010000000001     --mode ping
```

HTTP:

```bash
testing/ue/scripts/traffic-ue.sh     --imsi 001010000000001     --mode http
```

iperf:

```bash
testing/ue/scripts/traffic-ue.sh     --imsi 001010000000001     --mode iperf
```

## 7. Start many UEs

```bash
testing/ue/scripts/run-many-ues.sh     --csv testing/ue/csv/SimPool.with-imsi.csv     --count 10
```

This starts the first 10 enabled IMSIs from the CSV.

## 8. Detach one UE

```bash
testing/ue/scripts/detach-ue.sh --imsi 001010000000001
```

## 9. Cleanup all UEs

```bash
testing/ue/scripts/cleanup-ues.sh
```

## Validation checklist

A good E2E run should show:

```text
ue-agent attaches IMSI successfully
epcemu status shows attached UE
ping/http/iperf traffic succeeds
PCRF or OVS counters increase for that UE
no traffic bypasses the PCRF datapath
```

The last point is important. A successful ping alone is not enough; the test is
only meaningful if PCRF/OVS counters prove that traffic crossed the intended
node datapath.

