import { commonData, snackbarMessage } from '@/app-recoil';
import {
  Draft,
  Site,
  useAddDraftMutation,
  useDeleteDraftMutation,
  useGetDraftsQuery,
  useUpdateDraftNameMutation,
  useUpdateEventMutation,
  useUpdateSiteLocationMutation,
  useUpdateSiteMutation,
} from '@/generated/planning-tool';
import styles from '@/styles/Site_Planning.module.css';
import { PageContainer } from '@/styles/global';
import { colors } from '@/styles/theme';
import { TCommonData, TSnackMessage } from '@/types';
import DraftDropdown from '@/ui/molecules/DraftDropdown';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import Map from '@/ui/molecules/Map';
import {
  LeftOverlayUI,
  RightOverlayUI,
  SiteSummary,
} from '@/ui/molecules/MapOverlayUI';
import { formatSecondsToDuration } from '@/utils';
import { AlertColor, Popover, Stack, Typography } from '@mui/material';
import { LatLngLiteral } from 'leaflet';
import { useEffect, useState } from 'react';
import { useRecoilValue, useSetRecoilState } from 'recoil';

const ZOOM = 4;
const MATKER_INIT = { lat: 0, lng: 0 };
const DEFAULT_CENTER = { lat: 37.7780627, lng: -121.9822475 };
const SITE_INIT = {
  name: '',
  height: 0,
  solarUptime: 95,
  apOption: 'ONE_TO_ONE',
  isSetlite: true,
  location: {
    lat: '',
    lng: '',
    address: '',
  },
};

const Page = () => {
  const [zoom, setZoom] = useState<number>(ZOOM);
  const [site, setSite] = useState<Site>(SITE_INIT);
  const [selectedDraft, setSelectedDraft] = useState<Draft | undefined>();
  const [search, setSearch] = useState('');
  const [addSite, setAddSite] = useState(false);
  const [addLink, setAddLink] = useState(false);
  const [togglePower, setTogglePower] = useState(false);
  const _commonData = useRecoilValue<TCommonData>(commonData);
  const setSnackbarMessage = useSetRecoilState<TSnackMessage>(snackbarMessage);
  const [marker, setMarker] = useState<LatLngLiteral>(MATKER_INIT);
  const [anchorSiteInfo, setAnchorSiteInfo] =
    useState<HTMLButtonElement | null>(null);
  const showAlert = (
    id: string,
    message: string,
    type: AlertColor,
    show: boolean,
  ) =>
    setSnackbarMessage({
      id,
      message,
      type,
      show,
    });
  const {
    data: getDraftsData,
    loading: getDraftsLoading,
    refetch: refetchDrafts,
  } = useGetDraftsQuery({
    fetchPolicy: 'network-only',
    variables: {
      userId: _commonData.userId,
    },
    onCompleted: (data) => {
      if (!selectedDraft && data.getDrafts.length > 0) {
        setSelectedDraft(
          data.getDrafts.length > 0 ? data.getDrafts[0] : undefined,
        );
        setSite({ ...data.getDrafts[0].site });
        if (
          data.getDrafts[0].site.location.lat &&
          data.getDrafts[0].site.location.lng
        ) {
          setMarker({
            lat: parseFloat(data.getDrafts[0].site.location.lat),
            lng: parseFloat(data.getDrafts[0].site.location.lng),
          });
        }
      } else if (data.getDrafts.length === 0) {
        setSelectedDraft(undefined);
        setSite(SITE_INIT);
        setMarker(MATKER_INIT);
      }
    },
    onError: (error) => {
      showAlert('get-drafts-error', error.message, 'error', true);
    },
  });

  const [addDraftCall, { loading: addDraftLoading }] = useAddDraftMutation({
    onCompleted: (data) => {
      setSelectedDraft({
        ...data.addDraft,
      });
      setSite({
        ...data.addDraft.site,
      });
      refetchDrafts();
      showAlert(
        'add-drafts-success',
        'Draft added successfully.',
        'success',
        true,
      );
    },
    onError: (error) => {
      showAlert('add-drafts-error', error.message, 'error', true);
    },
  });
  const [updateDraftCall, { loading: updateDraftLoading }] =
    useUpdateDraftNameMutation({
      onCompleted: (data) => {
        if (data.updateDraftName.id === selectedDraft?.id) {
          setSelectedDraft({
            ...selectedDraft,
            name: data.updateDraftName.name,
          });
        } else {
          refetchDrafts();
        }
        showAlert(
          'update-drafts-success',
          'Draft updated successfully',
          'success',
          true,
        );
      },
      onError: (error) => {
        showAlert('update-drafts-error', error.message, 'error', true);
      },
    });
  const [updateSiteCall, { loading: updateSiteLoading }] =
    useUpdateSiteMutation({
      onCompleted: () => {
        showAlert(
          'update-site-success',
          'Site updated successfully',
          'success',
          true,
        );
      },
      onError: (error) => {
        showAlert('update-site-error', error.message, 'error', true);
      },
    });
  const [updateSiteLocationCall, { loading: updateSiteLocationLoading }] =
    useUpdateSiteLocationMutation({
      onCompleted: (data) => {
        setSelectedDraft({
          ...data?.updateSiteLocation,
        });
      },
      onError: (error) => {
        showAlert('update-site-location-error', error.message, 'error', true);
      },
    });
  const [updateEventCall, { loading: updateEventLoading }] =
    useUpdateEventMutation({
      onCompleted: (data) => {
        /* Show success message */
      },
      onError: (error) => {
        showAlert('update-event-error', error.message, 'error', true);
      },
    });

  const [deleteDraftCall, { loading: deleteDraftLoading }] =
    useDeleteDraftMutation({
      onCompleted: () => {
        setSelectedDraft(undefined);
        setSite(SITE_INIT);
        refetchDrafts();
        showAlert(
          'delte-drafts-success',
          'Draft  deleted successfully.',
          'success',
          true,
        );
      },
      onError: (error) => {
        showAlert('delete-drafts-error', error.message, 'error', true);
      },
    });

  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorSiteInfo(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorSiteInfo(null);
  };

  const open = Boolean(anchorSiteInfo);
  const id = open ? 'site-info-popover' : undefined;

  useEffect(() => {
    if (selectedDraft?.id) {
      const loc = {
        lat: marker.lat.toFixed(10).toString(),
        lng: marker.lng.toFixed(10).toString(),
        address: selectedDraft.site.location.address,
      };
      setSelectedDraft({
        ...selectedDraft,
        site: {
          ...selectedDraft.site,
          location: loc,
        },
      });
      setSite({
        ...site,
        location: loc,
      });
    }
  }, [marker]);

  const handleMarkerDrag = (e: LatLngLiteral) => {
    setMarker({ lat: e.lat, lng: e.lng });
    updateSiteLocationCall({
      variables: {
        draftId: selectedDraft?.id || '',
        data: {
          address: '',
          lat: e.lat.toString(),
          lng: e.lng.toString(),
        },
      },
    });
  };

  const handleMarkerAdd = (e: LatLngLiteral) => {
    if (addSite && marker.lat === 0 && marker.lng === 0) {
      setAddSite(false);
      setMarker(e);
    }
  };

  const handleSiteAction = () => {
    if (selectedDraft?.id)
      updateSiteCall({
        variables: {
          data: {
            siteName: site.name,
            apOption: site.apOption,
            isSetlite: site.isSetlite,
            address: site.location.address,
            lat: site.location.lat.toString(),
            lng: site.location.lng.toString(),
            height: parseFloat(site.height.toString()),
            lastSaved: Math.floor(new Date().getTime() / 1000),
            solarUptime: parseFloat(site.solarUptime.toString()),
          },
          draftId: selectedDraft?.id,
        },
      });
  };

  const handleAddSite = () => setAddSite(true);
  const handleAddLink = () => setAddLink(true);
  const handleOnOff = () => setTogglePower(!togglePower);
  const handleAddDraft = () => {
    addDraftCall({
      variables: {
        data: {
          name: 'New Draft',
          userId: _commonData.userId,
          lastSaved: Math.floor(new Date().getTime() / 1000),
        },
      },
    });
  };

  const handleDraftSelected = (draftId: string) => {
    const newDraft = getDraftsData?.getDrafts.find(({ id }) => id === draftId);
    setSelectedDraft(newDraft);
    setSite({
      ...(newDraft?.site || SITE_INIT),
    });
    const { lat, lng } = newDraft?.site.location || MATKER_INIT;
    setMarker({
      lat: parseFloat(lat.toString()),
      lng: parseFloat(lng.toString()),
    });
  };

  const handleDraftUpdated = (id: string, draft: string) => {
    updateDraftCall({
      variables: {
        name: draft,
        updateDraftNameId: id,
      },
    });
  };

  const handleDeleteDraft = (id: string) =>
    deleteDraftCall({
      variables: {
        draftId: id,
      },
    });

  return (
    <>
      <Popover
        id={id}
        open={open}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
        onClose={handleClose}
        anchorEl={anchorSiteInfo}
        sx={{ top: 4, left: -40 }}
        PaperProps={{
          sx: {
            width: '204px',
            padding: '16px 24px',
          },
        }}
      >
        <SiteSummary
          title={'Site summary'}
          subtitle={''}
          siteSummary={[
            {
              id: 1,
              name: 'Site Name 1',
              status: 'up',
            },
            {
              id: 2,
              name: 'Site Name 2',
              status: 'down',
            },
            {
              id: 3,
              name: 'Site Name 3',
              status: 'unknown',
            },
          ]}
        />
      </Popover>
      <Stack
        direction="row"
        alignItems={'center'}
        justifyContent={'space-between'}
      >
        <DraftDropdown
          draft={selectedDraft}
          isLoading={getDraftsLoading}
          handleAddDraft={handleAddDraft}
          handleDeleteDraft={handleDeleteDraft}
          drafts={getDraftsData?.getDrafts || []}
          handleDraftUpdated={handleDraftUpdated}
          handleDraftSelected={handleDraftSelected}
        />
        {selectedDraft && selectedDraft?.lastSaved > 0 && (
          <Typography variant="caption" sx={{ color: colors.black54 }}>
            {`Saved ${formatSecondsToDuration(
              Math.floor(new Date().getTime() / 1000) -
                selectedDraft?.lastSaved,
            )} ago.`}
          </Typography>
        )}
      </Stack>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={false}
        cstyle={{
          backgroundColor: false ? colors.white : 'transparent',
        }}
      >
        <PageContainer sx={{ padding: 0, mt: '12px' }}>
          <Map
            data={site}
            width={800}
            height={418}
            marker={marker}
            setData={setSite}
            setZoom={setZoom}
            isAddSite={addSite}
            id={'site-planning-map'}
            className={styles.homeMap}
            handleAction={handleSiteAction}
            handleAddMarker={handleMarkerAdd}
            handleDragMarker={handleMarkerDrag}
            zoom={marker.lat === 0 && marker.lng === 0 ? ZOOM : zoom}
            center={
              marker.lat === 0 && marker.lng === 0 ? DEFAULT_CENTER : marker
            }
          >
            {({ TileLayer }: any) => (
              <>
                <LeftOverlayUI
                  search={search}
                  setSearch={setSearch}
                  handleAddSite={handleAddSite}
                  handleAddLink={handleAddLink}
                  isCurrentDraft={!!selectedDraft}
                />
                <RightOverlayUI
                  id={id}
                  handleClick={handleClick}
                  handleTogglePower={handleOnOff}
                  isCurrentDraft={!!selectedDraft}
                />

                <TileLayer
                  maxZoom={19}
                  tileSize={270}
                  url="https://tiles.stadiamaps.com/tiles/alidade_smooth/{z}/{x}/{y}{r}.png"
                />
              </>
            )}
          </Map>
        </PageContainer>
      </LoadingWrapper>
    </>
  );
};

export default Page;
