import { useGetOrgsQuery } from '@/generated';
import {
  Box,
  Button,
  Card,
  Chip,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from '@mui/material';
import { format } from 'date-fns';

const OwnerOf = () => {
  const data = [
    {
      certificate: 'ukama',
      createdAt: '2023-12-05T09:15:13.364059Z',
      id: 'fbc2a80d-339e-4d3c-acaa-329dd3d1b221',
      isDeactivated: false,
      name: 'ukama-test',
      owner: '6d890d67-9850-41b6-a680-eb201dfd748f',
    },
    {
      certificate: 'ukama',
      createdAt: '2023-12-05T09:15:13.364059Z',
      id: 'fbc2a80d-339e-4d3c-acaa-329dd3d1b222',
      isDeactivated: false,
      name: 'ukama-test',
      owner: '6d890d67-9850-41b6-a680-eb201dfd718f',
    },
    {
      certificate: 'ukama',
      createdAt: '2023-12-05T09:15:13.364059Z',
      id: 'fbc2a80d-339e-4d3c-acaa-329dd3d1b223',
      isDeactivated: false,
      name: 'ukama-test',
      owner: '6d890d67-9850-41b6-a680-eb201dfd248f',
    },
    {
      certificate: 'ukama',
      createdAt: '2023-12-05T09:15:13.364059Z',
      id: 'fbc2a80d-339e-4d3c-acaa-329dd3d1b224',
      isDeactivated: false,
      name: 'ukama-test',
      owner: '6d890d67-9850-41b6-a680-eb201dfd748f',
    },
    {
      certificate: 'ukama',
      createdAt: '2023-12-05T09:15:13.364059Z',
      id: 'fbc2a80d-339e-4d3c-acaa-329dd3d1b225',
      isDeactivated: false,
      name: 'ukama-test',
      owner: '6d890d67-9850-41b6-a680-eb201dfd748f',
    },
    {
      certificate: 'ukama',
      createdAt: '2023-12-05T09:15:13.364059Z',
      id: 'fbc2a80d-339e-4d3c-acaa-329dd3d1b226',
      isDeactivated: false,
      name: 'ukama-test',
      owner: '6d890d67-9850-41b6-a680-eb201dfd748f',
    },
    {
      certificate: 'ukama',
      createdAt: '2023-12-05T09:15:13.364059Z',
      id: 'fbc2a80d-339e-4d3c-acaa-329dd3d1b227',
      isDeactivated: false,
      name: 'ukama-test',
      owner: '6d890d67-9850-41b6-a680-eb201dfd748f',
    },
    {
      certificate: 'ukama',
      createdAt: '2023-12-05T09:15:13.364059Z',
      id: 'fbc2a80d-339e-4d3c-acaa-329dd3d1b228',
      isDeactivated: false,
      name: 'ukama-test',
      owner: '6d890d67-9850-41b6-a680-eb201dfd748f',
    },
  ];
  return (
    <Stack
      direction={'column'}
      sx={{
        height: '300px',
        overflow: 'scroll',
      }}
    >
      <TableContainer>
        <Table
          sx={{ minWidth: 650, overflow: 'scroll' }}
          aria-label="simple table"
        >
          <TableHead>
            <TableRow>
              <TableCell align="left">Name</TableCell>
              <TableCell align="left">Role</TableCell>
              <TableCell align="left">Is active</TableCell>
              <TableCell align="left">Created at</TableCell>
              <TableCell align="left">Action</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data.map((row) => (
              <TableRow
                key={row.id}
                sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
              >
                <TableCell component="th" scope="row" align="left">
                  {row.name}
                </TableCell>
                <TableCell align="left">
                  <Chip size="small" variant="outlined" label={'admin'} />
                </TableCell>
                <TableCell align="left">
                  {`${!row.isDeactivated}`.toUpperCase()}
                </TableCell>
                <TableCell align="left">
                  {format(new Date(row.createdAt), 'MMM dd, yyyy')}
                </TableCell>
                <TableCell align="left">
                  <Button variant="contained" size="small">
                    Select
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Stack>
  );
};

const OnBoarding = () => {
  const {
    data: orgsData,
    error: orgsError,
    loading: orgsLoading,
  } = useGetOrgsQuery({
    onCompleted: (data) => {},
  });

  return (
    <Box
      sx={{
        width: '100%',
        height: 'calc(100vh - 20vh)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}
    >
      <Card sx={{ width: 'fit-content', height: 'fit-content', p: 3 }}>
        <Stack direction={'column'} spacing={1}>
          <Typography variant="h6" textAlign={'start'} fontWeight={500} pb={1}>
            Organizations
          </Typography>
          <OwnerOf />
        </Stack>
      </Card>
    </Box>
  );
};

export default OnBoarding;
