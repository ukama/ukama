import { colors } from '@/theme';
import EditIcon from '@mui/icons-material/Edit';
import { IconButton, InputAdornment, TextField } from '@mui/material';
import { useState } from 'react';

type EditableTextFieldProps = {
  label: string;
  type?: string;
  value: string;
  isEditable?: boolean;
  handleOnChange?: Function;
};

const EditableTextField = ({
  label,
  value,
  type = 'text',
  isEditable = true,
  handleOnChange = () => {},
}: EditableTextFieldProps) => {
  const [iseditable, setIsEditable] = useState(false);
  return (
    <TextField
      fullWidth
      id={label}
      name={label}
      label={label}
      value={value}
      variant="standard"
      disabled={!iseditable}
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
      inputRef={(input) => iseditable && input?.focus()}
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
              onClick={() => setIsEditable(!iseditable)}
              sx={{
                svg: {
                  path: {
                    fill: iseditable ? colors.primaryMain : '-moz-initial',
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
