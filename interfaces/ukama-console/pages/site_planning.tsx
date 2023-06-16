import { commonData, snackbarMessage } from '@/app-recoil';
import {
  Draft,
  useAddDraftMutation,
  useGetDraftsQuery,
  useUpdateDraftNameMutation,
  useUpdateEventMutation,
  useUpdateSiteMutation,
} from '@/generated/planning-tool';
import styles from '@/styles/Site_Planning.module.css';
import { PageContainer } from '@/styles/global';
import { colors } from '@/styles/theme';
import { TCommonData, TSite, TSnackMessage } from '@/types';
import SitePopup from '@/ui/SitePopup';
import DraftDropdown from '@/ui/molecules/DraftDropdown';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import Map from '@/ui/molecules/Map';
import {
  LeftOverlayUI,
  RightOverlayUI,
  SiteSummary,
} from '@/ui/molecules/MapOverlayUI';
import { AlertColor, Popover, Stack } from '@mui/material';
import { useState } from 'react';
import { useRecoilValue, useSetRecoilState } from 'recoil';

const DEFAULT_CENTER = [38.907132, -77.036546];

const Page = () => {
  const [site, setSite] = useState<TSite>({
    id: '',
    name: '',
    height: 0,
    solarUptime: 95,
    ap: 'ONE_TO_ONE',
    isBackhaul: true,
    location: {
      lat: '',
      lng: '',
      address: '',
    },
  });
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
    variables: {
      userId: _commonData.userId,
    },
    onCompleted: (data) => {
      if (data.getDrafts.length > 0) setSelectedDraft(data.getDrafts[0]);
    },
    onError: (error) => {
      showAlert('get-drafts-error', error.message, 'error', true);
    },
  });

  const [addDraftCall, { loading: addDraftLoading }] = useAddDraftMutation({
    onCompleted: () => {
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

  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorSiteInfo(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorSiteInfo(null);
  };

  const open = Boolean(anchorSiteInfo);
  const id = open ? 'site-info-popover' : undefined;

  const handleMarkerDrag = (e: any) => {
    setMarker(e.latlng);
    if (selectedDraft)
      setSelectedDraft({
        ...selectedDraft,
        site: {
          ...selectedDraft.site,
          location: {
            lat: e.latlng.lat,
            lng: e.latlng.lng,
            address: site.location.address,
          },
        },
      });
  };

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
            apOption: site.ap,
            siteName: site.name,
            height: site.height,
            lat: site.location.lat,
            lng: site.location.lng,
            isSetlite: site.isBackhaul,
            solarUptime: site.solarUptime,
            address: site.location.address,
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
      <Stack direction="row" justifyContent={'space-between'}>
        <DraftDropdown
          draft={selectedDraft}
          isLoading={getDraftsLoading}
          handleAddDraft={handleAddDraft}
          drafts={getDraftsData?.getDrafts || []}
          handleDraftUpdated={handleDraftUpdated}
          handleDraftSelected={handleDraftSelected}
        />
        {/* <Typography variant='caption' sx={{ color: colors.grey }}>
          {getDraftsData?.getDrafts[]} */}
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
                <Marker draggable position={marker} ondrag={handleMarkerDrag}>
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
