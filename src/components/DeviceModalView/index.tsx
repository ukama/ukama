import * as THREE from "three";
import { Suspense } from "react";
import { DDSLoader } from "three-stdlib";
import { Canvas } from "@react-three/fiber";
import { OrbitControls, useFBX } from "@react-three/drei";

THREE.DefaultLoadingManager.addHandler(/\.dds$/i, new DDSLoader());

const Scene = () => {
    const fbx = useFBX("node.fbx");
    return (
        <primitive
            position={[0, -2, 0]}
            scale={8}
            object={fbx}
            rotation={[Math.PI / 1, 0, 0]}
        />
    );
};

const DeviceModalView = () => {
    return (
        <div style={{ height: "70vh" }}>
            <Canvas>
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
            </Canvas>
        </div>
    );
};

export default DeviceModalView;
