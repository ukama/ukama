/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import * as Yup from 'yup';

// Validation rules
const nameValidationRule = Yup.string().required('Name is required.');
const networkNameValidationRule = Yup.string().required(
  'Network name is required.',
);
const emailValidatonRule = Yup.string()
  .required('Email is required.')
  .email('Please enter a valid email.');
const iccidValidator = Yup.string()
  .notRequired()
  .nullable()
  .matches(/^[0-9]+$/, 'Must be only digits')
  .min(18, 'Iccid must be 18 digits')
  .max(18, 'Iccid must be 18 digits');
const securitycodeValidator = Yup.string().required(
  'Security code is required.',
);

// Stepper form validation schema
const STEPPER_FORM_SCHEMA = [
  Yup.object().shape({
    switch: Yup.string().required('Switch is required'),
    power: Yup.string().required('Power is required'),
    backhaul: Yup.string().required('Backhaul is required'),
    access: Yup.string().required('Spectrum Band is required'),
  }),
  Yup.object().shape({
    siteName: Yup.string().required('Site Name is required'),
    network: Yup.string().required('Network is required'),
    latitude: Yup.number()
      .required('Latitude is required')
      .min(-90, 'Latitude must be between -90 and 90')
      .max(90, 'Latitude must be between -90 and 90'),
    longitude: Yup.number()
      .required('Longitude is required')
      .min(-180, 'Longitude must be between -180 and 180')
      .max(180, 'Longitude must be between -180 and 180'),
  }),
];
const UpdateSiteSchema = Yup.object().shape({
  siteName: Yup.string()
    .required('Site name is required')
    .matches(
      /^[a-z0-9-]*$/,
      'Site name must be lowercase alphanumeric and should not contain spaces, "-" are allowed.',
    ),
});

// Validation schemas
const ESIM_FORM_SCHEMA = Yup.object().shape({
  email: emailValidatonRule,
  name: nameValidationRule,
  simiccid: iccidValidator,
});
const NETWORK_NAME_SCHEMA_VALIDATOR = Yup.object().shape({
  name: networkNameValidationRule,
});
const PHYSICAL_SIM_FORM_SCHEMA = Yup.object().shape({
  iccid: iccidValidator,
  securityCode: securitycodeValidator,
});

const AddSiteValidationSchema = [
  Yup.object().shape({
    switch: Yup.string().required('Switch is required'),
    power: Yup.string().required('Power is required'),
    backhaul: Yup.string().required('Backhaul is required'),
    access: Yup.string().required('Access is required'),
    spectrum: Yup.string().required('Spectrum is required'),
  }),
  Yup.object().shape({
    siteName: Yup.string()
      .required('Site name is required')
      .matches(
        /^[a-z0-9-]*$/,
        'Site name must be lowercase alphanumeric and should not contain spaces, "-" are allowed.',
      ),
    network: Yup.string().required('Network is required'),
    latitude: Yup.number()
      .required('Latitude is required')
      .min(-90, 'Invalid latitude')
      .max(90, 'Invalid latitude'),
    longitude: Yup.number()
      .required('Longitude is required')
      .min(-180, 'Invalid longitude')
      .max(180, 'Invalid longitude'),
  }),
];

export {
  AddSiteValidationSchema,
  ESIM_FORM_SCHEMA,
  NETWORK_NAME_SCHEMA_VALIDATOR,
  PHYSICAL_SIM_FORM_SCHEMA,
  STEPPER_FORM_SCHEMA,
  UpdateSiteSchema,
};
