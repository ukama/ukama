/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Minimal local typings for react-simple-maps v3 (only the API surface we
 * use). A local shim instead of @types/react-simple-maps, which pins
 * @types/react 18 and conflicts with React 19 (BUILD-PLAN §7.1 caveat).
 */
declare module 'react-simple-maps' {
  import type * as React from 'react';

  export interface ProjectionConfig {
    center?: [number, number];
    scale?: number;
    rotate?: [number, number, number];
  }

  export interface ComposableMapProps {
    projection?: string;
    projectionConfig?: ProjectionConfig;
    width?: number;
    height?: number;
    style?: React.CSSProperties;
    children?: React.ReactNode;
  }
  export const ComposableMap: React.FC<ComposableMapProps>;

  export interface GeographyShape {
    rsmKey: string;
  }
  export interface GeographiesProps {
    geography: unknown;
    children: (args: { geographies: GeographyShape[] }) => React.ReactNode;
  }
  export const Geographies: React.FC<GeographiesProps>;

  export interface GeographyProps {
    geography: GeographyShape;
    style?: {
      default?: React.CSSProperties;
      hover?: React.CSSProperties;
      pressed?: React.CSSProperties;
    };
    tabIndex?: number;
  }
  export const Geography: React.FC<GeographyProps>;

  export interface MarkerProps {
    coordinates: [number, number];
    onClick?: (event: React.MouseEvent) => void;
    style?: {
      default?: React.CSSProperties;
      hover?: React.CSSProperties;
      pressed?: React.CSSProperties;
    };
    children?: React.ReactNode;
  }
  export const Marker: React.FC<MarkerProps>;

  export interface ZoomableGroupProps {
    center?: [number, number];
    zoom?: number;
    minZoom?: number;
    maxZoom?: number;
    onMoveEnd?: (position: { coordinates: [number, number]; zoom: number }) => void;
    /** Return false to ignore a d3-zoom source event (e.g. wheel). */
    filterZoomEvent?: (event: { type?: string }) => boolean;
    children?: React.ReactNode;
  }
  export const ZoomableGroup: React.FC<ZoomableGroupProps>;
}
