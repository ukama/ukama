import os
files = [f for f in os.listdir('.') if os.path.isfile(f)]

    # do something

# list = ["bootstrap-bootstrap-api.yaml", "bootstrap-lookup.yaml","cloud-api-gateway.yaml",
# "cloud-lwm2m-gateway.yaml", "cloud-registry.yaml", "common-lib.yml"]

for l in files:
    if l.endswith(".yaml")  or l.endswith(".yml"):   
        print("[![%s-ci](https://github.com/ukama/ukamaX/actions/workflows/%s/badge.svg)](https://github.com/ukama/ukamaX/actions/workflows/%s)" 
            % (l[:l.index(".")], l,l))


