rm domains.ext MwServerCA.* MwClientCA.*

echo -n "authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = localhost
DNS.2 = 127.0.0.1
DNS.3 = 0.0.0.0
DNS.4 = ::1
DNS.5 = *.mw.io
DNS.6 = *.middleware.io" > domains.ext

openssl req -x509 -nodes -new -sha256 -days 1024 -newkey rsa:2048 -keyout MwClientCA.key -out MwClientCA.pem -subj "/C=US/O=Middleware Pvt. Ltd./O=Middleware-Client-CA/CN=middleware.client"
openssl x509 -outform pem -in MwClientCA.pem -out MwClientCA.crt

openssl req -new -nodes -newkey rsa:2048 -keyout MwServerCA.key -out MwServerCA.csr -subj "/C=US/O=Middleware Pvt. Ltd./O=Middleware-Server-Certificate/CN=middleware.server"
openssl x509 -req -sha256 -days 1024 -in MwServerCA.csr -CA MwClientCA.pem -CAkey MwClientCA.key -CAcreateserial -extfile domains.ext -out MwServerCA.crt

# To verify Certificate:
# openssl verify -CAfile cert/MwClientCA.pem -verify_hostname *.middleware.io cert/MwServerCA.crt