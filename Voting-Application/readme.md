# Final Project - Voting Application

## 1. Where is the Go code?
In /db, /poll-api, /voter-api and /votes-api. /db is a module shared by the other three modules. It provides a unified interface for api to communicate with redis.

## 2. Whare are the Dockerfile and Compose file?
They are all in the same directory, /Voting-Application.

## 3. How to host the containers?
You **don't** need to build containers one by one before compose. You can directly run 'docker compose up' in /Voting-Application, or run the script **setup.sh**.

## 4. How to test the containers?
After all the containers are running, you can run the script **tests.sh**. It will output test results into standard output. There is an example output for referencing, **example-test-log.txt**.

I also make the Redis GUI port public (8001). You can run the tests in **tests.sh** one by one, and see how they change the data in redis database.
