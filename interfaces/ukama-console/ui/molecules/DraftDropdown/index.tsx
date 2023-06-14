import { colors } from '@/styles/theme';
import { hexToRGB } from '@/utils';
import { Edit } from '@mui/icons-material';
import AddIcon from '@mui/icons-material/Add';
import DoneIcon from '@mui/icons-material/Done';
import DotIcon from '@mui/icons-material/FiberManualRecord';
import {
  Divider,
  IconButton,
  MenuItem,
  Select,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { useState } from 'react';
import {
  ICON_STYLE,
  PaperProps,
  SelectDisplayProps,
  useStyles,
} from './styles';

const DraftDropdown = () => {
  const [newName, setNewName] = useState('');
  const [editName, setEditName] = useState(false);
  const classes = useStyles();
  const handleNameUpdate = (n: string) => {
    setNewName('');
  };
  return (
    <Select
      value={'Draft 1'}
      disableUnderline
      variant="standard"
      onChange={() => {}}
      SelectDisplayProps={SelectDisplayProps}
      MenuProps={{
        disablePortal: true,
        anchorOrigin: {
          vertical: 'bottom',
          horizontal: 'left',
        },
        transformOrigin: {
          vertical: 'top',
          horizontal: 'left',
        },
        PaperProps: {
          sx: {
            ...PaperProps,
          },
        },
      }}
      className={classes.selectStyle}
      renderValue={(selected) => selected}
    >
      {[{ id: 1, name: 'Draft 1' }].map(({ id, name }) => (
        <MenuItem
          key={id}
          value={name}
          sx={{
            m: 0,
            p: '6px 16px',
            justifyContent: 'space-between',
            ':hover': {
              backgroundColor: `${hexToRGB(
                colors.secondaryLight,
                0.25,
              )} !important`,
            },
          }}
        >
          <Stack direction={'row'} alignItems={'center'}>
            <DotIcon color="success" sx={ICON_STYLE} />

            {!editName ? (
              <Typography variant="body1">{name}</Typography>
            ) : (
              <TextField
                autoFocus
                variant="standard"
                value={newName}
                onClick={(e) => e.stopPropagation()}
                onChange={(event) => {
                  event.stopPropagation();
                  setNewName(event.target.value);
                }}
              />
            )}
          </Stack>
          <IconButton
            onClick={(e) => {
              e.stopPropagation();
              setEditName(!editName);
              if (editName) handleNameUpdate(newName);
              else setNewName(name);
            }}
          >
            {editName ? (
              <DoneIcon color="success" sx={ICON_STYLE} />
            ) : (
              <Edit color="action" sx={ICON_STYLE} />
            )}
          </IconButton>
        </MenuItem>
      ))}
      <Divider />
      <MenuItem
        onClick={(e) => {
          e.stopPropagation();
        }}
      >
        <AddIcon color="action" sx={ICON_STYLE} />
        <Typography variant="body1">Add new draft</Typography>
      </MenuItem>
    </Select>
  );
};
export default DraftDropdown;
