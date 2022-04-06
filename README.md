# OneLogin Auth CLI Tool

## Usage

### List all profiles:
```bash
onelogin-auth list
```

### Configuration

The onelogin auth CLI expects a file config.yaml:

```yaml
onelogin:
  clientID: clientID of API credential with "Authentication only"
  clientSecret: client Secret of API credential
  accountName: onelogin account name
  durationSeconds: 28800 # duration of the credentials in seconds (or remove for the default of 3600)
accounts:
  - name: myapp-prod
    appID: onelogin app id (e.g. 123456)
    accountID: AWS account ID
    profileName: AWS IAM profile to store credentials in (in ~/.aws/credentials)
roles:
  - iam-role-1 # role that is configured in onelogin and IAM to use with the onelogin identity provider
  - iam-role-2
defaultRegion: us-east-1
```

### Login

```
onelogin-auth login
```

You can also list the roles and accounts

Example:
```
$ onelogin-auth list

Roles:
[0] admin
[1] readonly
Accounts:
[0] myapp-prod

$ onelogin-auth login 1 0 
```
This example will make you login into the `myapp-prod` account with the `readonly` role.
