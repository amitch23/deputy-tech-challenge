# Deputy Technical Challenge Response

## Prereqs

Please install Go version 1.16.3

Clone this repo and navigate to this directory.

## What this tool does

This script takes as a command line flag a user Id and outputs a list of all the subordinate users of that user. 

It reads in the user roles and users data as JSON and assigns them to structs in order to handle the searching. 

## Tests

To run tests, cd into the directory this README is in, and run:

    go test

## To run

To run the tool and find subordinates of a user with id=1 (the default flag value if nothing is passed in on the command line), cd into the directory this README is in, and run:

    go run main.go 

To pass in a user ID via the command line, run:

    go run main.go -userid=<userId>

Your output (if successful :)) for finding all subordinates for user with Id=1 should look something like:

    [{5 Mary Manager 2} {4 Sam Supervisor 3} {2 Emily Employee 4} {3 Maya Employee 4} {6 Steve Trainer 5}]

If a userID that does not exist is passed in, you should see something like:

    "User not found..."

