import { commonData, snackbarMessage } from '@/app-recoil';
import {
  Draft,
  Site,
  useAddDraftMutation,
  useAddSiteMutation,
  useDeleteDraftMutation,
  useGetDraftsQuery,
  useUpdateDraftNameMutation,
  useUpdateLocationMutation,
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
import { calculateCenterLatLng, formatSecondsToDuration } from '@/utils';
import { AlertColor, Popover, Stack, Typography } from '@mui/material';
import { LatLngLiteral } from 'leaflet';
import { useState } from 'react';
import { useRecoilValue, useSetRecoilState } from 'recoil';

const ZOOM = 3;
const SITE_INIT = {
  id: '',
  name: '',
  height: 0,
  solarUptime: 95,
  apOption: 'ONE_TO_ONE',
  isSetlite: true,
  location: {
    id: '',
    lat: '',
    lng: '',
    address: '',
    lastSaved: 0,
  },
};

const getMarkers = (sites: Site[]) => {
  return sites.map((site) => ({
    lat: parseFloat(site.location.lat),
    lng: parseFloat(site.location.lng),
  }));
};

const getLastSavedInt = () => Math.floor(new Date().getTime() / 1000);

const Page = () => {
  const [zoom, setZoom] = useState<number>(ZOOM);
  const [selectedDraft, setSelectedDraft] = useState<Draft | undefined>(
    undefined,
  );
  const [search, setSearch] = useState('');
  const [addSite, setAddSite] = useState(false);
  const [addLink, setAddLink] = useState(false);
  const [togglePower, setTogglePower] = useState(false);
  const _commonData = useRecoilValue<TCommonData>(commonData);
  const setSnackbarMessage = useSetRecoilState<TSnackMessage>(snackbarMessage);
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
      if (data.getDrafts.length > 0) {
        if (!selectedDraft) {
          setSelectedDraft(data.getDrafts[0]);
        } else {
          setSelectedDraft(
            data.getDrafts.find((d) => d.id === selectedDraft?.id),
          );
        }
      } else {
        setSelectedDraft(undefined);
      }
    },
    onError: (error) => {
      showAlert('get-drafts-error', error.message, 'error', true);
    },
  });

  const [addDraftCall, { loading: addDraftLoading }] = useAddDraftMutation({
    onCompleted: (data) => {
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
      onCompleted: () => {
        refetchDrafts();
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
        refetchDrafts();
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
  const [addSiteCall, { loading: addSiteLoading }] = useAddSiteMutation({
    onCompleted: () => {
      refetchDrafts();
      showAlert(
        'add-site-success',
        'Site updated successfully',
        'success',
        true,
      );
    },
    onError: (error) => {
      showAlert('add-site-error', error.message, 'error', true);
    },
  });
  const [updateLocationCall, { loading: updateLocationLoading }] =
    useUpdateLocationMutation({
      onCompleted: () => {
        refetchDrafts();
      },
      onError: (error) => {
        showAlert('update-site-location-error', error.message, 'error', true);
      },
    });

  const [deleteDraftCall, { loading: deleteDraftLoading }] =
    useDeleteDraftMutation({
      onCompleted: () => {
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

  const handleMarkerDrag = (e: LatLngLiteral, id: string) => {
    updateLocationCall({
      variables: {
        locationId: id,
        draftId: selectedDraft?.id || '',
        data: {
          address: '',
          lastSaved: getLastSavedInt(),
          lat: e.lat.toFixed(9).toString(),
          lng: e.lng.toFixed(9).toString(),
        },
      },
    });
  };

  const handleMarkerAdd = (e: LatLngLiteral) => {
    if (addSite) {
      setAddSite(false);
      addSiteCall({
        variables: {
          draftId: selectedDraft?.id || '',
          data: {
            siteName: SITE_INIT.name,
            lastSaved: getLastSavedInt(),
            apOption: SITE_INIT.apOption,
            isSetlite: SITE_INIT.isSetlite,
            lat: e.lat.toFixed(9).toString(),
            lng: e.lng.toFixed(9).toString(),
            address: SITE_INIT.location.address,
            height: parseFloat(SITE_INIT.height.toString()),
            solarUptime: parseFloat(SITE_INIT.solarUptime.toString()),
          },
        },
      });
    }
  };

  const handleSiteAction = (s: Site) => {
    updateSiteCall({
      variables: {
        siteId: s.id,
        draftId: selectedDraft?.id || '',
        data: {
          siteName: s.name,
          apOption: s.apOption,
          isSetlite: s.isSetlite,
          lastSaved: getLastSavedInt(),
          address: s.location.address,
          lat: s.location.lat.toString(),
          lng: s.location.lng.toString(),
          height: parseFloat(s.height.toString()),
          solarUptime: parseFloat(s.solarUptime.toString()),
        },
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
          lastSaved: getLastSavedInt(),
        },
      },
    });
  };

  const handleDraftSelected = (draftId: string) => {
    const d = getDraftsData?.getDrafts.find(({ id }) => id === draftId);
    setSelectedDraft(d);
  };

  const handleDraftUpdated = (id: string, draft: string) => {
    updateDraftCall({
      variables: {
        name: draft,
        draftId: id,
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
          subtitle={selectedDraft?.sites.length || 0}
          siteSummary={selectedDraft?.sites || []}
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
        {selectedDraft && selectedDraft?.sites.length > 0 && (
          <Typography variant="caption" sx={{ color: colors.black54 }}>
            {`Saved ${formatSecondsToDuration(
              Math.floor(new Date().getTime() / 1000) -
                selectedDraft?.lastSaved || 0,
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
            width={800}
            height={418}
            setZoom={setZoom}
            isAddSite={addSite}
            id={'site-planning-map'}
            className={styles.homeMap}
            handleAction={handleSiteAction}
            data={selectedDraft?.sites || []}
            handleAddMarker={handleMarkerAdd}
            handleDragMarker={handleMarkerDrag}
            zoom={
              selectedDraft && selectedDraft?.sites.length > 0 ? zoom : ZOOM
            }
            center={calculateCenterLatLng(
              getMarkers(selectedDraft?.sites || []),
            )}
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
                  maxZoom={16}
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
