/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Skeleton rows mirroring a table's column count (table-kit.jsx). */
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
    <table className="tbl skeleton-tbl">
      <tbody>
        {Array.from({ length: rows }).map((_, r) => (
          <tr key={r} className="static">
            {Array.from({ length: cols }).map((_, c) => (
              <td key={c}>
                {lead && c === 0 ? (
                  <span style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
                    <span className="sk-dot" />
                    <span className="sk-bar" style={{ width: '58%' }} />
                  </span>
                ) : (
                  <span
                    className="sk-bar"
                    style={{
                      width: c === 0 ? '62%' : c === cols - 1 ? '38%' : '74%',
                    }}
                  />
                )}
              </td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  );
}
