# APOD Service

APOD Service at start checks if there is APOD data for today (if not, then it loads it). Then every 24 hours it will download APOD data.

## Usage

To start the postgresql container with [flyway](https://documentation.red-gate.com/fd/command-line-184127404.html) migrations and go service container.

```bash
make start
```

Go to browser

- `http://localhost:8080/today`
- `http://localhost:8080/list`
- `http://localhost:8080/bydate?date=` (paste date in format yyyy-mm-dd)
- `http://localhost:8080/storage/`

Or send requests in Swagger (without media storage)

- `http://localhost:8080/swagger/index.html#`

Service can be stopped 
```bash
make stop
```
or restarted.
```bash
make restart
```

## Running Tests
```bash
make test
```

## Running Linter
```bash
make lint
```
