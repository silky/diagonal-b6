import type { SVGProps } from 'react';
const SvgCarRepair = (props: SVGProps<SVGSVGElement>) => (
    <svg
        xmlns="http://www.w3.org/2000/svg"
        width={18}
        height={18}
        fill="none"
        viewBox="0 0 20 20"
        {...props}
    >
        <circle cx={10} cy={10} r={9.5} fill={props.fill} stroke="#fff" />
        <path
            fill="#fff"
            d="M6.477 7.16c.63 0 1.181-.343 1.476-.853h6.195a.852.852 0 0 0 0-1.705H7.953a1.704 1.704 0 0 0-2.953 0h1.477v1.705H5c.295.51.846.852 1.477.852m-.728 3.71a.43.43 0 0 0-.125.301v3.233c0 .236.191.426.426.426H7.33c.235 0 .426-.19.426-.426v-.426h5.114v.426c0 .236.191.426.426.426h1.279c.235 0 .426-.19.426-.426v-3.233a.43.43 0 0 0-.125-.301l-.727-.727-1.155-1.924a.43.43 0 0 0-.365-.207H7.996a.43.43 0 0 0-.365.207l-1.154 1.924zm7.66-.301H7.215l1.023-1.705h4.149zm-4.802 1.558v.317a.256.256 0 0 1-.255.255h-1.62a.256.256 0 0 1-.255-.255v-.711c0-.161.147-.283.305-.25l1.551.31a.34.34 0 0 1 .274.334m5.54-.29v.607a.256.256 0 0 1-.255.255h-1.62a.256.256 0 0 1-.255-.255v-.387c0-.122.086-.227.205-.25l1.518-.304a.34.34 0 0 1 .408.334"
        />
    </svg>
);
export default SvgCarRepair;
