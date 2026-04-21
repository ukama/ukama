/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import {
    Button,
    ClickAwayListener,
    MenuList,
    Paper,
    Popper,
    Stack,
    Switch,
    Typography
} from '@mui/material';
import React from 'react';
import BasicDialog from '../BasicDialog';

interface IToggleButtonsMenu {
  title: string;
  values: any[];
  options: any[];
  handleToggle: (id: string, value: boolean) => void;
}

const ToggleButtonsMenu = ({
  title,
  values,
  options,
  handleToggle,
}: IToggleButtonsMenu) => {
  const [open, setOpen] = React.useState(false);
  const anchorRef = React.useRef<HTMLButtonElement>(null);
  const [consentDialog, setConsentDialog] = React.useState(false);
  const [valuesState, setValuesState] = React.useState<any[]>(values);

  const handleOpen = () => {
    setOpen((open) => !open);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleToggleChange = (id: string, value: boolean) => {
    setValuesState(valuesState.map((option) => option.id === id ? { ...option, value } : option));
    handleToggle(id, valuesState.find((i: any) => i.id === id)?.value);
  };

  return (
    <>
      <Stack
        spacing={1}
        direction="row"
        sx={{
          px: 2,
          py: 1,
          pr: 1,
          borderRadius: '4px',
          bgcolor: 'primary.main',
        }}
      >
        <Typography
          variant="subtitle1"
          sx={{
            color: 'white',
            fontWeight: 500,
          }}
        >
          {title}
        </Typography>
        <Button
          ref={anchorRef}
          size="small"
          aria-haspopup="menu"
          aria-expanded="true"
          aria-controls="toggle-buttons-menu"
          aria-label="toggle buttons menu"
          onClick={handleOpen}
        >
          <ArrowDropDownIcon sx={{ color: 'white' }} />
        </Button>
      </Stack>
      <Popper
        open={open}
        role={undefined}
        placement='bottom-end'
        anchorEl={anchorRef.current}
        sx={{
          zIndex: 1000,
        }}
      >
        <Paper>
          <ClickAwayListener onClickAway={handleClose}>
            <MenuList id="node-action-menu">
              {options.map(({ id, name }) => (
                <Stack
                  key={id}
                  direction="row"
                  alignItems="center"
                  justifyContent="space-between"
                  sx={{ px: 2, py: 1, minWidth: 220 }}
                >
                  <Typography variant="body2">{name}</Typography>
                  <Switch
                    name={id}
                    size="small"
                    checked={valuesState.find((i: any) => i.id === id)?.value}
                    onChange={(e) => handleToggleChange(id, e.target.checked)}
                  />
                </Stack>
              ))}
            </MenuList>
          </ClickAwayListener>
        </Paper>
      </Popper>
      <BasicDialog
        isOpen={consentDialog}
        labelSuccessBtn={'Confirm'}
        labelNegativeBtn={'Cancel'}
        title={options.find((i) => i.id === values.find((j: any) => j.id === i.id)?.id)?.name}
        description={options.find((i) => i.id === values.find((j: any) => j.id === i.id)?.id)?.consent}
        handleCloseAction={() => setConsentDialog(false)}
        handleSuccessAction={() => {
          setConsentDialog(false);
        }}
      ></BasicDialog>
    </>
  );
};

export default ToggleButtonsMenu;
