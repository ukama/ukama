# Foo Service-boilerplate

## How to use boilerplate 

1. Copy content to the service directory
2. Run to replace dummy servcice name with your service name  in code and files:
`find . -type f -name '*' | xargs sed -i 's/foo/$SERVICE_NAME_LOWERCASE/g'` 
`find . -type f -name '*' | xargs sed -i 's/Foo/$SERVICE_NAME_CAMELCASE/g'`  
`find . -depth -name '*foo*' -execdir bash -c 'mv -i "$1" "${1//foo/$SERVICE_NAME_LOWERCASE}"' bash {} \;`
3. Run 'make gen' and then 'make build' adding all missing dependencies

