#!/bin/sh

# RUN THIS SCRIPT FROM THE ROOT OF THE PROJECT TO GIVE PERMISSIONS
#  chmod +x build_all.sh 
echo "Building registry..."
cd ../registry/org && make
cd ../users && make
cd ../node && make
cd ../network && make
cd ../api-gateway && make

cd ../

echo "Building subscriber..."
cd ../subscriber/sim-pool && make
cd ../sim-manager && make
cd ../registry && make
cd ../test-agent && make
cd ../api-gateway && make

cd ../

echo "Building data-plan..."
cd ../data-plan/base-rate && make
cd ../package && make
cd ../rate && make
cd ../api-gateway && make

cd ../

echo "Building services..."
cd ../services/msgClient && make
cd ../services/initClient && make