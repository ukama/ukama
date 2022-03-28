import { Paper, Grid, Typography, Card, Stack } from "@mui/material";
import { ContainerHeader } from "../../components";
import { LoadingWrapper } from "..";
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
                        <Grid xs={12} container spacing={1}>
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
                                    <Grid item xs={4} key={id}>
                                        <Card
                                            variant="outlined"
                                            sx={{
                                                padding: "15px 18px 8px 18px",
                                                borderRadius: "10px",
                                            }}
                                        >
                                            <Stack direction="column">
                                                <Typography
                                                    variant="h4"
                                                    sx={{
                                                        fontSize: "16px",
                                                    }}
                                                >
                                                    {title}
                                                </Typography>
                                                <Typography variant="caption">
                                                    {readingTime}
                                                </Typography>
                                                <img
                                                    src="https://i.ibb.co/BgfbTsP/1835cf7a15bd359317e492f4ea67166a.png"
                                                    alt="1835cf7a15bd359317e492f4ea67166a"
                                                    width="300"
                                                    height="160"
                                                />
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
