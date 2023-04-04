FROM continuumio/miniconda3


EXPOSE 8000

RUN : "---------- install generic build container deps ----------" \
    && set -x \
    && apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get -qy install --no-install-recommends \
    build-essential libgdal-dev git g++ cmake libbz2-dev imagemagick libz-dev

RUN mkdir -p /var/www/server \      
    && cd /var/www/server \        
    && git clone https://github.com/Cloud-RF/Signal-Server.git
RUN cd /var/www/server/Signal-Server/src && cmake . && make

COPY . /var/www/app
WORKDIR /var/www/app
RUN conda create --name myenv 
SHELL ["conda", "run", "-n", "myenv", "/bin/bash", "-c"]
RUN conda install --file requirements.txt -c conda-forge
RUN ["chmod", "+x", "/var/www/app/docker/api/startup.sh"]
