redash-exporter
=====

prometheus exporter for <a href="https://redash.io/">Redash</a>

# settings
__required__  
Set your redash api-key to the environment variable `REDASH_API_KEY`  

# defaults
listen addr `:9295`  
redash server scheme `http`  
redash server host `localhost`  
redash server port `5000`  

can see command line flags

```
./redash-exporter -h
```

docker image  
https://hub.docker.com/r/masahata/redash-exporter/
