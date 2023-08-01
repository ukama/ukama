import { colors } from '@/styles/theme';
import EditIcon from '@mui/icons-material/Edit';
import { IconButton, InputAdornment, TextField } from '@mui/material';
import { useState } from 'react';

type EditableTextFieldProps = {
  label: string;
  type?: string;
  value: any;
  isEditable?: boolean;
  handleOnChange?: Function;
};

const EditableTextField = ({
  label,
  value = '-djkdjflkdjlk',
  type = 'text',
  isEditable = true,
  // eslint-disable-next-line no-empty-function
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
      sx={{ width: { xs: '100%' } }}
      InputLabelProps={{
        shrink: true,
      }}
      onChange={(e) => handleOnChange(e.target.value)}
      inputRef={(input) => iseditable && input?.focus()}
      InputProps={{
        type: type,
        disableUnderline: true,
        color: 'primary',
        endAdornment: (
          <InputAdornment
            position="end"
            sx={{
              display: isEditable ? 'flex' : 'none',
            }}
          >
            <IconButton
              edge="end"
              onClick={() => setIsEditable(!iseditable)}
              sx={{
                svg: {
                  path: {
                    fill: `${iseditable ? colors.primaryMain : colors.silver}`,
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
