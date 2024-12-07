type CustomNetworkInfoProps = {
  networkIconColor?: string;
  siteOneIconColor?: string;
  siteTwoIconColor?: string;
  siteThreeIconColor?: string;
  networkColor?: string;
  siteOneColor?: string;
  siteTwoColor?: string;
  siteThreeColor?: string;
  networkName?: string;
  siteOneName?: string;
  siteTwoName?: string;
  siteThreeName?: string;
};

const CustomNetworkInfo = ({
  networkIconColor = '#6F7979',
  siteOneIconColor = '#6F7979',
  siteTwoIconColor = '#6F7979',
  siteThreeIconColor = '#6F7979',
  networkColor = '#333333',
  siteOneColor = '#333333',
  siteTwoColor = '#333333',
  siteThreeColor = '#333333',
  networkName = 'NETWORK',
  siteOneName = 'SITE 1',
  siteTwoName = 'SITE 2',
  siteThreeName = 'SITE 3',
}: CustomNetworkInfoProps) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="110.1561 150.5602 275.092 132.5885"
    width="364px"
    height="172px"
    preserveAspectRatio="none"
  >
    <mask
      id="mask0_345_2815"
      mask-type="luminance"
      maskUnits="userSpaceOnUse"
      x="151"
      y="0"
      width="35"
      height="35"
    >
      <path d="M186 0H151V35H186V0Z" fill="white" />
    </mask>
    <g
      id="object-18"
      transform="matrix(1, 0, 0, 1, 80.15612792968749, 146.33074951171875)"
    >
      <g mask="url(#mask0_345_2815)" id="object-2">
        <path
          d="M161.647 21.4378L163.397 19.6878C161.939 18.2295 161.21 16.3337 161.21 14.5837C161.21 12.6878 161.939 10.792 163.397 9.47949L161.647 7.72949C159.751 9.62533 158.73 12.1045 158.73 14.5837C158.73 17.0628 159.751 19.542 161.647 21.4378Z"
          fill={networkIconColor}
        />
        <path
          d="M178.855 4.22949L177.105 5.97949C179.439 8.31283 180.605 11.5212 180.605 14.5837C180.605 17.6462 179.439 20.8545 177.105 23.1878L178.855 24.9378C181.772 22.0212 183.085 18.3753 183.085 14.5837C183.085 10.792 181.626 7.14616 178.855 4.22949Z"
          fill={networkIconColor}
        />
        <path
          d="M159.897 5.97949L158.147 4.22949C155.376 7.14616 153.918 10.792 153.918 14.5837C153.918 18.3753 155.376 22.0212 158.147 24.9378L159.897 23.1878C157.564 20.8545 156.397 17.6462 156.397 14.5837C156.397 11.5212 157.564 8.31283 159.897 5.97949Z"
          fill={networkIconColor}
        />
        <path
          d="M175.355 21.4378C177.251 19.542 178.272 17.0628 178.272 14.5837C178.126 12.1045 177.251 9.62533 175.355 7.72949L173.605 9.47949C175.064 10.9378 175.793 12.8337 175.793 14.5837C175.793 16.4795 175.064 18.3753 173.605 19.6878L175.355 21.4378Z"
          fill={networkIconColor}
        />
        <path
          d="M172.148 14.5833C172.148 12.5708 170.515 10.9375 168.503 10.9375C166.49 10.9375 164.857 12.5708 164.857 14.5833C164.857 15.6917 165.353 16.6542 166.126 17.325L161.211 32.0833H164.128L165.105 29.1667H171.915L172.878 32.0833H175.794L170.88 17.325C171.653 16.6542 172.148 15.6917 172.148 14.5833ZM166.067 26.25L168.503 18.9583L170.938 26.25H166.067Z"
          fill={networkIconColor}
        />
      </g>
      <path
        d="M 167.559 59 L 167.559 86 L 167.559 59 Z"
        fill="black"
        id="object-3"
      />
      <path
        d="M 167.559 59 L 167.559 86"
        stroke="#B0B9C6"
        strokeWidth="0.999999"
        strokeLinecap="round"
        id="object-4"
      />
      <path d="M 305 71 L 30 71 L 305 71 Z" fill="black" id="object-5" />
      <path
        d="M 284.928 71 L 50.03 71"
        stroke="#B0B9C6"
        strokeWidth="0.999999"
        strokeLinecap="round"
        id="object-6"
      />
      <path d="M305 71L305.092 85.952Z" fill="black" id="object-7" />
      <path
        d="M 285 71 L 285.092 85.952"
        stroke="#B0B9C6"
        strokeWidth="0.999999"
        strokeLinecap="round"
        id="object-8"
      />
      <path
        d="M 30.008 71.231 L 30 85.97 L 30.008 71.231 Z"
        fill="black"
        id="object-9"
      />
      <path
        d="M 50.008 71.231 L 50 85.97"
        stroke="#B0B9C6"
        strokeWidth="0.999999"
        strokeLinecap="round"
        id="object-10"
      />
      <path
        d="M 50.188 91.246 C 44.544 91.246 39.979 95.811 39.979 101.454 C 39.979 109.111 50.188 120.413 50.188 120.413 C 50.188 120.413 60.396 109.111 60.396 101.454 C 60.396 95.811 55.831 91.246 50.188 91.246 Z M 50.188 105.1 C 48.175 105.1 46.542 103.467 46.542 101.454 C 46.542 99.442 48.175 97.808 50.188 97.808 C 52.2 97.808 53.833 99.442 53.833 101.454 C 53.833 103.467 52.2 105.1 50.188 105.1 Z"
        fill={siteOneIconColor}
        id="object-11"
      />
      <path
        d="M 285.209 91 C 279.565 91 275 95.565 275 101.208 C 275 108.865 285.209 120.167 285.209 120.167 C 285.209 120.167 295.417 108.865 295.417 101.208 C 295.417 95.565 290.852 91 285.209 91 Z M 285.209 104.854 C 283.196 104.854 281.563 103.221 281.563 101.208 C 281.563 99.196 283.196 97.562 285.209 97.562 C 287.221 97.562 288.854 99.196 288.854 101.208 C 288.854 103.221 287.221 104.854 285.209 104.854 Z"
        fill={siteTwoIconColor}
        id="object-12"
      />
      <path
        d="M168.209 91C162.565 91 158 95.565 158 101.208C158 108.865 168.209 120.167 168.209 120.167C168.209 120.167 178.417 108.865 178.417 101.208C178.417 95.565 173.852 91 168.209 91ZM168.209 104.854C166.196 104.854 164.563 103.221 164.563 101.208C164.563 99.196 166.196 97.562 168.209 97.562C170.221 97.562 171.854 99.196 171.854 101.208C171.854 103.221 170.221 104.854 168.209 104.854Z"
        fill={siteThreeIconColor}
        id="object-13"
      />
      <text
        y="52"
        x="136"
        id="object-14"
        style={{
          fontSize: '12px',
          whiteSpace: 'pre',
          fontWeight: 'normal',
          fill: networkColor,
          textAlign: 'center',
          transformBox: 'fill-box',
          transformOrigin: '50% 50%',
          fontFamily: 'Work Sans, sans-serif',
        }}
        textAnchor="middle"
        dominantBaseline="middle"
      >
        <tspan x="168" dy="-0.5em">
          {networkName}
        </tspan>
      </text>
      <text
        x="30"
        y="136"
        id="object-15"
        style={{
          fontSize: '12px',
          whiteSpace: 'pre',
          fill: siteOneColor,
          fontFamily: 'Work Sans, sans-serif',
        }}
      >
        {siteOneName}
      </text>
      <text
        x="148"
        y="136"
        id="object-16"
        style={{
          fontSize: '12px',
          whiteSpace: 'pre',
          fill: siteTwoColor,
          fontFamily: 'Work Sans, sans-serif',
        }}
      >
        {siteTwoName}
      </text>
      <text
        x="264"
        y="136"
        id="object-17"
        style={{
          width: '200px',
          fontSize: '12px',
          whiteSpace: 'pre',
          fill: siteThreeColor,
          fontFamily: 'Work Sans, sans-serif',
        }}
      >
        {siteThreeName}
      </text>
    </g>
  </svg>
);

export default CustomNetworkInfo;
