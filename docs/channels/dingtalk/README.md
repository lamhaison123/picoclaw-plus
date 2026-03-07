# DingTalk

DingTalk is Alibaba's enterprise communication platform, widely popular in Chinese workplaces. It uses a streaming SDK to maintain persistent connections.

## Configuration

```json
{
  "channels": {
    "dingtalk": {
      "enabled": true,
      "client_id": "YOUR_CLIENT_ID",
      "client_secret": "YOUR_CLIENT_SECRET",
      "allow_from": []
    }
  }
}
```

| Field         | Type   | Required | Description                                               |
| ------------- | ------ | -------- | --------------------------------------------------------- |
| enabled       | bool   | Yes      | Whether to enable the DingTalk channel                    |
| client_id     | string | Yes      | DingTalk Application Client ID                            |
| client_secret | string | Yes      | DingTalk Application Client Secret                        |
| allow_from    | array  | No       | Whitelist of user IDs; empty means allow all users        |

## Setup Process

1. Go to the [DingTalk Open Platform](https://open.dingtalk.com/).
2. Create an internal enterprise application.
3. Obtain the Client ID and Client Secret from the application settings.
4. Configure OAuth and event subscriptions (as needed).
5. Enter the Client ID and Client Secret into the configuration file.
