# wait-for-dependencies

Scenario: Wait for all the dependencies to be up before starting the application.

## Approach

Use [Dockerize](https://github.com/jwilder/dockerize#waiting-for-other-dependencies) to wait for the dependencies to be up and then start the application.

## Running

```bash
docker-compose up
```
