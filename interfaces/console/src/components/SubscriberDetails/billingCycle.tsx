/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import {
  TableContainer,
  Table,
  TableHead,
  TableRow,
  TableCell,
} from '@mui/material';

const BillingCycle: React.FC = () => (
  <TableContainer>
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>
            <strong style={{ fontWeight: 'bold' }}> Billing cycle</strong>
          </TableCell>
          <TableCell>
            <strong style={{ fontWeight: 'bold' }}> Data usage</strong>
          </TableCell>
          <TableCell>
            <strong style={{ fontWeight: 'bold' }}> Data plan</strong>
          </TableCell>
        </TableRow>
      </TableHead>
    </Table>
  </TableContainer>
);

export default BillingCycle;
