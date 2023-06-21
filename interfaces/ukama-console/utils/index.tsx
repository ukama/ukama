import colors from '@/styles/theme/colors';
import { TObject } from '@/types';
import { Typography } from '@mui/material';
import { format, intervalToDuration } from 'date-fns';
import { Alert_Type, Graphs_Tab, NodeDto, Node_Type } from '../generated';
const getTitleFromPath = (path: string) => {
  switch (path) {
    case '/home':
      return 'Home';
    case '/settings':
      return 'Settings';
    case '/site':
      return 'Sites';
    case '/node':
      return 'Nodes';
    case '/subscriber':
      return 'Subscribers';
    case '/site_planning':
      return 'Site Planning';
    case '/manage':
      return 'Manage';
    case '/uidev':
      return 'UIDEV';
    default:
      return '404';
  }
};

const getColorByType = (type: Alert_Type) => {
  switch (type) {
    case Alert_Type.Error:
      return 'error';
    case Alert_Type.Warning:
      return 'warning';
    default:
      return 'success';
  }
};
const NetworkStatusLabel = ({ message }: any) => {
  return (
    <Typography variant="h6">
      Your <span style={{ color: colors.primaryMain }}>network</span> {message}
    </Typography>
  );
};

const getStatusByType = (status: string) => {
  if (status === 'DOWN') return <NetworkStatusLabel message="is down." />;
  else if (status === 'ONLINE')
    return <NetworkStatusLabel message="is online and well:" />;
  else return <NetworkStatusLabel message=" is offline; no nodes connected." />;
};

const parseObjectInNameValue = (obj: any) => {
  let updatedObj: TObject[] = [];
  if (obj) {
    updatedObj = Object.keys(obj).map((key) => {
      return {
        name: key,
        value:
          key === 'timestamp' ? format(obj[key], 'MMM dd HH:mm:ss') : obj[key],
      };
    });

    let removeIndex = updatedObj
      .map((item) => item?.name)
      .indexOf('__typename');
    ~removeIndex && updatedObj.splice(removeIndex, 1);
    removeIndex = updatedObj.map((item) => item?.name).indexOf('id');
    ~removeIndex && updatedObj.splice(removeIndex, 1);
  }

  return updatedObj;
};

const uniqueObjectsArray = (name: string, list: TObject[]): TObject[] | [] => {
  const last =
    list.length > 0 ? list.filter((item: TObject) => item.name !== name) : [];
  return last;
};

const hexToRGB = (hex: string, alpha: number): string => {
  var h = '0123456789ABCDEF';
  var r = h.indexOf(hex[1]) * 16 + h.indexOf(hex[2]);
  var g = h.indexOf(hex[3]) * 16 + h.indexOf(hex[4]);
  var b = h.indexOf(hex[5]) * 16 + h.indexOf(hex[6]);
  if (alpha) {
    return `rgba(${r}, ${g}, ${b}, ${alpha})`;
  }

  return `rgba(${r}, ${g}, ${b})`;
};

const random = (min: number, max: number) => Math.random() * (max - min) + min;

const getRandomData = () => {
  const data = [];
  for (let i = 0; i < 10; i++) {
    data.push({
      x: Date.now() / 1000 - (10 - i),
      y: random(-2, 2),
    });
  }
  return data;
};

const getGraphFilterByType = (type: string) => {
  switch (type) {
    case 'DAY':
      return {
        to: Math.round(Date.now() / 1000),
        from: Math.round(Date.now() / 1000) - 86400,
      };
    case 'WEEK':
      return {
        to: Math.round(Date.now() / 1000),
        from: Math.round(Date.now() / 1000) - 604800,
      };
    case 'MONTH':
      return {
        to: Math.round(Date.now() / 1000),
        from: Math.round(Date.now() / 1000) - 2628002,
      };
  }
};

const getTabByIndex = (index: number) => {
  switch (index) {
    case 0:
      return Graphs_Tab.Overview;
    case 1:
      return Graphs_Tab.Network;
    case 2:
      return Graphs_Tab.Resources;
    case 3:
      return Graphs_Tab.Radio;
    case 4:
      return Graphs_Tab.Home;
    default:
      return Graphs_Tab.Overview;
  }
};

const getMetricPayload = ({
  tab = 0,
  nodeId = '',
  regPolling = true,
  nodeType = Node_Type.Home,
  to = Math.floor(Date.now() / 1000),
  from = Math.floor(Date.now() / 1000),
}: {
  to?: number;
  tab: number;
  from?: number;
  nodeId?: string;
  nodeType: string;
  regPolling?: boolean;
}) => {
  return {
    data: {
      step: 1,
      nodeId: nodeId,
      regPolling: regPolling,
      to: to,
      from: from, //20sec
      nodeType: nodeType as Node_Type,
      tab: getTabByIndex(tab),
    },
  };
};

const isMetricData = (metric: any) => {
  let isData = false;
  metric.forEach((item: any) => {
    if (item.data.length > 0) {
      isData = true;
    }
  });
  return isData;
};

const isContainNodeUpdate = (list: NodeDto[] = []): boolean => {
  let isUpdate = false;
  for (const ele of list) {
    if (ele.isUpdateAvailable) {
      isUpdate = true;
      break;
    }
  }

  return isUpdate;
};

const getDefaultMetricList = (name: string) => {
  return {
    name: name,
    data: [],
  };
};

const getTitleByKey = (key: string) => {
  switch (key) {
    case 'uptimetrx':
      return 'Uptime TRX';
    case 'temperaturetrx':
      return 'Temp. (TRX)';
    case 'temperaturerfe':
      return 'Temp. (RFE)';
    case 'subscribersactive':
      return 'Active';
    case 'subscribersattached':
      return 'Attached';
    case 'temperaturectl':
      return 'Temp. (CTL)';
    case 'temperaturecom':
      return 'Temp. (COM)';
    case 'rrc':
      return 'RRC CNX success';
    case 'rlc':
      return 'RLS  drop rate';
    case 'erab':
      return 'ERAB drop rate';
    case 'throughputuplink':
      return 'Throughput (U/L)';
    case 'throughputdownlink':
      return 'Throughput (D/L)';
    case 'cputrxusage':
      return 'CPU-TRX';
    case 'memorytrxused':
      return 'Memory-TRX';
    case 'disktrxused':
      return 'DISK-TRX';
    case 'cpuctlused':
      return 'CPU-CTL';
    case 'diskctlused':
      return 'DISK-CTL';
    case 'memoryctlused':
      return 'Memory-CTL';
    case 'powerlevel':
      return 'Power';
    case 'cpucomusage':
      return 'CPU-COM';
    case 'diskcomused':
      return 'DISK-COM';
    case 'memorycomused':
      return 'Memory-COM';
    case 'txpower':
      return 'TX Power';
    case 'rxpower':
      return 'RX Power';
    case 'papower':
      return 'PA Power';
    default:
      return '';
  }
};

const getMetricsInitObj = () => {
  return {
    temperaturetrx: null,
    temperaturerfe: null,
    subscribersactive: null,
    subscribersattached: null,
    temperaturectl: null,
    temperaturecom: null,
    rrc: null,
    rlc: null,
    erab: null,
    throughputuplink: null,
    throughputdownlink: null,
    cputrxusage: null,
    memorytrxused: null,
    disktrxused: null,
    cpuctlused: null,
    diskctlused: null,
    memoryctlused: null,
    powerlevel: null,
    cpucomusage: null,
    diskcomused: null,
    memorycomused: null,
    txpower: null,
    rxpower: null,
    papower: null,
    uptimetrx: null,
  };
};

const getMetricObjectByKey = (key: string) => {
  return { name: getTitleByKey(key), data: [] };
};

const formatBytes = (bytes = 0): string => {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const dm = 3;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ` ${sizes[i]}`;
};

const formatBytesToMB = (bytes = 0): string => {
  if (bytes === 0) return '0';
  return (bytes / (1024 * 1024)).toFixed(2);
};

const secondsToDuration = (end: any) => {
  // const units = ["year", "month", "day", "hour", "minute"];
  const units = ['hour', 'minute'];

  const duration: any = intervalToDuration({
    start: 0,
    end: end * 1000,
  });

  const response: any = [];
  let required = false;

  units.forEach((unit, index) => {
    if (
      duration[`${unit}s`] > 0 ||
      required === true ||
      index === units.length - 1
    ) {
      response.push(
        `${duration[`${unit}s`]} ${unit}${
          duration[`${unit}s`] === 1 ? '' : 's'
        }`,
      );
      required = true;
    }
  });

  return response.join(' ');
};

const secToHoursNMints = (seconds: number, separator: string) => {
  return (
    [Math.floor(seconds / 60 / 60), Math.floor((seconds / 60) % 60)]
      .join(separator ? separator : ':')
      .replace(/\b(\d)\b/g, '0$1') + ' minutes'
  );
};

const formatSecondsToDuration = (totalSeconds: number) => {
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds - hours * 3600 - minutes * 60;

  return [`${hours}h`, `${minutes}m`, `${seconds}s`]
    .filter((item) => item[0] !== '0')
    .join(' ');
};

const isEmailValid = (email: string): boolean =>
  /^\w+@[a-zA-Z_]+?\.[a-zA-Z]{2,20}$/.test(email);

const doesHttpOnlyCookieExist = (cookiename: string): boolean => {
  var d = new Date();
  d.setTime(d.getTime() + 1000);
  var expires = 'expires=' + d.toUTCString();

  document.cookie = cookiename + '=new_value;path=/;' + expires;
  return document.cookie.indexOf(cookiename + '=') == -1;
};

const getTowerNodeFromNodes = (nodes: NodeDto[]): string => {
  if (nodes.length > 0) {
    for (const node of nodes) {
      if (node.type === Node_Type.Tower) return node.id;
    }
    for (const node of nodes) {
      if (node.type === Node_Type.Home) return node.id;
    }
  }
  return '';
};

const getDataUsageSymbol = (dataUnit: string): string => {
  switch (dataUnit) {
    case 'GigaBytes':
      return 'GB';
    case 'MegaBytes':
      return 'MB';
    case 'KiloBytes':
      return 'KB';
    default:
      return 'MB';
  }
};

const getDataPlanUsage = (
  duration: string,
  currency: string,
  amount: string,
  dataVolume: string,
  dataUnit: string,
): string => {
  const symbol = currency === 'Dollar' ? '$' : '$';
  return `${symbol}${amount} / ${dataVolume} ${getDataUsageSymbol(
    dataUnit,
  )} / ${duration}`;
};

const fileToBase64 = (file: File): Promise<string> => {
  return new Promise<string>((resolve, reject) => {
    const reader = new FileReader();
    reader.onloadend = () => {
      const base64String = reader.result as string;
      resolve(base64String.split(',')[1]);
    };
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
};

function calculateCenterLatLng(coordinates: any) {
  var totalCoords = coordinates.length;
  if (totalCoords === 0) return { lat: 37.7780627, lng: -121.9822475 };
  if (totalCoords === 1)
    return { lat: coordinates[0].lat, lng: coordinates[0].lng };

  var sumLat = 0;
  var sumLng = 0;

  // Calculate the sum of latitudes and longitudes
  for (var i = 0; i < totalCoords; i++) {
    sumLat += coordinates[i].lat;
    sumLng += coordinates[i].lng;
  }

  // Calculate the average latitudes and longitudes
  var avgLat = sumLat / totalCoords;
  var avgLng = sumLng / totalCoords;

  // Return the center latitude and longitude
  return { lat: avgLat, lng: avgLng };
}

export {
  calculateCenterLatLng,
  doesHttpOnlyCookieExist,
  fileToBase64,
  formatBytes,
  formatBytesToMB,
  formatSecondsToDuration,
  getColorByType,
  getDataPlanUsage,
  getDefaultMetricList,
  getGraphFilterByType,
  getMetricObjectByKey,
  getMetricPayload,
  getMetricsInitObj,
  getRandomData,
  getStatusByType,
  getTitleFromPath,
  getTowerNodeFromNodes,
  hexToRGB,
  isContainNodeUpdate,
  isEmailValid,
  isMetricData,
  parseObjectInNameValue,
  secToHoursNMints,
  secondsToDuration,
  uniqueObjectsArray,
};
