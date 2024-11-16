# Real-Time-Data-Streaming-API-with-Golang-and-Redpanda

### Set Up Redpanda (Kafka) with Docker:

docker run -d --name=redpanda -p 9092:9092 vectorized/redpanda redpanda start --overprovisioned --smp 1 --memory 1G --reserve-memory 0M --node-id 0 --check=false --kafka-addr PLAINTEXT://0.0.0.0:9092 --advertise-kafka-addr PLAINTEXT://localhost:9092

It will pulls the Redpanda Docker image if it is not already downloaded.   
And run a redpanda container optimized for development with minimal resource consumption.  
Configures Redpanda to expose Kafka functionally on localhost:9092.
Bind Kafka client connections to the container's port 9092.

### Start the Real-Time Data Streaming API

Start this WSL terminal and run the follong code:  
go run main.go

For the security consideration, you shoud use set the header of the request as the Valid api key (Default is billyae) to use the API or you will be denied the access to all apis.
### Start a new stream (With the Valid_api_key billyae)

curl -X POST http://localhost:8080/stream/start -H "X-API-Key: billyae"

### Send data to a stream

curl -X POST http://localhost:8080/stream/{streamid}/send -H "X-API-Key: billyae" -d '{"key": "value"}'

### Receive data to a stream

curl -X GET http://localhost:8080/stream/{streamid}/results -H "X-API-Key: billyae"

It will create a process that will not terminate automatically. 

After this process is created, use the command above the send data to this stream and in the main WSL terminal you can see the data sent to this stream. 

![alt text](image-1.png)

### Check the metrics

curl http://localhost:8080/metrics -H "X-API-Key: billyae"

### Load Testing

ab -n 10000 -c 100 -H "X-API-Key: billyae" -p post_data.json http://localhost:8080/stream/{streamid}/send

### Test Data

#### unit_test

go test ./unit_test/kafka_test.go -v

#### integrated_test

Start the stream: go test -v -run TestStartStream ./integrated_test

Send the data via the stream: streamid = {streamid}  go test -v -run TestSendData ./integrated_test

Start a receiving stream: streamid = {streamid}  go test -v -run TestGetResults ./integrated_test

