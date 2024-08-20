'use client';
import InstallSiteLoading from '@/components/InstallSiteLoading';
import { CenterContainer } from '@/styles/global';
import colors from '@/theme/colors';
import GradientWrapper from '@/wrappers/gradiantWrapper';
import {
  Button,
  Checkbox,
  FormControlLabel,
  Paper,
  Stack,
  SvgIcon,
  Typography,
} from '@mui/material';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import NetworkInfo from '../../../public/svg/NetworkInfo';

const Page = () => {
  const router = useRouter();
  const [isChecked, setIsChecked] = useState(false);
  const [checkForInstallation, setCheckForInstallation] = useState(false);
  const handleNext = () => {
    if (isChecked) {
      setCheckForInstallation(true);
    }
  };
  const handleSkip = () => router.push('/console/home');
  const handleOnInstalled = (isChecked: boolean) => {
    setIsChecked(isChecked);
  };
  const onInstallProgressComplete = () => router.push(`/configure/123`);

  return (
    <CenterContainer>
      {checkForInstallation ? (
        <InstallSiteLoading
          duration={10}
          onCompleted={onInstallProgressComplete}
        />
      ) : (
        <GradientWrapper>
          <Paper elevation={0} sx={{ px: 4, py: 2 }}>
            <Typography variant="h6" fontWeight={500}>
              Install site -{' '}
              <span style={{ color: colors.black70, fontWeight: 400 }}>
                <i>optional</i> (1/6)
              </span>
            </Typography>
            <Stack mt={3} mb={3} direction={'column'} alignItems={'center'}>
              <Typography variant="body1" color={colors.vulcan}>
                If you would like to set up your network later, or if someone
                else will set up your network for you, skip this step.
                <br />
                <br />
                Install your node at the intended location, and ensure it is
                connected to a power and backhaul source. These three elements
                form a site, an abstracted representation of the aforementioned
                components. Each site can also hold up to three nodes for a
                stronger connection.
                <br />
                <br />
                You can follow the installation instructions in the provided
                manual, or in the PDF <Link href={''}>here</Link>.
              </Typography>
              <SvgIcon sx={{ width: 240, height: 176, mt: 2, mb: 4 }}>
                {NetworkInfo}
              </SvgIcon>
              <FormControlLabel
                sx={{ alignSelf: 'baseline' }}
                control={
                  <Checkbox
                    sx={{ p: 0, px: 1.3 }}
                    onChange={(e) => handleOnInstalled(e.target.checked)}
                  />
                }
                label="I have installed my node, power, and backhaul components"
              />
            </Stack>
            <Stack
              mb={1}
              spacing={2}
              direction={'row'}
              justifyContent={'space-between'}
            >
              <Button
                variant="text"
                onClick={handleSkip}
                sx={{ color: colors.black70, p: 0 }}
              >
                Skip
              </Button>
              <Button
                variant="contained"
                onClick={handleNext}
                disabled={!isChecked}
              >
                Next
              </Button>
            </Stack>
          </Paper>
        </GradientWrapper>
      )}
    </CenterContainer>
  );
};
export default Page;
