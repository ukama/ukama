import React, { useState } from 'react';
import {
  Checkbox,
  FormControlLabel,
  Slider,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { IMaskInput } from 'react-imask';

interface CustomProps {
  // eslint-disable-next-line no-unused-vars
  onChange: (event: { target: { name: string; value: string } }) => void;
}

const TextMaskCustom = React.forwardRef<HTMLInputElement, CustomProps>(
  function TextMaskCustom(props, _ref) {
    const { onChange, ...other } = props;
    return (
      <IMaskInput
        {...other}
        overwrite
        unmask={false}
        mask={'0#'}
        placeholder={'1'}
        definitions={{
          '#': /[0]/,
        }}
        onAccept={(value: any) => onChange({ target: { name: 'name', value } })}
      />
    );
  },
);

const CustomizePref = () => {
  const [data, setDate] = useState('2');
  const [sliderValue, setSliderValue] = useState(20);
  const [isSendAlerts, setIsSendAlerts] = useState(true);
  return (
    <Stack mt={1}>
      <Typography variant="subtitle2">Alert me when my data reaches</Typography>
      <Stack direction={'row'} spacing={2} alignItems="center">
        <Slider
          min={0}
          step={10}
          max={100}
          sx={{ height: 2 }}
          value={sliderValue}
          onChange={(_, value) => {
            setSliderValue(value as number);
            setDate(`${(value as number) / 10}`);
          }}
        />
        <TextField
          value={data}
          variant="outlined"
          sx={{
            width: 128,
            fontSize: 16,
            padding: '0px 8px',
            '.MuiOutlinedInput-input': {
              padding: '8px 14px',
            },
          }}
          InputLabelProps={{ shrink: false }}
          InputProps={{
            inputComponent: TextMaskCustom as any,
            endAdornment: (
              <Typography variant="body1" fontSize={16}>
                GB
              </Typography>
            ),
          }}
          onChange={(e) => {
            setDate(e.target.value);
            setSliderValue(Number(e.target.value) * 10);
          }}
        />
      </Stack>
      <FormControlLabel
        control={
          <Checkbox
            defaultChecked
            value={isSendAlerts}
            onChange={(e) => {
              setIsSendAlerts(e.target.checked);
            }}
          />
        }
        label="Send alert to email as well"
      />
    </Stack>
  );
};

export default CustomizePref;
