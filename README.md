redash-exporter
=====

prometheus exporter for <a href="https://redash.io/">Redash</a>

# Overview

Convert redash monitoring endpoint `/status.json` to pormetheus metircs.  
See also: https://redash.io/help/open-source/admin-guide/maintenance#Monitoring  

# Settings

## environment variable
__required__  
Set your redash api-key to the environment variable `REDASH_API_KEY`  

## Command-line flags
```
-listen-address string
    The address to listen HTTP requests. (default ":9295")
-metricsInterval int
    Interval to scrape status. (default 30)
-redashHost string
    target Redash host. (default "localhost")
-redashPort string
    target Redash port. (default "5000")
-redashScheme string
    target Redash scheme. (default "http")
-redashVersion int
    redash version. (default 8)
```

# docker image  
https://hub.docker.com/r/masahata/redash-exporter/
