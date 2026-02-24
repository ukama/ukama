/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import colors from "@/theme/colors";
import { Paper, SxProps, Theme } from "@mui/material";

interface ICardUI {
  children: React.ReactNode;
  sx?: SxProps<Theme>;
  isBorderLeft?: boolean;
}

export default function CardUI({ children, sx, isBorderLeft = false }: Readonly<ICardUI>) {
  return (
    <Paper
      elevation={1}
      sx={{
        p: 2,
        borderRadius: 1,
        height: "-webkit-fill-available",
        borderLeft: isBorderLeft ? `5px solid ${colors.primaryMain}` : "none",
        ...sx,
      }}
    >
      {children}
    </Paper>
  );
}