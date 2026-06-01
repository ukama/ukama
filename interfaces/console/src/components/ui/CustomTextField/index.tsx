/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { GlobalInput } from '@/styles/global';
import { useField } from 'formik';
import React from 'react';

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

  return (
    <GlobalInput
      {...field}
      label={label}
      fullWidth
      onChange={(e) => {
        field.onChange(e);
        if (onChange) {
          onChange(e);
        }
      }}
      helperText={meta.touched && meta.error ? meta.error : ''}
      error={meta.touched && Boolean(meta.error)}
      spellCheck={false}
      slotProps={{
        inputLabel: {
          shrink: true,
        },
      }}
    />
  );
};

export default CustomTextField;
