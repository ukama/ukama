import React, { useEffect, useState } from 'react';
import { colors } from '@/styles/theme';

import {
  Dialog,
  Button,
  DialogActions,
  Grid,
  DialogContent,
  Typography,
  DialogContentText,
  DialogTitle,
  IconButton,
  Select,
  MenuItem,
  FormControl,
  OutlinedInput,
  InputLabel,
  SelectChangeEvent,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import { makeStyles } from '@mui/styles';
import { PackageDto, SubcriberToSimDto } from '@/generated';

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
interface TopUpProps {
  onCancel: () => void;
  isToPup: boolean;
  subscriberId: string;
  handleTopUp: (planId: string, simId: string) => void;
  packages: PackageDto[];
  loadingTopUp: boolean;
  sims: SubcriberToSimDto[];
}

const TopUpData: React.FC<TopUpProps> = ({
  handleTopUp,
  onCancel,
  isToPup,
  subscriberId,
  packages,
  sims,
  loadingTopUp = false,
}) => {
  const [isToppingUp, setIsToppingUp] = useState(false);
  const [plan, setPlan] = useState<string>('');
  const [sim, setSim] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const classes = useStyles();

  useEffect(() => {
    if (sim && plan) {
      setLoading(false);
    } else {
      setLoading(true);
    }
  }, [sim, plan]);

  const onTopUp = () => {
    setIsToppingUp(true);
    onTopUp();
  };

  const handleClose = () => {
    if (!isToppingUp) {
      onCancel();
      setPlan('');
      setSim('');
    }
  };

  const handleselectPLan = (e: SelectChangeEvent) => {
    setPlan(e.target.value);
  };
  const handleselectSim = (e: SelectChangeEvent) => {
    setSim(e.target.value);
  };
  const handleClick = () => {
    handleTopUp(plan, sim);
  };

  return (
    <Dialog
      open={isToPup}
      onClose={handleClose}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <DialogTitle id="alert-dialog-title">
        <Typography variant="h6">{`Top up data - ${subscriberId.slice(
          0,
          10,
        )}...`}</Typography>
      </DialogTitle>
      <IconButton
        aria-label="close"
        onClick={handleClose}
        sx={{
          position: 'absolute',
          right: 8,
          top: 8,
        }}
      >
        <CloseIcon />
      </IconButton>
      <DialogContent>
        <DialogContentText id="alert-dialog-description">
          <Typography variant="body1" sx={{ color: colors.black }}>
            {` Add more data for ${subscriberId} for the rest of this month. Note: just like other data, it expires at the end of the month.`}
          </Typography>
        </DialogContentText>
        <Grid item xs={12} sx={{ pt: 2 }}>
          <FormControl variant="outlined" className={classes.formControl}>
            <InputLabel
              shrink
              variant="outlined"
              required
              htmlFor="outlined-age-always-notched"
            >
              SIMS
            </InputLabel>

            <Select
              variant="outlined"
              onChange={handleselectSim}
              value={sim}
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
              {sims &&
                sims.map((sim) => (
                  <MenuItem
                    key={sim.id}
                    value={sim.id}
                    sx={{
                      m: 0,
                      p: '6px 16px',
                    }}
                  >
                    <Typography variant="body1">{`${sim.iccid}`}</Typography>
                  </MenuItem>
                ))}
            </Select>
          </FormControl>
        </Grid>
        <Grid item xs={12} sx={{ pt: 2 }}>
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
              {packages &&
                packages.map((pkg) => (
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
      </DialogContent>
      <DialogActions>
        <Button
          onClick={handleClose}
          color="primary"
          autoFocus
          size="medium"
          disabled={loadingTopUp}
        >
          Cancel
        </Button>
        <Button
          variant="contained"
          onClick={handleClick}
          disabled={loadingTopUp || loading}
        >
          TOP UP
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default TopUpData;
