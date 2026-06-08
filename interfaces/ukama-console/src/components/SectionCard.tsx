/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/** Card with display-font section title + right slot (biz-common.jsx).
 *  When `selectable`, the card reads as a button with a left accent bar that
 *  turns blue while `active` — used by the node detail rail to pick a section. */
export default function SectionCard({
  title,
  count,
  right,
  children,
  style,
  bodyStyle,
  selectable,
  active,
  onClick,
}: {
  title?: React.ReactNode;
  count?: React.ReactNode;
  right?: React.ReactNode;
  children: React.ReactNode;
  style?: React.CSSProperties;
  bodyStyle?: React.CSSProperties;
  selectable?: boolean;
  active?: boolean;
  onClick?: () => void;
}) {
  const cls = [
    'card',
    'card-pad',
    selectable ? 'card-selectable' : '',
    selectable && active ? 'is-active' : '',
  ]
    .filter(Boolean)
    .join(' ');
  return (
    <div
      className={cls}
      style={style}
      role={selectable ? 'button' : undefined}
      tabIndex={selectable ? 0 : undefined}
      aria-pressed={selectable ? !!active : undefined}
      onClick={onClick}
      onKeyDown={
        selectable && onClick
          ? (e) => {
              if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                onClick();
              }
            }
          : undefined
      }
    >
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
