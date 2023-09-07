/* C4 Model for WIMC - Where Is My Content. */

workspace "Ukama-wimc" "Mode of WIMC - Where Is My Content" {

    model {
        client = softwareSystem "micro-services" "Device Containers"
        
        wimc = softwareSystem "WIMC" "daemon"{
            wimcd = container "wimc.d" "API for local clients, manage Agents and updates status CB" "daemon"
            agent = container "Agents" "One per method type (ftp, chunk etc.). Communicate with cloud provider to download the content" "process"
            db    = container "content-db" "Maintain database of content's local storage" "Database" 
        }

        ukamaCloud = softwareSystem "Ukama Cloud" "Ukama Cloud"{
            nodeRegistry = container "NodeRegistry" "NodeRegistry"
        }
        
        client      -> wimc "Ask for content" "HTTP/JSON"
        wimc        -> client "Content available" "HTTP/JSON"
        wimc        -> wimc "Lookup locally" "SQL"
        wimc        -> ukamaCloud "Request" "HTTP/JSON"
        ukamaCloud  -> wimc "Get Content" "FTP, Chunk, others"
                
    }

    views {
        systemLandscape {
            include * 
            autoLayout
        } 

        container wimc "WimcContainers" {
            include *          
            autoLayout
        }

        dynamic wimc "WIMC" {

          //  agent       -> wimcd "Register with supported method"
            client      -> wimcd "Where is my content?"
            wimcd       -> db "Lookup in its db"
            db          -> wimcd "Local path or 404"
            wimcd       -> client "Response with ID"
            wimcd       -> ukamaCloud "Ask for content"
            ukamaCloud  -> wimcd "CB URL & supporting method"
            wimcd       -> agent "Send CB URL to matching agent"
            agent       -> wimcd "Confirm transfer"
            agent       -> wimcd "Send periodic updates (ID)"
            wimcd       -> client "Content storage location"

        autoLayout
        }

    }    
}
