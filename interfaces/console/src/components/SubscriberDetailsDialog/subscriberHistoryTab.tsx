/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import {
  Box,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  CircularProgress,
} from '@mui/material';
import { colors } from '@/theme';
import { PackagesResDto } from '@/client/graphql/generated';
import { formatBytesToGB } from '@/utils';

interface SubscriberHistoryTabProps {
  packageHistories?: any[];
  packagesData?: PackagesResDto;
  loadingPackageHistories?: boolean;
}

const SubscriberHistoryTab: React.FC<SubscriberHistoryTabProps> = ({
  packageHistories,
  packagesData,
  loadingPackageHistories,
}) => {
  return (
    <>
      {loadingPackageHistories ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', py: 2 }}>
          <CircularProgress size={24} />
        </Box>
      ) : (
        <Box sx={{ overflowX: 'auto' }}>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Time active</TableCell>
                <TableCell>Data plan</TableCell>
                <TableCell>Data usage</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {packageHistories && packageHistories.length > 0 ? (
                [...packageHistories]
                  .sort(
                    (a, b) =>
                      new Date(b.start_date).getTime() -
                      new Date(a.start_date).getTime(),
                  )
                  // TODO: show only 5 need to discussion if we need to show more than 5 or add pagination
                  .slice(0, 5)
                  .map((pkg, idx) => {
                    const packageDetails = packagesData?.packages?.find(
                      (p) => p.uuid === pkg.package_id,
                    );
                    return (
                      <TableRow
                        key={pkg.id}
                        sx={{
                          backgroundColor:
                            idx % 2 === 1 ? colors.black10 : 'transparent',
                        }}
                      >
                        <TableCell sx={{ color: colors.black70 }}>
                          {`${new Date(pkg.start_date).toLocaleDateString()} - ${new Date(pkg.end_date).toLocaleDateString()}`}
                        </TableCell>
                        <TableCell sx={{ color: colors.black70 }}>
                          {packageDetails
                            ? `${packageDetails.dataVolume} ${packageDetails.dataUnit} /${packageDetails.duration} ${packageDetails.duration === 1 ? 'day' : 'days'}`
                            : 'Unknown plan'}
                        </TableCell>
                        <TableCell sx={{ color: colors.black70 }}>
                          {isNaN(Number(pkg.dataUsage))
                            ? '0 GB'
                            : `${formatBytesToGB(Number(pkg.dataUsage))} GB`}
                        </TableCell>
                      </TableRow>
                    );
                  })
              ) : (
                <TableRow>
                  <TableCell colSpan={3} align="center">
                    No package history available
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </Box>
      )}
    </>
  );
};

export default SubscriberHistoryTab;
