#!/usr/bin/env python3
"""
VE.Direct TEXT Protocol Emulator for controller.d testing.

This tool creates a virtual serial port and emits VE.Direct TEXT frames
that look like a Victron SmartSolar MPPT controller. It supports both a
compressed scenario simulator for fast testing and a real-time,
stateful plant emulator driven by wall-clock time, site location,
sunrise/sunset, PV array sizing, battery state, and site load.

Examples:

  python3 vedirect_sim.py --mode scenario --scenario sunny

  python3 vedirect_sim.py --mode scenario --scenario all

  python3 vedirect_sim.py --mode realtime --city "Goma, DRC"

  python3 vedirect_sim.py --mode realtime --lat -1.6792 \
      --lon 29.2228 --tz-offset 2

  python3 vedirect_sim.py --mode realtime --city "Nairobi, Kenya" \
      --pv-count 6 --pv-watts 550 --battery-ah 280 \
      --load-day-w 220 --load-night-w 150 --cloud-factor 0.90

  python3 vedirect_sim.py --mode realtime --city "Goma, DRC" \
      --use-api

Notes:
- The script creates a PTY pair and prints the slave path.
- Point controller.d to the printed serial port path.
- API/network lookups are optional and only happen at startup or on
  local date rollover.
- If API lookup fails, the emulator falls back to local approximations.

Requires:
- Python 3.9+ recommended
"""

from __future__ import annotations

from dataclasses import dataclass
from datetime import date
from datetime import datetime
from datetime import timedelta
from datetime import timezone
from typing import Callable
from typing import Optional
import argparse
import json
import math
import os
import pty
import signal
import sys
import time
import urllib.error
import urllib.parse
import urllib.request

PRODUCT_ID = "0xA042"
FIRMWARE = "159"
SERIAL = "HQ2242ABCDE"
INTERVAL_S = 1.0

SUNNY_CYCLE_S = 180.0
NIGHT_CYCLE_S = 120.0

SUNNY_DAY_YIELD_KWH = 5.6
YESTERDAY_YIELD_KWH = 5.2
TOTAL_YIELD_BASE_KWH = 1248.4
SUNNY_PEAK_POWER_W = 980
YESTERDAY_MAX_POWER_W = 910

CITY_PRESETS = {
    "goma": {
        "display_name": "Goma, DRC",
        "lat": -1.6792,
        "lon": 29.2228,
        "tz_offset_h": 2.0,
        "sunrise_h": 5.67,
        "sunset_h": 18.17,
    },
    "nairobi": {
        "display_name": "Nairobi, Kenya",
        "lat": -1.2864,
        "lon": 36.8172,
        "tz_offset_h": 3.0,
        "sunrise_h": 6.30,
        "sunset_h": 18.30,
    },
    "mbuji-mayi": {
        "display_name": "Mbuji-Mayi, DRC",
        "lat": -6.1360,
        "lon": 23.5898,
        "tz_offset_h": 2.0,
        "sunrise_h": 5.85,
        "sunset_h": 17.95,
    },
    "kinshasa": {
        "display_name": "Kinshasa, DRC",
        "lat": -4.4419,
        "lon": 15.2663,
        "tz_offset_h": 1.0,
        "sunrise_h": 5.95,
        "sunset_h": 18.00,
    },
}


@dataclass(frozen=True)
class ScenarioSpec:
    title: str
    summary: str
    duration_s: float
    frame_func: Callable[[float], tuple[bytes, dict[str, object]]]


@dataclass
class SiteInfo:
    name: str
    lat: float
    lon: float
    tz_offset_h: float
    sunrise_h: float
    sunset_h: float


@dataclass
class PlantConfig:
    site: SiteInfo
    pv_count: int = 4
    pv_watts_each: int = 450
    pv_derate: float = 0.82
    battery_nominal_v: float = 48.0
    battery_capacity_ah: float = 200.0
    soc_init: float = 0.72
    load_day_w: float = 180.0
    load_night_w: float = 120.0
    load_base_output_enabled: bool = True
    cloud_factor: float = 0.95
    ambient_day_c: float = 31.0
    ambient_night_c: float = 22.0
    total_yield_base_kwh: float = TOTAL_YIELD_BASE_KWH
    yesterday_yield_kwh: float = YESTERDAY_YIELD_KWH
    yesterday_max_power_w: int = YESTERDAY_MAX_POWER_W
    use_api: bool = False
    geocode_query: Optional[str] = None


@dataclass
class PlantState:
    soc: float
    batt_v: float
    batt_a: float = 0.0
    pv_w: int = 0
    pv_mv: int = 0
    temp_c: float = 25.0
    load_a: float = 0.0
    cs: int = 0
    mppt: int = 0
    err: str = "0"
    today_yield_kwh: float = 0.0
    total_yield_kwh: float = TOTAL_YIELD_BASE_KWH
    today_max_power_w: int = 0
    yesterday_yield_kwh: float = YESTERDAY_YIELD_KWH
    yesterday_max_power_w: int = YESTERDAY_MAX_POWER_W
    last_update_ts: float = 0.0
    last_local_date: Optional[date] = None
    sunrise_h: float = 6.0
    sunset_h: float = 18.0


def _safe_unlink(path: str) -> None:
    try:
        os.unlink(path)
    except FileNotFoundError:
        pass


def publish_serial_endpoint(slave_path: str,
                            serial_link: Optional[str],
                            ready_file: Optional[str]) -> None:
    if serial_link:
        parent = os.path.dirname(serial_link)
        if parent:
            os.makedirs(parent, exist_ok=True)

        if os.path.lexists(serial_link):
            _safe_unlink(serial_link)

        os.symlink(slave_path, serial_link)

    if ready_file:
        parent = os.path.dirname(ready_file)
        if parent:
            os.makedirs(parent, exist_ok=True)

        with open(ready_file, "w", encoding="utf-8") as file_desc:
            file_desc.write(slave_path + "\n")


def cleanup_serial_endpoint(serial_link: Optional[str],
                            ready_file: Optional[str]) -> None:
    if serial_link and os.path.islink(serial_link):
        _safe_unlink(serial_link)

    if ready_file and os.path.exists(ready_file):
        _safe_unlink(ready_file)


def clamp(value: float, low: float, high: float) -> float:
    return max(low, min(value, high))


def lerp(start: float, end: float, factor: float) -> float:
    return start + (end - start) * factor


def smoothstep01(value: float) -> float:
    value = clamp(value, 0.0, 1.0)
    return value * value * (3.0 - 2.0 * value)


def cycle_ratio(t: float, cycle_s: float) -> float:
    return (t % cycle_s) / cycle_s


def compressed_hour(t: float, cycle_s: float,
                    start_h: float, end_h: float) -> float:
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
    Build a valid VE.Direct TEXT frame from a list of (label, value)
    pairs. Appends a Checksum field so that sum of all frame bytes is
    0 mod 256.
    """
    body = b"".join(f"{k}\t{v}\r\n".encode("ascii") for k, v in fields)
    prefix = body + b"Checksum\t"
    current_sum = sum(prefix) % 256
    checksum_byte = (-current_sum) % 256
    return prefix + bytes([checksum_byte]) + b"\r\n"


def scenario_frame_sunny(t: float) -> tuple[bytes, dict[str, object]]:
    hour = compressed_hour(t, SUNNY_CYCLE_S, 6.0, 18.5)
    day_factor = math.sin(math.pi * clamp((hour - 6.0) / 12.0, 0.0, 1.0))
    day_factor = max(0.0, day_factor)

    pv_w = int(round(SUNNY_PEAK_POWER_W * day_factor))
    batt_v = lerp(49.2, 56.4, smoothstep01(day_factor))
    batt_a = pv_w / batt_v if batt_v > 0 else 0.0
    temp_c = lerp(23.0, 36.0, smoothstep01(day_factor))
    mppt = 2 if pv_w > 80 else 0
    cs = 5 if pv_w > 220 else (3 if pv_w > 30 else 0)
    yield_today = SUNNY_DAY_YIELD_KWH * smoothstep01(day_factor)
    pmax = pv_w

    summary = {
        "site_name": "Sunny scenario",
        "local_time": hhmm_string(hour),
        "pv_w": pv_w,
        "batt_v": batt_v,
        "temp_c": temp_c,
        "yield_today_kwh": yield_today,
        "state": charge_state_label(cs),
        "err": "0",
    }

    frame = build_frame([
        ("PID", PRODUCT_ID),
        ("FW", FIRMWARE),
        ("SER#", SERIAL),
        ("V", str(int(round(batt_v * 1000.0)))),
        ("I", str(int(round(batt_a * 1000.0)))),
        ("VPV", str(int(round((batt_v + 4.0 + 8.0 * day_factor) * 1000.0)))),
        ("PPV", str(pv_w)),
        ("CS", str(cs)),
        ("MPPT", str(mppt)),
        ("OR", "0"),
        ("ERR", "0"),
        ("LOAD", "ON"),
        ("IL", "0"),
        ("H19", f"{yield_today:.2f}"),
        ("H20", f"{yield_today:.2f}"),
        ("H21", str(pmax)),
        ("H22", f"{YESTERDAY_YIELD_KWH:.2f}"),
        ("H23", str(YESTERDAY_MAX_POWER_W)),
        ("T", str(int(round(temp_c)))),
    ])
    return frame, summary


def scenario_frame_night(t: float) -> tuple[bytes, dict[str, object]]:
    hour = compressed_hour(t, NIGHT_CYCLE_S, 19.0, 5.5)
    phase = cycle_ratio(t, NIGHT_CYCLE_S)
    load_breath = 0.5 + 0.5 * math.sin(2.0 * math.pi * phase)

    pv_w = 0
    batt_v = lerp(52.4, 48.2, phase)
    batt_a = -lerp(1.8, 4.5, load_breath)
    temp_c = lerp(24.0, 20.5, smoothstep01(phase))
    cs = 0
    mppt = 0
    yield_today = 0.0
    pmax = 0
    err = "0"

    if 0.55 < phase < 0.68:
        err = "18"

    summary = {
        "site_name": "Night scenario",
        "local_time": hhmm_string(hour),
        "pv_w": pv_w,
        "batt_v": batt_v,
        "temp_c": temp_c,
        "yield_today_kwh": yield_today,
        "state": charge_state_label(cs),
        "err": err,
    }

    frame = build_frame([
        ("PID", PRODUCT_ID),
        ("FW", FIRMWARE),
        ("SER#", SERIAL),
        ("V", str(int(round(batt_v * 1000.0)))),
        ("I", str(int(round(batt_a * 1000.0)))),
        ("VPV", "0"),
        ("PPV", "0"),
        ("CS", str(cs)),
        ("MPPT", str(mppt)),
        ("OR", "0"),
        ("ERR", err),
        ("LOAD", "ON"),
        ("IL", "0"),
        ("H19", f"{yield_today:.2f}"),
        ("H20", f"{yield_today:.2f}"),
        ("H21", str(pmax)),
        ("H22", f"{YESTERDAY_YIELD_KWH:.2f}"),
        ("H23", str(YESTERDAY_MAX_POWER_W)),
        ("T", str(int(round(temp_c)))),
    ])
    return frame, summary


SCENARIOS = {
    "sunny": ScenarioSpec(
        title="Sunny daytime",
        summary="Healthy production arc with charging and harvest",
        duration_s=SUNNY_CYCLE_S,
        frame_func=scenario_frame_sunny,
    ),
    "night": ScenarioSpec(
        title="Night discharge",
        summary="No PV, site running from battery, occasional low-input alarm",
        duration_s=NIGHT_CYCLE_S,
        frame_func=scenario_frame_night,
    ),
}


def scenario_order(name: str) -> list[str]:
    if name == "all":
        return list(SCENARIOS.keys())
    return [name]


def print_banner_scenario(requested_name: str,
                          slave_path: str,
                          advertised_path: Optional[str] = None) -> None:
    order = scenario_order(requested_name)
    first = SCENARIOS[order[0]]
    display_path = advertised_path or slave_path

    print(f"\n{'=' * 76}")
    title = first.title if requested_name != "all" else "All scenarios"
    print(f"  VE.Direct Emulator — {title}")
    print(f"{'=' * 76}")
    print("\n  Mode: scenario (compressed synthetic test mode)")
    if requested_name == "all":
        print(f"  Scenario order: {', '.join(order)}")
    else:
        print(f"  Profile: {first.summary}")
    print(f"\n  Virtual serial port: {slave_path}")
    if advertised_path:
        print(f"  Stable serial link:  {advertised_path}")
    print("\n  Run the daemon with:")
    print("    CONTROLLER_DRIVER=victron "
          f"CONTROLLER_SERIAL_PORT={display_path} ./controllerd")
    print("\n  Then test the API:")
    print("    curl http://localhost:18021/v1/ping")
    print("    curl http://localhost:18021/v1/status  | python3 -m json.tool")
    print("    curl http://localhost:18021/v1/metrics | python3 -m json.tool")
    print("    curl http://localhost:18021/v1/alarms  | python3 -m json.tool")
    print("\n  Press Ctrl+C to stop.\n")
    print(f"{'=' * 76}\n")


def local_dt_and_hour(now_ts: float, tz_offset_h: float) -> tuple[datetime, float]:
    tz = timezone(timedelta(hours=tz_offset_h))
    dt = datetime.fromtimestamp(now_ts, tz=tz)
    hour = dt.hour + (dt.minute / 60.0) + (dt.second / 3600.0)
    return dt, hour


def estimate_sunrise_sunset_offline(site: SiteInfo, local_day: date) -> tuple[float, float]:
    day_of_year = local_day.timetuple().tm_yday
    seasonal = math.sin((2.0 * math.pi / 365.0) * (day_of_year - 80))
    shift = 0.30 * seasonal * math.cos(math.radians(site.lat))
    sunrise = clamp(site.sunrise_h - shift, 4.5, 7.5)
    sunset = clamp(site.sunset_h + shift, 16.5, 19.5)
    return sunrise, sunset


def geocode_location(query: str) -> Optional[dict[str, object]]:
    base = "https://geocoding-api.open-meteo.com/v1/search"
    params = urllib.parse.urlencode({"name": query, "count": 1, "language": "en", "format": "json"})
    url = f"{base}?{params}"

    try:
        with urllib.request.urlopen(url, timeout=10) as response:
            data = json.loads(response.read().decode("utf-8"))
    except (urllib.error.URLError, TimeoutError, json.JSONDecodeError):
        return None

    results = data.get("results") or []
    if not results:
        return None

    first = results[0]
    return {
        "name": first.get("name"),
        "country": first.get("country"),
        "latitude": first.get("latitude"),
        "longitude": first.get("longitude"),
        "timezone": first.get("timezone"),
        "utc_offset_seconds": first.get("utc_offset_seconds"),
    }


def fetch_sunrise_sunset_api(lat: float, lon: float,
                             local_day: date,
                             tz_offset_h: float) -> Optional[tuple[float, float]]:
    start = local_day.isoformat()
    end = local_day.isoformat()
    tz_name = "auto"

    params = urllib.parse.urlencode({
        "latitude": f"{lat:.6f}",
        "longitude": f"{lon:.6f}",
        "daily": "sunrise,sunset",
        "start_date": start,
        "end_date": end,
        "timezone": tz_name,
    })
    url = f"https://api.open-meteo.com/v1/forecast?{params}"

    try:
        with urllib.request.urlopen(url, timeout=10) as response:
            data = json.loads(response.read().decode("utf-8"))
    except (urllib.error.URLError, TimeoutError, json.JSONDecodeError):
        return None

    daily = data.get("daily") or {}
    sunrises = daily.get("sunrise") or []
    sunsets = daily.get("sunset") or []
    if not sunrises or not sunsets:
        return None

    try:
        sunrise_dt = datetime.fromisoformat(sunrises[0])
        sunset_dt = datetime.fromisoformat(sunsets[0])
    except ValueError:
        return None

    sunrise_h = sunrise_dt.hour + sunrise_dt.minute / 60.0 + sunrise_dt.second / 3600.0
    sunset_h = sunset_dt.hour + sunset_dt.minute / 60.0 + sunset_dt.second / 3600.0

    return sunrise_h, sunset_h


def resolve_city_preset(name: str) -> Optional[SiteInfo]:
    key = name.strip().lower()
    if key in CITY_PRESETS:
        item = CITY_PRESETS[key]
        return SiteInfo(
            name=item["display_name"],
            lat=item["lat"],
            lon=item["lon"],
            tz_offset_h=item["tz_offset_h"],
            sunrise_h=item["sunrise_h"],
            sunset_h=item["sunset_h"],
        )

    for _, item in CITY_PRESETS.items():
        if key == item["display_name"].lower():
            return SiteInfo(
                name=item["display_name"],
                lat=item["lat"],
                lon=item["lon"],
                tz_offset_h=item["tz_offset_h"],
                sunrise_h=item["sunrise_h"],
                sunset_h=item["sunset_h"],
            )
    return None


def resolve_site(args: argparse.Namespace) -> SiteInfo:
    if args.lat is not None or args.lon is not None:
        if args.lat is None or args.lon is None:
            raise SystemExit("--lat and --lon must be provided together")
        if args.tz_offset is None:
            raise SystemExit("--tz-offset is required with --lat/--lon")
        return SiteInfo(
            name="Custom coordinates",
            lat=args.lat,
            lon=args.lon,
            tz_offset_h=args.tz_offset,
            sunrise_h=6.0,
            sunset_h=18.0,
        )

    if args.city:
        preset = resolve_city_preset(args.city)
        if preset:
            return preset

        if args.use_api:
            geo = geocode_location(args.city)
            if geo is not None:
                offset_seconds = geo.get("utc_offset_seconds")
                if offset_seconds is None and args.tz_offset is None:
                    raise SystemExit("Geocoder did not return UTC offset; pass --tz-offset")
                tz_offset_h = (
                    float(offset_seconds) / 3600.0
                    if offset_seconds is not None else args.tz_offset
                )
                return SiteInfo(
                    name=f"{geo.get('name')}, {geo.get('country')}",
                    lat=float(geo["latitude"]),
                    lon=float(geo["longitude"]),
                    tz_offset_h=float(tz_offset_h),
                    sunrise_h=6.0,
                    sunset_h=18.0,
                )

        print(f"Warning: unknown city preset '{args.city}', "
              "using offline fallback. Pass --use-api or exact coords "
              "for a better site lookup.",
              file=sys.stderr)
        return SiteInfo(
            name=args.city,
            lat=0.0,
            lon=0.0,
            tz_offset_h=args.tz_offset if args.tz_offset is not None else 0.0,
            sunrise_h=6.0,
            sunset_h=18.0,
        )

    return resolve_city_preset("goma")


def estimate_batt_v_from_soc(soc: float,
                             cs: int,
                             batt_a: float,
                             nominal_v: float) -> float:
    if nominal_v < 20.0:
        v_empty = 11.8
        v_full = 14.4
    elif nominal_v < 40.0:
        v_empty = 23.6
        v_full = 28.8
    else:
        v_empty = 47.2
        v_full = 57.6

    base = lerp(v_empty, v_full, clamp(soc, 0.0, 1.0))

    if batt_a > 0.2:
        base += clamp(batt_a * 0.015, 0.0, 0.8)
    elif batt_a < -0.2:
        base += clamp(batt_a * 0.020, -1.0, 0.0)

    if cs == 5:
        base = max(base, nominal_v * 1.12)
    elif cs == 4:
        base = max(base, nominal_v * 1.15)

    return base


def realtime_tick(config: PlantConfig,
                  state: PlantState,
                  now_ts: float) -> tuple[bytes, dict[str, object]]:
    dt, local_hour = local_dt_and_hour(now_ts, config.site.tz_offset_h)
    local_day = dt.date()

    if state.last_local_date != local_day:
        if state.last_local_date is not None:
            state.yesterday_yield_kwh = state.today_yield_kwh
            state.yesterday_max_power_w = state.today_max_power_w

        state.today_yield_kwh = 0.0
        state.today_max_power_w = 0
        state.last_local_date = local_day

        if config.use_api:
            api_sun = fetch_sunrise_sunset_api(
                config.site.lat,
                config.site.lon,
                local_day,
                config.site.tz_offset_h,
            )
            if api_sun is not None:
                state.sunrise_h, state.sunset_h = api_sun
            else:
                state.sunrise_h, state.sunset_h = estimate_sunrise_sunset_offline(
                    config.site, local_day
                )
        else:
            state.sunrise_h, state.sunset_h = estimate_sunrise_sunset_offline(
                config.site, local_day
            )

    if state.last_update_ts <= 0.0:
        state.last_update_ts = now_ts

    dt_hours = clamp((now_ts - state.last_update_ts) / 3600.0, 0.0, 0.05)
    state.last_update_ts = now_ts

    day_start = state.sunrise_h
    day_end = state.sunset_h

    if local_hour < day_start or local_hour > day_end:
        solar_shape = 0.0
    else:
        span = max(0.1, day_end - day_start)
        solar_phase = (local_hour - day_start) / span
        solar_shape = math.sin(math.pi * clamp(solar_phase, 0.0, 1.0))
        solar_shape = max(0.0, solar_shape)

    cloud_wave = 0.90 + 0.10 * math.sin(now_ts / 900.0)
    cloud_factor = clamp(config.cloud_factor * cloud_wave, 0.25, 1.0)

    pv_capacity_w = config.pv_count * config.pv_watts_each * config.pv_derate
    pv_w = int(round(pv_capacity_w * solar_shape * cloud_factor))
    pv_w = max(0, pv_w)

    state.today_max_power_w = max(state.today_max_power_w, pv_w)

    is_day = pv_w > 20
    load_w = config.load_day_w if is_day else config.load_night_w
    net_w = pv_w - load_w

    batt_capacity_wh = max(1.0, config.battery_nominal_v * config.battery_capacity_ah)
    delta_soc = (net_w * dt_hours) / batt_capacity_wh
    state.soc = clamp(state.soc + delta_soc, 0.05, 1.00)

    if pv_w <= 5 and state.soc < 0.20:
        state.cs = 2
        state.err = "2"
        state.mppt = 0
    elif pv_w <= 5:
        state.cs = 0
        state.err = "0"
        state.mppt = 0
    elif state.soc < 0.85:
        state.cs = 3
        state.err = "0"
        state.mppt = 2
    elif state.soc < 0.97:
        state.cs = 4
        state.err = "0"
        state.mppt = 2
    else:
        state.cs = 5
        state.err = "0"
        state.mppt = 2

    batt_v_est = estimate_batt_v_from_soc(
        state.soc, state.cs,
        (net_w / max(config.battery_nominal_v, 1.0)),
        config.battery_nominal_v,
    )
    state.batt_v = batt_v_est
    state.batt_a = net_w / max(state.batt_v, 1.0)

    if is_day:
        temp_base = lerp(config.ambient_night_c, config.ambient_day_c, smoothstep01(solar_shape))
        state.temp_c = temp_base + (pv_w / max(pv_capacity_w, 1.0)) * 4.0
    else:
        night_phase = 0.5 + 0.5 * math.sin(now_ts / 3600.0)
        state.temp_c = lerp(config.ambient_night_c - 1.0, config.ambient_night_c + 1.0, night_phase)

    state.load_a = load_w / max(state.batt_v, 1.0)
    state.pv_w = pv_w
    state.pv_mv = int(round((state.batt_v + 3.0 + 10.0 * solar_shape) * 1000.0))

    if pv_w > 0:
        state.today_yield_kwh += (pv_w * dt_hours) / 1000.0
        state.total_yield_kwh += (pv_w * dt_hours) / 1000.0

    summary = {
        "site_name": config.site.name,
        "local_time": dt.strftime("%H:%M:%S"),
        "pv_w": state.pv_w,
        "batt_v": state.batt_v,
        "temp_c": state.temp_c,
        "yield_today_kwh": state.today_yield_kwh,
        "state": charge_state_label(state.cs),
        "err": state.err,
    }

    frame = build_frame([
        ("PID", PRODUCT_ID),
        ("FW", FIRMWARE),
        ("SER#", SERIAL),
        ("V", str(int(round(state.batt_v * 1000.0)))),
        ("I", str(int(round(state.batt_a * 1000.0)))),
        ("VPV", str(state.pv_mv)),
        ("PPV", str(state.pv_w)),
        ("CS", str(state.cs)),
        ("MPPT", str(state.mppt)),
        ("OR", "0"),
        ("ERR", state.err),
        ("LOAD", "ON"),
        ("IL", str(int(round(state.load_a * 1000.0)))),
        ("H19", f"{state.today_yield_kwh:.2f}"),
        ("H20", f"{state.total_yield_kwh:.2f}"),
        ("H21", str(state.today_max_power_w)),
        ("H22", f"{state.yesterday_yield_kwh:.2f}"),
        ("H23", str(state.yesterday_max_power_w)),
        ("T", str(int(round(state.temp_c)))),
    ])
    return frame, summary


def print_banner_realtime(config: PlantConfig,
                          slave_path: str,
                          advertised_path: Optional[str] = None) -> None:
    display_path = advertised_path or slave_path

    print(f"\n{'=' * 76}")
    print("  VE.Direct Emulator — Real-time plant emulation")
    print(f"{'=' * 76}")
    print("\n  Mode: realtime")
    print(f"  Site: {config.site.name}")
    print("  Coordinates: "
          f"lat={config.site.lat:.5f}, lon={config.site.lon:.5f}")
    print(f"  Timezone offset: UTC{config.site.tz_offset_h:+.1f}")
    print(f"  PV array: {config.pv_count} x {config.pv_watts_each}W "
          f"derate={config.pv_derate:.2f}")
    print(f"  Battery: {config.battery_nominal_v:.1f}V "
          f"{config.battery_capacity_ah:.0f}Ah "
          f"soc_init={config.soc_init:.2f}")
    print(f"  Load: day={config.load_day_w:.0f}W "
          f"night={config.load_night_w:.0f}W")
    print(f"  Cloud factor: {config.cloud_factor:.2f}")
    source = "API+fallback" if config.use_api else "offline fallback/preset"
    print(f"  Sunrise/sunset source: {source}")
    print(f"\n  Virtual serial port: {slave_path}")
    if advertised_path:
        print(f"  Stable serial link:  {advertised_path}")
    print("\n  Run the daemon with:")
    print("    CONTROLLER_DRIVER=victron "
          f"CONTROLLER_SERIAL_PORT={display_path} ./controllerd")
    print("\n  Then test the API:")
    print("    curl http://localhost:18021/v1/ping")
    print("    curl http://localhost:18021/v1/status  | python3 -m json.tool")
    print("    curl http://localhost:18021/v1/metrics | python3 -m json.tool")
    print("    curl http://localhost:18021/v1/alarms  | python3 -m json.tool")
    print("\n  Press Ctrl+C to stop.\n")
    print(f"{'=' * 76}\n")


def run_scenario(scenario_name: str,
                 interval: float,
                 serial_link: Optional[str] = None,
                 ready_file: Optional[str] = None) -> None:
    master_fd, slave_fd = pty.openpty()
    slave_path = os.ttyname(slave_fd)
    publish_serial_endpoint(slave_path, serial_link, ready_file)

    order = scenario_order(scenario_name)
    display_path = serial_link or slave_path

    print_banner_scenario(
        scenario_name,
        slave_path,
        advertised_path=display_path,
    )

    scenario_idx = 0
    t0 = time.time()
    scenario_start = t0
    frame_count = 0

    def handle_exit(_sig, _frame):
        print(f"\n\nSent {frame_count} frames. Bye.")
        cleanup_serial_endpoint(serial_link, ready_file)
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
            f"site={summary['site_name']} "
            f"time={summary['local_time']} "
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
            print("\n  --- switching to scenario: "
                  f"{next_name} ({next_spec.title}) ---\n")

        time.sleep(interval)

    cleanup_serial_endpoint(serial_link, ready_file)


def run_realtime(config: PlantConfig,
                 interval: float,
                 serial_link: Optional[str] = None,
                 ready_file: Optional[str] = None) -> None:
    master_fd, slave_fd = pty.openpty()
    slave_path = os.ttyname(slave_fd)
    publish_serial_endpoint(slave_path, serial_link, ready_file)

    display_path = serial_link or slave_path

    print_banner_realtime(
        config,
        slave_path,
        advertised_path=display_path,
    )

    init_batt_v = estimate_batt_v_from_soc(
        config.soc_init, 0, 0.0, config.battery_nominal_v
    )
    state = PlantState(
        soc=config.soc_init,
        batt_v=init_batt_v,
        total_yield_kwh=config.total_yield_base_kwh,
        yesterday_yield_kwh=config.yesterday_yield_kwh,
        yesterday_max_power_w=config.yesterday_max_power_w,
        sunrise_h=config.site.sunrise_h,
        sunset_h=config.site.sunset_h,
    )

    t0 = time.time()
    frame_count = 0

    def handle_exit(_sig, _frame):
        print(f"\n\nSent {frame_count} frames. Bye.")
        cleanup_serial_endpoint(serial_link, ready_file)
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

        frame, summary = realtime_tick(config, state, now)

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
            f"site={summary['site_name']} "
            f"time={summary['local_time']} "
            f"batt={summary['batt_v']:.2f}V "
            f"pv={summary['pv_w']:4.0f}W "
            f"temp={summary['temp_c']:.0f}C "
            f"yield={summary['yield_today_kwh']:.2f}kWh "
            f"state={summary['state']}{err_str}"
        )

        time.sleep(interval)

    cleanup_serial_endpoint(serial_link, ready_file)


def main() -> None:
    parser = argparse.ArgumentParser(
        description="VE.Direct SmartSolar emulator for controller.d",
    )
    parser.add_argument(
        "--mode",
        choices=("scenario", "realtime"),
        default="scenario",
        help="Emulation mode",
    )
    parser.add_argument(
        "--scenario",
        choices=("sunny", "night", "all"),
        default="sunny",
        help="Scenario profile for scenario mode",
    )
    parser.add_argument(
        "--interval",
        type=float,
        default=INTERVAL_S,
        help="Seconds between frames",
    )

    parser.add_argument(
        "--serial-link",
        default=None,
        help="Stable symlink for the PTY slave, e.g. /tmp/victron-tty",
    )
    parser.add_argument(
        "--ready-file",
        default=None,
        help="Write the resolved PTY slave path here when ready",
    )

    parser.add_argument(
        "--city",
        help='Preset city or free-form location, e.g. "Goma, DRC"',
    )
    parser.add_argument("--lat", type=float, help="Latitude")
    parser.add_argument("--lon", type=float, help="Longitude")
    parser.add_argument(
        "--tz-offset",
        type=float,
        default=None,
        help="UTC offset in hours for realtime mode",
    )
    parser.add_argument(
        "--use-api",
        action="store_true",
        help="Use online geocoding and sunrise/sunset lookup",
    )

    parser.add_argument(
        "--pv-count",
        type=int,
        default=4,
        help="Number of PV panels",
    )
    parser.add_argument(
        "--pv-watts",
        type=int,
        default=450,
        help="Watts per PV panel",
    )
    parser.add_argument(
        "--pv-derate",
        type=float,
        default=0.82,
        help="PV derate factor 0..1",
    )
    parser.add_argument(
        "--battery-v",
        type=float,
        default=48.0,
        help="Nominal battery voltage",
    )
    parser.add_argument(
        "--battery-ah",
        type=float,
        default=200.0,
        help="Battery capacity Ah",
    )
    parser.add_argument(
        "--soc-init",
        type=float,
        default=0.72,
        help="Initial battery SOC 0..1",
    )
    parser.add_argument(
        "--load-day-w",
        type=float,
        default=180.0,
        help="Daytime site load W",
    )
    parser.add_argument(
        "--load-night-w",
        type=float,
        default=120.0,
        help="Nighttime site load W",
    )
    parser.add_argument(
        "--cloud-factor",
        type=float,
        default=0.95,
        help="Base cloud factor 0..1",
    )
    parser.add_argument(
        "--ambient-day-c",
        type=float,
        default=31.0,
        help="Day ambient temp C",
    )
    parser.add_argument(
        "--ambient-night-c",
        type=float,
        default=22.0,
        help="Night ambient temp C",
    )

    args = parser.parse_args()

    if args.mode == "scenario":
        run_scenario(
            args.scenario,
            args.interval,
            serial_link=args.serial_link,
            ready_file=args.ready_file,
        )
        return

    site = resolve_site(args)
    config = PlantConfig(
        site=site,
        pv_count=args.pv_count,
        pv_watts_each=args.pv_watts,
        pv_derate=args.pv_derate,
        battery_nominal_v=args.battery_v,
        battery_capacity_ah=args.battery_ah,
        soc_init=clamp(args.soc_init, 0.05, 1.0),
        load_day_w=max(0.0, args.load_day_w),
        load_night_w=max(0.0, args.load_night_w),
        cloud_factor=clamp(args.cloud_factor, 0.2, 1.0),
        ambient_day_c=args.ambient_day_c,
        ambient_night_c=args.ambient_night_c,
        use_api=args.use_api,
        geocode_query=args.city,
    )
    run_realtime(
        config,
        args.interval,
        serial_link=args.serial_link,
        ready_file=args.ready_file,
    )


if __name__ == "__main__":
    main()
