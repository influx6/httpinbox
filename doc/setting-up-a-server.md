# Setting up a server
With the massive interest and love for the Go language, there has being a massive
shift of providers which provide each their own mix in how and what approach you
use to deploy your new shiny apps.

Altough this is not an exhaustive list of providers, it does showcase three providers
that I do believe provide quiet the interesting package for deployment.

- [DigitalOcean](https://www.digitalocean.com)
  DigitalOcean builds on the popularity and beauty that Docker provides by letting
  you leverage their great support and cloud platform to have your docker images
  deployed with ease. They allow you deploy your docker images as `droplets` to
  custom servers within desired regions with a simple one click solution, these
  lets you take care of all needed setup with in your image and have that running
  within minutes on their cloud servers anywhere.


- [Amazon AWS](https://aws.amazon.com)
  Amazon AWS provides a complete cloud computing solution which allows you to deploy
  your code up on to a AWS server which can either take advantage of their core
  services or simple use their computing capabilities only. But regardless they
  do provide bang for your buck and lets you create your own servers and deploy
  as you see fit.

- [Heroku](https://www.heroku.com)
  Heroku provides a one push, one stop deployment solution which takes care of all
  the cruft of deployment and server setup and just lets you push and have your
  code running with as little setup needed. By providing a free and paid packages
  which allows you to take advantage of more features they are a wildly popular
  solution for all.


We will b using Heroku has our deployment server because this allows us more control
and lets us see how Harrow can seriously remove alot of headaches from our deployment
lifes.
