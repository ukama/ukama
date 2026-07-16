/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';
import { useGetPaymentsQuery } from '@/client/graphql/billing.generated';
import AppModal from '@/components/AppModal';
import { useToast } from '@/components/ToastProvider';
import { useAuth } from '@/lib/auth/context';
import { useCurrency } from '@/lib/currency';
import { buildReceiptView, pickPaymentForSim } from '@/lib/receipt';
import { downloadReceiptPdf } from '@/lib/receiptPdf';
import CheckCircleRounded from '@mui/icons-material/CheckCircleRounded';
import DownloadRounded from '@mui/icons-material/DownloadRounded';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Skeleton from '@mui/material/Skeleton';
import { useMemo, useState } from 'react';

const muted = 'var(--uk-ink-3)';
const ink = 'var(--uk-ink)';

function Meta({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <div style={{ fontSize: 11, color: muted, marginBottom: 2 }}>{label}</div>
      <div style={{ fontSize: 13, color: ink }}>{value}</div>
    </div>
  );
}

export default function ReceiptDialog({
  simId,
  packageId,
  planName,
  onClose,
}: {
  simId: string;
  packageId: string;
  planName: string;
  onClose: () => void;
}) {
  const toast = useToast();
  const { symbol } = useCurrency();
  const orgName = useAuth()?.orgName ?? 'Ukama';
  const [saving, setSaving] = useState(false);

  const { data, loading } = useGetPaymentsQuery({
    variables: {
      data: {
        itemId: packageId,
        type: 'package',
        paymentMethod: 'cash',
        status: 'completed',
      },
    },
    fetchPolicy: 'cache-and-network',
  });

  const view = useMemo(() => {
    const payment = pickPaymentForSim(data?.getPayments.payments ?? [], simId);
    return payment ? buildReceiptView(payment, { planName, symbol, orgName }) : null;
  }, [data, simId, planName, symbol, orgName]);

  const download = async () => {
    if (!view) return;
    try {
      setSaving(true);
      await downloadReceiptPdf(view);
    } catch {
      toast('Could not generate the receipt PDF.');
    } finally {
      setSaving(false);
    }
  };

  return (
    <AppModal
      title="Payment receipt"
      width={460}
      onClose={onClose}
      footer={
        <>
          <Button color="inherit" sx={{ color: 'var(--uk-ink-3)' }} onClick={onClose}>
            Close
          </Button>
          <Button
            variant="contained"
            disabled={!view || saving}
            startIcon={
              saving ? <CircularProgress size={16} color="inherit" /> : <DownloadRounded />
            }
            onClick={download}
          >
            {saving ? 'Preparing…' : 'Download'}
          </Button>
        </>
      }
    >
      {loading && !view ? (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 10, padding: '4px 0 12px' }}>
          <Skeleton variant="rounded" height={44} />
          <Skeleton variant="rounded" height={60} />
          <Skeleton variant="rounded" height={40} />
        </div>
      ) : !view ? (
        <div style={{ fontSize: 13, color: muted, padding: '8px 0 16px' }}>
          No receipt found for this package. It may have been allocated without a
          recorded payment.
        </div>
      ) : (
        <div style={{ padding: '2px 0 8px' }}>
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              marginBottom: 14,
            }}
          >
            <span style={{ fontSize: 13, color: muted }}>{view.orgName}</span>
            <span
              style={{
                display: 'inline-flex',
                alignItems: 'center',
                gap: 5,
                fontSize: 12,
                fontWeight: 500,
                color: 'var(--uk-success)',
                background: 'color-mix(in srgb, var(--uk-success) 12%, transparent)',
                padding: '4px 10px',
                borderRadius: 20,
              }}
            >
              <CheckCircleRounded sx={{ fontSize: 15 }} />
              {view.status}
            </span>
          </div>

          <div
            style={{
              display: 'flex',
              gap: 22,
              paddingBottom: 14,
              borderBottom: '0.5px solid var(--uk-line)',
            }}
          >
            <Meta label="Receipt no" value={view.id.slice(0, 8)} />
            <Meta label="Paid on" value={view.paidOn} />
            <Meta label="Method" value={view.method} />
          </div>

          <div style={{ padding: '14px 0 6px' }}>
            <div style={{ fontSize: 11, color: muted, marginBottom: 4 }}>Billed to</div>
            <div style={{ fontSize: 13, color: ink }}>Walk-in customer</div>
            <div style={{ fontSize: 12, color: 'var(--uk-ink-2)', marginTop: 1 }}>
              SIM {simId.slice(0, 8)} · package top-up
            </div>
          </div>

          <div
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              fontSize: 11,
              color: muted,
              padding: '8px 0',
              borderBottom: '0.5px solid var(--uk-line)',
            }}
          >
            <span>Description</span>
            <span>Amount</span>
          </div>
          <div
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'flex-start',
              padding: '12px 0',
            }}
          >
            <div style={{ minWidth: 0 }}>
              <div style={{ fontSize: 13.5, color: ink }}>{view.planName}</div>
              <div style={{ fontSize: 12, color: 'var(--uk-ink-2)', marginTop: 2 }}>
                Data package · qty 1
              </div>
            </div>
            <div style={{ fontSize: 13.5, color: ink }}>{view.amountLabel}</div>
          </div>

          <div
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'baseline',
              padding: '10px 0 4px',
              borderTop: '0.5px solid var(--uk-line)',
            }}
          >
            <span style={{ fontSize: 14, fontWeight: 500, color: ink }}>Total paid</span>
            <span style={{ fontSize: 20, fontWeight: 500, color: ink }}>{view.amountLabel}</span>
          </div>

          <div
            style={{
              marginTop: 12,
              background: 'var(--uk-hover)',
              borderRadius: 10,
              padding: '10px 12px',
            }}
          >
            <div style={{ fontSize: 11, color: muted }}>Payment ID</div>
            <div
              style={{
                fontSize: 11.5,
                color: 'var(--uk-ink-2)',
                wordBreak: 'break-all',
                marginTop: 1,
              }}
            >
              {view.id}
            </div>
          </div>

          <div style={{ fontSize: 11.5, color: muted, marginTop: 12 }}>
            Auto-generated · not a tax invoice
          </div>
        </div>
      )}
    </AppModal>
  );
}
