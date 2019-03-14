# post-boot-actions

**Scenario:** Fetch configuration after starting the application

> Note: Start of application is different from start of container

## Approach

- Have a new container which will wait for application to start. We can use [dockerize](https://github.com/jwilder/dockerize#waiting-for-other-dependencies) to wait for our application to start.
- Once the application starts, run the post-boot action in the new container.

## Running

```
docker-compose up
```
