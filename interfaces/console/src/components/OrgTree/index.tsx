/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { OrgTreeRes } from '@/client/graphql/generated';
import { convertToNamesJson } from '@/utils/getOrgTreeStructure';
import BatteryChargingFullOutlinedIcon from '@mui/icons-material/BatteryChargingFullOutlined';
import CellTowerOutlinedIcon from '@mui/icons-material/CellTowerOutlined';
import CorporateFareSharpIcon from '@mui/icons-material/CorporateFareSharp';
import DataSaverOffOutlinedIcon from '@mui/icons-material/DataSaverOffOutlined';
import DeveloperBoardOutlinedIcon from '@mui/icons-material/DeveloperBoardOutlined';
import Diversity1OutlinedIcon from '@mui/icons-material/Diversity1Outlined';
import LanOutlinedIcon from '@mui/icons-material/LanOutlined';
import PeopleOutlinedIcon from '@mui/icons-material/PeopleOutlined';
import PermDataSettingOutlinedIcon from '@mui/icons-material/PermDataSettingOutlined';
import PersonOutlinedIcon from '@mui/icons-material/PersonOutlined';
import RouterOutlinedIcon from '@mui/icons-material/RouterOutlined';
import SatelliteAltOutlinedIcon from '@mui/icons-material/SatelliteAltOutlined';
import SimCardOutlinedIcon from '@mui/icons-material/SimCardOutlined';
import { useCallback, useState } from 'react';
import Tree from 'react-d3-tree';

const useCenteredTree = (defaultTranslate = { x: 0, y: 0 }) => {
  const [translate, setTranslate] = useState(defaultTranslate);
  const containerRef = useCallback((containerElem: any) => {
    if (containerElem !== null) {
      const { width, height } = containerElem.getBoundingClientRect();
      setTranslate({ x: width / 2, y: height / 5 });
    }
  }, []);
  return { translate, containerRef };
};

const containerStyles = {
  width: '100vw',
  height: '100vh',
};

const textStyle = {
  fontSize: '14px',
  whiteSpace: 'pre',
  fontFamily: 'Work Sans, sans-serif',
};

const Wrapper = ({ children, style }: any) => (
  <div
    style={{
      display: 'flex',
      alignItems: 'center',
      stroke: 'transparent',
      flexDirection: 'column',
      ...style,
    }}
  >
    {children}
  </div>
);

const DefaultElement = ({ nodeDatum, toggleNode }: any) => (
  <g>
    <circle r="10" onClick={toggleNode} />
    <text fill="black" strokeWidth="1" x="20">
      {nodeDatum.name}
    </text>
  </g>
);

const renderCustomNodeElement = ({ nodeDatum, toggleNode }: any) => {
  switch (nodeDatum.elementType) {
    case 'ORG':
      return (
        <foreignObject width="100" height="50" x="-72" y="-20">
          <Wrapper>
            <CorporateFareSharpIcon />
            <span style={{ ...textStyle }}>{nodeDatum.name}</span>
          </Wrapper>
        </foreignObject>
      );
    case 'DATAPLAN':
      return (
        <g>
          <rect
            x="-50"
            y="-25"
            width="100"
            height="50"
            fill="white"
            stroke="lightgray"
          />
          <foreignObject width="100" height="50" x="-50" y="-18">
            <Wrapper>
              <PermDataSettingOutlinedIcon />
              <span style={{ ...textStyle }}>{nodeDatum.name}</span>
            </Wrapper>
          </foreignObject>
        </g>
      );
    case 'PLAN':
      return (
        <g>
          <rect
            x="-50"
            y="-25"
            width="100"
            height="50"
            fill="white"
            stroke="lightgray"
          />
          <foreignObject width="100" height="50" x="-50" y="-18">
            <Wrapper>
              <DataSaverOffOutlinedIcon />
              <span style={{ ...textStyle }}>{nodeDatum.name}</span>
            </Wrapper>
          </foreignObject>
        </g>
      );
    case 'SIMS':
      return (
        <g>
          <rect
            x="-50"
            y="-25"
            width="100"
            height="50"
            fill="white"
            stroke="lightgray"
          />
          <foreignObject width="100" height="50" x="-50" y="-18">
            <Wrapper>
              <SimCardOutlinedIcon />
              <span style={{ ...textStyle }}>{nodeDatum.name}</span>
            </Wrapper>
          </foreignObject>
        </g>
      );
    case 'SIM_STATS':
      return (
        <g>
          <foreignObject width="150" height="100" x="-75" y="-36">
            <div style={{ padding: '10px', fontSize: '12px' }}>
              <table
                style={{
                  width: '100%',
                  backgroundColor: 'white',
                  borderCollapse: 'collapse',
                  border: '1px solid lightgray',
                }}
              >
                <tbody>
                  <tr key={'availableSims'}>
                    <td>Available</td>
                    <td>{nodeDatum.availableSims}</td>
                  </tr>
                  <tr key={'consumed'}>
                    <td>Consumed</td>
                    <td>{nodeDatum.consumed}</td>
                  </tr>
                  <tr key={'totalSims'}>
                    <td>Total</td>
                    <td>{nodeDatum.totalSims}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </foreignObject>
        </g>
      );
    case 'MEMBERS':
      return (
        <g>
          <rect
            x="-50"
            y="-25"
            width="100"
            height="50"
            fill="white"
            stroke="lightgray"
          />
          <foreignObject width="100" height="50" x="-50" y="-18">
            <Wrapper>
              <PeopleOutlinedIcon />
              <span style={{ ...textStyle }}>{nodeDatum.name}</span>
            </Wrapper>
          </foreignObject>
        </g>
      );
    case 'MEMBER_STATS':
      return (
        <g>
          <foreignObject width="150" height="100" x="-75" y="-36">
            <div style={{ padding: '10px', fontSize: '12px' }}>
              <table
                style={{
                  width: '100%',
                  border: '1px solid lightgray',
                  borderCollapse: 'collapse',
                  backgroundColor: 'white',
                }}
              >
                <tbody>
                  <tr key={'activeMembers'}>
                    <td>Active</td>
                    <td>{nodeDatum.activeMembers}</td>
                  </tr>
                  <tr key={'inactiveMembers'}>
                    <td>Inactive</td>
                    <td>{nodeDatum.inactiveMembers}</td>
                  </tr>
                  <tr key={'totalMembers'}>
                    <td>Total</td>
                    <td>{nodeDatum.totalMembers}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </foreignObject>
        </g>
      );
    case 'SUBSCRIBERS':
      return (
        <g>
          <rect
            x="-50"
            y="-25"
            width="100"
            height="50"
            fill="white"
            stroke="lightgray"
          />
          <foreignObject width="100" height="50" x="-50" y="-18">
            <Wrapper>
              <PersonOutlinedIcon />
              <span style={{ ...textStyle }}>{nodeDatum.name}</span>
            </Wrapper>
          </foreignObject>
        </g>
      );
    case 'SUBSCRIBER_STATS':
      return (
        <g>
          <foreignObject width="150" height="100" x="-75" y="-36">
            <div style={{ padding: '10px', fontSize: '12px' }}>
              <table
                style={{
                  width: '100%',
                  backgroundColor: 'white',
                  borderCollapse: 'collapse',
                  border: '1px solid lightgray',
                }}
              >
                <tbody>
                  <tr key={'activeSubscribers'}>
                    <td>Active</td>
                    <td>{nodeDatum.activeSubscribers}</td>
                  </tr>
                  <tr key={'inactiveSubscribers'}>
                    <td>Inactive</td>
                    <td>{nodeDatum.inactiveSubscribers}</td>
                  </tr>
                  <tr key={'totalSubscribers'}>
                    <td>Total</td>
                    <td>{nodeDatum.totalSubscribers}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </foreignObject>
        </g>
      );
    case 'NETWORKS':
      return (
        <g>
          <rect
            x="-50"
            y="-25"
            width="100"
            height="50"
            fill="white"
            stroke="lightgray"
          />
          <foreignObject width="100" height="50" x="-50" y="-18">
            <Wrapper>
              <Diversity1OutlinedIcon />
              <span style={{ ...textStyle }}>{nodeDatum.name}</span>
            </Wrapper>
          </foreignObject>
        </g>
      );
    case 'NETWORK':
      return (
        <g>
          <rect
            x="-50"
            y="-25"
            width="100"
            height="50"
            fill="white"
            stroke="lightgray"
          />
          <foreignObject width="100" height="50" x="-50" y="-18">
            <Wrapper>
              <LanOutlinedIcon />
              <span style={{ ...textStyle }}>{nodeDatum.name}</span>
            </Wrapper>
          </foreignObject>
        </g>
      );
    case 'SITE':
      return (
        <g>
          <rect
            x="-50"
            y="-25"
            width="100"
            height="50"
            fill="white"
            stroke="lightgray"
          />
          <foreignObject width="100" height="50" x="-50" y="-18">
            <Wrapper>
              <CellTowerOutlinedIcon />
              <span style={{ ...textStyle }}>{nodeDatum.name}</span>
            </Wrapper>
          </foreignObject>
        </g>
      );
    case 'ACCESS':
      return (
        <g>
          <rect
            x="0"
            y="-25"
            width="364"
            height="50"
            fill="white"
            stroke="lightgray"
          />
          <foreignObject width="364" height="50" x="0" y="-18">
            <Wrapper>
              <RouterOutlinedIcon />
              <span style={{ ...textStyle }}>{nodeDatum.name}</span>
            </Wrapper>
          </foreignObject>
        </g>
      );
    case 'SWITCH':
      return (
        <g>
          <rect
            x="0"
            y="-25"
            width="364"
            height="50"
            fill="white"
            stroke="lightgray"
          />
          <foreignObject width="364" height="50" x="0" y="-18">
            <Wrapper>
              <DeveloperBoardOutlinedIcon />
              <span style={{ ...textStyle }}>{nodeDatum.name}</span>
            </Wrapper>
          </foreignObject>
        </g>
      );
    case 'POWER':
      return (
        <g>
          <rect
            x="0"
            y="-25"
            width="364"
            height="50"
            fill="white"
            stroke="lightgray"
          />
          <foreignObject width="364" height="50" x="0" y="-18">
            <Wrapper>
              <BatteryChargingFullOutlinedIcon />
              <span style={{ ...textStyle }}>{nodeDatum.name}</span>
            </Wrapper>
          </foreignObject>
        </g>
      );
    case 'BACKHAUL':
      return (
        <g>
          <rect
            x="0"
            y="-25"
            width="364"
            height="50"
            fill="white"
            stroke="lightgray"
          />
          <foreignObject width="364" height="50" x="0" y="-18">
            <Wrapper>
              <SatelliteAltOutlinedIcon />
              <span style={{ ...textStyle }}>{nodeDatum.name}</span>
            </Wrapper>
          </foreignObject>
        </g>
      );
    default:
      return <DefaultElement nodeDatum={nodeDatum} toggleNode={toggleNode} />;
  }
};

interface IOrgTree {
  data: OrgTreeRes | undefined;
}

export const OrgTree = ({ data }: IOrgTree) => {
  const { translate, containerRef } = useCenteredTree();
  return (
    <div style={containerStyles} ref={containerRef}>
      {data && (
        <Tree
          zoom={0.4}
          initialDepth={5}
          orientation="horizontal"
          data={convertToNamesJson(data)}
          translate={{ x: 100, y: translate?.y }}
          renderCustomNodeElement={renderCustomNodeElement}
        />
      )}
    </div>
  );
};
