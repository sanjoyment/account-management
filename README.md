# Account Management Tool
Mini project using Golang + Postgres.

# Help Link
https://gist.github.com/cecilemuller/9492b848eb8fe46d462abeb26656c4f8

# About Cert & Run Process
1. Run ```make cert``` to generate certificates.
2. Use server certs into following line: ```r.RunTLS(":8080", "cert/MwServerCA.crt", "cert/MwServerCA.key")```
3. Import ```cert/MwClientCA.crt``` into chrome-browser's SSL certificate Authorities list. 
4. Now run ```make run``` & hit ```http://localhost:8080/``` on browser.