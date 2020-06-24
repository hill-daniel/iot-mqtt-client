# Test client for processing protocol buffers messages with AWS IoT Core

## Part of the codecentric blog post ["Processing protocol buffers messages with AWS IoT Core"](https://blog.codecentric.de/en/2020/06/processing-protobufs-with-iot-core) 

This is the client for sending protocol buffers messages to the AWS IoT Core MQTT endpoint.

## Usage
Will send a protocol buffers message to MQTT endpoint every 10 seconds.
This example code comes with an .env file. You need to supply following values (see blog article chapter "Configure environment variables
" for further reference)
* `IOT_ENDPOINT` - enter the device endpoint supplied by AWS IoT (see “Get your endpoint”). Looks something like this <value>.iot.eu-central-1.amazonaws.com
* `CA_CERT_PATH` - the absolute path of the downloaded AWS Root certificate e.g. AmazonRootCA1.pem
* `DEVICE_CERT_PATH` - absolute path for the device certificate e.g. <value>-certificate.pem.crt
* `DEVICE_KEY_PATH` - the private key for the certificate e.g. <value>-private.pem.key
* `TOPIC` - the topic you want to send the protobuf messages to and where the IoT Core rule should listen to (use the one we configured in the policy, see “Create policy”)

### Build and run
* `go build cmd/iot-prototype/main.go` build main
* `./main` execute main 
