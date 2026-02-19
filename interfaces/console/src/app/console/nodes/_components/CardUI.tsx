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
}

export default function CardUI({ children, sx }: ICardUI) {
  return (
    <Paper
      elevation={1}
      sx={{
        p: 2,
        borderRadius: 1,
        height: "-webkit-fill-available",
        border: `1px solid ${colors.white}`,
        ...sx,
      }}
    >
      {children}
    </Paper>
  );
}