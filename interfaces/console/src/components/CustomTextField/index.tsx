/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
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
