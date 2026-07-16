/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import type { ReceiptView } from '@/lib/receipt';

const INK: [number, number, number] = [28, 30, 34];
const MUTED: [number, number, number] = [137, 141, 149];
const LINE: [number, number, number] = [231, 233, 238];

/**
 * Generate and download a receipt as receipt-<shortId>.pdf. jsPDF is imported
 * lazily so it is only bundled/loaded when a receipt is actually downloaded.
 */
export async function downloadReceiptPdf(v: ReceiptView): Promise<void> {
  const { jsPDF } = await import('jspdf');
  const doc = new jsPDF({ unit: 'pt', format: 'a5' });

  const pageW = doc.internal.pageSize.getWidth();
  const mx = 40;
  const right = pageW - mx;
  let y = 52;

  const rule = () => {
    doc.setDrawColor(...LINE);
    doc.setLineWidth(0.8);
    doc.line(mx, y, right, y);
  };
  const label = (t: string, x: number, yy: number) => {
    doc.setFontSize(8);
    doc.setTextColor(...MUTED);
    doc.text(t.toUpperCase(), x, yy);
  };
  const value = (t: string, x: number, yy: number, size = 11) => {
    doc.setFontSize(size);
    doc.setTextColor(...INK);
    doc.text(t || '—', x, yy);
  };

  doc.setFont('helvetica', 'bold');
  doc.setFontSize(15);
  doc.setTextColor(...INK);
  doc.text('Payment receipt', mx, y);
  doc.setFont('helvetica', 'normal');
  doc.setFontSize(9);
  doc.setTextColor(...MUTED);
  doc.text(v.orgName, mx, y + 15);

  doc.setFont('helvetica', 'bold');
  doc.setFontSize(10);
  doc.setTextColor(3, 116, 75);
  doc.text(v.status || 'Completed', right, y, { align: 'right' });
  doc.setFont('helvetica', 'normal');

  y += 40;
  rule();
  y += 22;

  label('Receipt no', mx, y);
  label('Paid on', mx + 130, y);
  label('Method', right - 70, y);
  y += 15;
  value(v.id.slice(0, 8), mx, y);
  value(v.paidOn, mx + 130, y, 10);
  value(v.method, right - 70, y);

  y += 30;
  label('Billed to', mx, y);
  y += 15;
  value('Walk-in customer', mx, y);
  if (v.simId) {
    doc.setFontSize(9);
    doc.setTextColor(...MUTED);
    doc.text(`SIM ${v.simId.slice(0, 8)} · package top-up`, mx, y + 13);
  }

  y += 34;
  rule();
  y += 18;
  label('Description', mx, y);
  doc.text('AMOUNT', right, y, { align: 'right' });
  y += 20;
  value(v.planName, mx, y, 11);
  doc.setFontSize(9);
  doc.setTextColor(...MUTED);
  doc.text('Data package · qty 1', mx, y + 13);
  value(v.amountLabel, right, y, 11);
  doc.text(v.amountLabel, right, y, { align: 'right' });

  y += 34;
  rule();
  y += 22;
  doc.setFont('helvetica', 'bold');
  doc.setFontSize(12);
  doc.setTextColor(...INK);
  doc.text('Total paid', mx, y);
  doc.setFontSize(15);
  doc.text(v.amountLabel, right, y, { align: 'right' });
  doc.setFont('helvetica', 'normal');

  y += 34;
  label('Payment ID', mx, y);
  y += 14;
  doc.setFontSize(9);
  doc.setTextColor(...MUTED);
  doc.text(v.id, mx, y);

  const bottom = doc.internal.pageSize.getHeight() - 34;
  doc.setFontSize(8);
  doc.setTextColor(...MUTED);
  doc.text('Auto-generated · not a tax invoice', mx, bottom);

  doc.save(`receipt-${v.id.slice(0, 8)}.pdf`);
}
