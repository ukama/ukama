/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { colors } from '@/theme';
import EditIcon from '@mui/icons-material/Edit';
import { IconButton, InputAdornment, TextField } from '@mui/material';
import { useState } from 'react';

type EditableTextFieldProps = {
  label: string;
  type?: string;
  value: string;
  editable?: boolean;
  handleOnChange?: (value: string) => void;
};

const EditableTextField = ({
  label,
  value,
  type = 'text',
  editable = true,
  handleOnChange = () => {},
}: EditableTextFieldProps) => {
  const [isEditable, setIsEditable] = useState<boolean>(false);
  return (
    <TextField
      fullWidth
      id={label}
      name={label}
      label={label}
      value={value}
      variant="standard"
      disabled={!editable}
      sx={{
        width: '300px',
        '& input': {
          color: colors.primaryMain,
        },
        '& input:disabled': {
          color: colors.black54,
          WebkitTextFillColor: colors.black54,
        },
      }}
      InputLabelProps={{ shrink: true }}
      onChange={(e) => handleOnChange(e.target.value)}
      inputRef={(input) => isEditable && input?.focus()}
      InputProps={{
        type: type,
        disableUnderline: true,
        endAdornment: (
          <InputAdornment
            position="end"
            sx={{ display: isEditable ? 'flex' : 'none' }}
          >
            <IconButton
              edge="end"
              onClick={() => setIsEditable(!isEditable)}
              sx={{
                svg: {
                  path: {
                    fill: isEditable ? colors.primaryMain : '-moz-initial',
                  },
                },
              }}
            >
              <EditIcon />
            </IconButton>
          </InputAdornment>
        ),
      }}
    />
  );
};

export default EditableTextField;
