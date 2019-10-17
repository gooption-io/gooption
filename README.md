# gooption

![Image of gopher option trader](gopher-gooption.png)

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
- [ ] QuantLib option pricer
     - [ ] Service definition using protobuf
     - [ ] goql service (price, greeks and implied vol)
     - [ ] Testing / Number validation
- [ ] Service discovery
     - [ ] Routing
     - [ ] Metadata for GUI discovery
- [ ] Front End
     - [ ] Angular front end to consume gobs service
     - [ ] Integration with dgraph
     - [ ] Integration with goql

### Setup

#### proto compiler

Download latest protobuf compiler stable version eg. protobuf-all-3.5.1.tar.gz from https://github.com/google/protobuf/releases/tag/v3.5.1
Extract content, cd in the folder then :

```
./configure
make
make check
sudo make install
protoc --version
```

#### gogo compiler

```
cd ~/go/src/github.com/gogo/protobuf
make
```

#### grpc gateway

```
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go
```

### gobs Docker image build
Ensure you're within the gooption root directory to get the right Docker build context
```
docker build -f gobs/Dockerfile . -t gooption/gobs:<BUILD> --no-cache
```