import { commonData, snackbarMessage } from '@/app-recoil';
import {
  Draft,
  Site,
  useAddDraftMutation,
  useAddLinkMutation,
  useAddSiteMutation,
  useDeleteDraftMutation,
  useDeleteSiteMutation,
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
  PlanningSummary,
  RightOverlayUI,
} from '@/ui/molecules/MapOverlayUI';
import { calculateCenterLatLng, formatSecondsToDuration } from '@/utils';
import { AlertColor, Popover, Stack, Typography } from '@mui/material';
import { LatLngLiteral } from 'leaflet';
import { useEffect, useState } from 'react';
import { useRecoilValue, useSetRecoilState } from 'recoil';

const POWER_SUMMARY = {
  sites: [
    {
      id: 'site-id-1',
      name: 'Site 1',
      status: 'up',
      usage: 10,
      panels: 1,
      battries: 2,
    },
    {
      id: 'site-id-2',
      name: 'Site 2',
      status: 'up',
      usage: 20,
      panels: 1,
      battries: 4,
    },
    {
      id: 'site-id-3',
      name: 'Site 3',
      status: 'up',
      usage: 30,
      panels: 1,
      battries: 3,
    },
    {
      id: 'site-id-4',
      name: 'Site 4',
      status: 'up',
      usage: 5,
      panels: 1,
      battries: 2,
    },
  ],
};

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

const INIT_LINK = {
  siteA: '',
  siteB: '',
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
  const [center, setCenter] = useState<LatLngLiteral>({
    lat: 37.7780627,
    lng: -121.9822475,
  });
  const [mapInteraction, setMapInteraction] = useState({
    isAddLink: false,
    isAddSite: false,
  });
  const [linkSites, setLinkSites] = useState(INIT_LINK);
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
      setSelectedDraft(data.addDraft);
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
      onCompleted: () => refetchDrafts(),
      onError: (error) => {
        showAlert('update-site-location-error', error.message, 'error', true);
      },
    });

  const [addLinkCall, { loading: addlinkLoading }] = useAddLinkMutation({
    onCompleted: () => {
      setLinkSites(INIT_LINK);
      refetchDrafts();
      showAlert('add-link-success', `Link created`, 'success', true);
    },
    onError: (error) => {
      showAlert('add-link-error', error.message, 'error', true);
    },
  });

  const [deleteDraftCall, { loading: deleteDraftLoading }] =
    useDeleteDraftMutation({
      onCompleted: () => {
        setSelectedDraft(undefined);
        refetchDrafts();
        showAlert(
          'delte-drafts-success',
          'Draft deleted successfully.',
          'success',
          true,
        );
      },
      onError: (error) => {
        showAlert('delete-drafts-error', error.message, 'error', true);
      },
    });

  const [deleteSiteCall, { loading: deleteSiteLoading }] =
    useDeleteSiteMutation({
      onCompleted: () => {
        refetchDrafts();
        showAlert(
          'delte-site-success',
          'Site deleted successfully.',
          'success',
          true,
        );
      },
      onError: (error) => {
        showAlert('delete-site-error', error.message, 'error', true);
      },
    });

  useEffect(() => {
    if (selectedDraft) {
      setCenter(
        calculateCenterLatLng(
          getMarkers(selectedDraft.sites || []),
        ) as LatLngLiteral,
      );
    } else {
      setCenter({ lat: 37.7780627, lng: -121.9822475 } as LatLngLiteral);
    }
  }, [selectedDraft]);

  useEffect(() => {
    if (linkSites.siteA && linkSites.siteB) {
      addLinkCall({
        variables: {
          draftId: selectedDraft?.id || '',
          data: {
            lastSaved: getLastSavedInt(),
            siteA: linkSites.siteA,
            siteB: linkSites.siteB,
          },
        },
      });
    }
  }, [linkSites]);

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

  const handleMarkerAdd = (e: LatLngLiteral, id: string) => {
    if (mapInteraction.isAddSite) {
      setMapInteraction({ ...mapInteraction, isAddSite: false });
      addSiteCall({
        variables: {
          draftId: selectedDraft?.id || '',
          data: {
            locationId: id,
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
          locationId: s.location.id,
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

  const handleDeleteSite = (id: string) =>
    deleteSiteCall({
      variables: {
        siteId: id,
      },
    });

  const handleLocationSelected = (loc: LatLngLiteral) => {
    setZoom(6);
    setCenter(loc);
  };

  const handleAddSite = () =>
    setMapInteraction({
      isAddLink: false,
      isAddSite: !mapInteraction.isAddSite,
    });
  const handleAddLink = () =>
    setMapInteraction({
      isAddLink: !mapInteraction.isAddLink,
      isAddSite: false,
    });
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

  const handleAddLinkToSite = (siteId: string) => {
    if (mapInteraction.isAddLink) {
      if (!linkSites.siteA) {
        setLinkSites({ ...linkSites, siteA: siteId });
      } else if (!linkSites.siteB) {
        setLinkSites({ ...linkSites, siteB: siteId });
        setMapInteraction({ ...mapInteraction, isAddLink: false });
      }
    }
  };
  console.log(mapInteraction);
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
            width: '220px',
            padding: '16px 24px',
          },
        }}
      >
        <PlanningSummary
          subtitleOne={`Sites`}
          subtitleTwo={`Power`}
          siteSummary={selectedDraft?.sites || []}
          powerSummary={POWER_SUMMARY}
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
            zoom={zoom}
            height={418}
            center={center}
            setZoom={setZoom}
            id={'site-planning-map'}
            linkSites={linkSites}
            className={styles.homeMap}
            handleAction={handleSiteAction}
            data={selectedDraft?.sites || []}
            handleAddMarker={handleMarkerAdd}
            links={selectedDraft?.links || []}
            handleDeleteSite={handleDeleteSite}
            handleDragMarker={handleMarkerDrag}
            isAddSite={mapInteraction.isAddSite}
            isAddLink={mapInteraction.isAddLink}
            handleAddLinkToSite={handleAddLinkToSite}
          >
            {() => (
              <>
                <LeftOverlayUI
                  handleAddSite={handleAddSite}
                  handleAddLink={handleAddLink}
                  isCurrentDraft={!!selectedDraft}
                  isAddSite={mapInteraction.isAddSite}
                  isAddLink={mapInteraction.isAddLink}
                  handleLocationSelected={handleLocationSelected}
                />
                <RightOverlayUI
                  id={id}
                  handleClick={handleClick}
                  handleTogglePower={handleOnOff}
                  isCurrentDraft={!!selectedDraft}
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
