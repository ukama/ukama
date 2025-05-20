/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Typography, Box, styled, Stack } from '@mui/material';
import { colors } from '@/theme';
import { PackagesResDto } from '@/client/graphql/generated';
import { formatBytesToGB } from '@/utils';

interface SubscriberDataPlansTabProps {
  packageHistories?: any[];
  packagesData?: PackagesResDto;
  dataUsage: string;
}

const FieldLabel = styled(Typography)(({ theme }) => ({
  color: colors.black38,
  fontSize: theme.typography.caption.fontSize,
  lineHeight: theme.typography.caption.lineHeight,
  marginBottom: theme.spacing(0.5),
}));

const SubscriberDataPlansTab: React.FC<SubscriberDataPlansTabProps> = ({
  packageHistories,
  packagesData,
  dataUsage,
}) => {
  return (
    <>
      <Box>
        <Box sx={{ display: 'flex', mb: 1 }}>
          <Typography
            variant="body1"
            color={colors.black70}
            sx={{ width: 150 }}
          >
            Date plan
          </Typography>
          <Typography variant="subtitle1" color={colors.black70}>
            {(() => {
              const currentPackage =
                packageHistories && packageHistories.length > 0
                  ? packagesData?.packages?.find(
                      (p) => p.uuid === packageHistories[0]?.package_id,
                    )
                  : null;

              return currentPackage
                ? `${currentPackage.dataVolume} ${currentPackage.dataUnit} / ${currentPackage.duration} days / $${currentPackage.amount}`
                : 'No active plan';
            })()}
          </Typography>
        </Box>

        <Box sx={{ display: 'flex', mb: 1 }}>
          <Typography
            variant="body1"
            color={colors.black70}
            sx={{ width: 150 }}
          >
            Data usage
          </Typography>
          <Typography variant="subtitle1" color={colors.black70}>
            {isNaN(Number(dataUsage))
              ? '0 GB'
              : `${formatBytesToGB(Number(dataUsage))} GB`}
          </Typography>
        </Box>

        <Box sx={{ display: 'flex', mb: 1 }}>
          <Typography
            variant="body1"
            color={colors.black70}
            sx={{ width: 150, alignSelf: 'flex-start' }}
          >
            Upcoming
          </Typography>
          <Box>
            {packageHistories && packageHistories.length > 0 ? (
              (() => {
                const upcomingPackages = packageHistories
                  .filter((pkg) => new Date(pkg.start_date) > new Date())
                  .sort(
                    (a, b) =>
                      new Date(a.start_date).getTime() -
                      new Date(b.start_date).getTime(),
                  );

                return upcomingPackages.length > 0 ? (
                  upcomingPackages.map((pkg, index) => {
                    const packageDetails = packagesData?.packages?.find(
                      (p) => p.uuid === pkg.package_id,
                    );

                    return (
                      <Typography
                        variant="subtitle1"
                        key={pkg.id}
                        sx={{ mb: 2 }}
                        color={colors.black70}
                      >
                        {packageDetails
                          ? `${packageDetails.dataVolume} ${packageDetails.dataUnit} / ${packageDetails.duration} days / $${packageDetails.amount}`
                          : 'Unknown plan'}
                      </Typography>
                    );
                  })
                ) : (
                  <Typography variant="subtitle1" color={colors.black70}>
                    No upcoming plans
                  </Typography>
                );
              })()
            ) : (
              <Typography variant="subtitle1" color={colors.black70}>
                No upcoming plans
              </Typography>
            )}
          </Box>
        </Box>
      </Box>
    </>
  );
};

export default SubscriberDataPlansTab;
