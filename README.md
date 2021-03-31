### Setup

`make build` or `make all` to run tests then build
`chmod +x ./server`
`./server`

#### Config

Configuration defaults can be found in `config.json` and altered as desired
Default port is 8080, default storage directory is `./messages`

### Accessing API

Using basic auth:

**POST:** `curl -u <username>:<password> localhost:8080/api/saveMessage -d '{"message": "This is a test message"} -X POST`

**GET:** `curl -u <same username>:<password> localhost:8080/api/retrieveMessage`

### Comments

I did not get around to adding integration tests, as I prioritized thorough config options and good unit testing. Thank you for this opportunity, and I hope to hear from you soon!