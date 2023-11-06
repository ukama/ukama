import React, { useCallback } from 'react';
import {
  Button,
  Typography,
  InputLabel,
  FormControl,
  OutlinedInput,
  MenuItem,
  Select,
  Stack,
  Grid,
  SelectChangeEvent,
} from '@mui/material';
import { PackageDto, SimDto } from '@/generated';
import colors from '@/styles/theme/colors';
import { useState } from 'react';
import { makeStyles } from '@mui/styles';

const useStyles = makeStyles(() => ({
  selectStyle: () => ({
    width: '100%',
    height: '48px',
  }),
  formControl: {
    width: '100%',
    height: '48px',
  },
}));

interface SubscriberDialogProps {
  onClose: () => void;
  handlePlanInstallation: Function;
  goBack: () => void;
  submitButtonState: boolean;
  packages: PackageDto[];
  sims: SimDto[];

  selectedSimType: string;
}

const Step2: React.FC<SubscriberDialogProps> = React.memo(
  ({
    onClose,
    handlePlanInstallation,
    goBack,
    submitButtonState = false,
    packages,
    sims,
    selectedSimType,
  }) => {
    const classes = useStyles();
    const [plan, setPlan] = useState<string>('');
    const [simIccid, setSimIccid] = useState<string>('');

    const handleselectPLan = useCallback((e: SelectChangeEvent) => {
      setPlan(e.target.value);
    }, []);
    const handleselectSim = useCallback((e: SelectChangeEvent) => {
      setSimIccid(e.target.value);
    }, []);

    const handleButtonClick = useCallback(() => {
      handlePlanInstallation(plan, simIccid);
    }, [handlePlanInstallation, plan, simIccid]);

    const getButtonState = () => {
      let isButtonDisabled;
      if (selectedSimType == 'eSim') {
        isButtonDisabled = !plan;
        return isButtonDisabled;
      }
      isButtonDisabled = !plan || !simIccid;
      return isButtonDisabled;
    };
    const isButtonDisabled = getButtonState();
    return (
      <>
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <FormControl variant="outlined" className={classes.formControl}>
              <InputLabel
                shrink
                variant="outlined"
                required
                htmlFor="outlined-age-always-notched"
              >
                {selectedSimType == 'pSim' ? `pSIM ICCID` : `eSIM ICCID`}
              </InputLabel>

              <Select
                variant="outlined"
                onChange={handleselectSim}
                value={simIccid}
                required
                sx={{
                  '& legend': { width: '93px' },
                }}
                input={
                  <OutlinedInput
                    notched
                    label="ICCID"
                    name={'iccid'}
                    id="outlined-age-always-notched"
                    fullWidth
                  />
                }
                disabled={selectedSimType == 'eSim' ? true : false}
                MenuProps={{
                  disablePortal: false,
                  PaperProps: {
                    sx: {
                      boxShadow:
                        '0px 5px 5px -3px rgba(0, 0, 0, 0.2), 0px 8px 10px 1px rgba(0, 0, 0, 0.14), 0px 3px 14px 2px rgba(0, 0, 0, 0.12)',
                      borderRadius: '4px',
                    },
                  },
                }}
                className={classes.selectStyle}
              >
                {sims.map((sim) => (
                  <MenuItem
                    key={sim.id}
                    value={sim.iccid}
                    sx={{
                      m: 0,
                      p: '6px 16px',
                    }}
                  >
                    <Typography variant="body1">{sim.iccid}</Typography>
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          <Grid item xs={12}>
            <FormControl variant="outlined" className={classes.formControl}>
              <InputLabel
                shrink
                variant="outlined"
                required
                htmlFor="outlined-age-always-notched"
              >
                DATA PLAN
              </InputLabel>

              <Select
                variant="outlined"
                onChange={handleselectPLan}
                value={plan}
                required
                sx={{
                  '& legend': { width: '93px' },
                }}
                input={
                  <OutlinedInput
                    fullWidth
                    notched
                    label="Plan"
                    name={'plan'}
                    id="outlined-age-always-notched"
                  />
                }
                MenuProps={{
                  disablePortal: false,
                  PaperProps: {
                    sx: {
                      boxShadow:
                        '0px 5px 5px -3px rgba(0, 0, 0, 0.2), 0px 8px 10px 1px rgba(0, 0, 0, 0.14), 0px 3px 14px 2px rgba(0, 0, 0, 0.12)',
                      borderRadius: '4px',
                    },
                  },
                }}
                className={classes.selectStyle}
              >
                {packages.map((pkg) => (
                  <MenuItem
                    key={pkg.uuid}
                    value={pkg.uuid}
                    sx={{
                      m: 0,
                      p: '6px 16px',
                    }}
                  >
                    <Typography variant="body1">
                      {`${pkg.name} - $${pkg.amount}/${
                        Number(pkg.dataVolume) / 1024
                      } GB`}
                    </Typography>
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>

          <Grid item xs={12}>
            <Stack
              direction="row"
              justifyContent="space-between"
              mt={1}
              sx={{ mb: 2 }}
            >
              <Button
                variant="text"
                onClick={() => {
                  goBack();
                }}
              >
                {' Go Back'}
              </Button>

              <Stack direction="row" spacing={3}>
                <Button
                  variant="text"
                  onClick={() => {
                    onClose();
                  }}
                >
                  {' CANCEL'}
                </Button>

                <Button
                  variant="contained"
                  type="submit"
                  onClick={handleButtonClick}
                  disabled={submitButtonState || isButtonDisabled}
                  sx={{
                    color: submitButtonState ? colors.primaryLight : undefined,
                  }}
                >
                  <Typography variant="body1"> ADD SUBSCRIBER</Typography>
                </Button>
              </Stack>
            </Stack>
          </Grid>
        </Grid>
      </>
    );
  },
);
Step2.displayName = 'Step2';
export default Step2;
