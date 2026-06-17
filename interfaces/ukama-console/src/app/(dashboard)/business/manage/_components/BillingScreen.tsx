/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Billing — what you owe Ukama for running this network (screens-manage.jsx). */
import { useCallback, useMemo, useState } from 'react';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import BoltRounded from '@mui/icons-material/BoltRounded';
import DownloadRounded from '@mui/icons-material/DownloadRounded';
import { useBillingOverviewQuery } from '@/client/graphql/commerce.generated';
import {
  useGetPaymentsQuery,
  useUpdatePaymentMutation,
  useGetGeneratedPdfReportLazyQuery,
} from '@/client/graphql/billing.generated';
import { EmptyState } from '@/components/EmptyState';
import SkeletonTable from '@/components/data-table/SkeletonTable';
import TableFooter from '@/components/data-table/TableFooter';
import PageHeader from '@/components/PageHeader';
import { sectionValue } from '@/components/SectionFallback';
import StatusBadge from '@/components/StatusBadge';
import { useToast } from '@/components/ToastProvider';
import { useAuth } from '@/lib/auth/context';
import { useCurrency } from '@/lib/currency';
import { base64ToBlob, downloadBlob } from '@/lib/download';
import { useUiPrefs } from '@/lib/store';
import StripePaymentDialog from './StripePaymentDialog';

/** Invoice total comes from the provider in minor units (cents) as a string. */
const centsToMajor = (cents?: string | null): number | null => {
  const n = Number(cents);
  return Number.isFinite(n) ? n / 100 : null;
};

type Invoice = NonNullable<
  NonNullable<
    ReturnType<typeof useBillingOverviewQuery>['data']
  >['commerceView']['invoices']
>['reports'] extends (infer R)[] | null | undefined
  ? R
  : never;

export default function BillingScreen({ embed }: { embed?: boolean }) {
  const toast = useToast();
  const { money } = useCurrency();
  const user = useAuth();
  const networkId = useUiPrefs((s) => s.networkId);

  const {
    data,
    loading: invLoading,
    refetch: refetchOverview,
  } = useBillingOverviewQuery({ variables: { networkId }, skip: !networkId });
  const invoicesSection = data?.commerceView.invoices;
  const balanceSection = data?.commerceView.balance;
  const invoices = useMemo(() => invoicesSection?.reports ?? [], [invoicesSection]);

  // Phase B — payment records, matched to invoices by itemId. Direct query
  // (mirrors the legacy console); commerceView still supplies the invoices.
  const { data: paymentsData, refetch: refetchPayments } = useGetPaymentsQuery({
    variables: { data: { type: 'invoice' } },
    fetchPolicy: 'network-only',
    onError: (err) => toast(err.message),
  });
  const paymentByItemId = useMemo(() => {
    const map = new Map<string, NonNullable<typeof paymentsData>['getPayments']['payments'][number]>();
    paymentsData?.getPayments?.payments.forEach((p) => map.set(p.itemId, p));
    return map;
  }, [paymentsData]);

  // ---- Phase A: on-demand PDF download -------------------------------------
  const [downloadingId, setDownloadingId] = useState<string | null>(null);
  const [getPdf] = useGetGeneratedPdfReportLazyQuery({
    fetchPolicy: 'network-only',
    onCompleted: (pdf) => {
      const report = pdf?.getGeneratedPdfReport;
      if (report?.downloadUrl) {
        try {
          const blob = base64ToBlob(report.downloadUrl, report.contentType || 'application/pdf');
          downloadBlob(blob, report.filename || 'invoice.pdf');
        } catch {
          toast('Failed to download PDF');
        }
      }
      setDownloadingId(null);
    },
    onError: (err) => {
      toast(err.message);
      setDownloadingId(null);
    },
  });
  const handleDownload = useCallback(
    (reportId: string) => {
      setDownloadingId(reportId);
      getPdf({ variables: { Id: reportId } });
    },
    [getPdf],
  );

  // ---- Phase C: Stripe pay flow --------------------------------------------
  const [payingId, setPayingId] = useState<string | null>(null);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [clientSecret, setClientSecret] = useState('');
  const [activeInvoice, setActiveInvoice] = useState<Invoice | null>(null);
  const [updatePayment] = useUpdatePaymentMutation({
    onError: (err) => toast(err.message),
  });

  const handlePay = useCallback(
    async (invoice: Invoice) => {
      const payment = paymentByItemId.get(invoice.id);
      if (!payment) {
        toast('No payment record for this invoice yet');
        return;
      }
      setPayingId(invoice.id);
      try {
        // Already provisioned — open the dialog with the existing secret.
        if (payment.extra) {
          setActiveInvoice(invoice);
          setClientSecret(payment.extra);
          setDialogOpen(true);
          return;
        }
        // Otherwise provision a Stripe PaymentIntent, then read its secret
        // (returned in payment.extra) from a fresh payments fetch.
        const res = await updatePayment({
          variables: {
            data: {
              id: payment.id,
              paymentMethod: 'stripe',
              payerEmail: user?.email,
              payerName: user?.name,
            },
          },
        });
        const mutationError = res.errors?.[0];
        if (mutationError) {
          toast(mutationError.message);
          return;
        }
        const refreshed = await refetchPayments();
        const updated = refreshed.data?.getPayments?.payments.find(
          (p) => p.itemId === invoice.id,
        );
        if (updated?.extra) {
          setActiveInvoice(invoice);
          setClientSecret(updated.extra);
          setDialogOpen(true);
        } else {
          toast('Could not start payment — please try again');
        }
      } catch (err) {
        toast(err instanceof Error ? err.message : 'Could not start payment');
      } finally {
        setPayingId(null);
      }
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps -- `toast` is a stable context fn
    [paymentByItemId, updatePayment, refetchPayments, user?.email, user?.name],
  );

  const firstUnpaid = useMemo(() => invoices.find((i) => !i.isPaid), [invoices]);

  const closeDialog = useCallback(() => {
    setDialogOpen(false);
    setClientSecret('');
    setActiveInvoice(null);
  }, []);

  const handlePaymentSuccess = useCallback(async () => {
    toast('Payment completed successfully');
    await Promise.all([refetchOverview(), refetchPayments()]);
    closeDialog();
  }, [toast, refetchOverview, refetchPayments, closeDialog]);

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
            {sectionValue(
              balanceSection?.error ? null : money(balanceSection?.outstandingAmount ?? 0),
              balanceSection?.error ?? null,
            )}
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
              disabled={!firstUnpaid || payingId !== null}
              onClick={() => firstUnpaid && handlePay(firstUnpaid)}
            >
              {payingId && payingId === firstUnpaid?.id ? 'Starting…' : 'Pay now'}
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
              onCta={() => refetchOverview()}
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
                  <TableCell sx={{ width: 160 }} />
                </TableRow>
              </TableHead>
              <TableBody>
                {invoices.map((inv) => {
                  const canPay = !inv.isPaid && paymentByItemId.has(inv.id);
                  return (
                    <TableRow key={inv.id}>
                      <TableCell className="tnum" style={{ fontWeight: 600 }}>
                        {inv.id.slice(0, 8)}
                      </TableCell>
                      <TableCell>{inv.period}</TableCell>
                      <TableCell align="right" className="tnum">
                        {money(centsToMajor(inv.rawReport?.totalAmountCents))}
                      </TableCell>
                      <TableCell>
                        <StatusBadge status={inv.isPaid ? 'paid' : 'pending'} variant="pill" />
                      </TableCell>
                      <TableCell align="right">
                        <div style={{ display: 'flex', gap: 4, justifyContent: 'flex-end' }}>
                          {canPay && (
                            <Button
                              size="small"
                              variant="outlined"
                              disabled={payingId !== null}
                              onClick={() => handlePay(inv)}
                            >
                              {payingId === inv.id ? <CircularProgress size={16} /> : 'Pay'}
                            </Button>
                          )}
                          <Button
                            size="small"
                            startIcon={
                              downloadingId === inv.id ? (
                                <CircularProgress size={16} />
                              ) : (
                                <DownloadRounded />
                              )
                            }
                            sx={{ color: 'var(--uk-ink-3)' }}
                            disabled={downloadingId === inv.id}
                            onClick={() => handleDownload(inv.id)}
                          >
                            PDF
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          )}
        </div>
        {!invLoading && !invoicesSection?.error && (
          <TableFooter count={invoices.length} noun="invoices" />
        )}
      </div>

      <StripePaymentDialog
        open={dialogOpen}
        onClose={closeDialog}
        clientSecret={clientSecret}
        amountCents={activeInvoice?.rawReport?.totalAmountCents}
        periodLabel={activeInvoice?.period}
        onPaymentSuccess={handlePaymentSuccess}
        onPaymentError={(message) => toast(message)}
      />
    </>
  );

  if (embed) return body;

  return (
    <div className="page">
      <PageHeader
        title="Billing"
        sub="What you owe Ukama for running this network."
        actions={
          <Button
            variant="contained"
            startIcon={<BoltRounded />}
            disabled={!firstUnpaid || payingId !== null}
            onClick={() => firstUnpaid && handlePay(firstUnpaid)}
          >
            Pay now
          </Button>
        }
      />
      {body}
    </div>
  );
}
