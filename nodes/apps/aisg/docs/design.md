# AISG Design

## Architecture

```text
aisgd
  -> controller contract
  -> aisg-ctrl / aisg-emu
  -> raw-rs485 / stm-uart
  -> real RET or emulated RET
```

Production path:

```text
aisgd
  -> controller contract
  -> aisg-ctrl/stm-uart
  -> STM daughter card
  -> real RET
```

V1 is single RET only. Daisy-chain and firmware download are out of scope.

Overall design

                 +------------------+
                 |      aisgd       |
                 | REST/status/etc. |
                 +---------+--------+
                           |
                 southbound controller contract
                           |
      +--------------------+--------------------+
      |                    |                    |
      v                    v                    v
 aisg-ctrl            aisg-emu            aisg-ctrl
(raw-rs485 mode)     (software only)      (stm-uart mode)
      |                                       |
      v                                       v
USB-RS485 + PSU                           STM daughter card (cnode)
      |                                       |
      v                                       v
   real RET                               real RET

~                                                           
