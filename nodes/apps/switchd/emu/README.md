# switchemu.d
Tycon management-plane emulator in C.

Build:
    make

Run:
    ./switchemu.d --http-port 18088 --snmp-port 1161 --tftp-port 1069

Debug endpoints:
- GET /debug/state
- GET /debug/ports
- GET /debug/firmware
- POST /debug/scenario {"name":"normal"}
- POST /debug/ports/<id>/link {"value":"up|down"}
- POST /debug/ports/<id>/poe {"value":"on|off"}
- POST /debug/switch/reachable {"value":true|false}

SNMP:
- Minimal v2c GET/SET agent for key Tycon enterprise OIDs
- Simplified line protocol also accepted for easy manual testing:
  GET <oid>\n
  SET <oid> <int>\n
TFTP:
- Simple UDP write receiver for staging firmware files

Scenarios:
- scenarios/normal.json
- scenarios/high_temp.json
- scenarios/poe_fault_port4.json
- scenarios/firmware_fail.json
