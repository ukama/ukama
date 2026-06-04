/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Skeleton rows mirroring a table's column count — MUI Table + Skeleton. */
import Skeleton from '@mui/material/Skeleton';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableRow from '@mui/material/TableRow';

export default function SkeletonTable({
  cols = 5,
  rows = 6,
  lead,
}: {
  cols?: number;
  rows?: number;
  lead?: boolean;
}) {
  return (
    <Table>
      <TableBody>
        {Array.from({ length: rows }).map((_, r) => (
          <TableRow key={r}>
            {Array.from({ length: cols }).map((_, c) => (
              <TableCell key={c}>
                {lead && c === 0 ? (
                  <span style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
                    <Skeleton variant="circular" width={28} height={28} />
                    <Skeleton variant="rounded" height={11} sx={{ width: '58%', borderRadius: 1.5 }} />
                  </span>
                ) : (
                  <Skeleton
                    variant="rounded"
                    height={11}
                    sx={{
                      borderRadius: 1.5,
                      width: c === 0 ? '62%' : c === cols - 1 ? '38%' : '74%',
                    }}
                  />
                )}
              </TableCell>
            ))}
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
