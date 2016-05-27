# Deployment of HttpInbox
We will be using Harrow to manage our deployment of our application from creating
a runnable docker image to getting our image running on our Heroku server, by
having such a critical, central system to manage all these process we can ensure
simplicity in our lives and a fast, easy way to quickly grasp our deployment,testing
and integration steps.


## Deployment with Docker
We will using docker images for creating a self contained runnable containers that
take our custom docker image and have that deployed quickly to our servers, these
allows us reduce alot of issues with environment by simply letting us set things up
and have the same setup running on numerous servers around the world.

For this we need to write a `dockerfile` which will be the base for creating our
application image for our application.


```dockerfile
FROM `golang:1.6`

RUN mkdir /app

ADD ./app/views /app/views
ADD ./httpinbox /app/httpinbox

RUN chmod +x /app/httpinbox

ENTRYPOINT /app/httpinbox
```

We also need to write a build script which will create the needed binary for our
application, then create the docker image and have that deployed to the dockerhub
server which will then allow us to easily pull that anywhere and to any server
for running within a container.
