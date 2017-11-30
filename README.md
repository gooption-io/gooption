# gooption
Scalable, Distributed, Low Latency, High Throughput, Extensible option pricing system.

### roadmap
- [ ] Use case 1 : BS option pricer written in go
- [ ] Use case 2 : BS option pricer using quantlib
- [ ] Market data dependency management
- [ ] CLI for scaffolding new service
     - [ ] gprc stub / mock generation
     - [ ] Functional tests generation 
     - [ ] Technical / Infra tests generation 
     - [ ] Request / Response generation
     - [ ] Cloud deploiement automation
- [ ] Docker integration
- [ ] Service discovery / configuration using consul.io
- [ ] Dgraph integration 
     - [ ] Data loading to query (snapshot) option service
     - [ ] Crawler to update market data for each tick 
     - [ ] Data loading to query (realtime) option service
- [ ] Machine learning 
     - [ ] Crawler allowing to create tensorflow graph from dgraph
     - [ ] Tensorlow use case
- [ ] Trade events / life cycle 
     - [ ] Integrate hyperledger 
     - [ ] Distributed trade creation / management
     - [ ] Distributed trade fixing
     - [ ] Corporate action
- [ ] Web app for option trade booking and risk management 
