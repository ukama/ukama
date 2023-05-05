import { LoadingWrapper } from "..";
import { ContainerHeader } from "../../components";
import { Paper, Grid, Typography, Card, Stack, Box } from "@mui/material";

type ISchematicsProps = {
    schematicsSpecsData?: any;
    getSearchValue: Function;
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
                    height="100%"
                    radius={"small"}
                    isLoading={loading}
                >
                    <Paper sx={{ p: 2 }}>
                        <Grid container xs={12}>
                            <Grid
                                item
                                xs={12}
                                container
                                justifyContent="flex-start"
                            >
                                <Typography variant="h6">
                                    {nodeTitle}
                                </Typography>
                            </Grid>
                            <Grid
                                item
                                xs={12}
                                container
                                justifyContent="center"
                                sx={{ height: "300px" }}
                            >
                                <img
                                    src="https://i.ibb.co/d2cNd1d/Screen-Shot-2022-03-28-at-11-08-50.png"
                                    alt="1835cf7a15bd359317e492f4ea67166a"
                                    width="600"
                                    height="300"
                                />
                            </Grid>
                        </Grid>
                    </Paper>
                </LoadingWrapper>
                <Paper sx={{ p: 2 }}>
                    <Grid container xs={12}>
                        <Grid xs={12} container spacing={2}>
                            <Grid item xs={12}>
                                <ContainerHeader
                                    title="Resources"
                                    showButton={false}
                                    showSearchBox={true}
                                    handleSearchChange={getSearchValue}
                                />
                            </Grid>
                            {schematicsSpecsData.map(
                                ({ id, title, readingTime }: any) => (
                                    <Grid item key={id} xs md lg={4}>
                                        <Card
                                            variant="outlined"
                                            sx={{
                                                padding: "15px 18px 8px 18px",
                                                borderRadius: "10px",
                                            }}
                                        >
                                            <Stack
                                                spacing={1}
                                                direction="column"
                                            >
                                                <Typography
                                                    variant="h6"
                                                    sx={{
                                                        fontSize: "16px",
                                                    }}
                                                >
                                                    {title}
                                                </Typography>
                                                <Typography variant="caption">
                                                    {readingTime}
                                                </Typography>
                                                <Box
                                                    component="div"
                                                    sx={{
                                                        width: "100%",
                                                        display: "flex",
                                                        alignItems: "center",
                                                        justifyContent:
                                                            "center",
                                                    }}
                                                >
                                                    <img
                                                        src="https://i.ibb.co/BgfbTsP/1835cf7a15bd359317e492f4ea67166a.png"
                                                        alt="1835cf7a15bd359317e492f4ea67166a"
                                                        width="300"
                                                        height="160"
                                                    />
                                                </Box>
                                            </Stack>
                                        </Card>
                                    </Grid>
                                )
                            )}
                        </Grid>
                    </Grid>
                </Paper>
            </Stack>
        </>
    );
};

export default NodeSchematicTab;
