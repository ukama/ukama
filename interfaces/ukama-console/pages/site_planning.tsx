import { commonData, snackbarMessage } from '@/app-recoil';
import {
  Draft,
  Site,
  useAddDraftMutation,
  useDeleteDraftMutation,
  useGetDraftsQuery,
  useUpdateDraftNameMutation,
  useUpdateEventMutation,
  useUpdateSiteMutation,
} from '@/generated/planning-tool';
import styles from '@/styles/Site_Planning.module.css';
import { PageContainer } from '@/styles/global';
import { colors } from '@/styles/theme';
import { TCommonData, TSnackMessage } from '@/types';
import SitePopup from '@/ui/SitePopup';
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
import { useEffect, useState } from 'react';
import { useRecoilValue, useSetRecoilState } from 'recoil';

const DEFAULT_CENTER = [38.907132, -77.036546];
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
  const [site, setSite] = useState<Site>(SITE_INIT);
  const [selectedDraft, setSelectedDraft] = useState<Draft | undefined>();
  const [search, setSearch] = useState('');
  const [addSite, setAddSite] = useState(false);
  const [addLink, setAddLink] = useState(false);
  const [togglePower, setTogglePower] = useState(false);
  const _commonData = useRecoilValue<TCommonData>(commonData);
  const setSnackbarMessage = useSetRecoilState<TSnackMessage>(snackbarMessage);
  const [marker, setMarker] = useState([0, 0]);
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
      setSelectedDraft(
        data.getDrafts.length > 0 ? data.getDrafts[0] : undefined,
      );
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
        'update-drafts-success',
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
      onCompleted: (data) => {
        /* Show success message */
      },
      onError: (error) => {
        showAlert('update-site-error', error.message, 'error', true);
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
        lat: marker[0].toFixed(10).toString(),
        lng: marker[1].toFixed(10).toString(),
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

  const handleMarkerDrag = (e: any) => setMarker([e.lat, e.lng]);

  const handleMarkerAdd = (e: number[]) => {
    if (
      addSite &&
      (marker.length === 0 || (marker[0] === 0 && marker[1] === 0))
    ) {
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
            lat: site.location.lat,
            lng: site.location.lng,
            apOption: site.apOption,
            isSetlite: site.isSetlite,
            solarUptime: site.solarUptime,
            address: site.location.address,
            height: parseFloat(site.height.toString()),
            lastSaved: Math.floor(new Date().getTime() / 1000),
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
    setSelectedDraft(getDraftsData?.getDrafts.find(({ id }) => id === draftId));
    setSite({
      ...site,
    });
    setMarker([0, 0]);
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
        <SiteSummary />
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
            id={'site-planning-map'}
            zoom={12}
            width={800}
            height={418}
            isAddSite={addSite}
            center={DEFAULT_CENTER}
            className={styles.homeMap}
            onMapClick={handleMarkerAdd}
          >
            {({ TileLayer, Marker, Popup }: any) => (
              <>
                <LeftOverlayUI
                  search={search}
                  setSearch={setSearch}
                  handleAddSite={handleAddSite}
                  handleAddLink={handleAddLink}
                  isCurrentDraft={selectedDraft?.id !== ''}
                />
                <RightOverlayUI
                  id={id}
                  handleClick={handleClick}
                  handleTogglePower={handleOnOff}
                  isCurrentDraft={selectedDraft?.id !== ''}
                />
                <TileLayer url="https://tiles.stadiamaps.com/tiles/alidade_smooth/{z}/{x}/{y}{r}.png" />
                <Marker
                  autoPan
                  draggable
                  position={marker}
                  ondrag={handleMarkerDrag}
                  eventHandlers={{
                    moveend: (event: any) =>
                      handleMarkerDrag(event.target.getLatLng()),
                  }}
                >
                  <Popup>
                    <SitePopup
                      data={site}
                      setData={setSite}
                      handleAction={handleSiteAction}
                    />
                  </Popup>
                </Marker>
              </>
            )}
          </Map>
        </PageContainer>
      </LoadingWrapper>
    </>
  );
};

export default Page;
