import { isDarkmode } from '@/app-recoil';
import { NodeDto } from '@/generated';
import { colors } from '@/styles/theme';
import KeyboardArrowLeftIcon from '@mui/icons-material/KeyboardArrowLeft';
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight';
import { IconButton, Stack, Theme, useMediaQuery } from '@mui/material';
import { useEffect, useState } from 'react';
import Carousel from 'react-material-ui-carousel';
import { useRecoilValue } from 'recoil';
import { NodeCard } from '..';

interface INodeSlider {
  items: NodeDto[];
  handleItemAction: Function;
}

const NodeSlider = ({ items = [], handleItemAction }: INodeSlider) => {
  const _isDarkMod = useRecoilValue(isDarkmode);
  const small = useMediaQuery((theme: Theme) => theme.breakpoints.up('sm'));
  const medium = useMediaQuery((theme: Theme) => theme.breakpoints.up('md'));
  const [list, setList] = useState<any>([]);

  const NavButtonStyle = {
    p: 0.5,
    backgroundColor: _isDarkMod ? colors.nightGrey12 : colors.white,
    ':hover': {
      opacity: 1,
      backgroundColor: _isDarkMod ? colors.nightGrey12 : colors.white,
    },
  };

  useEffect(() => {
    const slides = [];
    const isSmall = small ? 2 : 1;
    const chunk = medium ? 3 : isSmall;
    for (let i = 0; i < items.length; i += chunk) {
      slides.push({
        cid: `chunk-${i}`,
        item: items.slice(i, i + chunk),
      });
    }
    setList(slides);
  }, [items, small, medium]);

  return (
    <Carousel
      swipe={true}
      animation="slide"
      autoPlay={false}
      indicators={false}
      cycleNavigation={false}
      navButtonsAlwaysVisible
      sx={{
        width: '100%',
        minHeight: '240px',
        pt: 3,
        pb: 0,
        div: {
          ':hover button': {
            backgroundColor: _isDarkMod ? colors.nightGrey12 : colors.white,
            opacity: '1 !important',
          },
        },
      }}
      NextIcon={
        <IconButton sx={NavButtonStyle}>
          <KeyboardArrowRightIcon />
        </IconButton>
      }
      PrevIcon={
        <IconButton sx={NavButtonStyle}>
          <KeyboardArrowLeftIcon />
        </IconButton>
      }
      navButtonsProps={{
        style: {
          margin: 0,
          padding: 0,
          height: 'fit-content',
          boxShadow:
            '0px 3px 1px -2px rgba(0, 0, 0, 0.1), 0px 2px 2px rgba(0, 0, 0, 0.07), 0px 1px 5px rgba(0, 0, 0, 0.06)',
        },
      }}
    >
      {list.map(({ cid, item }: any) => (
        <Stack
          key={cid}
          spacing={4}
          direction={'row'}
          sx={{
            justifyContent: {
              xs: 'center',
              md: items.length > 1 ? 'center' : 'flex-start',
            },
          }}
        >
          {item.map(
            (
              {
                id,
                type,
                name,
                description,
                updateShortNote,
                isUpdateAvailable,
              }: any,
              i: number,
            ) => (
              <NodeCard
                key={i}
                id={id}
                users={3}
                type={type}
                title={name}
                loading={false}
                subTitle={description}
                updateShortNote={updateShortNote}
                isUpdateAvailable={isUpdateAvailable}
                handleOptionItemClick={(type: string) =>
                  handleItemAction(id, type)
                }
              />
            ),
          )}
        </Stack>
      ))}
    </Carousel>
  );
};

export default NodeSlider;
