import { Suspense } from "react";
import { Canvas } from "@react-three/fiber";
import { OrbitControls, useFBX } from "@react-three/drei";

const Model = () => {
    const fbx = useFBX("/earth.fbx");
    return <primitive object={fbx} />;
};

const DeviceModalView = () => {
    return (
        <Canvas
            style={{ height: 400 }}
            camera={{ fov: 35, zoom: 4, near: 1, far: 1000 }}
        >
            <ambientLight intensity={1} />
            <Suspense fallback={null}>
                <Model />
            </Suspense>
            <OrbitControls
                makeDefault={true}
                minPolarAngle={0}
                maxPolarAngle={Math.PI / 1.75}
            />
        </Canvas>
    );
};

export default DeviceModalView;
