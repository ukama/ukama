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
  currencySymbol?: string;
}

const SubscriberDataPlansTab: React.FC<SubscriberDataPlansTabProps> = ({
  packageHistories,
  packagesData,
  dataUsage,
  currencySymbol,
}) => {
  return (
    <>
      <Box>
        <Box sx={{ display: 'flex', mb: 1 }}>
          <Typography
            variant="body2"
            color={colors.black70}
            sx={{ width: 150 }}
          >
            Date plan
          </Typography>
          <Typography variant="body2" color={colors.black70}>
            {(() => {
              const currentPackage =
                packageHistories && packageHistories.length > 0
                  ? packagesData?.packages?.find(
                      (p) => p.uuid === packageHistories[0]?.package_id,
                    )
                  : null;

              return currentPackage
                ? `${currentPackage.dataVolume} ${currentPackage.dataUnit} / ${currentPackage.duration} ${currentPackage.duration === 1 ? 'day' : 'days'} /${currencySymbol}   ${currentPackage.amount}`
                : 'No active plan';
            })()}
          </Typography>
        </Box>

        <Box sx={{ display: 'flex', mb: 1 }}>
          <Typography
            variant="body2"
            color={colors.black70}
            sx={{ width: 150 }}
          >
            Data usage
          </Typography>
          <Typography variant="body2" color={colors.black70}>
            {isNaN(Number(dataUsage))
              ? '0 GB'
              : `${formatBytesToGB(Number(dataUsage))} GB`}
          </Typography>
        </Box>

        <Box sx={{ display: 'flex', mb: 1 }}>
          <Typography
            variant="body2"
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
                // TODO: show only 3 need to discussion if we need to show more than 3
                const packagesToShow = upcomingPackages.slice(0, 3);

                return upcomingPackages.length > 0 ? (
                  <Stack spacing={1}>
                    {packagesToShow.map((pkg, index) => {
                      const packageDetails = packagesData?.packages?.find(
                        (p) => p.uuid === pkg.package_id,
                      );

                      return (
                        <Typography
                          variant="body2"
                          key={pkg.id}
                          color={colors.black70}
                        >
                          {packageDetails
                            ? `${packageDetails.dataVolume} ${packageDetails.dataUnit} / ${packageDetails.duration} ${packageDetails.duration === 1 ? 'day' : 'days'} / ${currencySymbol}  ${packageDetails.amount}`
                            : 'Unknown plan'}
                        </Typography>
                      );
                    })}
                  </Stack>
                ) : (
                  <Typography variant="body2" color={colors.black70}>
                    No upcoming plans
                  </Typography>
                );
              })()
            ) : (
              <Typography variant="body2" color={colors.black70}>
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
