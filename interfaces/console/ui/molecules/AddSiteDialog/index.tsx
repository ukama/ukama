import CloseIcon from '@mui/icons-material/Close';
import { globalUseStyles } from '@/styles/global';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  MenuItem,
  Select,
  Stack,
  TextField,
  Typography,
  Grid,
} from '@mui/material';
import React, { useState } from 'react';
import { Form, Formik } from 'formik';
import * as Yup from 'yup';
import { NetworkDto } from '@/generated';

type AddSiteDialogProps = {
  title: string;
  isOpen: boolean;
  description: string;
  handleCloseAction: any;
  handleAddSite: Function;
  networks: NetworkDto[];
};

interface SiteFormValues {
  site: string;
}
const initialValues: SiteFormValues = {
  site: '',
};
const validationSchema = Yup.object({
  site: Yup.string().required('Site name is required'),
});

const AddSiteDialog: React.FC<AddSiteDialogProps> = ({
  title,
  isOpen,
  description,
  handleCloseAction,
  handleAddSite,
  networks,
}: AddSiteDialogProps) => {
  const [selectedNetwork, setSelectedNetwork] = useState('');
  const isAddButtonDisabled = !selectedNetwork;
  const gclasses = globalUseStyles();

  return (
    <Dialog fullWidth open={isOpen} maxWidth="sm" onClose={handleCloseAction}>
      <Stack direction="row" alignItems="center" justifyContent="space-between">
        <DialogTitle>{title}</DialogTitle>
        <IconButton
          onClick={handleCloseAction}
          sx={{ position: 'relative', right: 8 }}
        >
          <CloseIcon fontSize="small" />
        </IconButton>
      </Stack>

      <Formik
        initialValues={initialValues}
        validationSchema={validationSchema}
        //Loading...
        onSubmit={async (values) => {
          await handleAddSite({
            ...values,
            selectedNetwork,
          });
        }}
      >
        {({
          values,
          errors,
          touched,
          handleChange,
          handleSubmit,
          handleBlur,
        }) => (
          <>
            <Form onSubmit={handleSubmit}>
              <DialogContent>
                <Typography variant="body2">{description}</Typography>
                <Grid container spacing={1}>
                  <Grid xs={12} item>
                    <TextField
                      required
                      fullWidth
                      label={'SITE'}
                      name={'site'}
                      InputLabelProps={{
                        shrink: true,
                      }}
                      onBlur={handleBlur}
                      onChange={handleChange}
                      value={values.site}
                      helperText={touched.site && errors.site}
                      error={touched.site && Boolean(errors.site)}
                      id={'site'}
                      spellCheck={false}
                      InputProps={{
                        classes: {
                          input: gclasses.inputFieldStyle,
                        },
                      }}
                    />
                  </Grid>
                  <Grid item xs={12}>
                    <Select
                      fullWidth
                      label="Network"
                      value={selectedNetwork}
                      onChange={(e) =>
                        setSelectedNetwork(e.target.value as string)
                      }
                    >
                      {networks.map((network) => (
                        <MenuItem key={network.id} value={network.id}>
                          {network.name}
                        </MenuItem>
                      ))}
                    </Select>
                  </Grid>
                </Grid>
              </DialogContent>
              <DialogActions>
                <Stack
                  width={'100%'}
                  spacing={2}
                  direction={'row'}
                  alignItems="center"
                  justifyContent={'space-between'}
                >
                  <Button type="submit" variant="contained">
                    Add site
                  </Button>
                </Stack>
              </DialogActions>
            </Form>
          </>
        )}
      </Formik>
    </Dialog>
  );
};

export default AddSiteDialog;
