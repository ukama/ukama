/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import React from 'react';
import { Box } from '@mui/material';

interface TabPanelProps {
  id: string;
  index: number;
  value: number;
  children?: React.ReactNode;
}

const TabPanel = ({ id, index, value, children }: TabPanelProps) => {
  return (
    <div id={id} role="tabpanel" aria-labelledby={id} hidden={value !== index}>
      {value === index && <Box component="div">{children}</Box>}
    </div>
  );
};

export default TabPanel;
