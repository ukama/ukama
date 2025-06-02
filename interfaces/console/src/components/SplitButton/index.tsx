/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import Button from '@mui/material/Button';
import ButtonGroup from '@mui/material/ButtonGroup';
import ClickAwayListener from '@mui/material/ClickAwayListener';
import Grow from '@mui/material/Grow';
import MenuItem from '@mui/material/MenuItem';
import MenuList from '@mui/material/MenuList';
import Paper from '@mui/material/Paper';
import Popper from '@mui/material/Popper';
import * as React from 'react';
import BasicDialog from '../BasicDialog';

type splitButtonProps = {
  options: any[];
  handleSplitActionClick: (id: string) => void;
};
const SplitButton = ({ options, handleSplitActionClick }: splitButtonProps) => {
  const isHaveOptions = options.length > 1;
  const [open, setOpen] = React.useState(false);
  const anchorRef = React.useRef<HTMLDivElement>(null);
  const [consentDialog, setConsentDialog] = React.useState(false);
  const [selectedIndex, setSelectedIndex] = React.useState(
    options[0]?.id || 'node-restart',
  );

  const handleOptionSelected = (
    event: React.MouseEvent<HTMLLIElement, MouseEvent>,
    id: string,
  ) => {
    setSelectedIndex(id);
    setOpen(false);
  };

  const handleToggle = () => {
    setOpen((prevOpen) => !prevOpen);
  };

  const handleClose = (event: Event) => {
    if (
      anchorRef.current &&
      anchorRef.current.contains(event.target as HTMLElement)
    ) {
      return;
    }

    setOpen(false);
  };

  return (
    <>
      <ButtonGroup
        variant="contained"
        ref={anchorRef}
        aria-label="split button"
        sx={{
          whiteSpace: 'nowrap',
          width: isHaveOptions ? '180px' : 'fit-content',
        }}
      >
        <Button fullWidth onClick={() => setConsentDialog(true)}>
          {options.find((i) => i.id === selectedIndex).name}
        </Button>
        {isHaveOptions && (
          <Button
            size="small"
            aria-controls={open ? 'split-button-menu' : undefined}
            aria-expanded={open ? 'true' : undefined}
            aria-label="select merge strategy"
            aria-haspopup="menu"
            onClick={handleToggle}
          >
            <ArrowDropDownIcon />
          </Button>
        )}
      </ButtonGroup>
      <Popper
        open={open}
        anchorEl={anchorRef.current}
        role={undefined}
        transition
      >
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            style={{
              transformOrigin:
                placement === 'bottom' ? 'center top' : 'center bottom',
            }}
          >
            <Paper>
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList id="node-action-menu">
                  {options.map(({ id, name }) => (
                    <MenuItem
                      key={id}
                      selected={id === selectedIndex}
                      onClick={(event) => handleOptionSelected(event, id)}
                    >
                      {name}
                    </MenuItem>
                  ))}
                </MenuList>
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
      <BasicDialog
        isOpen={consentDialog}
        labelSuccessBtn={'Confirm'}
        labelNegativeBtn={'Cancel'}
        title={options.find((i) => i.id === selectedIndex).name}
        description={options.find((i) => i.id === selectedIndex).consent}
        handleCloseAction={() => setConsentDialog(false)}
        handleSuccessAction={() => {
          setConsentDialog(false);
          handleSplitActionClick(
            options.find((i) => i.id === selectedIndex).id,
          );
        }}
      ></BasicDialog>
    </>
  );
};
export default SplitButton;
