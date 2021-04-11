<img src="https://alkatronic.focustronic.com/images/alkatronic_logo.png" width="300" alt="Alkatronic">

## Purpose
This projects creates a CLI application that allows scraping the Alkatronic API and saves
the data in a Postgresql database. You can run one time to prepopulate the last `7`, `30`, or `90` days worth
of metrics on the Alkatronic site. You can also run as a daemon to continuously scrape the Alkatronic site to
to collect metrics as they come in. This project relies on the Alkatronic client API [here](https://github.com/codyja/alkatronic)


## Examples
1. This example pulls down the last 90 days worth of test records from Alkatronic and saves to the database:
```shell
docker run -it \
-e ALKATRONIC_USERNAME="user" \
-e ALKATRONIC_PASSWORD="pass" \
-e DB_CONNECTION_STRING="postgresql://postgres:password@db-host_here:5432/db-name" \
codyja/alkalog:0.0.2 alkalog -days 90
```
2. Run in as a long running service that pulls new records every 30 minutes:
```shell
docker run -d --name alkalog --restart=always \
-e ALKATRONIC_USERNAME="user" \
-e ALKATRONIC_PASSWORD="pass" \
-e DB_CONNECTION_STRING="postgresql://postgres:password@db-host_here:5432/db-name" \
codyja/alkalog:0.0.2 alkalog -d
```
