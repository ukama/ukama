/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

export interface GlowieBgProps {
  children?: React.ReactNode;
  backgroundSize?: string;
}

export default function GlowieBg({
  children,
  backgroundSize = "110%",
}: Readonly<GlowieBgProps>) {
  return (
    <div
      style={{
        height: "100%",
        backgroundSize,
        display: "flex",
        paddingTop: "44px",
        position: "relative",
        alignItems: "center",
        justifyContent: "center",
        backgroundImage: "url(/images/vglow.png)",
        backgroundPosition: "center",
        backgroundRepeat: "no-repeat",
      }}
    >
        {children}
    </div>
  );
}