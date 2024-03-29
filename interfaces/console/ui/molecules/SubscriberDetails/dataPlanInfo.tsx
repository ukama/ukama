import React from 'react';
import Stack from '@mui/material/Stack';
import Typography from '@mui/material/Typography';
import Skeleton from '@mui/material/Skeleton';
import { colors } from '@/styles/theme';

interface Props {
  packageName: string | undefined;
  currentSite: string | undefined;
  bundle: string | undefined;
}

const DataPlanComponent: React.FC<Props> = ({
  packageName,
  currentSite,
  bundle,
}) => {
  return (
    <Stack direction="column" spacing={2}>
      <Stack direction="row" spacing={2}>
        <Typography variant="body1" sx={{ color: colors.black }}>
          Data plan
        </Typography>
        <Typography variant="subtitle1" sx={{ color: colors.black }}>
          {packageName && packageName.length ? (
            packageName
          ) : (
            <Skeleton
              variant="rectangular"
              width={120}
              height={24}
              sx={{ backgroundColor: colors.black10 }}
            />
          )}
        </Typography>
      </Stack>
      <Stack direction="row" spacing={2}>
        <Typography variant="body1" sx={{ color: colors.black }}>
          Current site
        </Typography>
        <Typography variant="subtitle1" sx={{ color: colors.black }}>
          {currentSite || ''}
        </Typography>
      </Stack>
      <Stack direction="row" spacing={2}>
        <Typography variant="body1" sx={{ color: colors.black }}>
          Month usage
        </Typography>
        <Typography variant="subtitle1" sx={{ color: colors.black }}>
          {bundle && bundle.length ? (
            bundle
          ) : (
            <Skeleton
              variant="rectangular"
              width={120}
              height={24}
              sx={{ backgroundColor: colors.black10 }}
            />
          )}
        </Typography>
      </Stack>
    </Stack>
  );
};

export default DataPlanComponent;
