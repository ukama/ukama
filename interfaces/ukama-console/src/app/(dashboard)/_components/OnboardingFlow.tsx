/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';
import InputAdornment from '@mui/material/InputAdornment';
import TextField from '@mui/material/TextField';

/**
 * First-run network setup — progressive-disclosure stepper:
 * welcome → network → first node → SIMs → review → done (onboarding.jsx).
 */
import { useState } from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import CheckRounded from '@mui/icons-material/CheckRounded';
import ChevronLeftRounded from '@mui/icons-material/ChevronLeftRounded';
import ChevronRightRounded from '@mui/icons-material/ChevronRightRounded';
import HubRounded from '@mui/icons-material/HubRounded';
import LocationOnRounded from '@mui/icons-material/LocationOnRounded';
import QrCodeScannerRounded from '@mui/icons-material/QrCodeScannerRounded';
import RouterRounded from '@mui/icons-material/RouterRounded';
import SettingsInputAntennaRounded from '@mui/icons-material/SettingsInputAntennaRounded';
import SimCardRounded from '@mui/icons-material/SimCardRounded';
import UMark from '@/components/UMark';
import { useToast } from '@/components/ToastProvider';
import { PLANS } from '@/data';
import { Ic } from './icons';

const OB_STEPS = [
  { key: 'network', label: 'Network' },
  { key: 'node', label: 'First node' },
  { key: 'sims', label: 'SIMs' },
  { key: 'review', label: 'Review' },
];
const OB_COUNTRIES = ['Zambia', 'DR Congo', 'Tanzania', 'Malawi', 'Kenya'];
const OB_NODE_TYPES = ['Tower node', 'Amplifier node', 'Indoor node'];

type Stage = 'welcome' | 0 | 1 | 2 | 3 | 'done';

function ObField({
  label,
  hint,
  children,
}: {
  label: string;
  hint?: string;
  children: React.ReactNode;
}) {
  return (
    <div>
      <label className="flabel">{label}</label>
      {children}
      {hint && <div style={{ fontSize: 12, color: 'var(--uk-ink-3)', marginTop: 5 }}>{hint}</div>}
    </div>
  );
}

export default function OnboardingFlow({ onClose }: { onClose: () => void }) {
  const toast = useToast();
  const [stage, setStage] = useState<Stage>('welcome');
  const [net, setNet] = useState({ name: '', country: 'Zambia' });
  const [node, setNode] = useState({ serial: '', type: 'Tower node' });
  const [sims, setSims] = useState({ count: '500', plan: PLANS[0]?.name ?? 'Starter' });

  const isNum = typeof stage === 'number';
  const canContinue =
    stage === 0
      ? net.name.trim().length > 1
      : stage === 1
        ? node.serial.trim().length > 2
        : stage === 2
          ? Number(sims.count) > 0
          : true;

  const next = () => {
    if (stage === 'welcome') return setStage(0);
    if (stage === 3) return setStage('done');
    if (isNum) setStage((stage + 1) as Stage);
  };
  const back = () => {
    if (stage === 0) return setStage('welcome');
    if (isNum) setStage((stage - 1) as Stage);
  };

  return (
    <Dialog
      open
      onClose={onClose}
      slotProps={{
        paper: { sx: { width: 620, maxWidth: '94vw', borderRadius: 3.5, overflow: 'hidden' } },
      }}
    >
      <div className="ob-grad" />

      {stage === 'welcome' && (
        <div className="ob-welcome">
          <div className="ob-mark">
            <UMark />
          </div>
          <div className="ob-title" style={{ fontSize: 25 }}>
            Let’s set up your network
          </div>
          <div className="ob-sub" style={{ maxWidth: 420, margin: '8px auto 0' }}>
            A few quick steps to stand up your private cellular network. You can change any
            of this later.
          </div>
          <div className="ob-checklist">
            {(
              [
                ['hub', 'Name your network and choose a region'],
                ['router', 'Register your first node'],
                ['sim_card', 'Allocate SIMs so subscribers can connect'],
              ] as const
            ).map(([icon, text]) => (
              <div key={text} className="ob-check">
                <span className="ob-check-ic">
                  <Ic name={icon} sx={{ fontSize: 19 }} />
                </span>
                <span>{text}</span>
              </div>
            ))}
          </div>
          <div style={{ display: 'flex', gap: 10, justifyContent: 'center', marginTop: 26 }}>
            <Button color="inherit" sx={{ color: 'var(--uk-ink-3)' }} onClick={onClose}>
              Maybe later
            </Button>
            <Button variant="contained" endIcon={<ChevronRightRounded />} onClick={next}>
              Get started
            </Button>
          </div>
        </div>
      )}

      {isNum && (
        <>
          <div className="ob-head">
            <div className="ob-title">Set up your network</div>
          </div>
          <div className="ob-steps">
            {OB_STEPS.map((s, i) => (
              <div
                key={s.key}
                className={`ob-step${i < stage ? ' done' : ''}${i === stage ? ' active' : ''}`}
              >
                <span className="ob-dot">
                  {i < stage ? <CheckRounded sx={{ fontSize: 17 }} /> : i + 1}
                </span>
                <span className="ob-step-label">{s.label}</span>
              </div>
            ))}
          </div>

          <div className="ob-body">
            {stage === 0 && (
              <div className="ob-field-grid">
                <ObField label="Network name" hint="The name operators and subscribers will see.">
                  <TextField fullWidth autoFocus value={net.name} onChange={(e) => setNet({ ...net, name: e.target.value })} placeholder="e.g. Kafue Valley Mobile" slotProps={{ input: { startAdornment: (<InputAdornment position="start"><HubRounded sx={{ fontSize: 19, color: 'var(--uk-ink-3)' }} /></InputAdornment>) } }} />
                </ObField>
                <ObField label="Country / region">
                  <TextField fullWidth select value={net.country} onChange={(e) => setNet({ ...net, country: e.target.value })} slotProps={{ select: { native: true }, input: { startAdornment: (<InputAdornment position="start"><LocationOnRounded sx={{ fontSize: 19, color: 'var(--uk-ink-3)' }} /></InputAdornment>) } }}>{OB_COUNTRIES.map((c) => (<option key={c}>{c}</option>))}</TextField>
                </ObField>
              </div>
            )}
            {stage === 1 && (
              <div className="ob-field-grid">
                <ObField label="Node serial" hint="Find this on the label or scan the QR code on the unit.">
                  <TextField fullWidth autoFocus value={node.serial} onChange={(e) => setNode({ ...node, serial: e.target.value })} placeholder="uk-tnode-a06-0000" slotProps={{ input: { startAdornment: (<InputAdornment position="start"><QrCodeScannerRounded sx={{ fontSize: 19, color: 'var(--uk-ink-3)' }} /></InputAdornment>) } }} />
                </ObField>
                <ObField label="Node type">
                  <div className="ob-segrow">
                    {OB_NODE_TYPES.map((tp) => {
                      const SegIcon = tp.startsWith('Amp')
                        ? SettingsInputAntennaRounded
                        : RouterRounded;
                      return (
                        <button
                          key={tp}
                          type="button"
                          className={`ob-seg${node.type === tp ? ' on' : ''}`}
                          onClick={() => setNode({ ...node, type: tp })}
                        >
                          <SegIcon sx={{ fontSize: 18 }} />
                          {tp}
                        </button>
                      );
                    })}
                  </div>
                </ObField>
              </div>
            )}
            {stage === 2 && (
              <div className="ob-field-grid">
                <ObField
                  label="SIMs to allocate"
                  hint="You can upload more SIM batches at any time from the SIM pool."
                >
                  <TextField fullWidth autoFocus type="number" value={sims.count} onChange={(e) => setSims({ ...sims, count: e.target.value })} placeholder="500" slotProps={{ input: { startAdornment: (<InputAdornment position="start"><SimCardRounded sx={{ fontSize: 19, color: 'var(--uk-ink-3)' }} /></InputAdornment>) } }} />
                </ObField>
                <ObField label="Default data plan">
                  <div className="ob-segrow ob-segrow-wrap">
                    {PLANS.map((p) => (
                      <button
                        key={p.id}
                        type="button"
                        className={`ob-seg${sims.plan === p.name ? ' on' : ''}`}
                        onClick={() => setSims({ ...sims, plan: p.name })}
                      >
                        <span className="ob-plan-dot" style={{ background: p.color }} />
                        {p.name}
                        <span className="ob-plan-price">${p.price}</span>
                      </button>
                    ))}
                  </div>
                </ObField>
              </div>
            )}
            {stage === 3 && (
              <div>
                <div className="ob-sub" style={{ marginBottom: 14 }}>
                  Review your setup, then create the network.
                </div>
                {(
                  [
                    ['hub', 'Network', `${net.name || 'Untitled network'} · ${net.country}`, 0],
                    ['router', node.type, node.serial || '—', 1],
                    [
                      'sim_card',
                      'SIMs',
                      `${Number(sims.count).toLocaleString()} allocated · ${sims.plan} plan`,
                      2,
                    ],
                  ] as const
                ).map(([icon, k, v, target]) => (
                  <div key={k} className="ob-review-row">
                    <span className="ob-review-ic">
                      <Ic name={icon} sx={{ fontSize: 20 }} />
                    </span>
                    <div style={{ flex: 1 }}>
                      <div style={{ fontSize: 12.5, color: 'var(--uk-ink-3)' }}>{k}</div>
                      <div style={{ fontWeight: 600, fontSize: 14 }}>{v}</div>
                    </div>
                    <button type="button" className="link" onClick={() => setStage(target as Stage)}>
                      Edit
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>

          <div className="ob-foot">
            <Button color="inherit" sx={{ color: 'var(--uk-ink-3)' }} startIcon={<ChevronLeftRounded />} onClick={back}>
              Back
            </Button>
            <Button
              variant="contained"
              disabled={!canContinue}
              endIcon={stage === 3 ? <CheckRounded /> : <ChevronRightRounded />}
              onClick={next}
            >
              {stage === 3 ? 'Create network' : 'Continue'}
            </Button>
          </div>
        </>
      )}

      {stage === 'done' && (
        <div className="ob-welcome">
          <div className="ob-success">
            <CheckRounded sx={{ fontSize: 34 }} />
          </div>
          <div className="ob-title" style={{ fontSize: 24 }}>
            {net.name || 'Your network'} is ready
          </div>
          <div className="ob-sub" style={{ maxWidth: 420, margin: '8px auto 0' }}>
            Your network is live with {Number(sims.count).toLocaleString()} SIMs and one node
            registered. Subscribers can connect now.
          </div>
          <div style={{ display: 'flex', gap: 10, justifyContent: 'center', marginTop: 26 }}>
            <Button
              variant="contained"
              endIcon={<ChevronRightRounded />}
              onClick={() => {
                onClose();
                toast(`${net.name || 'New network'} created`);
              }}
            >
              Go to dashboard
            </Button>
          </div>
        </div>
      )}
    </Dialog>
  );
}
