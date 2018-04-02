# gooption
Scalable, Distributed, Low Latency, High Throughput, Extensible option pricing system.

### Road map

- [ ] grpc option pricer 
     - [ ] Service definition using protobuf
     - [ ] gobs service (price, greeks and implied vol)
     - [ ] Testing / Number validation 
- [ ] Monitoring
     - [ ] Prometheus integration including healthchecks
     - [ ] Grafana dashboard for memory and performance statistics
     - [ ] Custom counters
- [ ] CI / CD pipeline
     - [ ] Containerization with Docker
     - [ ] Orchestration with Docker Compose
     - [ ] CI / CD pipeline with Drone
     - [ ] Orchestration with Kubernetes
- [ ] Data management
     - [ ] Integration with Dgraph
     - [ ] Client to retrieve data and call gobs service
- [ ] grpc middleware
     - [ ] Authentication
     - [ ] Logging
     - [ ] Data validation
     - [ ] Reverse proxy
- [ ] Service discovery
     - [ ] Routing
     - [ ] Metadata for GUI discovery
- [ ] Front End
     - [ ] Angular front end to consume gobs service
     - [ ] Integration with dgraph
     - [ ] Integration with goquantlib
     

     
 
     
