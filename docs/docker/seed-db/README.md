# seed-db

Scenario: Seed data after the Database is up.

## Approach

1. Start DB container
2. Start a helper container containg DB client and the data to be seeded
3. Make helper container wait for the DB to be up using [Dockerize](https://github.com/jwilder/dockerize)
4. Seed the data using the helper container

## Running

```bash
docker-compose up
```