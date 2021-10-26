redash-exporter
=====

prometheus exporter for <a href="https://redash.io/">Redash</a> version 9/10

# Overview

Convert redash monitoring endpoint `/status.json` to pormetheus metircs.  
See also: https://redash.io/help/open-source/admin-guide/maintenance#Monitoring  

# Settings

## required variable
Set your redash api-key to the environment variable `REDASH_API_KEY`  

## Command-line flags or environment variable
```
-listen_address string
    The address to listen HTTP requests. (default ":9295")
-metrics_interval int
    Interval to scrape status. (default 30)
-redash_host string
    target Redash host. (default "localhost")
-redash_port string
    target Redash port. (default "5000")
-redash_scheme string
    target Redash scheme. (default "http")
```

# docker image  
https://hub.docker.com/r/masahata/redash-exporter/
