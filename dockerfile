FROM `golang:1.6`

RUN mkdir /app

ADD ./app/views /app/views
ADD ./httpinbox /app/httpinbox

RUN chmod +x /app/httpinbox

ENTRYPOINT /app/httpinbox
