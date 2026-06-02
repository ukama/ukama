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
  Skeleton,
  Stack,
  Switch,
  Typography,
} from '@mui/material';
import React from 'react';
import BasicDialog from '@/components/ui/BasicDialog';

export interface ToggleOption {
  id: string;
  name: string;
  consent?: string;
}

export interface ToggleValue {
  id: string;
  value: boolean;
}

interface IToggleButtonsMenu {
  title: string;
  values: ToggleValue[];
  options: ToggleOption[];
  isLoading: boolean;
  handleToggle: (id: string, value: boolean) => void;
}

const ToggleButtonsMenu = ({
  title,
  values,
  options,
  isLoading,
  handleToggle,
}: IToggleButtonsMenu) => {
  const [open, setOpen] = React.useState(false);
  const anchorRef = React.useRef<HTMLButtonElement>(null);
  const [consentDialog, setConsentDialog] = React.useState(false);

  const handleOpen = () => {
    setOpen((open) => !open);
  };

  const handleClose = () => {
    setOpen(false);
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
        placement="bottom-end"
        // eslint-disable-next-line react-hooks/refs
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
                  sx={{ px: 2, py: 1, width: 220 }}
                >
                  <Typography variant="body2">{name}</Typography>

                  {isLoading ? (
                    <Skeleton
                      width={40}
                      height={20}
                      variant="rounded"
                      sx={{ borderRadius: '25%' }}
                    />
                  ) : (
                    <Switch
                      name={id}
                      size="small"
                      onChange={(e) => handleToggle(id, e.target.checked)}
                      checked={values.find((j) => j.id === id)?.value}
                    />
                  )}
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
        title={
          options.find(
            (i) => i.id === values.find((j) => j.id === i.id)?.id,
          )?.name ?? ''
        }
        description={
          options.find(
            (i) => i.id === values.find((j) => j.id === i.id)?.id,
          )?.consent ?? ''
        }
        handleCloseAction={() => setConsentDialog(false)}
        handleSuccessAction={() => {
          setConsentDialog(false);
        }}
      />
    </>
  );
};

export default ToggleButtonsMenu;
