#!/bin/bash

# There are three sections of tests. Some of the tests
# need previous to create necessary data. So please
# run them in order. Otherwise some tests might fail.

# reset database
echo '<<<--- Reset Database' &&
curl -i --request DELETE 'http://localhost/votes' &&
curl -i --request DELETE 'http://localhost:1080/voters' &&
curl -i --request DELETE 'http://localhost:1081/polls'

# Test Section 1 - Poll API --------------------------------
echo $'<<<--- Test Section 1 - Poll API ---------------------------------\n'

echo '--->>> test 1.1 create poll_1' &&
curl --silent --location 'http://localhost:1081/polls/1' \
--header 'Content-Type: application/json' \
--data '{
    "id": 1,
    "title": "pet type",
    "question": "which type of pet do you like?",
    "options" : [
        {
            "id": 1,
            "value": "dog"
        },
        {
            "id": 2,
            "value": "cat"
        }
    ]
}' && echo $'\n'

echo '--->>> test 1.2 get poll_1' &&
curl --silent --location 'http://localhost:1081/polls/1' && echo $'\n'

echo '--->>> test 1.3 update poll_1' &&
curl --silent --location --request PUT 'http://localhost:1081/polls' \
--header 'Content-Type: application/json' \
--data '{
    "id": 1,
    "title": "animal type",
    "question": "which type of animal do you like?",
    "options" : [
        {
            "id": 1,
            "value": "fish"
        },
        {
            "id": 2,
            "value": "bird"
        }
    ]
}' && echo $'\n'

echo '--->>> test 1.4 delete poll_1' &&
curl -i --location --request DELETE 'http://localhost:1081/polls/1'

echo '--->>> test 1.5 create poll_2' &&
curl --silent --location 'http://localhost:1081/polls/2' \
--header 'Content-Type: application/json' \
--data '{
    "id": 2,
    "title": "paper type",
    "question": "paper down or paper up?",
    "options" : [
        {
            "id": 1,
            "value": "down"
        },
        {
            "id": 2,
            "value": "up"
        }
    ]
}' && echo $'\n'

echo '--->>> test 1.6 add option_3' &&
curl --silent --location 'http://localhost:1081/polls/2/options' \
--header 'Content-Type: application/json' \
--data '{
   "options" : [
        {
            "id": 3,
            "value": "whatever"
        }
    ]
}' && echo $'\n'

echo '--->>> test 1.7 get option_3' &&
curl --silent --location 'http://localhost:1081/polls/2/options/3' && echo $'\n'

echo '--->>> test 1.8 update option_3' &&
curl --silent --location --request PUT 'http://localhost:1081/polls/2/options' \
--header 'Content-Type: application/json' \
--data '{
   "options" : [
        {
            "id": 3,
            "value": "I don'\''t use paper"
        }
    ]
}' && echo $'\n'

echo '--->>> test 1.9 delete all the options of poll_2' &&
curl -i --location --request DELETE 'http://localhost:1081/polls/2/options'

echo '--->>> test 1.10 get all the options of poll_2' &&
curl --silent --location 'http://localhost:1081/polls/2/options' && echo $'\n'

echo '--->>> test 1.11 delete all the polls' &&
curl -i --location --request DELETE 'http://localhost:1081/polls'

echo '--->>> test 1.12 get all the polls' &&
curl --silent --location 'http://localhost:1081/polls' && echo $'\n'

echo '--->>> test 1.13 check poll-api health' &&
curl --silent --location 'http://localhost:1081/polls/health' && echo $'\n'

# Test Section 2 - Voter API --------------------------------
echo $'<<<--- Test Section 2 - Voter API ---------------------------------\n'

echo '--->>> test 2.1 create voter_1' &&
curl --silent --location 'http://localhost:1080/voters/1' \
--header 'Content-Type: application/json' \
--data '{
    "id": 1,
    "firstName": "Mike"
}' && echo $'\n'

echo '--->>> test 2.2 get voter_1' &&
curl --silent --location 'http://localhost:1080/voters/1' && echo $'\n'

echo '--->>> test 2.3 update voter_1' &&
curl --silent --location --request PUT 'http://localhost:1080/voters' \
--header 'Content-Type: application/json' \
--data '{
    "id": 1,
    "firstName": "Mike",
    "lastName": "F",
    "history": [
        {
            "id": 20,
            "date": "2023-07-25T16:47:26.3570871-04:00"
        }
    ]
}' && echo $'\n'

echo '--->>> test 2.4 delete voter_1' &&
curl -i --location --request DELETE 'http://localhost:1080/voters/1'

echo '--->>> test 2.5 create voter_2' &&
curl --silent --location 'http://localhost:1080/voters/2' \
--header 'Content-Type: application/json' \
--data '{
    "id": 2,
    "firstName": "Xiao",
    "lastName": "Fang"
}' && echo $'\n'

echo '--->>> test 2.6 add voterPoll record' &&
curl --silent --location 'http://localhost:1080/voters/2/polls' \
--header 'Content-Type: application/json' \
--data '{
    "history": [
        {
            "id": 1,
            "date": "2021-07-25T16:47:26.3570871-04:00"
        }
    ]
}' && echo $'\n'

echo '--->>> test 2.7 get voterPoll record' &&
curl --silent --location 'http://localhost:1080/voters/2/polls/1' && echo $'\n'

echo '--->>> test 2.8 update voterPoll record' &&
curl --silent --location --request PUT 'http://localhost:1080/voters/2/polls' \
--header 'Content-Type: application/json' \
--data '{
    "history": [
        {
            "id": 1,
            "date": "2023-07-25T16:47:26.3570871-04:00"
        }
    ]
}' && echo $'\n'

echo '--->>> test 2.9 delete the voteHistory of voter_2' &&
curl -i --location --request DELETE 'http://localhost:1080/voters/2/polls'

echo '--->>> test 2.10 get the voteHistory of voter_2' &&
curl --silent --location 'http://localhost:1080/voters/2/polls' && echo $'\n'

echo '--->>> test 2.11 delete all the voters' &&
curl -i --location --request DELETE 'http://localhost:1080/voters'

echo '--->>> test 2.12 get all the voters' &&
curl --silent --location 'http://localhost:1080/voters' && echo $'\n'

echo '--->>> test 2.13 check voter-api health' &&
curl --silent --location 'http://localhost:1080/voters/health' && echo $'\n'

# Test Section 3 - Votes API --------------------------------
echo $'<<<--- Test Section 3 - Votes API ---------------------------------\n'

echo '--->>> test 3.1 vote can not be create if the associated voter does not exist' &&
curl --silent --location 'http://localhost/votes/1' \
--header 'Content-Type: application/json' \
--data '{
    "id": 1,
    "voterId": 1,
    "pollId": 2,
    "choiceId": 2
}' && echo $'\n'

echo '--->>> test 3.2 vote can not be create if the associated Poll does not exist' &&
curl --silent --location 'http://localhost:1080/voters/1' \
--header 'Content-Type: application/json' \
--data '{
    "id": 1,
    "firstName": "Xiao",
    "lastName": "Fang"
}' && echo '' &&
curl --silent --location 'http://localhost/votes/1' \
--header 'Content-Type: application/json' \
--data '{
    "id": 1,
    "voterId": 1,
    "pollId": 2,
    "choiceId": 2
}' && echo $'\n'

echo '--->>> test 3.3 vote can not be create if the associated vote option does not exist' &&
curl --silent --location 'http://localhost:1081/polls/2' \
--header 'Content-Type: application/json' \
--data '{
    "id": 2,
    "title": "paper type",
    "question": "paper down or paper up?",
    "options" : [
        {
            "id": 1,
            "value": "down"
        },
        {
            "id": 2,
            "value": "up"
        }
    ]
}' && echo '' &&
curl --silent --location 'http://localhost/votes/1' \
--header 'Content-Type: application/json' \
--data '{
    "id": 1,
    "voterId": 1,
    "pollId": 2,
    "choiceId": 3
}' && echo $'\n'

echo '--->>> test 3.4 vote can be created if all the associated entities exist' &&
curl --silent --location 'http://localhost/votes/1' \
--header 'Content-Type: application/json' \
--data '{
    "id": 1,
    "voterId": 1,
    "pollId": 2,
    "choiceId": 2
}' && echo $'\n'

echo '--->>> test 3.5 one voter can only has one voting in a poll' &&
curl --silent --location 'http://localhost/votes/2' \
--header 'Content-Type: application/json' \
--data '{
    "id": 2,
    "voterId": 1,
    "pollId": 2,
    "choiceId": 1
}' && echo $'\n'

echo '--->>> test 3.6 one voter can change their opinion by updating the vote' &&
curl --silent --location --request PUT 'http://localhost/votes' \
--header 'Content-Type: application/json' \
--data '{
    "id": 1,
    "voterId": 1,
    "pollId": 2,
    "choiceId": 1
}' && echo $'\n'


echo '--->>> test 3.7 get the vote_1' &&
curl --silent --location --location 'http://localhost/votes/1'

echo '--->>> test 3.8 delete the vote_1' &&
curl -i --location --request DELETE 'http://localhost/votes/1'


echo '--->>> test 3.9 create vote_2' &&
curl --silent --location 'http://localhost/votes/2' \
--header 'Content-Type: application/json' \
--data '{
    "id": 2,
    "voterId": 1,
    "pollId": 2,
    "choiceId": 1
}' && echo $'\n'

echo '--->>> test 3.10 delete all the votes' &&
curl -i --location --request DELETE 'http://localhost/votes'

echo '--->>> test 3.9 get all the votes' &&
curl --silent --location 'http://localhost/votes' && echo $'\n'

echo '--->>> test 3.10 check votes-api health' &&
curl --silent --location 'http://localhost/votes/health' && echo $'\n'

# reset database again
echo '<<<--- Reset Database Again' &&
curl -i --request DELETE 'http://localhost/votes' &&
curl -i --request DELETE 'http://localhost:1080/voters' &&
curl -i --request DELETE 'http://localhost:1081/polls'