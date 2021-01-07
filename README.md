# Shortname

Shortname is a tool built to make it easier to get to your favorite sites.

## Usage

First run the binary to create an initial config:


```
$ sudo ./shortname -f $HOME/.shortnamerc        # need sudo since we're on port 80
```

Now, edit your `~/.shortnamerc` file and give it some appropriate mappings:

```yaml
sites:
  drive: https://drive.google.com
  hn: https://news.ycombinator.com
```

Finally, go to `hn/` in your browser, and you should see a Hacker News redirect!

You can also go to paths on your links! For example if I have:

```yaml
sites:
  gp: https://github.com/pennlabs/
```

I can go to `gp/kittyhawk/` in my browser, and it'll take me to `github.com/pennlabs/kittyhawk`!

Finally, if you want a list of all your sites, just go to `sn/` in your browser to get a json-serialized list of all your sites and where they redirect to.

## Installation

Make sure you've got go installed, and go get the package:

```
$ go get github.com/pawalt/shortname
```

Once you've got the binary, make sure it's in your path, and you're ready to roll. I recommend starting it at login if your OS supports that.

## Enabling on startup

Since shortname has to run as root, we have to start it up with an init process like systemd. If you're using systemd, drop this file at /etc/systemd/system/shortname.service:

```
[Unit]
Description=shortname
After=network.target
StartLimitIntervalSec=0
[Service]
Type=simple
Restart=always
RestartSec=1
User=root
ExecStart=/home/<YOUR_USER_HERE>/bin/shortname -f /home/<YOUR_USER_HERE>/.shortnamerc

[Install]
WantedBy=multi-user.target
```

Now start it and enable it at startup:

```
$ sudo systemctl enable shortname
$ sudo systemctl start shortname
```
