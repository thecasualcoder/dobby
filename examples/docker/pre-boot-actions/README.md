# pre-boot-actions

**Scenario:** Fetch configuration before starting the application/container

## Approach

- Have a bootstrap script which will run the pre-boot action and then the run the user given (or default) command
- Copy that script into your image
- Use that script as the `entrypoint` for the container

## Running

```
docker-compose up
```
