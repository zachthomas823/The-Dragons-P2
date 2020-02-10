FROM ubuntu

EXPOSE 4000

WORKDIR /app

ADD logs/ /app/logs

ADD stats_firewall /app/stats_firewall

RUN mkdir web
ADD web/index.html /app/web/index.html

ADD nba_data.csv /app/nba_data.csv

CMD ["./stats_firewall"]