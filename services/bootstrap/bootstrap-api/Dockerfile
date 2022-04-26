FROM devopsfaith/krakend:1.4
COPY ./config /etc/krakend/config/
ENV FC_ENABLE=1
ENV FC_SETTINGS="/etc/krakend/config/settings" 


CMD ["run", "-d", "-c",  "/etc/krakend/config/krakend.json", "-p", "8080" ]