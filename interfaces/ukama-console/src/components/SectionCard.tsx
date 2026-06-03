/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Card with display-font section title + right slot (biz-common.jsx). */
export default function SectionCard({
  title,
  count,
  right,
  children,
  style,
  bodyStyle,
}: {
  title?: React.ReactNode;
  count?: React.ReactNode;
  right?: React.ReactNode;
  children: React.ReactNode;
  style?: React.CSSProperties;
  bodyStyle?: React.CSSProperties;
}) {
  return (
    <div className="card card-pad" style={style}>
      {(title || right) && (
        <div className="sec-head">
          {title ? (
            <div className="sec-title">
              {title}
              {count != null && <span className="cnt tnum">{count}</span>}
            </div>
          ) : (
            <span />
          )}
          {right}
        </div>
      )}
      <div style={bodyStyle}>{children}</div>
    </div>
  );
}
