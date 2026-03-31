/* macOS shim for <linux/if_ether.h> */
#pragma once
#include <net/ethernet.h>
#define ETH_ALEN  6
#define ETH_P_IP  0x0800
#define ETH_P_ARP 0x0806
