Accepts GitHub webhooks, and forwards them to ElasticSearch to do whatever you
like with.

## GitHub Setup

1. Define a webhook with payload URL of `<yourdomain>/webhook`.
2. Select a Content Type of "application/json"
3. Define a webhook secret, and note it down to set as an environment variable

## Environment Variables

| Name                        | Description                                                                                                                           |
| --------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| `GITHUB_APP_WEBHOOK_SECRET` | GitHub webhook secret. See https://docs.github.com/en/developers/webhooks-and-events/securing-your-webhooks#setting-your-secret-token |
| `ELASTICSEARCH_INDEX`       | ElasticSearch Index to send events to.                                                                                                |
| `ELASTICSEARCH_URL`         | ElasticSearch URL.                                                                                                                    |
