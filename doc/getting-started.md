# Building your own HttpInbox
This article is about showing how you can build your own application using Go and take it from completion to deployment using Harrow to simplify your life, this allows you a insight into how you can combine different technology with Harrow to seamlessly move from product iteration to product continuous deployment and integration.

We will be building a simple application that logs http requests coming
to it by allocating `inboxes` with URLs that allows you capture and store which requests are hitting your servers, this example will be very contrived but with this as a base more complex versions can indeed be built.

The concept is simple: A service that allows you to create a http request inbox where clients could hit and we can capture that requests
to allow us perform whatever metric or process we wish to do so with that data.

Our application really has a very simple structure and this is very
important, because the structure of our code has a massive effect on
our deployment strategy on a longer term.

```bash
~/.../httpinbox > tree -d .
.
├── app
│   ├── api.go
│   ├── datawriter.go
│   └── views
│       ├── all.tml
│       ├── layout.tml
│       ├── list.tml
│       └── single.tml
├── main.go
├── readme.md
└── vendor
    ├── github.com
```

We have a `app` directory where the controller code for our application with their view templates are stored and also a `main.go` file which will be used to both run our application and build our binary from when deploying.
