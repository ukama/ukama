import React from 'react';
import { TextField } from '@mui/material';
import { useField } from 'formik';
import { globalUseStyles } from '@/styles/global';

interface CustomTextFieldProps {
  label: string;
  name: string;
  onChange?: (event: any) => void;
}

const CustomTextField: React.FC<CustomTextFieldProps> = ({
  label,
  name,
  onChange,
}) => {
  const [field, meta] = useField(name);
  const gclasses = globalUseStyles();

  return (
    <TextField
      {...field}
      label={label}
      fullWidth
      onChange={(e) => {
        field.onChange(e);
        if (onChange) {
          onChange(e);
        }
      }}
      InputLabelProps={{
        shrink: true,
      }}
      helperText={meta.touched && meta.error ? meta.error : ''}
      error={meta.touched && Boolean(meta.error)}
      spellCheck={false}
      InputProps={{
        classes: {
          input: gclasses.inputFieldStyle,
        },
      }}
    />
  );
};

export default CustomTextField;
