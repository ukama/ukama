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
  Skeleton,
} from '@mui/material';
import { colors } from '@/theme';
import { PackagesResDto } from '@/client/graphql/generated';

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
  if (loadingPackageHistories || !packagesData) {
    return (
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
            {[1, 2, 3].map((index) => (
              <TableRow
                key={`skeleton-${index}`}
                sx={{
                  backgroundColor:
                    index % 2 === 0 ? colors.black10 : 'transparent',
                }}
              >
                <TableCell>
                  <Skeleton variant="text" width={120} height={20} />
                </TableCell>
                <TableCell>
                  <Skeleton variant="text" width={100} height={20} />
                </TableCell>
                <TableCell>
                  <Skeleton variant="text" width={60} height={20} />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Box>
    );
  }

  const sortedHistories =
    packageHistories && packageHistories.length > 0
      ? [...packageHistories]
          .sort(
            (a, b) =>
              new Date(b.start_date).getTime() -
              new Date(a.start_date).getTime(),
          )
          .slice(0, 5)
      : [];

  return (
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
          {sortedHistories.length > 0 ? (
            sortedHistories.map((pkg, idx) => {
              console.log('Package in table:', {
                id: pkg.id,
                isActive: pkg.isActive,
                simId: pkg.simId,
                package_id: pkg.package_id,
              });

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
                      ? `${packageDetails.dataVolume} ${packageDetails.dataUnit} / ${packageDetails.duration} ${
                          packageDetails.duration === 1 ? 'day' : 'days'
                        }`
                      : 'Unknown plan'}
                  </TableCell>
                  <TableCell sx={{ color: colors.black70 }}>
                    <span
                      style={{ color: colors.black38, fontStyle: 'italic' }}
                    >
                      N/A
                    </span>
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
  );
};

export default SubscriberHistoryTab;
