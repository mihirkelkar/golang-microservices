<h1> Instructions on building this service </h1>

1. Create a network.
`docker network myevents-network`

2. Start a mogodb container on the same network.
`docker container run -d --name events-db --network myevents-network mongo`

3. Start the kafka and zookeeper. Note how the network is already in the docker-compose file.
The docker compose file is in the ../docker  folder
`docker-compose up`

4. Change the events_config.json file to set the mongo-url to the container
and the kafka url to the kafka container

4. Build the docker file.
`docker image build -t myevents/eventservice .`

5. Run the container on this network
`docker run -p 5000:5000 --network myevents-network myevents/eventservice`
