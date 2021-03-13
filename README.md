<img src="https://alkatronic.focustronic.com/images/alkatronic_logo.png" width="300" alt="Alkatronic">

## Purpose
This projects creates a CLI application that allows scraping the Alkatronic API and saves
the data in a Postgresql database. You can run one time to prepopulate the last `7`, `30`, or `90` days worth
of metrics on the Alkatronic site. You can also run as a daemon to continuously scrape the Alkatronic site to
to collect metrics as they come in. This project relies on the Alkatronic client API [here](https://github.com/codyja/alkatronic)


## Examples