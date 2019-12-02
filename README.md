# Infinite Etudes
*infinite-etudes* generates ear training exercises for instrumentalists.

You can run it from the command line (cli mode) or as a web server (server mode).

In cli mode, *infinite-etudes* generates a set of 7 midi files for each of 12 key
signatures. Each set covers all possible combinations of 3 pitches within the
key. The files are generated in the current working directory.

In server mode, *infinite-etudes* is a high-performance self-contained web server
that provides a simple user interface that allows the user to choose a key, a
scale pattern and an instrument sound and play a freshly-generated etude in
the web browser. A public demo instance is running at 

https://etudes.ellisandgrant.com

## Installation
You need to have Go installed to build and test infinite-etudes. Get it from https://golang.org/dl/ .

After installing Go, do

```
  go get github.com/Michael-F-Ellis/infinite-etudes
  cd $GOPATH/src/github.com/Michael-F-Ellis/infinite-etudes
  go test
  go install
```

Then run `infinite-etudes -h` for options and usage instructions.

## Serving HTTPS
When serving on port 443, you'll need to set 2 environment variables, IETUDE_CERT_PATH and 
IETUDE_CERTKEY_PATH to point to the certificates. If you're serving from a linux with systemctl,
a typical service file looks like the following:
```
[Unit]
Description=Infinite Etudes server
After=network.target

[Service]
Type=simple
User=mellis
WorkingDirectory=/home/mellis/ietudes
# Always attempt to renew the certificate before (re)starting infinite-etudes
ExecStartPre=+/usr/bin/certbot renew
# infinite-etudes needs two environment variables that give full paths to the certificate
# fullchain and key files.
Environment="IETUDE_CERT_PATH=/etc/letsencrypt/live/etudes.ellisandgrant.com/fullchain.pem"
Environment="IETUDE_CERTKEY_PATH=/etc/letsencrypt/live/etudes.ellisandgrant.com/privkey.pem"
# run infinite-etudes as an https server
ExecStart=/home/mellis/go/bin/infinite-etudes -s -p :443
# Ensure that the process is always restarted on failure or if terminated by a signal
# A 5 second restart delay is used to reduce the possibility of thrashing if
# something is badly wrong.
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```
In the above, you'll need to replace `mellis` and `etudes.ellisandgrant.com` with your
user name and host name, respectively.

Also, you'll need to arrange for the infinite-etudes executable to have permission to
read the CERT files and bind to port 443.  If you use something like
```
sudo setcap sudo setcap 'cap_net_bind_service,cap_dac_read_search=+ep' /home/mellis/go/bin/infinite-etudes
```
you'll need to re-run that command any time you rebuild the executable.