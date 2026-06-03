/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Billing — what you owe Ukama for running this network (screens-manage.jsx). */
import Button from '@mui/material/Button';
import BoltRounded from '@mui/icons-material/BoltRounded';
import CreditCardRounded from '@mui/icons-material/CreditCardRounded';
import DownloadRounded from '@mui/icons-material/DownloadRounded';
import ReceiptLongRounded from '@mui/icons-material/ReceiptLongRounded';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import TableFooter from '@/components/data-table/TableFooter';
import PageHeader from '@/components/PageHeader';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import { BILLING } from '@/data';
import { useFirstLoad } from '@/lib/useFirstLoad';

export default function BillingScreen({ embed }: { embed?: boolean }) {
  const b = BILLING;
  const invLoading = useFirstLoad('billing-invoices');
  const toast = useToast();

  const body = (
    <>
      <div
        className="tile-grid"
        style={{ gridTemplateColumns: '1fr 1.3fr', marginBottom: 'var(--uk-gap)', alignItems: 'stretch' }}
      >
        <div className="card card-pad" style={{ display: 'flex', flexDirection: 'column' }}>
          <div style={{ fontSize: 13, color: 'var(--uk-ink-2)' }}>
            Current balance · {b.cycle}
          </div>
          <div
            className="tnum"
            style={{
              fontFamily: 'var(--font-display)',
              fontSize: 38,
              fontWeight: 500,
              margin: '6px 0 2px',
            }}
          >
            ${b.current.toLocaleString(undefined, { minimumFractionDigits: 2 })}
          </div>
          <div style={{ fontSize: 13, color: 'var(--uk-ink-2)' }}>
            Due {b.due} · {b.method}
          </div>
          <hr className="divider" style={{ margin: '16px 0' }} />
          <div style={{ display: 'grid', gap: 10 }}>
            {b.breakdown.map((r) => (
              <div
                key={r.label}
                style={{ display: 'flex', justifyContent: 'space-between', fontSize: 13.5 }}
              >
                <span style={{ color: 'var(--uk-ink-2)' }}>{r.label}</span>
                <span className="tnum" style={{ fontWeight: 600 }}>
                  ${r.amt.toFixed(2)}
                </span>
              </div>
            ))}
          </div>
          <div style={{ marginTop: 'auto', paddingTop: 16 }}>
            <Button
              fullWidth
              variant="contained"
              startIcon={<BoltRounded />}
              onClick={() => toast('Payment processed — thank you!')}
            >
              Pay now
            </Button>
          </div>
        </div>

        <div className="card card-pad" style={{ display: 'flex', flexDirection: 'column' }}>
          <div className="sec-title" style={{ marginBottom: 14 }}>
            Billing details
          </div>
          <div style={{ display: 'grid', gap: 14 }}>
            {(
              [
                ['Plan', b.plan],
                ['Payment method', b.method],
                ['Billing cycle', b.cycle],
                ['Next invoice due', b.due],
              ] as const
            ).map(([k, v]) => (
              <div key={k}>
                <div style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>{k}</div>
                <div style={{ fontWeight: 600, fontSize: 13.5, marginTop: 2 }}>{v}</div>
              </div>
            ))}
          </div>
          <hr className="divider" style={{ margin: '16px 0' }} />
          <div style={{ fontSize: 12.5, color: 'var(--uk-ink-3)', textWrap: 'pretty' }}>
            Usage-based billing. You’re charged monthly for active SIMs and data carried
            across your network.
          </div>
          <div style={{ marginTop: 'auto', paddingTop: 16, display: 'flex', gap: 10 }}>
            <Button
              fullWidth
              variant="outlined"
              startIcon={<CreditCardRounded />}
              onClick={() => toast('Update payment method')}
            >
              Update payment
            </Button>
          </div>
        </div>
      </div>

      <div className="card card-pad">
        <div className="sec-head">
          <div className="sec-title">Invoice history</div>
        </div>
        <div className="tbl-wrap">
          {invLoading ? (
            <SkeletonTable cols={5} rows={4} />
          ) : (
            <table className="tbl">
              <thead>
                <tr className="static">
                  <th>Invoice</th>
                  <th>Period</th>
                  <th className="num">Amount</th>
                  <th>Status</th>
                  <th style={{ width: 60 }} />
                </tr>
              </thead>
              <tbody>
                {b.invoices.map((inv) => (
                  <tr key={inv.id} className="static">
                    <td className="tnum" style={{ fontWeight: 600 }}>
                      {inv.id}
                    </td>
                    <td>{inv.period}</td>
                    <td className="num tnum">${inv.amt.toFixed(2)}</td>
                    <td>
                      <StatusBadge status="paid" variant="pill" />
                    </td>
                    <td>
                      <Button
                        size="small"
                        startIcon={<DownloadRounded />}
                        sx={{ color: 'var(--uk-ink-3)' }}
                        onClick={() => toast(`${inv.id}.pdf downloaded`)}
                      >
                        PDF
                      </Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
        {!invLoading && <TableFooter count={b.invoices.length} noun="invoices" />}
      </div>
    </>
  );

  if (embed) return body;

  return (
    <div className="page">
      <PageHeader
        title="Billing"
        sub="What you owe Ukama for running this network."
        actions={
          <>
            <Button
              variant="outlined"
              startIcon={<CreditCardRounded />}
              onClick={() => toast('Update payment method')}
            >
              Payment method
            </Button>
            <Button
              variant="contained"
              startIcon={<ReceiptLongRounded />}
              onClick={() => toast('Invoices downloaded')}
            >
              Download invoices
            </Button>
          </>
        }
      />
      {body}
    </div>
  );
}
