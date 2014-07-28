[![Pre-compiled binaries](http://img.shields.io/badge/Precompiled-Download-green.svg)](http://beta.gobuild.io/github.com/bjeanes/hk-deploy)
[![License](http://img.shields.io/badge/license-MIT-green.svg)](http://bjeanes.mit-license.org/)

# hk-deploy

A plugin to the fast Heroku CLI [`hk`](https://github.com/heroku/hk) for
deploying via the [Heroku Build API](https://devcenter.heroku.com/articles/build-and-release-using-the-api).

This allows deploying without a `git push` or even having Git installed.

## How?

It asks [the here-included backend service](/web) for a temporary set of
credentials for uploading to an S3 bucket, archives the requested directory,
uploads it, then kicks off a build by supplying a short-lived public URL for
the uploaded source archive.

## License

[MIT](bjeanes.mit-license.org)
