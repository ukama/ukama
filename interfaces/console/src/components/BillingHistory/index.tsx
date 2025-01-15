/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React, { useState } from 'react';
import {
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  Typography,
  TextField,
  InputAdornment,
  Box,
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';

interface BillingHistoryTableProps {
  data: {
    id: string;
    posted: string;
    billing: string;
    payment: string;
    description: string;
    pdf: string;
  }[];
}

const BillingHistoryTable: React.FC<BillingHistoryTableProps> = ({ data }) => {
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(5);
  const [searchQuery, setSearchQuery] = useState('');

  const handleChangePage = (event: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const filteredData = data.filter(
    (row) =>
      row.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
      row.payment.toLowerCase().includes(searchQuery.toLowerCase()) ||
      row.posted.toLowerCase().includes(searchQuery.toLowerCase()) ||
      row.billing.toLowerCase().includes(searchQuery.toLowerCase()),
  );

  return (
    <Paper sx={{ p: 2, minHeight: '300px', borderRadius: '10px' }}>
      <Typography variant="h6" sx={{ mb: 2 }}>
        Billing History
      </Typography>
      {data.length === 0 ? (
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            height: '150px',
          }}
        >
          <Typography variant="body1">No past invoices available.</Typography>
        </Box>
      ) : (
        <>
          <TextField
            label="Search"
            variant="outlined"
            sx={{ mb: 2, width: '30%' }}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            InputProps={{
              endAdornment: (
                <InputAdornment position="end">
                  <SearchIcon />
                </InputAdornment>
              ),
            }}
          />
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Billing Period</TableCell>
                  <TableCell>Posted</TableCell>
                  <TableCell>Payment Status</TableCell>
                  <TableCell>Description</TableCell>
                  <TableCell>PDF</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {filteredData
                  .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                  .map((row) => (
                    <TableRow key={row.id}>
                      <TableCell>{row.billing}</TableCell>
                      <TableCell>{row.posted}</TableCell>
                      <TableCell>{row.payment}</TableCell>
                      <TableCell>{row.description}</TableCell>
                      <TableCell>
                        <a
                          href={row.pdf}
                          target="_blank"
                          rel="noopener noreferrer"
                        >
                          Download
                        </a>
                      </TableCell>
                    </TableRow>
                  ))}
              </TableBody>
            </Table>
          </TableContainer>
          <TablePagination
            rowsPerPageOptions={[5, 10, 25]}
            component="div"
            count={filteredData.length}
            rowsPerPage={rowsPerPage}
            page={page}
            onPageChange={handleChangePage}
            onRowsPerPageChange={handleChangeRowsPerPage}
          />
        </>
      )}
    </Paper>
  );
};

export default BillingHistoryTable;
