# log-monitoring
This a solution to monitor w3c-formatted http access log
## Prerequisites (Tested with these versions)
Valid Go installation 1.11+  
GNU Make  
GO Dep 0.5+  
Docker 18.06+  - only if you want to utilize the docker build
### Build
`make build`
### Run
`log-monitor --logPath pathToLog --rpsThreshold yourRpsThreshold` or simply  
`log-monitor` to read from `/tmp/access.log` with rpsThreshold `10`  
### Docker build and run
`make docker`
`docker run -it ruseinov/log-monitoring:latest -v /your/log/file:/tmp/access.log`
### How the architecture of this app could be made better
1. Being able to read from STDIN would make the dockerized version easy to run on any system like this:   
`cat your_log | docker run -it ruseinov/log-monitoring:latest` and would also allow for monitoring several logs at a time   
using multitail and such
2. Introducing template config for printer would allow for more output flexbility
3. Support more alert types via more generic alerting approach, possibly passing in log entries for that
4. Current implementation ignores timestamps as it assumes that all the incoming data has consistent timestamps,  
instead we could validate these timestamps to see if the monitoring solution is falling behind or the timestamps  
are off