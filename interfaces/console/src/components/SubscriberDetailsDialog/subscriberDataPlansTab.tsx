/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import { Typography, Box, Stack, Skeleton } from '@mui/material';
import { colors } from '@/theme';
import { PackagesResDto } from '@/client/graphql/generated';
import { formatBytesToGB } from '@/utils';

interface SubscriberDataPlansTabProps {
  packageHistories?: any[];
  packagesData?: PackagesResDto;
  dataUsage: string;
  currencySymbol?: string;
  loadingPackageHistories?: boolean;
}

const SubscriberDataPlansTab: React.FC<SubscriberDataPlansTabProps> = ({
  packageHistories,
  packagesData,
  dataUsage,
  currencySymbol,
  loadingPackageHistories,
}) => {
  if (loadingPackageHistories || !packageHistories) {
    return (
      <Box>
        <Box sx={{ display: 'flex', mb: 1 }}>
          <Typography
            variant="body2"
            color={colors.black70}
            sx={{ width: 150 }}
          >
            Data plan
          </Typography>
          <Skeleton variant="text" width={200} height={20} />
        </Box>

        <Box sx={{ display: 'flex', mb: 1 }}>
          <Typography
            variant="body2"
            color={colors.black70}
            sx={{ width: 150 }}
          >
            Data usage
          </Typography>
          <Skeleton variant="text" width={80} height={20} />
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
            <Stack spacing={1}>
              <Skeleton variant="text" width={200} height={20} />
              <Skeleton variant="text" width={200} height={20} />
            </Stack>
          </Box>
        </Box>
      </Box>
    );
  }

  const getCurrentPlan = () => {
    const activePackage = packageHistories.find(
      (pkg) => pkg.is_active === true,
    );
    if (!activePackage) return 'No active plan';

    const basePackage = packagesData?.packages?.find(
      (p) => p.uuid === activePackage.package_id,
    );

    if (!basePackage) return 'Unknown plan';

    return `${basePackage.dataVolume} ${basePackage.dataUnit} / ${basePackage.duration} ${
      basePackage.duration === 1 ? 'day' : 'days'
    } / ${currencySymbol ?? ''}${basePackage.amount}`;
  };

  const getUpcomingPackages = () => {
    const now = new Date();

    return packageHistories
      .filter((pkg) => {
        const startDate = new Date(pkg.start_date);
        return pkg.is_active === false && startDate > now;
      })
      .sort(
        (a, b) =>
          new Date(a.start_date).getTime() - new Date(b.start_date).getTime(),
      );
  };

  const upcomingPackages = getUpcomingPackages();

  return (
    <Box>
      <Box sx={{ display: 'flex', mb: 1 }}>
        <Typography variant="body2" color={colors.black70} sx={{ width: 150 }}>
          Date plan
        </Typography>
        <Typography variant="body2" color={colors.black70}>
          {getCurrentPlan()}
        </Typography>
      </Box>

      <Box sx={{ display: 'flex', mb: 1 }}>
        <Typography variant="body2" color={colors.black70} sx={{ width: 150 }}>
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
          {upcomingPackages.length > 0 ? (
            <Stack spacing={1}>
              {upcomingPackages.map((pkg) => {
                const basePackage = packagesData?.packages?.find(
                  (p) => p.uuid === pkg.package_id,
                );

                return (
                  <Typography
                    variant="body2"
                    key={pkg.id}
                    color={colors.black70}
                  >
                    {basePackage
                      ? `${basePackage.dataVolume} ${basePackage.dataUnit} / ${
                          basePackage.duration
                        } ${
                          basePackage.duration === 1 ? 'day' : 'days'
                        } / ${currencySymbol ?? ''}${basePackage.amount}`
                      : 'Unknown plan'}
                  </Typography>
                );
              })}
            </Stack>
          ) : (
            <Typography variant="body2" color={colors.black70}>
              No upcoming plans
            </Typography>
          )}
        </Box>
      </Box>
    </Box>
  );
};

export default SubscriberDataPlansTab;
