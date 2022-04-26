export FC_ENABLE=1 
export FC_SETTINGS="$PWD/config/settings" 
export FC_PARTIALS="$PWD/config/partials"  
export FC_OUT="karken_final.json"
krakend check -t -d -c "$PWD/config/krakend.json" 