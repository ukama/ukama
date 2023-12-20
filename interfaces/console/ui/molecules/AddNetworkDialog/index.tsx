import CloseIcon from '@mui/icons-material/Close';
import {
  Button,
  Dialog,
  DialogContent,
  DialogTitle,
  IconButton,
  Stack,
  TextField,
  Typography,
} from '@mui/material';

import { Formik } from 'formik';
import * as Yup from 'yup';

type AddNetworkDialogProps = {
  title: string;
  isOpen: boolean;
  loading: boolean;
  // networks: NetworkDto[];
  description: string;
  isClosable?: boolean;
  handleCloseAction: any;
  labelSuccessBtn?: string;
  handleSuccessAction?: any;
  labelNegativeBtn?: string;
};

interface AddNetworkForm {
  name: string;
  budget: number;
  countries: { name: string; code: string }[];
  networks: { id: string; name: string }[];
}

const validationSchema = Yup.object({
  networks: Yup.array().optional().default([]),
  countries: Yup.array().optional().default([]),
  name: Yup.string()
    .required('Network name is required')
    .matches(
      /^[a-z0-9\-]*$/,
      'Network name must be lowercase alphanumeric and should not contain spaces, "-" are allowed.',
    ),
  budget: Yup.number().default(0),
});

const initialValues: AddNetworkForm = {
  name: '',
  budget: 0,
  countries: [],
  networks: [],
};

const AddNetworkDialog = ({
  title,
  isOpen,
  loading,
  // networks,
  description,
  labelSuccessBtn,
  labelNegativeBtn,
  handleCloseAction,
  handleSuccessAction,
}: AddNetworkDialogProps) => {
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
          <CloseIcon fontSize="small" />
        </IconButton>
      </Stack>

      <DialogContent sx={{ maxHeight: '400px', overflow: 'auto' }}>
        <Formik
          initialValues={initialValues}
          validationSchema={validationSchema}
          onSubmit={async (values) => {
            handleSuccessAction(values);
          }}
        >
          {({
            values,
            errors,
            touched,
            handleChange,
            handleSubmit,
            handleBlur,
            setFieldValue,
          }) => (
            <form onSubmit={handleSubmit}>
              <Stack spacing={2} direction={'column'} alignItems="start">
                {description && (
                  <Typography variant="body1">{description}</Typography>
                )}
                <TextField
                  fullWidth
                  name={'name'}
                  size="medium"
                  placeholder="network-name"
                  label={'Network name'}
                  InputLabelProps={{
                    shrink: true,
                  }}
                  onBlur={handleBlur}
                  onChange={handleChange}
                  value={values.name}
                  helperText={touched.name && errors.name}
                  error={touched.name && Boolean(errors.name)}
                  id={'name'}
                />
                <Stack
                  width="100%"
                  spacing={2}
                  direction={'row'}
                  alignItems="center"
                  justifyContent={'flex-end'}
                >
                  <Button
                    variant="text"
                    color={'primary'}
                    onClick={handleCloseAction}
                  >
                    {labelNegativeBtn}
                  </Button>

                  <Button type="submit" variant="contained" disabled={loading}>
                    {labelSuccessBtn}
                  </Button>
                </Stack>
              </Stack>
            </form>
          )}
        </Formik>
      </DialogContent>
    </Dialog>
  );
};

export default AddNetworkDialog;

{
  /* <Stack direction={'row'} width={'100%'} spacing={2} mt={1}>
      <TextField
        fullWidth
        name={'name'}
        size="medium"
        placeholder="Mesh"
        label={'Network name'}
        InputLabelProps={{
          shrink: true,
        }}
        onBlur={handleBlur}
        onChange={handleChange}
        value={values.name}
        helperText={touched.name && errors.name}
        error={touched.name && Boolean(errors.name)}
        id={'name'}
      />
      <TextField
        sx={{ width: '50%' }}
        name={'budget'}
        id={'budget'}
        size="medium"
        type="number"
        placeholder="100"
        label={'Network budget'}
        InputLabelProps={{
          shrink: true,
        }}
        onBlur={handleBlur}
        value={values.budget}
        onChange={handleChange}
        helperText={touched.budget && Boolean(errors.budget)}
        error={touched.budget && Boolean(errors.budget)}
      />
    </Stack>
    <Autocomplete
      multiple
      options={COUNTRIES}
      getOptionLabel={(option) => option.name}
      renderInput={(params) => (
        <TextField
          {...params}
          label="Allowed Countries"
          placeholder="Country"
          InputLabelProps={{
            shrink: true,
          }}
        />
      )}
      sx={{
        width: '100%',
        mt: 1,
        p: 0,
        '.MuiOutlinedInput-root': {
          p: '4px',
        },
      }}
      id={'countries'}
      onBlur={handleBlur}
      onChange={(_, value: any) =>
        setFieldValue('countries', value)
      }
    /> 
    <Autocomplete
      multiple
      options={networks.length > 0 ? networks : []}
      getOptionLabel={(option: NetworkDto) =>
        option ? option.name : 'No network available'
      }
      renderInput={(params) => (
        <TextField
          {...params}
          label="Allowed Networks"
          placeholder="Network"
          InputLabelProps={{
            shrink: true,
          }}
          value={values.networks}
        />
      )}
      sx={{
        width: '100%',
        mt: 1,
        '.MuiOutlinedInput-root': {
          p: '4px',
        },
      }}
      id={'networks'}
      onBlur={handleBlur}
      onChange={(_, value: any) => setFieldValue('networks', value)}
    /> */
}
