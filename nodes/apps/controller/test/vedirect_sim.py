#!/usr/bin/env python3
"""
VE.Direct TEXT Protocol Simulator for controller.d testing.

Creates a virtual serial port pair (PTY) and streams realistic
Victron SmartSolar MPPT 150/45 frames so the daemon can be tested
without physical hardware.

Usage:
    python3 vedirect_sim.py [--scenario SCENARIO] [--interval SECS]

Scenarios:
    sunny       Clear-sky daytime profile for Goma (default)
    low_batt    Battery drops below warning/critical thresholds
    fault       Controller reports error code 17 (charger temp high)
    night       Overnight profile for Goma with load-driven discharge
    all         Cycles through all scenarios in sequence
"""

from dataclasses import dataclass
from typing import Callable
import argparse
import math
import os
import pty
import signal
import sys
import time

PRODUCT_ID = "0xA042"  # SmartSolar MPPT 150/45
FIRMWARE = "159"
SERIAL = "HQ2242ABCDE"
INTERVAL_S = 1.0

GOMA_LOCATION = "Goma, DRC"
GOMA_TIMEZONE = "UTC+2"
GOMA_SUNRISE_H = 5.67   # ~05:40
GOMA_SUNSET_H = 18.17   # ~18:10
SUNNY_CYCLE_S = 180.0
NIGHT_CYCLE_S = 120.0

SUNNY_DAY_YIELD_KWH = 5.6
YESTERDAY_YIELD_KWH = 5.2
TOTAL_YIELD_BASE_KWH = 1248.4
SUNNY_PEAK_POWER_W = 980
YESTERDAY_MAX_POWER_W = 910


@dataclass(frozen=True)
class ScenarioSpec:
    title: str
    summary: str
    duration_s: float
    frame_func: Callable[[float], tuple[bytes, dict[str, object]]]


def clamp(value: float, low: float, high: float) -> float:
    return max(low, min(value, high))


def lerp(start: float, end: float, factor: float) -> float:
    return start + (end - start) * factor


def smoothstep01(value: float) -> float:
    value = clamp(value, 0.0, 1.0)
    return value * value * (3.0 - 2.0 * value)


def cycle_ratio(t: float, cycle_s: float) -> float:
    return (t % cycle_s) / cycle_s


def compressed_hour(t: float, cycle_s: float, start_h: float, end_h: float) -> float:
    return lerp(start_h, end_h, cycle_ratio(t, cycle_s))


def hhmm_string(hour: float) -> str:
    total_minutes = int(round((hour % 24.0) * 60.0))
    hh = (total_minutes // 60) % 24
    mm = total_minutes % 60
    return f"{hh:02d}:{mm:02d}"


def charge_state_label(cs_value: int) -> str:
    return {
        0: "Off",
        1: "Low power",
        2: "Fault",
        3: "Bulk",
        4: "Absorption",
        5: "Float",
        6: "Storage",
    }.get(cs_value, f"CS={cs_value}")


def build_frame(fields: list[tuple[str, str]]) -> bytes:
    """
    Build a valid VE.Direct TEXT frame from a list of (label, value) pairs.
    Appends a Checksum field so that sum of all frame bytes == 0 (mod 256).
    """
    frame = b""
    for label, value in fields:
        frame += f"{label}\t{value}\r\n".encode("latin-1")

    prefix = frame + b"Checksum\t"
    total = sum(prefix) % 256
    chk = (256 - total) % 256

    return prefix + bytes([chk]) + b"\r\n"


def frame_sunny(t: float) -> tuple[bytes, dict[str, object]]:
    """
    Compressed clear-sky day in Goma.

    Models an equatorial day with near-constant day length, a strong midday
    solar peak, bulk charging in the morning, absorption around noon, and
    float in the afternoon.
    """
    local_hour = compressed_hour(t, SUNNY_CYCLE_S, GOMA_SUNRISE_H, GOMA_SUNSET_H)
    daylight = (local_hour - GOMA_SUNRISE_H) / (GOMA_SUNSET_H - GOMA_SUNRISE_H)
    daylight = clamp(daylight, 0.0, 1.0)

    solar_shape = math.sin(math.pi * daylight)
    solar_shape = max(0.0, solar_shape) ** 1.15
    yield_ratio = smoothstep01(daylight)

    pv_w = int(round(SUNNY_PEAK_POWER_W * solar_shape))
    pv_mv = int(round(90000 + 38000 * solar_shape)) if pv_w > 0 else 0
    load_a = lerp(1.1, 1.8, solar_shape)
    temp_c = lerp(21.0, 34.0, solar_shape)

    if daylight < 0.55:
        stage = smoothstep01(daylight / 0.55)
        batt_v = lerp(50.8, 57.4, stage)
        batt_a = lerp(2.5, 18.0, solar_shape)
        cs, mppt = 3, 2   # BULK
    elif daylight < 0.78:
        stage = smoothstep01((daylight - 0.55) / 0.23)
        batt_v = lerp(57.4, 57.7, stage)
        batt_a = lerp(10.5, 4.0, stage)
        cs, mppt = 4, 2   # ABSORPTION
    else:
        stage = smoothstep01((daylight - 0.78) / 0.22)
        batt_v = lerp(56.2, 55.2, stage)
        batt_a = lerp(2.5, -0.4, stage)
        cs, mppt = 5, 1   # FLOAT

    yield_today = SUNNY_DAY_YIELD_KWH * yield_ratio
    max_power_today = max(pv_w, int(round(SUNNY_PEAK_POWER_W * yield_ratio)))
    total_yield = TOTAL_YIELD_BASE_KWH + yield_today

    fields = [
        ("PID", PRODUCT_ID),
        ("FW", FIRMWARE),
        ("SER#", SERIAL),
        ("V", str(int(round(batt_v * 1000)))),
        ("I", str(int(round(batt_a * 1000)))),
        ("VPV", str(pv_mv)),
        ("PPV", str(pv_w)),
        ("CS", str(cs)),
        ("MPPT", str(mppt)),
        ("ERR", "0"),
        ("T", str(int(round(temp_c)))),
        ("LOAD", "ON"),
        ("IL", str(int(round(load_a * 1000)))),
        ("Relay", "OFF"),
        ("H19", str(int(round(total_yield * 100)))),
        ("H20", str(int(round(yield_today * 100)))),
        ("H21", str(max_power_today)),
        ("H22", str(int(round(YESTERDAY_YIELD_KWH * 100)))),
        ("H23", str(YESTERDAY_MAX_POWER_W)),
    ]

    summary = {
        "local_time": hhmm_string(local_hour),
        "batt_v": batt_v,
        "pv_w": pv_w,
        "temp_c": temp_c,
        "state": charge_state_label(cs),
        "yield_today_kwh": yield_today,
        "err": "0",
    }
    return build_frame(fields), summary


def frame_low_batt(t: float) -> tuple[bytes, dict[str, object]]:
    """Battery drains through the warning and critical thresholds."""
    phase = clamp((t % 60.0) / 60.0, 0.0, 1.0)
    batt_v = lerp(52.0, 42.0, phase)
    batt_a = -2.0
    temp_c = 24.0

    fields = [
        ("PID", PRODUCT_ID),
        ("FW", FIRMWARE),
        ("SER#", SERIAL),
        ("V", str(int(round(batt_v * 1000)))),
        ("I", str(int(round(batt_a * 1000)))),
        ("VPV", "0"),
        ("PPV", "0"),
        ("CS", "0"),
        ("MPPT", "0"),
        ("ERR", "0"),
        ("T", str(int(round(temp_c)))),
        ("LOAD", "ON"),
        ("IL", "1800"),
        ("Relay", "OFF"),
        ("H19", str(int(round((TOTAL_YIELD_BASE_KWH + YESTERDAY_YIELD_KWH) * 100)))),
        ("H20", "0"),
        ("H21", "0"),
        ("H22", str(int(round(YESTERDAY_YIELD_KWH * 100)))),
        ("H23", str(YESTERDAY_MAX_POWER_W)),
    ]

    summary = {
        "local_time": "03:20",
        "batt_v": batt_v,
        "pv_w": 0,
        "temp_c": temp_c,
        "state": charge_state_label(0),
        "yield_today_kwh": 0.0,
        "err": "0",
    }
    return build_frame(fields), summary


def frame_fault(t: float) -> tuple[bytes, dict[str, object]]:
    """Controller reports ERR=17 on a hot afternoon."""
    pv_w = 120
    temp_c = 68.0
    fields = [
        ("PID", PRODUCT_ID),
        ("FW", FIRMWARE),
        ("SER#", SERIAL),
        ("V", "51200"),
        ("I", "0"),
        ("VPV", "145000"),
        ("PPV", str(pv_w)),
        ("CS", "2"),        # FAULT state
        ("MPPT", "0"),
        ("ERR", "17"),      # VERR_CHARGER_TEMP_HIGH
        ("T", str(int(round(temp_c)))),
        ("LOAD", "ON"),
        ("IL", "1500"),
        ("Relay", "OFF"),
        ("H19", str(int(round((TOTAL_YIELD_BASE_KWH + 4.8) * 100)))),
        ("H20", "480"),
        ("H21", "860"),
        ("H22", str(int(round(YESTERDAY_YIELD_KWH * 100)))),
        ("H23", str(YESTERDAY_MAX_POWER_W)),
    ]

    summary = {
        "local_time": "13:40",
        "batt_v": 51.2,
        "pv_w": pv_w,
        "temp_c": temp_c,
        "state": charge_state_label(2),
        "yield_today_kwh": 4.8,
        "err": "17",
    }
    return build_frame(fields), summary


def frame_night(t: float) -> tuple[bytes, dict[str, object]]:
    """
    Compressed Goma overnight profile.

    Covers post-sunset through pre-dawn, with site load slowly discharging the
    battery bank, no solar production, and the daily yield resetting at midnight.
    """
    start_h = GOMA_SUNSET_H + 0.35
    end_h = GOMA_SUNRISE_H + 24.0 - 0.25
    absolute_hour = compressed_hour(t, NIGHT_CYCLE_S, start_h, end_h)
    phase = clamp((absolute_hour - start_h) / (end_h - start_h), 0.0, 1.0)
    local_hour = absolute_hour % 24.0

    batt_v = lerp(54.8, 50.9, smoothstep01(phase))
    batt_a = -lerp(1.3, 2.4, phase)
    temp_c = lerp(28.0, 21.5, smoothstep01(phase))
    load_a = lerp(1.4, 2.2, phase)

    if phase < 0.12:
        cs, mppt = 5, 1   # Briefly still at float after sunset
    else:
        cs, mppt = 1, 0   # Low power overnight

    if absolute_hour < 24.0:
        yield_today = 5.4
    else:
        yield_today = 0.0

    total_yield = TOTAL_YIELD_BASE_KWH + 5.4

    fields = [
        ("PID", PRODUCT_ID),
        ("FW", FIRMWARE),
        ("SER#", SERIAL),
        ("V", str(int(round(batt_v * 1000)))),
        ("I", str(int(round(batt_a * 1000)))),
        ("VPV", "0"),
        ("PPV", "0"),
        ("CS", str(cs)),
        ("MPPT", str(mppt)),
        ("ERR", "0"),
        ("T", str(int(round(temp_c)))),
        ("LOAD", "ON"),
        ("IL", str(int(round(load_a * 1000)))),
        ("Relay", "OFF"),
        ("H19", str(int(round(total_yield * 100)))),
        ("H20", str(int(round(yield_today * 100)))),
        ("H21", "0"),
        ("H22", str(int(round(YESTERDAY_YIELD_KWH * 100)))),
        ("H23", str(YESTERDAY_MAX_POWER_W)),
    ]

    summary = {
        "local_time": hhmm_string(local_hour),
        "batt_v": batt_v,
        "pv_w": 0,
        "temp_c": temp_c,
        "state": charge_state_label(cs),
        "yield_today_kwh": yield_today,
        "err": "0",
    }
    return build_frame(fields), summary


SCENARIOS = {
    "sunny": ScenarioSpec(
        title="Sunny day in Goma",
        summary="Compressed clear-sky day with bulk -> absorption -> float and a ~5.6 kWh daily yield.",
        duration_s=SUNNY_CYCLE_S,
        frame_func=frame_sunny,
    ),
    "night": ScenarioSpec(
        title="Night in Goma",
        summary="Compressed post-sunset to pre-dawn profile with battery discharge under site load and no PV production.",
        duration_s=NIGHT_CYCLE_S,
        frame_func=frame_night,
    ),
    "low_batt": ScenarioSpec(
        title="Low battery",
        summary="Synthetic discharge case that crosses the 46V and 44V alarm thresholds.",
        duration_s=45.0,
        frame_func=frame_low_batt,
    ),
    "fault": ScenarioSpec(
        title="Controller fault",
        summary="Hot-controller fault injection with ERR=17 and elevated temperature.",
        duration_s=30.0,
        frame_func=frame_fault,
    ),
}


def scenario_order(name: str) -> list[str]:
    if name == "all":
        return ["sunny", "night", "low_batt", "fault"]
    return [name]


def print_banner(requested_name: str, slave_path: str) -> None:
    order = scenario_order(requested_name)
    first = SCENARIOS[order[0]]

    print(f"\n{'=' * 68}")
    print(f"  VE.Direct Simulator — {first.title if requested_name != 'all' else 'All scenarios'}")
    print(f"{'=' * 68}")
    print(f"\n  Location model: {GOMA_LOCATION} ({GOMA_TIMEZONE})")
    print(f"  Daylight model: sunrise ~{hhmm_string(GOMA_SUNRISE_H)}, sunset ~{hhmm_string(GOMA_SUNSET_H)}")
    if requested_name == "all":
        print(f"  Scenario order: {', '.join(order)}")
    else:
        print(f"  Profile: {first.summary}")
    print(f"\n  Virtual serial port: {slave_path}")
    print(f"\n  Run the daemon with:")
    print(f"    CONTROLLER_DRIVER=victron CONTROLLER_SERIAL_PORT={slave_path} ./controllerd")
    print(f"\n  Then test the API:")
    print(f"    curl http://localhost:8095/v1/controller/status | python3 -m json.tool")
    print(f"    curl http://localhost:8095/v1/controller/metrics | python3 -m json.tool")
    print(f"    curl http://localhost:8095/v1/controller/alarms  | python3 -m json.tool")
    print(f"    ./test_api.sh")
    print(f"\n  Press Ctrl+C to stop.\n")
    print(f"{'=' * 68}\n")


def run(scenario_name: str, interval: float) -> None:
    master_fd, slave_fd = pty.openpty()
    slave_path = os.ttyname(slave_fd)
    order = scenario_order(scenario_name)

    print_banner(scenario_name, slave_path)

    scenario_idx = 0
    t0 = time.time()
    scenario_start = t0
    frame_count = 0

    def handle_exit(_sig, _frame):
        print(f"\n\nSent {frame_count} frames. Bye.")
        try:
            os.close(master_fd)
        except OSError:
            pass
        try:
            os.close(slave_fd)
        except OSError:
            pass
        sys.exit(0)

    signal.signal(signal.SIGINT, handle_exit)
    signal.signal(signal.SIGTERM, handle_exit)

    while True:
        now = time.time()
        elapsed = now - t0
        rel_t = now - scenario_start
        current_name = order[scenario_idx]
        spec = SCENARIOS[current_name]

        frame, summary = spec.frame_func(rel_t)

        try:
            os.write(master_fd, frame)
        except OSError:
            print("Daemon closed the port.")
            break

        frame_count += 1

        err = summary["err"]
        err_str = f" err={err}" if err != "0" else ""
        print(
            f"  [{elapsed:6.1f}s] frame={frame_count:4d} "
            f"goma={summary['local_time']} "
            f"batt={summary['batt_v']:.2f}V "
            f"pv={summary['pv_w']:4.0f}W "
            f"temp={summary['temp_c']:.0f}C "
            f"yield={summary['yield_today_kwh']:.2f}kWh "
            f"state={summary['state']}{err_str}"
        )

        if scenario_name == "all" and rel_t >= spec.duration_s:
            scenario_idx = (scenario_idx + 1) % len(order)
            scenario_start = now
            next_name = order[scenario_idx]
            next_spec = SCENARIOS[next_name]
            print(f"\n  --- switching to scenario: {next_name} ({next_spec.title}) ---\n")

        time.sleep(interval)


def main() -> None:
    parser = argparse.ArgumentParser(description="VE.Direct frame simulator for controller.d")
    parser.add_argument(
        "--scenario",
        choices=list(SCENARIOS.keys()) + ["all"],
        default="sunny",
        help="Scenario to simulate (default: sunny)",
    )
    parser.add_argument(
        "--interval",
        type=float,
        default=INTERVAL_S,
        help="Seconds between frames (default: 1.0)",
    )
    args = parser.parse_args()
    run(args.scenario, args.interval)


if __name__ == "__main__":
    main()
