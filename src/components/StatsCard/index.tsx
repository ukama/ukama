import { LineChart, TimeFilter } from "..";
import { colors } from "../../theme";
import { useRecoilValue } from "recoil";
import { makeStyles } from "@mui/styles";
import { isDarkmode } from "../../recoil";
import { StatsItemType } from "../../types";
import { RoundedCard, SkeletonRoundedCard } from "../../styles";
import { Grid, Select, MenuItem, Typography } from "@mui/material";

type StyleProps = { color: string };

const useStyles = makeStyles(() => ({
    selectStyle: ({ color }: StyleProps) => ({
        width: "172px",
        "& p": {
            fontSize: "20px",
            fontWeight: "500",
            lineHeight: "160%",
            fontFamily: "Rubik",
            letterSpacing: "0.15px",
            color: color,
        },
    }),
}));

type StatsCardProps = {
    loading: boolean;
    handleSelect: any;
    selectOption: number;
    selectedButton: string;
    options: StatsItemType[];
    handleSelectedButton: any;
};

const StatsCard = ({
    options,
    loading,
    selectOption,
    handleSelect,
    selectedButton,
    handleSelectedButton,
}: StatsCardProps) => {
    const _isDarkMod = useRecoilValue(isDarkmode);
    const styleProps = { color: _isDarkMod ? colors.white : colors.black };
    const classes = useStyles(styleProps);
    return (
        <>
            {loading ? (
                <SkeletonRoundedCard variant="rectangular" height={337} />
            ) : (
                <RoundedCard>
                    <Grid container spacing={1}>
                        <Grid item container xs={12}>
                            <Grid item xs={12} sm={6}>
                                <Select
                                    sx={{
                                        width: "auto",
                                    }}
                                    disableUnderline
                                    variant="standard"
                                    value={selectOption}
                                    onChange={handleSelect}
                                    className={classes.selectStyle}
                                >
                                    {options.map(
                                        ({ id, label }: StatsItemType) => (
                                            <MenuItem key={id} value={id}>
                                                <Typography variant="body1">
                                                    {label}
                                                </Typography>
                                            </MenuItem>
                                        )
                                    )}
                                </Select>
                            </Grid>
                            <Grid
                                item
                                xs={12}
                                sm={6}
                                display="flex"
                                justifyContent={"flex-end"}
                            >
                                <TimeFilter
                                    filter={selectedButton}
                                    handleFilterSelect={handleSelectedButton}
                                />
                            </Grid>
                        </Grid>
                        <Grid item xs={12}>
                            <LineChart
                                title={""}
                                hasData={true}
                                showFilter={false}
                            />
                        </Grid>
                    </Grid>
                </RoundedCard>
            )}
        </>
    );
};
export default StatsCard;
