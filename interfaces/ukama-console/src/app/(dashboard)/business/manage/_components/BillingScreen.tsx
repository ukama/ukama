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
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import BoltRounded from '@mui/icons-material/BoltRounded';
import CreditCardRounded from '@mui/icons-material/CreditCardRounded';
import DownloadRounded from '@mui/icons-material/DownloadRounded';
import ReceiptLongRounded from '@mui/icons-material/ReceiptLongRounded';
import { useBillingOverviewQuery } from '@/client/graphql/commerce.generated';
import { EmptyState } from '@/components/EmptyState';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import TableFooter from '@/components/data-table/TableFooter';
import PageHeader from '@/components/PageHeader';
import { sectionValue } from '@/components/SectionFallback';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import { useUiPrefs } from '@/lib/store';

export default function BillingScreen({ embed }: { embed?: boolean }) {
  const toast = useToast();
  const networkId = useUiPrefs((s) => s.networkId);
  const { data, loading: invLoading, refetch } = useBillingOverviewQuery({
    variables: { networkId },
  });
  const invoicesSection = data?.commerceView.invoices;
  const balanceSection = data?.commerceView.balance;
  const invoices = invoicesSection?.reports ?? [];

  const body = (
    <>
      <div
        className="tile-grid"
        style={{ gridTemplateColumns: '1fr 1.3fr', marginBottom: 'var(--uk-gap)', alignItems: 'stretch' }}
      >
        <div className="card card-pad" style={{ display: 'flex', flexDirection: 'column' }}>
          <div style={{ fontSize: 13, color: 'var(--uk-ink-2)' }}>
            Current balance
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
            {/* TODO(backend-gap #11): outstanding amount needs an invoice
                amount field on the reports list */}
            {sectionValue(null, balanceSection?.error ?? null)}
          </div>
          <div style={{ fontSize: 13, color: 'var(--uk-ink-2)' }}>
            {balanceSection?.error
              ? '—'
              : `${balanceSection?.outstandingCount ?? 0} unpaid invoice${
                  (balanceSection?.outstandingCount ?? 0) === 1 ? '' : 's'
                }${
                  balanceSection?.latestUnpaidPeriod
                    ? ` · latest ${balanceSection.latestUnpaidPeriod}`
                    : ''
                }`}
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
                ['Billing cycle', 'Monthly'],
                ['Unpaid invoices', sectionValue(balanceSection?.outstandingCount, balanceSection?.error)],
                ['Latest unpaid period', balanceSection?.latestUnpaidPeriod ?? '—'],
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
          ) : invoicesSection?.error ? (
            <EmptyState
              art="error"
              title="Couldn't load invoices"
              sub={invoicesSection.error.message}
              cta="Try again"
              onCta={() => refetch()}
            />
          ) : invoices.length === 0 ? (
            <EmptyState art="invoice" title="No invoices yet" sub="Invoices appear here each cycle." />
          ) : (
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Invoice</TableCell>
                  <TableCell>Period</TableCell>
                  <TableCell align="right">Amount</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell sx={{ width: 44 }} />
                </TableRow>
              </TableHead>
              <TableBody>
                {invoices.map((inv) => (
                  <TableRow key={inv.id}>
                    <TableCell className="tnum" style={{ fontWeight: 600 }}>
                      {inv.id.slice(0, 8)}
                    </TableCell>
                    <TableCell>{inv.period}</TableCell>
                    {/* TODO(backend-gap #11): invoice amount field */}
                    <TableCell align="right" className="tnum">—</TableCell>
                    <TableCell>
                      <StatusBadge status={inv.isPaid ? 'paid' : 'pending'} variant="pill" />
                    </TableCell>
                    <TableCell>
                      <Button
                        size="small"
                        startIcon={<DownloadRounded />}
                        sx={{ color: 'var(--uk-ink-3)' }}
                        onClick={() => toast(`${inv.id}.pdf downloaded`)}
                      >
                        PDF
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </div>
        {!invLoading && !invoicesSection?.error && (
          <TableFooter count={invoices.length} noun="invoices" />
        )}
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
