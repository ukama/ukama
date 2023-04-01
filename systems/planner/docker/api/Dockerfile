# FROM tiangolo/uvicorn-gunicorn-fastapi:python3.11
FROM continuumio/miniconda3


EXPOSE 8000


RUN : "---------- install generic build container deps ----------" \
    && set -x \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
    g++ cmake libbz2-dev imagemagick


COPY . /var/www/app
WORKDIR /var/www/app
ENV ENV=dev
# Activate the environment, and make sure it's activated:
RUN cd /var/www/app/Signal-Server/src \
    && cmake . \
    && make
RUN conda create --name myenv 
SHELL ["conda", "run", "-n", "myenv", "/bin/bash", "-c"]
RUN conda install -y --file requirements.txt -c conda-forge



RUN ["chmod", "+x", "/var/www/app/docker/api/startup.sh"]
# ENTRYPOINT ["/bin/bash"] 
ENTRYPOINT ["/var/www/app/docker/api/startup.sh"]