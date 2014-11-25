Packhunter
===========

**Packhunter** is a web app that allows you to use your followers from Product Hunt to organize top products, rather than only using the global rankings.

# Deployment

To deploy the backend, you will need to compile it for your production environment and upload the executable file and the template files to the server.

## /etc/packhunter/packhunter.gcfg

You will need this file whether you're developing locally or deploying to a server.

```
[CtaApi]
key = xxxxxxxxxx
[Host]
port = :80
path = /path/to/packhunter
```

Get your API key [from the CTA](http://www.transitchicago.com/developers/traintracker.aspx).

## systemd config

Edit `/lib/systemd/system/packhunter.service`:

```
[Unit]
Description=packhunter.co web service
After=syslog.target network.target

[Service]
Type=simple
ExecStart=/path/to/packhunter/packhunter.linux

[Install]
WantedBy=multi-user.target
```

Then, make a symbolic link:

```
ln -s /lib/systemd/system/packhunter.service /etc/systemd/system/packhunter.service
```

And start/enable your service:

```
systemctl start packhunter.service
systemctl enable packhunter.service
```

## compile/upload/deploy script

I deploy with a script like this:

```
echo "Compiling"
cd packhunter; GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o packhunter.linux; cd ..
echo "Uploading"
pwd
tar cvf packhunter.tar packhunter/packhunter.linux packhunter/template packhunter/static
ssh -l sirsean packhunter.co mkdir -p packhunter
ssh -l sirsean packhunter.co rm packhunter/packhunter.linux
scp packhunter.tar sirsean@packhunter.co:packhunter.tar
ssh -l sirsean packhunter.co tar xvf packhunter.tar
echo "Restarting"
ssh -l root packhunter.co systemctl restart packhunter.service
echo "Deployed"

```

Note that you will need to set up your system to cross-compile for your target environment. (I deploy to a 64-bit Linux server.)

