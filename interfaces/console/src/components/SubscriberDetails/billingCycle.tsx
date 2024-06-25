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
