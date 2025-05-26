/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import ContainerHeader from '@/components/ContainerHeader';
import LoadingWrapper from '@/components/LoadingWrapper';
import { Box, Card, Grid, Paper, Stack, Typography } from '@mui/material';
import Image from 'next/image';

type ISchematicsProps = {
  schematicsSpecsData?: any;
  getSearchValue: (value: string) => void;
  nodeTitle: string | undefined;
  loading: boolean;
};

const NodeSchematicTab = ({
  schematicsSpecsData,
  getSearchValue,
  nodeTitle,
  loading,
}: ISchematicsProps) => {
  return (
    <>
      <Stack direction="column" spacing={2}>
        <LoadingWrapper
          width="100%"
          height="400px"
          radius={'small'}
          isLoading={loading}
        >
          <Paper sx={{ p: 2, height: '100%' }}>
            <Grid container>
              <Grid item xs={12} container justifyContent="flex-start">
                <Typography variant="h6">{nodeTitle}</Typography>
              </Grid>
              <Grid item xs={12} container justifyContent="center">
                <Image
                  src="https://ukama-site-assets.s3.amazonaws.com/images/schematic.png"
                  alt="trx_schematic"
                  width={0}
                  height={0}
                  sizes="100vw"
                  style={{
                    width: '100%',
                    height: 'auto',
                    maxWidth: '720px',
                  }}
                />
              </Grid>
            </Grid>
          </Paper>
        </LoadingWrapper>
        <Paper sx={{ p: 2 }}>
          <Grid container>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <ContainerHeader
                  title="Resources"
                  showButton={false}
                  showSearchBox={true}
                  handleSearchChange={getSearchValue}
                />
              </Grid>
              {schematicsSpecsData.map(({ id, title, readingTime }: any) => (
                <Grid item key={id} xs md lg={4}>
                  <Card
                    variant="outlined"
                    sx={{
                      cursor: 'pointer',
                      borderRadius: '10px',
                      padding: '15px 18px 8px 18px',
                      ':hover': {
                        boxShadow: '0px 3px 5px 2px #0000001F',
                      },
                    }}
                  >
                    <Stack spacing={1} direction="column">
                      <Typography
                        variant="h6"
                        sx={{
                          fontSize: '16px',
                        }}
                      >
                        {title}
                      </Typography>
                      <Typography variant="caption">{readingTime}</Typography>
                      <Box
                        component="div"
                        sx={{
                          width: '100%',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                        }}
                      >
                        <Image
                          src="https://i.ibb.co/BgfbTsP/1835cf7a15bd359317e492f4ea67166a.png"
                          alt="1835cf7a15bd359317e492f4ea67166a"
                          width="300"
                          height="160"
                        />
                      </Box>
                    </Stack>
                  </Card>
                </Grid>
              ))}
            </Grid>
          </Grid>
        </Paper>
      </Stack>
    </>
  );
};

export default NodeSchematicTab;
