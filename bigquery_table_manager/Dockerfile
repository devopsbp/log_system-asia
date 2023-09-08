FROM python:3.11.1-slim-bullseye
RUN apt update \
&& apt install -y curl \ 
dnsutils \
iputils-ping \
netcat \
net-tools \
telnet

WORKDIR /home/application/bigquery_table_manager

COPY . .
RUN pip install -r requirements.txt

ENTRYPOINT ["python"]
CMD ["create_or_update_bigquery_tables.py"]