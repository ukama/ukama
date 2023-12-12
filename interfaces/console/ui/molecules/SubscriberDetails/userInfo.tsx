import React from 'react';
import { Typography, TextField, IconButton, Stack } from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import { colors } from '@/styles/theme';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';

interface UserInfoProps {
  subscriberLoading: boolean;
  onEditName: boolean;
  firstName: string;
  handleEditName: () => void;
  onEditEmail: boolean;
  email: string;
  handleSimEdit: () => void;
  setOnEditName: any;
  setOnEditEmail: any;
}
const UserInfo: React.FC<UserInfoProps> = ({
  subscriberLoading,
  onEditName,
  firstName,
  handleEditName,
  onEditEmail,
  email,
  handleSimEdit,
  setOnEditName,
  setOnEditEmail,
}) => (
  <LoadingWrapper
    radius="small"
    width={'100%'}
    isLoading={subscriberLoading}
    cstyle={{
      overflow: 'auto',
      backgroundColor: false ? colors.white : 'transparent',
    }}
  >
    <Stack direction="column" spacing={2}>
      <Typography variant="body1" sx={{ color: colors.black }}>
        Name
      </Typography>
      <Stack
        direction="row"
        spacing={2}
        alignItems={'center'}
        justifyContent={'space-between'}
        sx={{ pr: 2 }}
      >
        <TextField
          id="outlined-basic"
          value={firstName}
          variant="standard"
          disabled={!onEditName}
          size="small"
          onChange={handleEditName}
          sx={{ width: '100%' }}
        />

        <IconButton
          size="small"
          color="primary"
          onClick={() => setOnEditName(!onEditName)}
        >
          <EditIcon fontSize="small" />
        </IconButton>
      </Stack>

      <Typography variant="body1" sx={{ color: colors.black }}>
        Email
      </Typography>
      <Stack
        direction="row"
        spacing={2}
        alignItems={'center'}
        justifyContent={'space-between'}
        sx={{ pr: 2 }}
      >
        <TextField
          id="outlined-basic"
          value={email}
          variant="standard"
          disabled={!onEditEmail}
          size="small"
          onChange={handleSimEdit}
          sx={{ width: '100%' }}
        />

        <IconButton
          size="small"
          color="primary"
          onClick={() => setOnEditEmail(!onEditEmail)}
        >
          <EditIcon fontSize="small" />
        </IconButton>
      </Stack>
    </Stack>
  </LoadingWrapper>
);

export default UserInfo;
