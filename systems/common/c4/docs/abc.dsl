workspace "Ukama" "This is a model of my software system." {

    model {
        user = person "User" "A user of my software system."
        
        dev = person "Ukama Developer" "Ukama dev team."
        
        node = softwareSystem "Node" "Node"{
            lwm2mClient = container "lwm2m Client"
        }

    

        containerRegistry = softwareSystem "ContainerRegistry" "Container Registry all Ukama images"{
            registry = container "OCI Images registry"
        }

        ukamaCloud = softwareSystem "UkamaCloud" "Ukama Cloud. System for managing and monitoring network of Ukama devices                - Must have static IP" "oss" {
            
            dashboard = container "Ukama dashboard" "Ukama Dashboard" "React" "WebFrontend"
            mobile = container "Ukama Mobile App" "Ukama Mobile App" "Iphone/Android" "MobileApp"
            nodeRegistry = container "NodeRegistry" "NodeRegistry"
            controller = container "Controller" "Node Controller"
            kpiStorate = container "Kpi Storage" "Node Kpi storage. Prometheus"
            lwm2m = container "Lwm2m service" "Lightwaight machine to machine service"
            lwm2mGateway = container "Lwm2m Gateway"
            queue = container "Queue Service" "RabbitMQ" "Database"
            apiGateway = container "API Gateway"
            identityGateway = container "Authentication service"
        }

        bootstrap = softwareSystem "BootstrapSystem" "Ukama Device Bootstrap"{
            url  "http:\\bootstrap.ukama.com"
            bootstrapService = container "Bootstrap service" "" "REST, HTTP"
            dmr = container "Device Manufacturing Registry" "" "Lambda"
            deviceManufacturingRegistryDB = container "Device Manufacturing DB" "" "Database" "Database"
            uuidToOrgDb = container "UUID to Org and Org to Server" "Stores mapping between Organisation and Server IPs(NodeRegistry) and UUID and Org" "RDBMS Database" "Database"
            bootstrapApiGateway = container "Interface for Clouds services" "" "REST, HTTP"
        }

        deviceFactory = softwareSystem "Device production facility" {
            ukamaFactoryClient = container "Tool that works on the assembly line"
            -> bootstrap "Adds uuid to" "Rest/GRPC"
        }
        
        user -> ukamaCloud "Uses"
        user -> node "Connects internet and power to"
        node -> bootstrap "comunicates with"
        node -> ukamaCloud "Sends data to"
        ukamaCloud -> node "Manages configuration of"
        node -> bootstrap "Requests server IPs and Certs from"
        ukamaCloud -> bootstrap "Adds node UUID and org IP and Certs to " 
        ukamaCloud -> containerRegistry  "Pulls images from"
        dev -> containerRegistry "Pushes images with container updates to"        
        user -> dashboard "Manages nodes in"
        dashboard -> apiGateway "Makes API Calls to" "JSON/HTTPS"
        apiGateway -> nodeRegistry "Makes API calls to" "GRPC"
        apiGateway -> identityGateway "Makes calls to"
        node -> bootstrapService "Makes API calls to" "JSON/HTTPS"
        bootstrapService -> dmr "Makes API calls to" "JSON/HTTPS"
        bootstrapService -> uuidToOrgDb "Queries"
        controller -> queue "Reads messages from" "RabbitMQ API"
        

    }

    views {
        systemContext ukamaCloud "SystemContext" "A System Context diagram." {
            include  *
            autoLayout
        }


        container bootstrap "BootstrapContainers" {
            include *          
            autoLayout
        }

        container bootstrap "UkamaCloudContainers" {
            include *          
            autoLayout
        }


        dynamic bootstrap "NodePowerOnNoOrg"{
            node -> bootstrapService "Requests Certs and Server IPs from"
            bootstrapService -> dmr "Requests node info by UUID from"
            dmr -> bootstrapService "Returns UUID and device manufacturing data"
            bootstrapService -> uuidToOrgDb "Request Server IP and Cert by UUID"
            uuidToOrgDb -> bootstrapService "Returns NoServerIP and Cert response"
            bootstrapService -> node "Returns DEVIECE_NOT_IN_ORG respose"
        }

        dynamic bootstrap "NodeClaimProcedure"{
            user -> dashboard "Creates an organisation"
            dashboard -> apiGateway "Sends request to add new network to"
            apiGateway -> nodeRegistry "Adds organisation and certs to"
            nodeRegistry -> bootstrapApiGateway  "Adds organisation id and certs to"
            bootstrapApiGateway -> uuidToOrgDb "Adds Org and server IP record"

            user -> dashboard "Claims device by UUID" 
            dashboard -> apiGateway "Send claim device request to"
            apiGateway -> nodeRegistry "Requests to add UUID to Org"
            nodeRegistry -> bootstrapApiGateway "Adds an UUID to Org"
            bootstrapApiGateway -> uuidToOrgDb "Adds UUID to Org record"

            node -> bootstrapService "Requests Certs and Server IPs from"
            bootstrapService -> node "Returns IP and certs"
             autoLayout
            
        
        }
        

        dynamic bootstrap "NodePowerOn"{
            
            node -> bootstrapService "Requests Server IPs and Certs by sending UUID from"
            bootstrapService -> dmr "Requests node info by UUID from" 
            dmr -> bootstrapService "Returns node info if exist" 
            bootstrapService -> uuidToOrgDb "Requests server IPs and Certs" 
            uuidToOrgDb -> bootstrapService "Returns server IPs(Cloud Node Registry by default) and certs"
            bootstrapService -> node "Returns server IPs and certs" 
            node -> nodeRegistry "Initiate connetion with"
            
            nodeRegistry -> node  "Accepts connetion"             
            autoLayout
        }
        

        styles {
            element "Software System" {
                background #1168bd
                color #ffffff
            }

            element "Person" {
                shape Person
                background #08427b
                color #ffffff
            }

            element "Database" {
                shape Cylinder
            }

            element "WebFrontend"{
                shape WebBrowser
            }

            element "MobileApp"{
                shape MobileDeviceLandscape
            }

            element "oss"{
                color #08420a
            }
            
        }
    }
    
}
