/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/theme';
import { Select, styled } from '@mui/material';

const SelectStyle = styled(Select)({
  width: 'fit-content',
  color: colors.primaryMain,
});

const SelectDisplayProps = {
  style: {
    fontWeight: 600,
    display: 'flex',
    fontSize: '20px',
    marginLeft: '4px',
    alignItems: 'center',
    minWidth: 'fit-content',
  },
};

const PaperProps = {
  boxShadow:
    '0px 5px 5px -3px rgba(0, 0, 0, 0.2), 0px 8px 10px 1px rgba(0, 0, 0, 0.14), 0px 3px 14px 2px rgba(0, 0, 0, 0.12)',
  borderRadius: '4px',
};

export { PaperProps, SelectDisplayProps, SelectStyle };
