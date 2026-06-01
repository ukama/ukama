/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { MetricsStateRes } from '@/client/graphql/generated/subscriptions';
import { SITE_KPI_TYPES } from '@/constants';

// ---------------------------------------------------------------------------
// Low-level value extraction helpers
// ---------------------------------------------------------------------------

type MetricValuePayload = number | [unknown, number] | unknown[];

/**
 * Extracts a numeric metric value from a raw payload element.
 * Handles both scalar numbers and 2-tuple [timestamp, value] arrays.
 */
export const extractMetricValue = (value: unknown): number | null => {
  if (Array.isArray(value) && value.length > 1) {
    return typeof value[1] === 'number' ? value[1] : null;
  }
  return typeof value === 'number' ? value : null;
};

/**
 * Extracts a numeric metric value from a PubSub payload, which may itself be
 * a 2-tuple where the second element is the actual metric payload.
 */
export const extractMetricFromPubSubPayload = (
  payload: MetricValuePayload | [unknown, MetricValuePayload] | unknown,
): number | null => {
  if (payload === null || payload === undefined) return null;
  if (Array.isArray(payload) && payload.length > 1) {
    return extractMetricValue(payload[1]);
  }
  return extractMetricValue(payload);
};

// ---------------------------------------------------------------------------
// Domain helpers
// ---------------------------------------------------------------------------

/**
 * Returns the current value for a single metric type scoped to a specific
 * site, or null if no matching metric exists in the snapshot.
 */
export const getSiteMetricValue = (
  metricId: string,
  metricsData: MetricsStateRes | undefined,
  siteId: string | undefined,
): number | null => {
  if (!metricsData?.metrics || !siteId) return null;

  const metric = metricsData.metrics.find(
    (m) => m.type === metricId && m.success === true && m.siteId === siteId,
  );

  return metric ? extractMetricValue(metric.value) : null;
};

/**
 * Aggregates active-subscriber counts across all metric entries for a site.
 * Returns null when no subscriber metrics exist yet.
 */
export const getSiteActiveSubscribers = (
  metricsData: MetricsStateRes | undefined,
  siteId: string | undefined,
): number | null => {
  if (!metricsData?.metrics || !siteId) return null;

  const subscriberMetrics = metricsData.metrics.filter(
    (m) =>
      m.type === SITE_KPI_TYPES.ACTIVE_SUBSCRIBERS &&
      m.success === true &&
      m.siteId === siteId,
  );

  if (subscriberMetrics.length === 0) return null;

  return subscriberMetrics.reduce((total, metric) => {
    const value = extractMetricValue(metric.value);
    return total + (value ?? 0);
  }, 0);
};

// ---------------------------------------------------------------------------
// UI helpers
// ---------------------------------------------------------------------------

/** Truncates text to maxLength characters, appending "…" when cut. */
export const truncateText = (text: string, maxLength: number): string => {
  if (!text || typeof text !== 'string') return '';
  if (text.length <= maxLength) return text;
  return `${text.substring(0, maxLength)}...`;
};
