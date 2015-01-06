TinyCI
======

Simplest possible CI server. Runs scripts in response to events.

All scripts should go in a scripts directory. Scripts directory can be defined by `TINYCI-SCRIPT-DIR` environment variable, or will simply be `scripts` directory located where TinyCI executable is located.

Event sources are:

1. Github webhooks:
-------
Simply point your github webhook at `http://yourTinyCIserver/gh`.

Will run scripts in this order (if they exist in scripts folder):

1. `gh-githubname.reponame.sh`
2. `gh-githubname.reponame~branch.sh`

For better security on webhooks you can set an environment variable called `github-hook-secret` with the same secret you supply to github when creating the webhook.
	
2. Docker hub webhooks:
----

Point docker hub webhooks to `http://yourTinyCIserver/dh`

Will run `dh-yourname.repo.sh` if it exists.

3. Git polling:
----