'use client';
import { Typography } from '@mui/material';

// TODO: Handle unauthorized access here
/**
 * Where user auth session is value valid but user token request failed
 */

const Unauthorized = () => {
  return <Typography variant="h1">Unauthorized</Typography>;
};

export default Unauthorized;
