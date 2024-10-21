# PullSync

### Description
Sync notifications of `Github Pull Request` events to your `Slack`.

### Tools
Monorepo powered by [NX](https://nx.dev/)

[Golang Getting Started](https://github.com/nx-go/nx-go)


✨ **Create a GO library** ✨

nx g @nx-go/nx-go:library `<name>` --directory=library/go



✨ **Create a GO application** ✨

nx g @nx-go/nx-go:application `<name>` --directory=app

```
To remove:
nx g rm `<name>`
```

```
To remove:
nx g rm library-backend-<name>
```

### Github Webhook events

* Check runs
* Check suites
* Commit comments
* Discussion comments
* Issue comments
* Pull request review comments
* Pull request review threads
* Pull request reviews
* Pull requests

### Slack Oath & Permissions (Scopes)

You need to create an [app](https://api.slack.com/apps) then add the following scopes:

* channels:history
* channels:join
* channels:read
* chat:write
* incoming-webhook
* reactions:read
* reactions:write


### Development

You must have `docker` and `docker compose` installed.
Prerequisite is we need to run `dynamodb` locally.

```
docker compose up -d
```

The `api` service can be run directly as `REST API` or through `lambda` function using `SAM CLI`.

As `API`:

```
export PORT=:3000
nx serve api
```

As `lambda` function:

```
nx lambda.build api (optional)
nx lambda.serve api
```