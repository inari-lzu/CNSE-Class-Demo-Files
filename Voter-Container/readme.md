# Cloud Native Packaging

## General

1. My updated Voter API code is in **/Voter-Container/voter-api**

2. The dockerfile that containerizes my API is also in **/Voter-Container/voter-api**

3. The script to build my API into a container is **/Voter-Container/build-voter-api.sh**

4. The docker-compose file that starts my API and enables it to work with an off-the-shelf redis container, is **/Voter-Container/docker-compose.yaml**.

> I set backend network as internal, and unforward the ports of 6379 and 8001 in redis container. So you would not be able to check whether redis is working as we expected through redis-gui. If you want to use it, you have to edit the compose file.

> Also, I added some mock data. So when the service launch up, there would be two records of voters already in redis.

> If you want to test the end points, you can use postman with local agent. I shared my workspace so you can join and use my testing endpoints: [**link**](https://galactic-robot-304045.postman.co/workspace/My-Workspace~8585261e-51a3-40d8-8e35-27d45097193d/collection/23204427-79b9cb4a-f9f0-497d-893a-c8673177cb78?action=share&creator=23204427)

## Extra Credit
1. I pushed both the basic version and multi-platform version image to dockerhub. I created a public repository for them: [**link**](https://hub.docker.com/repository/docker/xf2000/cst680/tags)

2. The script for multi-platform building is **/Voter-Container/build-multiplatform.sh**, and I have already push the built images to my repository: [**link**](https://hub.docker.com/layers/xf2000/cst680/multiplatform/images/sha256-3257eca302fed3ea39cbfb179086e7c31b274ca4bd16538324cb2eeada610613?context=repo)
