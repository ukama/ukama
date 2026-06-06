/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Guided "bring your site online" checklist. Four steps — power on each of
 * the three units, then confirm the location lock. Status is driven entirely
 * by the polled getNodes response (see computeSiteReadiness); the user just
 * follows along. The first incomplete step is marked active with a spinner.
 */
'use client';

import CircularProgress from '@mui/material/CircularProgress';
import IconButton from '@mui/material/IconButton';
import Tooltip from '@mui/material/Tooltip';
import CheckCircleRounded from '@mui/icons-material/CheckCircleRounded';
import RadioButtonUncheckedRounded from '@mui/icons-material/RadioButtonUncheckedRounded';
import RefreshRounded from '@mui/icons-material/RefreshRounded';

import type { SiteReadiness } from './detectSites';

const timeFmt = (d: Date): string =>
  d.toLocaleTimeString('en-GB', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });

interface Item {
  key: keyof Pick<
    SiteReadiness,
    'tower' | 'amplifier' | 'controller' | 'located'
  >;
  title: string;
  hint: string;
}

const ITEMS: Item[] = [
  {
    key: 'tower',
    title: 'Turn on your tower unit',
    hint: 'Power it on and connect it to the internet. This is the main unit of your site.',
  },
  {
    key: 'amplifier',
    title: 'Turn on your amplifier unit',
    hint: 'Power it on and connect it — it extends your coverage.',
  },
  {
    key: 'controller',
    title: 'Turn on your controller unit',
    hint: 'Power it on and connect it — it manages the site.',
  },
  {
    key: 'located',
    title: 'Confirming your site location',
    hint: 'Your tower reports its location automatically once it has a clear view of the sky.',
  },
];

export default function SiteReadinessChecklist({
  readiness,
  lastFetched,
  refreshing,
  onRefresh,
}: {
  readiness: SiteReadiness;
  lastFetched: Date | null;
  refreshing: boolean;
  onRefresh: () => void;
}) {
  const done = ITEMS.filter((i) => readiness[i.key]).length;
  // The active step is the first one not yet complete.
  const activeIndex = ITEMS.findIndex((i) => !readiness[i.key]);

  return (
    <div className="cfg-checklist">
      <div className="cfg-checklist-progress">
        <div className="cfg-checklist-bar">
          <span style={{ width: `${(done / ITEMS.length) * 100}%` }} />
        </div>
        <span className="cfg-checklist-count">{done}/{ITEMS.length} done</span>
      </div>

      <div className="cfg-checklist-sync">
        <span>
          {lastFetched
            ? `Last checked ${timeFmt(lastFetched)}`
            : 'Checking…'}
        </span>
        <Tooltip title="Check now">
          <span>
            <IconButton
              size="small"
              aria-label="Check now"
              onClick={onRefresh}
              disabled={refreshing}
              sx={{ color: 'var(--uk-ink-2)' }}
            >
              <RefreshRounded
                sx={{
                  fontSize: 18,
                  animation: refreshing
                    ? 'cfg-spin 0.8s linear infinite'
                    : 'none',
                }}
              />
            </IconButton>
          </span>
        </Tooltip>
      </div>

      <ul className="cfg-steps-list">
        {ITEMS.map((item, i) => {
          const complete = readiness[item.key];
          const active = !complete && i === activeIndex;
          return (
            <li
              key={item.key}
              className="cfg-step-item"
              data-state={complete ? 'done' : active ? 'active' : 'pending'}
            >
              <span className="cfg-step-ic">
                {complete ? (
                  <CheckCircleRounded
                    sx={{ fontSize: 22, color: 'var(--uk-success-bright)' }}
                  />
                ) : active ? (
                  <CircularProgress size={18} />
                ) : (
                  <RadioButtonUncheckedRounded
                    sx={{ fontSize: 22, color: 'var(--uk-ink-3)' }}
                  />
                )}
              </span>
              <span className="cfg-step-text">
                <span className="cfg-step-title">{item.title}</span>
                <span className="cfg-step-hint">
                  {complete
                    ? item.key === 'located'
                      ? 'Location confirmed.'
                      : 'Powered on and online.'
                    : active
                      ? item.key === 'located'
                        ? 'Waiting for a location fix…'
                        : 'Waiting for this unit to come online…'
                      : item.hint}
                </span>
              </span>
            </li>
          );
        })}
      </ul>
    </div>
  );
}
