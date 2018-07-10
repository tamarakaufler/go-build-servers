#!/bin/bash
#sudo apt-get install libssl1.0.0 -y

export FQDN="localhost"
echo -------------------
echo FQDN is $FQDN
echo -------------------

openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
-keyout $FQDN.key -out $FQDN.crt \
-subj "/C=GB/ST=HT/L=STA/O=bumpyride/CN=$FQDN"

cat $FQDN.crt $FQDN.key > $FQDN.pem
openssl x509 -noout -subject -in $FQDN.crt
