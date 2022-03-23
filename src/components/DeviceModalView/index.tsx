import * as THREE from "three";
import { Suspense } from "react";
import { DDSLoader } from "three-stdlib";
import { Canvas, useLoader } from "@react-three/fiber";
import { OrbitControls } from "@react-three/drei";
// import { OBJLoader } from "three/examples/jsm/loaders/OBJLoader";
// import { MTLLoader } from "three/examples/jsm/loaders/MTLLoader";

THREE.DefaultLoadingManager.addHandler(/\.dds$/i, new DDSLoader());

const Scene = () => {
    // const materials = useLoader(MTLLoader, "Poimandres.mtl");
    const OBJLoader = require("three/examples/jsm/loaders/OBJLoader").OBJLoader;
    const obj = useLoader(OBJLoader, "ukama_node.obj", () => {
        // materials.preload();
        // loader.setMaterials(materials);
    });

    return <primitive position={[0, -2, 0]} object={obj} scale={0.025} />;
};

const DeviceModalView = () => {
    return (
        <div style={{ height: "50vh" }}>
            <Canvas>
                <pointLight position={[0, 0, 100]} />
                <pointLight position={[0, 0, -100]} />
                <Suspense fallback={null}>
                    <Scene />
                    <OrbitControls
                        minDistance={1.5}
                        maxDistance={10}
                        maxPolarAngle={Math.PI / 2}
                    />
                </Suspense>
            </Canvas>
        </div>
    );
};

export default DeviceModalView;
