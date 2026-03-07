# WeCom Self-Built App

A WeCom self-built application is an app created by an enterprise within WeCom, primarily for internal use. These apps enable efficient communication and collaboration with employees, improving overall productivity.

## Configuration

```json
{
  "channels": {
    "wecom_app": {
      "enabled": true,
      "corp_id": "wwxxxxxxxxxxxxxxxx",
      "corp_secret": "YOUR_CORP_SECRET",
      "agent_id": 1000002,
      "token": "YOUR_TOKEN",
      "encoding_aes_key": "YOUR_ENCODING_AES_KEY",
      "webhook_path": "/webhook/wecom-app",
      "allow_from": [],
      "reply_timeout": 5
    }
  }
}
```

| Field            | Type   | Required | Description                                               |
| ---------------- | ------ | -------- | --------------------------------------------------------- |
| corp_id          | string | Yes      | Enterprise ID (CorpID)                                    |
| corp_secret      | string | Yes      | Application Secret                                        |
| agent_id         | int    | Yes      | Application Agent ID                                      |
| token            | string | Yes      | Callback Verification Token                               |
| encoding_aes_key | string | Yes      | 43-character AES Key                                      |
| webhook_path     | string | No       | Webhook Path (Default: /webhook/wecom-app)                |
| allow_from       | array  | No       | Whitelist of User IDs                                     |
| reply_timeout    | int    | No       | Reply Timeout (seconds)                                   |

## Setup Process

1. Log in to the [WeCom Management Backend](https://work.weixin.qq.com/).
2. Go to "App Management" -> "Create App".
3. Obtain the Enterprise ID (CorpID) and App Secret.
4. In the app settings, configure "Receive Messages" to obtain the Token and EncodingAESKey.
5. Set the Callback URL to `http://<your-server-ip>:<port>/webhook/wecom-app`.
6. Enter the CorpID, Secret, AgentID, etc., into the configuration file.

   Note: PicoClaw now uses a shared Gateway HTTP server to receive webhook callbacks for all channels, listening on 127.0.0.1:18790 by default. To receive callbacks from the public internet, please reverse proxy your external domain to the Gateway (default port 18790).
