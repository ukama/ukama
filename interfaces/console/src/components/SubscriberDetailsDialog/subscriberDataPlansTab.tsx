import React from 'react';
import { Typography, Box, styled } from '@mui/material';
import { colors } from '@/theme';
import { PackagesResDto } from '@/client/graphql/generated';

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

const FieldValue = styled(Typography)(({ theme }) => ({
  fontSize: theme.typography.body1.fontSize,
  marginBottom: theme.spacing(2),
}));

const SubscriberDataPlansTab: React.FC<SubscriberDataPlansTabProps> = ({
  packageHistories,
  packagesData,
  dataUsage,
}) => {
  return (
    <>
      <Box sx={{ mb: 3 }}>
        <FieldLabel>DATE PLAN</FieldLabel>
        {packageHistories && packageHistories.some((pkg) => pkg.is_active) ? (
          <FieldValue>
            {(() => {
              const activePackage = packageHistories.find(
                (pkg) => pkg.is_active,
              );
              if (!activePackage) return 'No active plan';

              const packageDetails = packagesData?.packages?.find(
                (p) => p.uuid === activePackage.package_id,
              );

              return packageDetails
                ? `${packageDetails.dataVolume} ${packageDetails.dataUnit} / ${packageDetails.duration} days / $${packageDetails.amount}`
                : 'Unknown plan';
            })()}
          </FieldValue>
        ) : (
          <FieldValue>No active plan</FieldValue>
        )}
      </Box>

      <Box sx={{ mb: 3 }}>
        <FieldLabel>DATA USAGE</FieldLabel>
        {packageHistories && packageHistories.some((pkg) => pkg.is_active) ? (
          <FieldValue>
            {dataUsage != undefined ? `${dataUsage}GB` : 'No Usage'}
          </FieldValue>
        ) : (
          <FieldValue>Not applicable</FieldValue>
        )}
      </Box>

      <Box>
        <FieldLabel>UPCOMING</FieldLabel>
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
                  <FieldValue key={pkg.id} sx={{ mb: 1 }}>
                    {packageDetails
                      ? `${packageDetails.dataVolume} ${packageDetails.dataUnit} / ${packageDetails.duration} days / $${packageDetails.amount}`
                      : 'Unknown plan'}
                  </FieldValue>
                );
              })
            ) : (
              <FieldValue>No upcoming plans</FieldValue>
            );
          })()
        ) : (
          <FieldValue>No upcoming plans</FieldValue>
        )}
      </Box>
    </>
  );
};

export default SubscriberDataPlansTab;
