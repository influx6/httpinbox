# HttpInbox
This demos a [WebHookInbox](http://webhookinbox.com/) app built in go.

## Install

```bash
go get -u github.com/influx6/httpinbox/...
```

## Deployment
The app uses environment variables as its core means of deployment which allows
us to store sensitive details away securely. It requires these environment
variables to be set before attempting to run its binary or main file.

  ```
  HTTPINBOX_HOST_ADDR
  ```
  Sets the address which the application http server will be started with.

  ```
  HTTPINBOX_MONGO_HOST
  ```
  Sets the mongo host address which will be used to communicate with the underline
  mongo database.

  ```
  HTTPINBOX_MONGO_AUTHDB
  ```
  Sets the mongo authentication database name which will be used to authenticate
  the provided credentials.

  ```
  HTTPINBOX_MONGO_DB
  ```
  Sets the mongo database name which will be used to getting the needed database.

  ```
  HTTPINBOX_MONGO_USER
  ```
  Sets the mongo username credential for the mongo authentication process.

  ```
  HTTPINBOX_MONGO_PASS
  ```
  Sets the mongo username password credential for the mongo authentication process.

## Running
Once all these are set within the host environment or within the deployed container
the application can be started either by:

- Running the `main.go` file

```bash
> go run main.go
```

- Running the generated binary for the project.

```bash
> ./httpinbox
```
