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
                            ? `${packageDetails.dataVolume} ${packageDetails.dataUnit} / ${packageDetails.duration} days`
                            : 'Unknown plan'}
                        </TableCell>
                        <TableCell>
                          {(pkg.dataUsage || 'N/A') + ' GB'}
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
