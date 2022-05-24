import { Box } from "@mui/material";
import { Node_Type } from "../../generated";
import { NODE_IMAGES } from "../../constants";
// import * as THREE from "three";
// import { Suspense } from "react";
// import { DDSLoader } from "three-stdlib";
// import { Canvas } from "@react-three/fiber";
// import { OrbitControls, useFBX } from "@react-three/drei";

interface IDeviceModalView {
    nodeType: Node_Type | undefined;
}

// THREE.DefaultLoadingManager.addHandler(/\.dds$/i, new DDSLoader());

// const Scene = () => {
//     const fbx = useFBX("node.fbx");
//     return (
//         <primitive
//             position={[0, -2, 0]}
//             scale={8}
//             object={fbx}
//             rotation={[Math.PI / 1, 0, 0]}
//         />
//     );
// };

const DeviceModalView = ({ nodeType = Node_Type.Home }: IDeviceModalView) => {
    return (
        <Box
            component={"div"}
            sx={{
                height: { xs: "80vh", md: "62vh" },
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                marginTop: 6,
            }}
        >
            {/* <Canvas>
                <pointLight
                    color="white"
                    intensity={1}
                    position={[10, 10, 10]}
                />
                <pointLight
                    color="white"
                    intensity={0.5}
                    position={[10, -10, -10]}
                />
                <Suspense fallback={null}>
                    <Scene />
                    <OrbitControls minDistance={1.5} maxDistance={10} />
                </Suspense>
            </Canvas> */}
            <img
                src={NODE_IMAGES[nodeType]}
                alt="node-img"
                style={{ maxWidth: "100%", maxHeight: "500px" }}
            />
        </Box>
    );
};

export default DeviceModalView;
