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
    """
    Publish a stable path for the dynamically allocated PTY slave.

    serial_link:
        Symlink path such as /tmp/victron-tty -> /dev/pts/N

    ready_file:
        Small text file containing the resolved PTY slave path.
        Useful for supervisor/container startup sequencing.
    """
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

        with open(ready_file, "w", encoding="utf-8") as f:
            f.write(slave_path + "\n")


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
    frame = b""
    for label, value in fields:
        line = f"{label}\t{value}\r\n".encode("latin-1")
        frame += line

    prefix = frame + b"Checksum\t"
    total = sum(prefix) % 256
    chk = (256 - total) % 256

    return prefix + bytes([chk]) + b"\r\n"


def frame_sunny(t: float) -> tuple[bytes, dict[str, object]]:
    local_hour = compressed_hour(t, SUNNY_CYCLE_S, 5.67, 18.17)
    daylight = (local_hour - 5.67) / (18.17 - 5.67)
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
        cs, mppt = 3, 2
    elif daylight < 0.78:
        stage = smoothstep01((daylight - 0.55) / 0.23)
        batt_v = lerp(57.4, 57.7, stage)
        batt_a = lerp(10.5, 4.0, stage)
        cs, mppt = 4, 2
    else:
        stage = smoothstep01((daylight - 0.78) / 0.22)
        batt_v = lerp(56.2, 55.2, stage)
        batt_a = lerp(2.5, -0.4, stage)
        cs, mppt = 5, 1

    yield_today = SUNNY_DAY_YIELD_KWH * yield_ratio
    max_power_today = max(pv_w, int(round(SUNNY_PEAK_POWER_W *
                                          yield_ratio)))
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
        "site_name": "Goma, DRC",
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
        ("H19", str(int(round((TOTAL_YIELD_BASE_KWH +
                               YESTERDAY_YIELD_KWH) * 100)))),
        ("H20", "0"),
        ("H21", "0"),
        ("H22", str(int(round(YESTERDAY_YIELD_KWH * 100)))),
        ("H23", str(YESTERDAY_MAX_POWER_W)),
    ]

    summary = {
        "site_name": "Synthetic",
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
        ("CS", "2"),
        ("MPPT", "0"),
        ("ERR", "17"),
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
        "site_name": "Synthetic",
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
    start_h = 18.17 + 0.35
    end_h = 5.67 + 24.0 - 0.25
    absolute_hour = compressed_hour(t, NIGHT_CYCLE_S, start_h, end_h)
    phase = clamp((absolute_hour - start_h) / (end_h - start_h), 0.0, 1.0)
    local_hour = absolute_hour % 24.0

    batt_v = lerp(54.8, 50.9, smoothstep01(phase))
    batt_a = -lerp(1.3, 2.4, phase)
    temp_c = lerp(28.0, 21.5, smoothstep01(phase))
    load_a = lerp(1.4, 2.2, phase)

    if phase < 0.12:
        cs, mppt = 5, 1
    else:
        cs, mppt = 1, 0

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
        "site_name": "Goma, DRC",
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
        summary="Compressed clear-sky day with bulk, absorption, "
                "float, and about 5.6 kWh daily yield.",
        duration_s=SUNNY_CYCLE_S,
        frame_func=frame_sunny,
    ),
    "night": ScenarioSpec(
        title="Night in Goma",
        summary="Compressed post-sunset to pre-dawn profile with "
                "battery discharge under site load and no PV.",
        duration_s=NIGHT_CYCLE_S,
        frame_func=frame_night,
    ),
    "low_batt": ScenarioSpec(
        title="Low battery",
        summary="Synthetic discharge case that crosses the 46V and "
                "44V alarm thresholds.",
        duration_s=45.0,
        frame_func=frame_low_batt,
    ),
    "fault": ScenarioSpec(
        title="Controller fault",
        summary="Hot-controller fault injection with ERR=17 and "
                "elevated temperature.",
        duration_s=30.0,
        frame_func=frame_fault,
    ),
}


def scenario_order(name: str) -> list[str]:
    if name == "all":
        return ["sunny", "night", "low_batt", "fault"]
    return [name]

def city_key(value: str) -> str:
    value = value.strip().lower()
    value = value.replace(",", " ")
    value = value.replace("_", " ")
    value = " ".join(value.split())
    return value.replace(" ", "-")

def resolve_city_preset(value: str) -> Optional[dict[str, object]]:
    key = city_key(value)

    if key in CITY_PRESETS:
        return CITY_PRESETS[key]

    for preset_key, preset in CITY_PRESETS.items():
        display_key = city_key(preset["display_name"])
        if key == display_key:
            return preset

        city_only = city_key(preset["display_name"].split(",")[0])
        if key == city_only:
            return preset

    return None

def hour_from_dt(dt: datetime) -> float:
    return dt.hour + dt.minute / 60.0 + dt.second / 3600.0

def now_at_offset(tz_offset_h: float) -> datetime:
    tz = timezone(timedelta(hours=tz_offset_h))
    return datetime.now(tz)

def fallback_sunrise_sunset(lat: float,
                            day_of_year: int) -> tuple[float, float]:
    """
    Cheap local approximation for offline use.
    """
    lat_factor = clamp(abs(lat) / 60.0, 0.0, 1.0)
    seasonal = math.cos((2.0 * math.pi * (day_of_year - 172)) / 365.25)
    daylight_h = 12.1 + 2.2 * lat_factor * seasonal
    daylight_h = clamp(daylight_h, 10.7, 13.4)
    solar_noon = 12.0
    sunrise = solar_noon - daylight_h / 2.0
    sunset = solar_noon + daylight_h / 2.0
    return sunrise, sunset


def http_get_json(url: str, timeout_s: float = 8.0) -> dict:
    req = urllib.request.Request(
        url,
        headers={
            "User-Agent": "vedirect-sim/1.0 (+controller.d test utility)"
        },
    )
    with urllib.request.urlopen(req, timeout=timeout_s) as resp:
        return json.loads(resp.read().decode("utf-8"))


def geocode_city_nominatim(query: str) -> Optional[tuple[float, float, str]]:
    params = urllib.parse.urlencode({
        "q": query,
        "format": "jsonv2",
        "limit": 1,
    })
    url = f"https://nominatim.openstreetmap.org/search?{params}"
    try:
        data = http_get_json(url)
        if not data:
            return None
        row = data[0]
        return float(row["lat"]), float(row["lon"]), row.get(
            "display_name", query
        )
    except Exception:
        return None


def lookup_sunrise_sunset_api(lat: float, lon: float,
                              local_dt: datetime) -> Optional[
                                  tuple[float, float]]:
    day_str = local_dt.date().isoformat()
    params = urllib.parse.urlencode({
        "lat": f"{lat:.6f}",
        "lng": f"{lon:.6f}",
        "date": day_str,
        "formatted": 0,
    })
    url = f"https://api.sunrise-sunset.org/json?{params}"
    try:
        data = http_get_json(url)
        if data.get("status") != "OK":
            return None

        results = data["results"]
        sunrise_utc = datetime.fromisoformat(
            results["sunrise"].replace("Z", "+00:00")
        )
        sunset_utc = datetime.fromisoformat(
            results["sunset"].replace("Z", "+00:00")
        )

        local_tz = local_dt.tzinfo
        sunrise_local = sunrise_utc.astimezone(local_tz)
        sunset_local = sunset_utc.astimezone(local_tz)
        return hour_from_dt(sunrise_local), hour_from_dt(sunset_local)
    except Exception:
        return None

def resolve_site(args: argparse.Namespace) -> SiteInfo:
    if args.city:
        city = resolve_city_preset(args.city)
        if city is not None:
            tz_offset_h = city["tz_offset_h"]
            if args.tz_offset is not None:
                tz_offset_h = args.tz_offset

            return SiteInfo(
                name=city["display_name"],
                lat=city["lat"],
                lon=city["lon"],
                tz_offset_h=tz_offset_h,
                sunrise_h=city["sunrise_h"],
                sunset_h=city["sunset_h"],
            )

        if not args.use_api:
            raise SystemExit(
                f"Unknown city preset: {args.city!r}. Use --use-api "
                f"for online geocoding, or provide --lat/--lon."
            )

        geo = geocode_city_nominatim(args.city)
        if geo is None:
            raise SystemExit(f"Could not geocode city: {args.city!r}")

        lat, lon, display_name = geo
        tz_offset_h = 0.0 if args.tz_offset is None else args.tz_offset
        day_of_year = now_at_offset(tz_offset_h).timetuple().tm_yday
        sunrise_h, sunset_h = fallback_sunrise_sunset(
            lat,
            day_of_year,
        )

        return SiteInfo(
            name=display_name,
            lat=lat,
            lon=lon,
            tz_offset_h=tz_offset_h,
            sunrise_h=sunrise_h,
            sunset_h=sunset_h,
        )

    if args.lat is not None and args.lon is not None:
        tz_offset_h = 0.0 if args.tz_offset is None else args.tz_offset
        day_of_year = now_at_offset(tz_offset_h).timetuple().tm_yday
        sunrise_h, sunset_h = fallback_sunrise_sunset(
            args.lat,
            day_of_year,
        )
        return SiteInfo(
            name=f"Lat {args.lat:.4f}, Lon {args.lon:.4f}",
            lat=args.lat,
            lon=args.lon,
            tz_offset_h=tz_offset_h,
            sunrise_h=sunrise_h,
            sunset_h=sunset_h,
        )

    city = CITY_PRESETS["goma"]
    tz_offset_h = city["tz_offset_h"]
    if args.tz_offset is not None:
        tz_offset_h = args.tz_offset

    return SiteInfo(
        name=city["display_name"],
        lat=city["lat"],
        lon=city["lon"],
        tz_offset_h=tz_offset_h,
        sunrise_h=city["sunrise_h"],
        sunset_h=city["sunset_h"],
    )

def pv_peak_w(config: PlantConfig) -> float:
    return config.pv_count * config.pv_watts_each * config.pv_derate


def update_site_sun_times_if_needed(config: PlantConfig,
                                    state: PlantState,
                                    local_dt: datetime) -> None:
    if state.last_local_date == local_dt.date():
        return

    prior_date = state.last_local_date
    state.last_local_date = local_dt.date()

    sunrise_h = config.site.sunrise_h
    sunset_h = config.site.sunset_h

    if config.use_api:
        online = lookup_sunrise_sunset_api(
            config.site.lat, config.site.lon, local_dt
        )
        if online is not None:
            sunrise_h, sunset_h = online
        else:
            day_of_year = local_dt.timetuple().tm_yday
            sunrise_h, sunset_h = fallback_sunrise_sunset(
                config.site.lat, day_of_year
            )
    else:
        day_of_year = local_dt.timetuple().tm_yday
        sunrise_h, sunset_h = fallback_sunrise_sunset(
            config.site.lat, day_of_year
        )

    state.sunrise_h = sunrise_h
    state.sunset_h = sunset_h

    if prior_date is not None:
        state.yesterday_yield_kwh = state.today_yield_kwh
        state.yesterday_max_power_w = state.today_max_power_w
        state.today_yield_kwh = 0.0
        state.today_max_power_w = 0


def solar_shape_for_hour(local_hour: float,
                         sunrise_h: float,
                         sunset_h: float) -> float:
    if local_hour <= sunrise_h or local_hour >= sunset_h:
        return 0.0

    daylight = (local_hour - sunrise_h) / (sunset_h - sunrise_h)
    shape = math.sin(math.pi * daylight)
    return max(0.0, shape) ** 1.18


def cloud_multiplier(now_ts: float, base_cloud_factor: float) -> float:
    ripple = (
        0.04 * math.sin(now_ts / 37.0) +
        0.02 * math.sin(now_ts / 11.0) +
        0.01 * math.sin(now_ts / 5.0)
    )
    factor = base_cloud_factor + ripple
    return clamp(factor, 0.25, 1.02)


def ambient_temp_for_hour(local_hour: float,
                          sunrise_h: float,
                          sunset_h: float,
                          ambient_night_c: float,
                          ambient_day_c: float) -> float:
    shape = solar_shape_for_hour(local_hour, sunrise_h, sunset_h)
    return lerp(ambient_night_c, ambient_day_c, smoothstep01(shape))


def load_power_for_hour(local_hour: float,
                        sunrise_h: float,
                        sunset_h: float,
                        day_w: float,
                        night_w: float) -> float:
    shape = solar_shape_for_hour(local_hour, sunrise_h, sunset_h)
    return lerp(night_w, day_w, smoothstep01(shape))


def estimate_batt_v_from_soc(soc: float, cs: int,
                             batt_a: float,
                             nominal_v: float) -> float:
    """
    Simple, believable approximation for a 48V-class system.
    """
    base_v = lerp(46.0, 54.2, soc)

    if cs == 3:
        base_v += 2.0
    elif cs == 4:
        base_v += 3.2
    elif cs == 5:
        base_v += 1.2

    if batt_a < 0:
        base_v += batt_a * 0.08
    else:
        base_v += batt_a * 0.03

    return clamp(base_v, nominal_v - 8.0, nominal_v + 10.0)


def choose_charge_state(soc: float, pv_w: float,
                        net_w: float) -> tuple[int, int]:
    if pv_w <= 5:
        if net_w < -10:
            return 1, 0
        return 0, 0

    if soc < 0.85:
        return 3, 2
    if soc < 0.98:
        return 4, 2
    return 5, 1


def build_realtime_frame(config: PlantConfig,
                         state: PlantState) -> bytes:
    fields = [
        ("PID", PRODUCT_ID),
        ("FW", FIRMWARE),
        ("SER#", SERIAL),
        ("V", str(int(round(state.batt_v * 1000)))),
        ("I", str(int(round(state.batt_a * 1000)))),
        ("VPV", str(int(round(state.pv_mv)))),
        ("PPV", str(int(round(state.pv_w)))),
        ("CS", str(state.cs)),
        ("MPPT", str(state.mppt)),
        ("ERR", str(state.err)),
        ("T", str(int(round(state.temp_c)))),
        ("LOAD", "ON" if config.load_base_output_enabled else "OFF"),
        ("IL", str(int(round(state.load_a * 1000)))),
        ("Relay", "OFF"),
        ("H19", str(int(round(state.total_yield_kwh * 100)))),
        ("H20", str(int(round(state.today_yield_kwh * 100)))),
        ("H21", str(int(round(state.today_max_power_w)))),
        ("H22", str(int(round(state.yesterday_yield_kwh * 100)))),
        ("H23", str(int(round(state.yesterday_max_power_w)))),
    ]
    return build_frame(fields)


def realtime_tick(config: PlantConfig,
                  state: PlantState,
                  now_ts: float) -> tuple[bytes, dict[str, object]]:
    local_dt = now_at_offset(config.site.tz_offset_h)
    update_site_sun_times_if_needed(config, state, local_dt)

    if state.last_update_ts <= 0:
        dt_h = INTERVAL_S / 3600.0
    else:
        dt_h = clamp(
            (now_ts - state.last_update_ts) / 3600.0,
            0.0,
            10.0 / 3600.0,
        )

    local_hour = hour_from_dt(local_dt)
    sunrise_h = state.sunrise_h
    sunset_h = state.sunset_h

    sun_shape = solar_shape_for_hour(local_hour, sunrise_h, sunset_h)
    cloud = cloud_multiplier(now_ts, config.cloud_factor)
    ambient = ambient_temp_for_hour(
        local_hour,
        sunrise_h,
        sunset_h,
        config.ambient_night_c,
        config.ambient_day_c,
    )

    max_pv = pv_peak_w(config)
    pv_w = max_pv * sun_shape * cloud

    temp_derate = 1.0 - max(0.0, ambient - 25.0) * 0.0035
    temp_derate = clamp(temp_derate, 0.82, 1.02)
    pv_w *= temp_derate

    if pv_w < 2.0:
        pv_w = 0.0

    load_w = load_power_for_hour(
        local_hour,
        sunrise_h,
        sunset_h,
        config.load_day_w,
        config.load_night_w,
    )

    batt_v_guess = state.batt_v if state.batt_v > 1.0 else \
        config.battery_nominal_v
    net_w = pv_w - load_w
    batt_a = net_w / batt_v_guess
    batt_a = clamp(batt_a, -20.0, 30.0)

    state.soc = clamp(
        state.soc + (batt_a * dt_h) / config.battery_capacity_ah,
        0.05,
        1.0,
    )

    cs, mppt = choose_charge_state(state.soc, pv_w, net_w)
    batt_v = estimate_batt_v_from_soc(
        state.soc, cs, batt_a, config.battery_nominal_v
    )

    batt_a = clamp(net_w / max(batt_v, 1.0), -20.0, 30.0)

    controller_temp = ambient + 4.0 * sun_shape + 0.007 * pv_w
    controller_temp = clamp(
        controller_temp,
        config.ambient_night_c - 2.0,
        72.0,
    )

    err = "17" if controller_temp >= 67.0 else "0"
    if err != "0":
        cs, mppt = 2, 0

    pv_mv = 0.0
    if pv_w > 0:
        pv_mv = lerp(90000.0, 145000.0, smoothstep01(sun_shape))

    load_a = load_w / max(batt_v, 1.0)

    pv_kwh = (pv_w * dt_h) / 1000.0
    state.today_yield_kwh += max(0.0, pv_kwh)
    state.total_yield_kwh += max(0.0, pv_kwh)
    state.today_max_power_w = max(state.today_max_power_w,
                                  int(round(pv_w)))

    state.batt_v = batt_v
    state.batt_a = batt_a
    state.pv_w = int(round(pv_w))
    state.pv_mv = int(round(pv_mv))
    state.temp_c = controller_temp
    state.load_a = load_a
    state.cs = cs
    state.mppt = mppt
    state.err = err
    state.last_update_ts = now_ts

    summary = {
        "site_name": config.site.name,
        "local_time": local_dt.strftime("%H:%M:%S"),
        "batt_v": state.batt_v,
        "pv_w": state.pv_w,
        "temp_c": state.temp_c,
        "state": charge_state_label(state.cs),
        "yield_today_kwh": state.today_yield_kwh,
        "err": state.err,
    }

    return build_realtime_frame(config, state), summary


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
        description="VE.Direct frame emulator for controller.d"
    )

    parser.add_argument(
        "--mode",
        choices=["scenario", "realtime"],
        default="scenario",
        help="scenario = compressed synthetic mode, realtime = "
             "wall-clock plant emulation",
    )
    parser.add_argument(
        "--scenario",
        choices=list(SCENARIOS.keys()) + ["all"],
        default="sunny",
        help="Scenario to simulate in scenario mode",
    )
    parser.add_argument(
        "--interval",
        type=float,
        default=INTERVAL_S,
        help="Seconds between frames",
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
        soc_init=args.soc_init,
        load_day_w=args.load_day_w,
        load_night_w=args.load_night_w,
        cloud_factor=args.cloud_factor,
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
