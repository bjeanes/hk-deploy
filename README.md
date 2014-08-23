[![Pre-compiled binaries](http://img.shields.io/badge/Precompiled-Download-green.svg)](http://beta.gobuild.io/github.com/bjeanes/hk-deploy)
[![License](http://img.shields.io/badge/license-MIT-green.svg)](http://bjeanes.mit-license.org/)

# hk-deploy

A plugin to the fast Heroku CLI [`hk`](https://github.com/heroku/hk) for
deploying via the [Heroku Build API](https://devcenter.heroku.com/articles/build-and-release-using-the-api).

This allows deploying without a `git push` or even having Git installed.

If you want something like this for the classic Heroku Toolbelt, see @jesperfj's
[heroku-git-flow](https://github.com/jesperfj/heroku-git-flow).

## Usage

```sh-session
$ hk deploy .
Creating source...done.
Compressing source...done.
Uploading source...done.
Building:

-----> Fetching custom git buildpack... done
...
```

## Install

### Pre-compiled Binaries

Pre-compiled binaries are available
[here](http://beta.gobuild.io/github.com/bjeanes/hk-deploy) but I've had mixed
luck with them (some memory-related `panic()`s are happening).

After unarchiving, stick the `deploy` binary in `/usr/local/lib/hk/plugin` or
your custom `$HKPATH`.

### Source install

Make sure you have Go (only 1.3 has been tested) installed.

```sh-session
$ go get github.com/bjeanes/hk-deploy
$ cd (go env GOPATH)/src/github.com/bjeanes/hk-deploy
$ go build
$ mkdir -p /usr/local/lib/hk/plugin # or any custom $HKPATH
$ mv ./hk-deploy /usr/local/lib/hk/plugin/deploy
$ hk help deploy
hk deploy: Deploy a directory of code to Heroku using the Build API.

Run "hk deploy DIRECTORY" to deploy the specified directory to Heroku.
```

## Hacking

```sh-session
$ go get github.com/bjeanes/hk-deploy
$ cd (go env GOPATH)/src/github.com/bjeanes/hk-deploy
```

I created a file `$HKPATH/deploy` with the following contents and `chmod +x`'d
it for easy testing:

```sh
#!/usr/bin/env sh

cd (go env GOPATH)/src/github.com/bjeanes/hk-deploy
go build *.go $*
```

YMMV

## License

[MIT](bjeanes.mit-license.org)
