import { Draft } from '@/generated/planning-tool';
import { colors } from '@/styles/theme';
import { hexToRGB } from '@/utils';
import { Edit } from '@mui/icons-material';
import AddIcon from '@mui/icons-material/Add';
import DeleteOutlineIcon from '@mui/icons-material/DeleteOutlineOutlined';
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
  draft: Draft | undefined;
  drafts: Draft[];
  isLoading: boolean;
  handleAddDraft: Function;
  handleDeleteDraft: Function;
  handleDraftUpdated: Function;
  handleDraftSelected: Function;
}

const DraftDropdown = ({
  draft,
  drafts = [],
  isLoading,
  handleAddDraft,
  handleDeleteDraft,
  handleDraftUpdated,
  handleDraftSelected,
}: IDraftDropdown) => {
  const { id = '', name = '' } = draft || {};
  const [newName, setNewName] = useState('');
  const [editNameId, setEditNameId] = useState('-1');
  const classes = useStyles({ isEmpty: !id });
  const handleNameUpdate = (i: string, n: string) => {
    setNewName('');
    setEditNameId('-1');
    handleDraftUpdated(i, n);
  };
  return (
    <LoadingWrapper isLoading={isLoading} width={'200px'} height={'32px'}>
      <Select
        disableUnderline
        variant="standard"
        value={name}
        onChange={(e) => {
          e.stopPropagation();
          handleDraftSelected(e.target.value);
        }}
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
        displayEmpty
        className={classes.selectStyle}
        renderValue={(selected) => (selected ? selected : 'Add/Select a draft')}
        onClose={() => {
          setEditNameId('-1');
          setNewName('');
        }}
      >
        {drafts.map(({ id: _id, name: _name }: any) => (
          <MenuItem
            key={_id}
            value={_id}
            sx={{
              m: 0,
              p: '6px 16px',
              justifyContent: 'space-between',
              backgroundColor: `${
                id === _id ? hexToRGB(colors.secondaryLight, 0.25) : 'inherit'
              } !important`,
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

              {editNameId !== _id ? (
                <Typography variant="body1" ml={1}>
                  {_name}
                </Typography>
              ) : (
                <TextField
                  autoFocus
                  variant="standard"
                  onClick={(e) => e.stopPropagation()}
                  value={editNameId === _id ? newName : _name}
                  onChange={(event) => {
                    event.stopPropagation();
                    setNewName(event.target.value);
                  }}
                />
              )}
            </Stack>
            <Stack direction={'row'}>
              <IconButton
                onClick={(e) => {
                  e.stopPropagation();
                  if (editNameId !== '-1')
                    handleNameUpdate(editNameId, newName);
                  else {
                    setEditNameId(_id);
                    setNewName(_name);
                  }
                }}
              >
                {editNameId === _id ? (
                  <DoneIcon color="success" sx={ICON_STYLE} />
                ) : (
                  <Edit color="action" sx={ICON_STYLE} />
                )}
              </IconButton>
              <IconButton
                onClick={() => {
                  handleDeleteDraft(id);
                }}
              >
                <DeleteOutlineIcon color="error" sx={ICON_STYLE} />
              </IconButton>
            </Stack>
          </MenuItem>
        ))}
        {drafts.length > 0 && <Divider />}
        <MenuItem
          onClick={(e) => {
            e.stopPropagation();
            handleAddDraft();
          }}
        >
          <AddIcon color="action" sx={ICON_STYLE} />
          <Typography variant="body2" ml={1}>
            Add new draft
          </Typography>
        </MenuItem>
      </Select>
    </LoadingWrapper>
  );
};
export default DraftDropdown;
