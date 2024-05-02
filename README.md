# go-nethttp

Test project to demonstrate new `net.http` features in go version 1.22. It shows how to use new features of the HTTP router and how to implement middleware without external dependencies.

## Todo

- ~~DONE git initial docker container support~~
- IN PROGRESS move main.go starting file into cmd/server folder
    - for docker and /cmd/app support, review code from github.com/valantonini/go-coffee-service repo and make adjustments.
- ~~DONE docker container for development based on info from https://www.youtube.com/watch?v=zfNqp85g5JM ~~
    - ~~add `docker compose watch` support for improved Dev Experience.~~
- customer web UI 
- admin web UI

## Notes 

From vido: youtube.com/watch?v=zfNqp85g5JM

```shell
docker run -it -v ./terraform:/app infra
root@contID:/app# ls
main.tf terraform.tf
root@contID:/app# terraform init

Initializing the backend
....
``` 

Better approach

```shell
$ docker compose run
```

```yml
# file compose.yml
version: '3.8'
services:
  terraform:
    image: hashicorp/terraform:1.4.0
    volumes: 
      - ./terraform:/infra
    working_dir: /infra

  aws-cli:
    image: amazon/aws-cli:2.2.20
    volumes: 
      - /aws:/aws
```

```shell
docker compose run --rm terraform init
```

This approach elevates compose.yml file to becomes the source of truth for all of the tooling your team is using.

Based on previous statement probable db clients and other tools also has to be added the same way.

## Docker compose develop feature

Docker compose yml file can have optional develop/watch/action section where watch for changes in the file system will done and listed action will be performed.

## Docker compose watch feature

Watch build context for service and rebuild/refresh containers when files are updated.
Watches for any file updates on the current working directory. When the file is updated then 
application will be automatically redeployed into your stack.

Starting docker compose in the watch mode.

```shell
docker compose watch
```


```shell
docker compose watch product-service
```

```yaml
# file compose.yml
Services:
  frontend:
    build:
      context: ./chatly-web
      dockerfile: docker/dev.Dockerfile
    ports:
      - 5173:5173
    develop:
      watch:
        - action: sync
          path: ./chatly-web
          target: /usr/src/app
          ignore: 
            - node_modules/
        - action: rebuild
          path: ./chatly-web/package.json
  server:
    build:
      context: ./
    ports:
      - 3000:3000
    develop:
      watch:
        - action: rebuild
          path: ./chaty
          ignore: 
            - target/
    
```

# References

The project is inspired and based on an excellent video intro and git repository.

- https://www.youtube.com/watch?v=H7tbjKFSg58
- https://github.com/dreamsofcode-io/nethttp
