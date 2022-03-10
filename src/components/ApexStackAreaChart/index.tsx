import Chart from "react-apexcharts";
import { format } from "date-fns";
import { GraphTitleWrapper } from "..";

const TIME_RANGE_IN_MILLISECONDS = 100;

interface IApexStackAreaChart {
    data: any;
    name: string;
    filter?: string;
    hasData: boolean;
    refreshInterval?: number;
    onRefreshData?: Function;
    onFilterChange?: Function;
}

const StackAreaChart = (props: any) => {
    const options: any = {
        chart: {
            type: "area",
            height: 350,
            stacked: true,
            zoom: {
                type: "x",
                enabled: false,
                autoScaleYaxis: true,
            },
        },
        colors: ["#008FFB"],
        plotOptions: {
            area: {
                fillTo: "end",
            },
        },
        dataLabels: {
            enabled: false,
        },
        stroke: {
            curve: "smooth",
        },
        fill: {
            type: "gradient",
            gradient: {
                opacityFrom: 0.6,
                opacityTo: 0.8,
            },
        },
        xaxis: {
            type: "datetime",
            range: props.range,
            labels: {
                formatter: (val: any) =>
                    val ? format(new Date(val * 1000), "hh:mm:ss") : "",
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
            height={"300px"}
            options={options}
            series={props.dataList}
        />
    );
};

const ApexStackAreaChart = ({
    name,
    data = [],
    filter = "LIVE",
    hasData = false,
    onFilterChange = () => {
        /*DEFAULT FUNCTION*/
    },
}: IApexStackAreaChart) => {
    return (
        <GraphTitleWrapper
            key={name}
            title={name}
            filter={filter}
            hasData={hasData}
            variant="subtitle1"
            handleFilterChange={onFilterChange}
        >
            <StackAreaChart
                key={name}
                name={name}
                dataList={data}
                range={TIME_RANGE_IN_MILLISECONDS}
            />
        </GraphTitleWrapper>
    );
};

export default ApexStackAreaChart;
