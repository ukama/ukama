import { StatsItemType } from '@/types';
import {
  Checkbox,
  FormControl,
  FormControlLabel,
  FormGroup,
  FormLabel,
} from '@mui/material';
import React from 'react';

type FormControlCheckboxesProp = {
  values: any;
  handleChange: Function;
  checkboxList: StatsItemType[];
};

const FormControlCheckboxes = ({
  values,
  handleChange,
  checkboxList,
}: FormControlCheckboxesProp) => {
  return (
    <FormControl
      component="div"
      variant="standard"
      sx={{
        display: 'flex',
        alignItems: 'center',
        flexDirection: 'row',
      }}
      onChange={(e: any) => handleChange(e.target.name, e.target.checked)}
    >
      <FormLabel component="legend" sx={{ mr: '44px' }}>
        Server Down
      </FormLabel>
      <FormGroup row>
        {checkboxList.map(({ id, label, value }: StatsItemType) => (
          <FormControlLabel
            label={label}
            key={`${value}-${id}`}
            control={
              <Checkbox
                name={value}
                checked={values[value] ? values[value] : false}
              />
            }
          />
        ))}
      </FormGroup>
    </FormControl>
  );
};

export default React.memo(FormControlCheckboxes);
