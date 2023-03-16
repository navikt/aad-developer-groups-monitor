# Azure AD developer groups monitor

Provide member count metrics for the developer groups in Azure AD.

## Developer groups

The so called developer groups are dynamic groups that is used to grant access to for instance GitHub and/or Google Cloud Platform through a set of enterprise applications in Azure AD. The queries these dynamic groups use rely on information synced from various other systems, and these values have previously been changed without notification, which results in numerous developers losing access to critical systems.

The metrics are used to alerts the NAIS team if the member count changes with an abnormal amount, so that we can remediate the issue as quick as possible.

Currently we are monitoring the following groups:

| Object ID                              | Name                                |
|:---------------------------------------|:------------------------------------|
| `eb5c5556-6c9a-4e54-83fc-f70cae25358d` | Utviklere i IT utvikling (interne)  |
| `76e9ee7e-2cd1-4814-b199-6c0be007d7b4` | Utviklere i IT utvikling (eksterne) |
| `48120347-8582-4329-8673-7beb3ed6ca06` | Utviklere i Data (interne)          |
| `15f9ea54-1987-475c-a0d5-f0e1a0e3f811` | Utviklere i Data (eksterne)         |

## License

Console is licensed under the MIT License, see [LICENSE.md](LICENSE.md).
