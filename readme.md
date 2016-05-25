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
  HTTPINBOX_LISTEN
  ```
  Sets the address which the application http server will be started with.

  ```
  HTTPINBOX_DATA
  ```
  Sets the directory where all the inbox and their data
  will be stored on disk

  ```
  HTTPINBOX_VIEWS
  ```
  Sets the directory path where the views/template files
  are located


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
