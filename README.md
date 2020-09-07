# Shortname

Shortname is a tool built to make it easier to get to your favorite sites.

## Usage

First run the binary to create an initial config:


```
$ sudo ./shortname                       # need sudo since we're on port 80
```

Now, edit your `~/.shortnamerc` file and give it some appropriate mappings:

```
sites:
  drive: https://drive.google.com
  hn: https://news.ycombinator.com
```

Finally, go to `hn/` in your browser, and you should see a Hacker News redirect!

## Installation

Make sure you've got go installed, and go get the package:

```
$ go get github.com/pawalt/shortname
```

Once you've got the binary, make sure it's in your path, and you're ready to roll. I recommend starting it at login if your OS supports that.