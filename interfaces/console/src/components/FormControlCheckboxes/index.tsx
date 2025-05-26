/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

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
  handleChange: (name: string, checked: boolean) => void;
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
