import React from "react";
import { format } from "date-fns";
import Chart from "react-apexcharts";
import { colors } from "../../theme";
import { GraphTitleWrapper } from "..";
import { useRecoilValue } from "recoil";
import { makeStyles } from "@mui/styles";
import { isDarkmode } from "../../recoil";
import { Box, Theme } from "@mui/material";
import { isMetricData } from "../../utils";

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
    name: string;
    filter?: string;
    hasData?: boolean;
    refreshInterval?: number;
    onRefreshData?: Function;
    onFilterChange?: Function;
}

const ApexLineChart = (props: any) => {
    const _isDarkMod = useRecoilValue(isDarkmode);
    const options: any = {
        stroke: {
            lineCap: "butt",
            curve: "smooth",
            width: 5,
        },
        chart: {
            minHeight: "200px",
            height: "100%",
            width: "100%",
            zoom: {
                type: "x",
                enabled: false,
                autoScaleYaxis: true,
            },
            animations: {
                enabled: true,
                easing: "linear",
                dynamicAnimation: {
                    speed: 1000,
                },
            },
            dropShadow: {
                enabled: true,
                top: 1,
                left: 1,
                bottom: 1,
                blur: 3,
                opacity: 0.2,
            },
        },
        grid: {
            borderColor: _isDarkMod ? colors.vulcan60 : colors.white60,
            opacity: 0.3,
        },
        xaxis: {
            type: "datetime",
            range: props.range,
            labels: {
                formatter: (val: any) =>
                    val ? format(new Date(val * 1000), "mm:ss") : "",
            },
            tooltip: {
                enabled: false,
                offsetX: 0,
            },
        },

        yaxis: {
            labels: {
                formatter: (val: any) => val.toFixed(2),
            },
            // min: 0,
            // max: 100,
            // tooltip: {
            //     enabled: true,
            // },
            tickAmount: 8,
        },
    };
    return (
        <Chart
            type="line"
            key={props.name}
            options={options}
            height={"300px"}
            series={props.dataList}
        />
    );
};

const ApexLineChartIntegration = ({
    name,
    data = [],
    onRefreshData,
    filter = "LIVE",
    refreshInterval = 10000,
    onFilterChange = () => {
        /*DEFAULT FUNCTION*/
    },
}: IApexLineChartIntegration) => {
    const classes = useStyles();
    React.useEffect(() => {
        const interval = setInterval(() => {
            onRefreshData && onRefreshData();
        }, refreshInterval);

        return () => clearInterval(interval);
    });

    return (
        <GraphTitleWrapper
            key={name}
            title={name}
            filter={filter}
            variant="subtitle1"
            hasData={isMetricData(data)}
            handleFilterChange={onFilterChange}
        >
            <Box component="div" className={classes.chartStyle}>
                <ApexLineChart
                    key={name}
                    name={name}
                    dataList={data}
                    range={TIME_RANGE_IN_MILLISECONDS}
                />
            </Box>
        </GraphTitleWrapper>
    );
};

export default ApexLineChartIntegration;
