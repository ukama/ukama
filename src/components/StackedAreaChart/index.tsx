import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
    Filler,
} from "chart.js";
import "chartjs-adapter-luxon";
import { Chart } from "react-chartjs-2";
import zoomPlugin from "chartjs-plugin-zoom";
import ChartStreaming from "chartjs-plugin-streaming";
import { GraphTitleWrapper } from "..";

ChartJS.register(
    zoomPlugin,
    ChartStreaming,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
    Filler
);

const data = {
    datasets: [
        {
            data: [],
            label: "My First dataset",
            borderColor: "rgb(255,99,132)",
            backgroundColor: "rgba(255,99,132, 0.5)",
            fill: true,
        },
        {
            data: [],
            label: "My Second dataset",
            borderColor: "rgb(76,192,192)",
            backgroundColor: "rgba(76,192,192, 0.5)",
            fill: true,
        },
        {
            data: [],
            label: "My Third dataset",
            borderColor: "rgb(53, 162, 235)",
            backgroundColor: "rgba(53, 162, 235, 0.5)",
            fill: true,
        },
        {
            data: [],
            label: "My Fourth dataset",
            borderColor: "rgb(255,205,87)",
            backgroundColor: "rgba(255,205,87, 0.5)",
            fill: true,
        },
    ],
};

const onRefresh = (chart: any) => {
    const now = Date.now();
    chart.data.datasets.forEach((dataset: any) => {
        dataset.data.push({
            x: now,
            y: Math.floor(Math.random() * 100) + -10,
        });
    });
};

const zoomOptions = {
    pan: {
        enabled: true,
        mode: "x",
    },
    zoom: {
        pinch: {
            enabled: true,
        },
        wheel: {
            enabled: true,
        },
        mode: "x",
    },
    limits: {
        x: {
            minDelay: 0,
            maxDelay: 4000,
            minDuration: 10000, //ZoomIn
            maxDuration: 80000, //ZoomOut
        },
    },
};

const config: any = {
    type: "line",
    data: data,
    options: {
        responsive: true,
        scales: {
            x: {
                type: "realtime",
                realtime: {
                    duration: 20000,
                    refresh: 1000,
                    delay: 2000,
                    onRefresh: onRefresh,
                },
                title: {
                    display: true,
                    text: "Month",
                },
            },
            y: {
                stacked: true,
                title: {
                    display: true,
                    text: "Value",
                },
            },
        },
        interaction: {
            mode: "nearest",
            axis: "x",
            intersect: false,
        },
        plugins: {
            annotation: false,
            datalabels: false,
            zoom: zoomOptions,
            legend: { position: "bottom" },
        },
    },
};

interface IStackedAreaChart {
    title: string;
    hasData: boolean;
    height?: number;
}

const StackedAreaChart = ({
    title,
    hasData,
    height = 80,
}: IStackedAreaChart) => {
    return (
        <GraphTitleWrapper hasData={hasData} variant="subtitle1" title={title}>
            <Chart
                type="line"
                key={title}
                id={title}
                height={height}
                data={data}
                options={{ ...config.options, id: title }}
            />
        </GraphTitleWrapper>
    );
};

export default StackedAreaChart;
