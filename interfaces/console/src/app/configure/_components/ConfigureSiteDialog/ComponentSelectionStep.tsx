/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { AddSiteValidationSchema } from '@/helpers/formValidators';
import { GlobalInput } from '@/styles/global';
import { TSiteForm } from '@/types';
import { Button, MenuItem, Stack, TextField } from '@mui/material';
import { Form, Formik } from 'formik';

export interface Component {
  id: string;
  inventory_id: string;
  category: string;
  type: string;
  user_id: string;
  description: string;
  datasheet_url: string;
  images_url: string;
  part_number: string;
  manufacturer: string;
  managed: string;
  warranty: number;
  specification: string;
}

interface ComponentSelectionStepProps {
  initialValues: TSiteForm;
  switchComponents: Component[];
  powerComponents: Component[];
  backhaulComponents: Component[];
  accessComponents: Component[];
  onComplete: (values: Partial<TSiteForm>) => void;
  onCancel: () => void;
}

/**
 * Step 0 of ConfigureSiteDialog.
 * Collects Switch, Power, Backhaul, Access, and Spectrum band selections.
 */
const ComponentSelectionStep: React.FC<ComponentSelectionStepProps> = ({
  initialValues,
  switchComponents,
  powerComponents,
  backhaulComponents,
  accessComponents,
  onComplete,
  onCancel,
}) => (
  <Formik
    initialValues={initialValues}
    onSubmit={onComplete}
    validationSchema={AddSiteValidationSchema[0]}
  >
    {({ errors, touched, isValid, handleReset }) => (
      <Form>
        <Stack>
          <GlobalInput
            as={TextField}
            fullWidth
            select
            required
            name="switch"
            label="SWITCH"
            margin="normal"
            slotProps={{ inputLabel: { shrink: true } }}
            error={touched.switch && Boolean(errors.switch)}
            helperText={touched.switch && errors.switch}
          >
            {switchComponents.map((c) => (
              <MenuItem key={c.id} value={c.id}>
                {c.description}
              </MenuItem>
            ))}
          </GlobalInput>

          <GlobalInput
            fullWidth
            select
            required
            name="power"
            label="POWER"
            margin="normal"
            slotProps={{ inputLabel: { shrink: true } }}
            error={touched.power && Boolean(errors.power)}
            helperText={touched.power && errors.power}
          >
            {powerComponents.map((c) => (
              <MenuItem key={c.id} value={c.id}>
                {c.description}
              </MenuItem>
            ))}
          </GlobalInput>

          <GlobalInput
            fullWidth
            select
            required
            name="backhaul"
            label="BACKHAUL"
            margin="normal"
            slotProps={{ inputLabel: { shrink: true } }}
            error={touched.backhaul && Boolean(errors.backhaul)}
            helperText={touched.backhaul && errors.backhaul}
          >
            {backhaulComponents.map((c) => (
              <MenuItem key={c.id} value={c.id}>
                {c.description}
              </MenuItem>
            ))}
          </GlobalInput>

          <GlobalInput
            fullWidth
            select
            required
            name="access"
            label="ACCESS"
            margin="normal"
            slotProps={{ inputLabel: { shrink: true } }}
            error={touched.access && Boolean(errors.access)}
            helperText={touched.access && errors.access}
          >
            {accessComponents.map((c) => (
              <MenuItem key={c.id} value={c.id}>
                {c.description}
              </MenuItem>
            ))}
          </GlobalInput>

          <GlobalInput
            fullWidth
            select
            required
            name="spectrum"
            label="SPECTRUM BAND"
            margin="normal"
            slotProps={{ inputLabel: { shrink: true } }}
            error={touched.spectrum && Boolean(errors.spectrum)}
            helperText={touched.spectrum && errors.spectrum}
          >
            {/* Spectrum uses access components per existing business logic */}
            {accessComponents.map((c) => (
              <MenuItem key={c.id} value={c.id}>
                {c.description}
              </MenuItem>
            ))}
          </GlobalInput>
        </Stack>

        <Stack
          direction="row"
          spacing={1}
          justifyContent="flex-end"
          sx={{ mt: 1 }}
        >
          <Button onClick={() => { handleReset(); onCancel(); }}>Cancel</Button>
          <Button type="submit" variant="contained" color="primary" disabled={!isValid}>
            Next
          </Button>
        </Stack>
      </Form>
    )}
  </Formik>
);

export default ComponentSelectionStep;
