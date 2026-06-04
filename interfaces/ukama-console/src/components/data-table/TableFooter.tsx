/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Table footing — anchors a list with a count (+ optional pager slot). */
export default function TableFooter({
  count,
  noun,
  showing,
  total,
  children,
}: {
  count?: number;
  noun?: string;
  showing?: number;
  total?: number;
  children?: React.ReactNode;
}) {
  return (
    <div className="tbl-foot">
      <span className="tnum">
        {showing != null && total != null
          ? `Showing ${showing} of ${total.toLocaleString()}`
          : `${count != null ? count.toLocaleString() + ' ' : ''}${noun ?? ''}`}
      </span>
      {children && <div className="pager">{children}</div>}
    </div>
  );
}
