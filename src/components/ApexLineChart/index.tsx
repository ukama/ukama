import LineChart from "./LineChart";
import { GraphTitleWrapper } from "..";
import { makeStyles } from "@mui/styles";
import { Box, Theme } from "@mui/material";

const TIME_RANGE_IN_MILLISECONDS = 100;

const useStyles = makeStyles<Theme>(theme => ({
    chartStyle: {
        "& .apexcharts-yaxis-label tspan": {
            fill: theme.palette.text.primary,
        },
        "& .apexcharts-xaxis-label tspan": {
            fill: theme.palette.text.primary,
        },
    },
}));

interface IApexLineChartIntegration {
    data: any;
    name?: string;
    filter?: string;
    hasData?: boolean;
    onFilterChange?: Function;
}

const ApexLineChart = ({
    data = { name: "-", data: [] },
    filter = "LIVE",
    onFilterChange = () => {
        /*DEFAULT FUNCTION*/
    },
}: IApexLineChartIntegration) => {
    const classes = useStyles();
    return (
        <GraphTitleWrapper
            filter={filter}
            variant="subtitle1"
            title={data?.name || ""}
            handleFilterChange={onFilterChange}
            hasData={data?.data.length > 0 || false}
        >
            <Box component="div" className={classes.chartStyle}>
                <LineChart
                    name={data?.name || ""}
                    dataList={[data]}
                    range={TIME_RANGE_IN_MILLISECONDS}
                />
            </Box>
        </GraphTitleWrapper>
    );
};

export default ApexLineChart;
