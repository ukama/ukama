/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * The three steps of site creation (see useCreateSite), shown while the run
 * is in flight. Same list styling as the readiness checklist so the flow
 * looks continuous from one step to the next.
 */
import CircularProgress from '@mui/material/CircularProgress';
import CheckCircleRounded from '@mui/icons-material/CheckCircleRounded';
import ErrorRounded from '@mui/icons-material/ErrorRounded';
import RadioButtonUncheckedRounded from '@mui/icons-material/RadioButtonUncheckedRounded';

import type { CreateStep } from './useCreateSite';

const icon = (state: CreateStep['state']) => {
  if (state === 'done')
    return (
      <CheckCircleRounded
        sx={{ fontSize: 22, color: 'var(--uk-success-bright)' }}
      />
    );
  if (state === 'failed')
    return <ErrorRounded sx={{ fontSize: 22, color: 'var(--uk-error)' }} />;
  if (state === 'active') return <CircularProgress size={18} />;
  return (
    <RadioButtonUncheckedRounded
      sx={{ fontSize: 22, color: 'var(--uk-ink-3)' }}
    />
  );
};

export default function CreateSiteProgress({ steps }: { steps: CreateStep[] }) {
  return (
    <ul className="cfg-steps-list">
      {steps.map((step) => (
        <li
          key={step.key}
          className="cfg-step-item"
          data-state={step.state === 'failed' ? 'active' : step.state}
        >
          <span className="cfg-step-ic">{icon(step.state)}</span>
          <span className="cfg-step-text">
            <span className="cfg-step-title">{step.title}</span>
            {step.detail && (
              <span className="cfg-step-hint">{step.detail}</span>
            )}
          </span>
        </li>
      ))}
    </ul>
  );
}
