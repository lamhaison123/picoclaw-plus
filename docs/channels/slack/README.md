# Slack

Slack is a leading global enterprise-grade instant messaging platform. PicoClaw utilizes Slack's Socket Mode for real-time, bidirectional communication without the need for public Webhook endpoints.

## Configuration

```json
{
  "channels": {
    "slack": {
      "enabled": true,
      "bot_token": "xoxb-...",
      "app_token": "xapp-...",
      "allow_from": []
    }
  }
}
```

| Field      | Type   | Required | Description                                               |
| ---------- | ------ | -------- | --------------------------------------------------------- |
| enabled    | bool   | Yes      | Whether to enable the Slack channel                       |
| bot_token  | string | Yes      | Slack Bot User OAuth Token (starts with `xoxb-`)          |
| app_token  | string | Yes      | Slack Socket Mode App Level Token (starts with `xapp-`)   |
| allow_from | array  | No       | Whitelist of user IDs; empty means allow all users        |

## Setup Process

1. Go to [Slack API](https://api.slack.com/) and create a new Slack App.
2. Enable Socket Mode and obtain the App Level Token.
3. Add Bot Token Scopes (e.g., `chat:write`, `im:history`, etc.).
4. Install the app to your workspace and obtain the Bot User OAuth Token.
5. Enter the Bot Token and App Token into the configuration file.
