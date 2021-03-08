# String encryptor - Exercise Task

### String Randomizer Service

Generate the list of random strings and sending them via POST request like JSON to Encryptor service, then get response with JSON of list of encrypted strings, output random encrypted string on the page.

### String Encryptor Service

Receive the list of random strings, encrypt them using SHA256 and the list of encrypted strings to Randomizer Service



### Running

You can configure the environment variables in `docker-compose.yml`

```
docker-compose build
docker-compose up
```

### Example of usage


Open link in browser `http://localhost:8001/getRandomHash?size=23` where `size` - length of random generated list of strings

