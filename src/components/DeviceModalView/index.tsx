import * as THREE from "three";
import { Suspense } from "react";
import { DDSLoader } from "three-stdlib";
import { Canvas } from "@react-three/fiber";
import { OrbitControls, useFBX } from "@react-three/drei";

THREE.DefaultLoadingManager.addHandler(/\.dds$/i, new DDSLoader());

const Scene = () => {
    const fbx = useFBX("node.fbx");
    return <primitive position={[0, 1, 0]} scale={8} object={fbx} />;
};

const DeviceModalView = () => {
    return (
        <div style={{ height: "70vh" }}>
            <Canvas>
                <pointLight position={[0, 50, 100]} intensity={0.5} />
                <pointLight position={[0, -50, -100]} intensity={0.2} />
                <Suspense fallback={null}>
                    <Scene />
                    <OrbitControls minDistance={1.5} maxDistance={10} />
                </Suspense>
            </Canvas>
        </div>
    );
};

export default DeviceModalView;
