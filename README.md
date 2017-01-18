### glb
A configurable Round Robbin Load Balancer/Reverse proxy written in Go.
#### Configuration File
The configuration file should be saved in the directory of the executable as `glb.json`

Document structure is defined in the `glb.schema.json` file. Below is a sample configuration. 
```json
{
  "Basic": false,
  "DisableKeepAlives": false,
  "IdleConnTimeoutSeconds": 10,
  "Host": {
    "Addr": "localhost",
    "Port": ":9090",
    "SslPort": ":8443"
  },
  "Registry": {
    "s1": {
      "v1": [
        "localhost:8080",
        "localhost:8081"
      ]
    }
  }
}
```
##### Explanation of Fields
* Basic determines the semantics of the registry object. If enabled only one service/version 
paring will be available. That service must be defined with the service name of "default" and 
the service version of "default". This allows the service/version qualifier requirement to be 
bypassed allowing for the original service URL to be utilized. 
* DisabledKeepAlives forces the transport object to dial each time a request is made. This 
effectively set the IdleConnTimeoutSeconds to a true zero. This will also cause a the load 
balancer to alternate on each request.
* IdleConnTimeoutSeconds is the number of seconds that an inactive transport should keep its 
connection alive. This will affect the behavior of the load balancing algorithm causing requests
to be directed to a single host until the timeout or other qualifing condition is met. 
* Host describes the load balancer properties.
    * Addr is the address that the server should bind listeners to. 
    * Port is the HTTP port the server will use. 
    * SslPort is the HTTPS port the server will use. If blank only HTTP will be used. 
* Registry is the data store that handles the service/name to address mappings. this is represented 
by a map of maps whose values are a slice of strings representing the addresses. The Keys are 
strings of the service and version. 
#### Operation and Functionality
The load balancer mechanism's behavior will changee based on the connection settings in the config. 
For the most part the operation is consistent with what is expected of a round robin balancer with 
the exception a high or unlimited IdleConnTimeoutSeconds value. If the timeout is unlimited then the 
dial will only be called as the Go HTTP package sees necessary (until the connection pool is full). 
This is because the balancing logic is contained within the dial function. 

The registry object is keyed by the service name first, then the service version and, the a list of 
arrayss containing the addresses of the appropriate service/verison mapping. No guarentee is made of
the ordering of the proxy target addresses.

The registry can be overridden with a struct that implements Registry.
#### To Do List
1) Multiplier on round robin counter threshhold.
2) Service endpoint operations
3) Endpoints for managing the registry
4) Endpoint to write and reload a new configuration


