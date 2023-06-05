import { DATA_DURATION, DATA_UNIT } from '@/constants';
import CloseIcon from '@mui/icons-material/Close';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControl,
  Grid,
  IconButton,
  InputLabel,
  MenuItem,
  Select,
  Stack,
  TextField,
} from '@mui/material';
import { useState } from 'react';

type AddDataPlanDialogProps = {
  title: string;
  isOpen: boolean;
  handleCloseAction: any;
  labelSuccessBtn?: string;
  handleSuccessAction: any;
  labelNegativeBtn?: string;
};

const AddDataPlanDialog = ({
  title,
  isOpen,
  labelSuccessBtn,
  labelNegativeBtn,
  handleCloseAction,
  handleSuccessAction,
}: AddDataPlanDialogProps) => {
  const [dataplan, setDataplan] = useState({
    name: '',
    dataVolume: 0,
    dataUnit: '',
    amount: 0,
    duration: 0,
  });

  return (
    <Dialog
      fullWidth
      open={isOpen}
      maxWidth="sm"
      onClose={handleCloseAction}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <Stack direction="row" alignItems="center" justifyContent="space-between">
        <DialogTitle>{title}</DialogTitle>
        <IconButton
          onClick={handleCloseAction}
          sx={{ position: 'relative', right: 8 }}
        >
          <CloseIcon />
        </IconButton>
      </Stack>

      <DialogContent>
        <Grid
          container
          rowSpacing={2}
          gridAutoRows={2}
          columnSpacing={2}
          gridAutoColumns={1}
          alignItems={'center'}
          justifyContent={'center'}
        >
          <Grid item xs={12}>
            <TextField
              fullWidth
              required
              label="DATA PLAN NAME"
              value={dataplan.name}
              id={'data-plan-name'}
              InputLabelProps={{
                shrink: true,
              }}
              onChange={(e) =>
                setDataplan({ ...dataplan, name: e.target.value })
              }
            />
          </Grid>
          <Grid item container xs={12} sm={6} columnSpacing={1} rowSpacing={2}>
            <Grid item xs={6}>
              <TextField
                fullWidth
                required
                type="number"
                label="DATA LIMIT"
                value={dataplan.dataVolume}
                id={'data-plan-limit'}
                InputLabelProps={{
                  shrink: true,
                }}
                onChange={(e) =>
                  setDataplan({
                    ...dataplan,
                    dataVolume: parseInt(e.target.value as string),
                  })
                }
              />
            </Grid>
            <Grid item xs={6}>
              <TextField
                fullWidth
                required
                type="number"
                label="PACKAGE PRICE"
                value={dataplan.amount}
                id={'data-plan-price'}
                InputLabelProps={{
                  shrink: true,
                }}
                onChange={(e) =>
                  setDataplan({
                    ...dataplan,
                    amount: parseInt(e.target.value as string),
                  })
                }
              />
            </Grid>
          </Grid>
          <Grid item container xs={12} sm={6} columnSpacing={1} rowSpacing={2}>
            <Grid item xs={4}>
              <FormControl fullWidth>
                <InputLabel id={'data-plan-unit-label'} shrink>
                  Unit
                </InputLabel>
                <Select
                  notched
                  required
                  label="UNIT"
                  value={dataplan.dataUnit}
                  id={'data-plan-unit'}
                  labelId="data-plan-unit-label"
                  onChange={(e) =>
                    setDataplan({
                      ...dataplan,
                      dataUnit: e.target.value as string,
                    })
                  }
                >
                  {DATA_UNIT.map(({ id, label, value }) => (
                    <MenuItem key={id} value={value}>
                      {label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={8}>
              <FormControl fullWidth>
                <InputLabel
                  id={'data-plan-unit-price-label'}
                  shrink
                  sx={{
                    '& legend': {
                      letterSpacing: 0.6,
                    },
                  }}
                >
                  PRICE PER
                </InputLabel>
                <Select
                  notched
                  required
                  label="UNIT"
                  value={dataplan.duration}
                  id={'data-plan-unit'}
                  labelId="data-plan-unit-price-label"
                  onChange={(e) =>
                    setDataplan({
                      ...dataplan,
                      duration: parseInt(e.target.value as string),
                    })
                  }
                >
                  {DATA_DURATION.map(({ id, label, value }) => (
                    <MenuItem key={id} value={value}>
                      {label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
          </Grid>
        </Grid>
      </DialogContent>

      <DialogActions>
        <Stack direction={'row'} alignItems="center" spacing={2}>
          {labelNegativeBtn && (
            <Button
              variant="text"
              color={'primary'}
              onClick={handleCloseAction}
            >
              {labelNegativeBtn}
            </Button>
          )}
          {labelSuccessBtn && (
            <Button
              variant="contained"
              onClick={() => handleSuccessAction(dataplan)}
            >
              {labelSuccessBtn}
            </Button>
          )}
        </Stack>
      </DialogActions>
    </Dialog>
  );
};

export default AddDataPlanDialog;
