/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Reusable tile map (Leaflet) — used by Home, Site detail and onboarding.
 * Street / Satellite / Terrain base layers (free, no API key), colored status
 * pins, and optional per-marker popup content. Render via the dynamic
 * `UkamaMap` wrapper (ssr:false) so Leaflet never runs on the server.
 *
 * The base layer comes from the shared `useUiPrefs().mapView` preference, and
 * switching layers in-map updates it, so every console map shows one view.
 */
import 'leaflet/dist/leaflet.css';

import { useEffect, useRef } from 'react';
import L from 'leaflet';
import { useColorScheme } from '@mui/material/styles';
import {
  LayersControl,
  MapContainer,
  Marker,
  Popup,
  TileLayer,
  useMap,
  useMapEvent,
} from 'react-leaflet';
import { useUiPrefs } from '@/lib/store';
import { MAP_VIEWS, basemapUrl, viewFromLabel, type MapView } from './basemaps';

export interface UkamaMapMarker {
  id: string;
  lat: number;
  lng: number;
  /** Pin colour (CSS color). Defaults to the accent. */
  color?: string;
  /** Popup content shown on marker click. */
  popup?: React.ReactNode;
}

export interface UkamaMapProps {
  markers?: UkamaMapMarker[];
  /** Fallback view when there are no markers. */
  center?: [number, number];
  zoom?: number;
  height?: number | string;
  /** Fit the viewport to all markers (default true when >1 marker). */
  fitToMarkers?: boolean;
  /** Disable pan/zoom interactions (e.g. small onboarding preview). */
  interactive?: boolean;
  /** Zoom on mouse-wheel scroll. Off by default so page scroll isn't trapped. */
  scrollWheelZoom?: boolean;
  /** Show the in-map base layer switcher (off for small previews). */
  layersControl?: boolean;
  /** Marker click handler (e.g. select a site on the home map). */
  onSelect?: (id: string) => void;
}

/** Site marker — Material CellTower icon tinted by status. */
const siteIcon = (color: string) =>
  L.divIcon({
    className: 'uk-map-pin',
    html: `<svg viewBox="0 0 24 24" width="30" height="30" fill="${color}"
      style="filter:drop-shadow(0 1px 2px rgba(0,0,0,.45));">
      <path d="m7.3 14.7 1.2-1.2c-1-1-1.5-2.3-1.5-3.5 0-1.3.5-2.6 1.5-3.5L7.3 5.3c-1.3 1.3-2 3-2 4.7s.7 3.4 2 4.7M19.1 2.9l-1.2 1.2c1.6 1.6 2.4 3.8 2.4 5.9s-.8 4.3-2.4 5.9l1.2 1.2c2-2 2.9-4.5 2.9-7.1s-1-5.1-2.9-7.1"/>
      <path d="M6.1 4.1 4.9 2.9C3 4.9 2 7.4 2 10s1 5.1 2.9 7.1l1.2-1.2c-1.6-1.6-2.4-3.8-2.4-5.9s.8-4.3 2.4-5.9m10.6 10.6c1.3-1.3 2-3 2-4.7-.1-1.7-.7-3.4-2-4.7l-1.2 1.2c1 1 1.5 2.3 1.5 3.5 0 1.3-.5 2.6-1.5 3.5zM14.5 10c0-1.38-1.12-2.5-2.5-2.5S9.5 8.62 9.5 10c0 .76.34 1.42.87 1.88L7 22h2l.67-2h4.67l.66 2h2l-3.37-10.12c.53-.46.87-1.12.87-1.88m-4.17 8L12 13l1.67 5z"/>
    </svg>`,
    iconSize: [30, 30],
    iconAnchor: [15, 15],
    popupAnchor: [0, -15],
  });

/** A whole-world view used when there's nothing to show yet (no sites). */
const WORLD_CENTER: [number, number] = [25, 0];
const WORLD_ZOOM = 2;

/** Keep the view in sync with markers/center, and re-measure the container
 *  (Leaflet paints partial tiles when its size isn't final at init). */
function ViewSync({
  markers,
  center,
  zoom,
  fitToMarkers,
}: {
  markers: UkamaMapMarker[];
  center?: [number, number];
  zoom: number;
  fitToMarkers: boolean;
}) {
  const map = useMap();
  useEffect(() => {
    if (fitToMarkers && markers.length > 1) {
      map.fitBounds(
        markers.map((m) => [m.lat, m.lng] as [number, number]),
        { padding: [40, 40] },
      );
    } else if (markers.length === 1) {
      map.setView([markers[0]!.lat, markers[0]!.lng], zoom);
    } else if (center) {
      map.setView(center, zoom);
    } else {
      // No sites and no explicit center → show the zoomed-out world.
      map.setView(WORLD_CENTER, WORLD_ZOOM);
    }
  }, [map, markers, center, zoom, fitToMarkers]);

  useEffect(() => {
    const invalidate = () => map.invalidateSize();
    const timers = [0, 150, 400].map((d) => setTimeout(invalidate, d));
    const ro = new ResizeObserver(invalidate);
    ro.observe(map.getContainer());
    return () => {
      timers.forEach(clearTimeout);
      ro.disconnect();
    };
  }, [map]);
  return null;
}

/** Base layers bound to the shared `mapView` preference. The switcher writes
 *  the choice back to the store; a change from Settings swaps the active tile
 *  layer in place, with no remount. */
function BaseLayers({
  view,
  dark,
  withControl,
}: {
  view: MapView;
  dark: boolean;
  withControl: boolean;
}) {
  const map = useMap();
  const setMapView = useUiPrefs((s) => s.setMapView);
  const layers = useRef<Partial<Record<MapView, L.TileLayer>>>({});

  useMapEvent('baselayerchange', (e) => {
    const next = viewFromLabel(e.name);
    if (next) setMapView(next);
  });

  useEffect(() => {
    const target = layers.current[view];
    if (!target || map.hasLayer(target)) return;
    MAP_VIEWS.forEach((v) => {
      const layer = layers.current[v.id];
      if (layer && layer !== target && map.hasLayer(layer)) {
        map.removeLayer(layer);
      }
    });
    map.addLayer(target);
  }, [map, view]);

  // Previews without the switcher only ever mount the selected layer.
  if (!withControl) return <TileLayer url={basemapUrl(view, dark)} />;

  return (
    <LayersControl position="topright">
      {MAP_VIEWS.map((v) => (
        <LayersControl.BaseLayer
          key={v.id}
          name={v.label}
          checked={v.id === view}
        >
          <TileLayer
            url={basemapUrl(v.id, dark)}
            ref={(layer) => {
              layers.current[v.id] = layer ?? undefined;
            }}
          />
        </LayersControl.BaseLayer>
      ))}
    </LayersControl>
  );
}

const DEFAULT_CENTER: [number, number] = [0, 20];

export default function UkamaMapImpl({
  markers = [],
  center,
  zoom = 6,
  height = 300,
  fitToMarkers = true,
  interactive = true,
  scrollWheelZoom = false,
  layersControl = true,
  onSelect,
}: UkamaMapProps) {
  const { mode, systemMode } = useColorScheme();
  const dark = (mode === 'system' ? systemMode : mode) === 'dark';
  const mapView = useUiPrefs((s) => s.mapView);

  // With no sites and no explicit center, open on the whole-world view so the
  // first paint isn't an empty ocean at a tight zoom.
  const emptyWorld = markers.length === 0 && !center;
  const start: [number, number] = emptyWorld
    ? WORLD_CENTER
    : (center ??
      (markers[0] ? [markers[0].lat, markers[0].lng] : DEFAULT_CENTER));
  const startZoom = emptyWorld ? WORLD_ZOOM : zoom;

  return (
    <MapContainer
      center={start}
      zoom={startZoom}
      style={{ height, width: '100%' }}
      scrollWheelZoom={interactive && scrollWheelZoom}
      dragging={interactive}
      zoomControl={interactive}
      doubleClickZoom={interactive}
      attributionControl={false}
    >
      <BaseLayers view={mapView} dark={dark} withControl={layersControl} />

      {markers.map((m) => (
        <Marker
          key={m.id}
          position={[m.lat, m.lng]}
          icon={siteIcon(m.color ?? '#2190f6')}
          eventHandlers={onSelect ? { click: () => onSelect(m.id) } : undefined}
        >
          {m.popup && <Popup>{m.popup}</Popup>}
        </Marker>
      ))}

      <ViewSync
        markers={markers}
        center={center}
        zoom={zoom}
        fitToMarkers={fitToMarkers}
      />
    </MapContainer>
  );
}
