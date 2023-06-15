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
import LoadingWrapper from '../LoadingWrapper';
import {
  ICON_STYLE,
  PaperProps,
  SelectDisplayProps,
  useStyles,
} from './styles';

interface IDraftDropdown {
  drafts: any;
  isLoading: boolean;
  handleAddDraft: Function;
  handleDraftUpdated: Function;
  handleDraftSelected: Function;
  currentDraft: { id: string; name: string };
}

const DraftDropdown = ({
  drafts,
  isLoading,
  currentDraft,
  handleAddDraft,
  handleDraftUpdated,
  handleDraftSelected,
}: IDraftDropdown) => {
  const [newName, setNewName] = useState('');
  const [editName, setEditName] = useState(false);
  const classes = useStyles();
  const handleNameUpdate = (n: string) => {
    setNewName('');
    handleDraftUpdated(currentDraft.id, n);
  };
  return (
    <LoadingWrapper isLoading={isLoading} width={'200px'} height={'32px'}>
      <Select
        disableUnderline
        variant="standard"
        value={currentDraft.name}
        onChange={(e) => handleDraftSelected(e.target.value)}
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
        onClose={() => {
          setEditName(false);
          setNewName('');
        }}
      >
        {drafts.map(({ id, name }: any) => (
          <MenuItem
            key={id}
            value={id}
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
            handleAddDraft();
          }}
        >
          <AddIcon color="action" sx={ICON_STYLE} />
          <Typography variant="body1">Add new draft</Typography>
        </MenuItem>
      </Select>
    </LoadingWrapper>
  );
};
export default DraftDropdown;
