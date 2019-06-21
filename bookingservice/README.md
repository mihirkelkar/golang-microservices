1. Start the kafka and zookeeper. Note how the network is already in the docker-compose file.
The docker compose file is in the ../docker  folder
`docker-compose up`

2. Change the events_config.json file to set the mongo-url to the container
and the kafka url to the kafka container

3. Build the docker file.
`docker image build -t booking/bookingservice .`

4. Run the container on this network
`docker run -p 5000:5000 --network myevents-network booking/bookingservice`
