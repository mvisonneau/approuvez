# Quickstart - An example usage of 'approuvez' for experimentation purposes

## Requirements

- **~15 min of your time**
- [ngrok](https://ngrok.com/) or something equivalent to be able to access the approuvez's server from internet
- Sufficient privileges over your Slack workspace to be able to create a new app

## Get your external endpoint information

I wont paraphrase the installation / usage of [ngrok](https://ngrok.com/) here but this is basically all you need
to do once you are settled:

```bash
~$ ngrok approuvez 8443
ngrok by @inconshreveable                                                                               (Ctrl+C to quit)

Session Status                online
Version                       2.3.35
Region                        Europe (eu)
Web Interface                 http://127.0.0.1:4040
Forwarding                    http://1bafa1e508da.eu.ngrok.io -> http://localhost:8443
Forwarding                    https://1bafa1e508da.eu.ngrok.io -> http://localhost:8443

Connections                   ttl     opn     rt1     rt5     p50     p90
                              39      0       0.00    0.00    115.00  312.19

HTTP Requests
-------------

POST /callbacks/slack          200 OK
```

## Slack app creation

You will need to create a new Slack application, it is a 4 steps process starting at https://api.slack.com/apps

- Create your app

![slack_app_create](./images/slack_app_create.png)

- Configure the interactivity request URL (Features / Interactivity & Shortcuts pane)

![slack_app_interactivity_url](./images/slack_app_interactivity_url.png)

- Configure the capabilities / scopes of the app (Features / Oauth & Permissions)
  - `chat:write`
  - `users:read`
  - `users:read.email`

![slack_app_config_scopes](./images/slack_app_config_scopes.png)

- From the same page, you also need to write down the oauth-token of your app

![slack_app_oauth_token](./images/slack_app_oauth_token.png)

## Let's roll

If you haven't installed the `approuvez` binary yet, you can refer to this [install section first](/README.md#Install).

```bash
# Disable mTLS for the time being
~$ export APPROUVEZ_TLS_DISABLE=true

# Configure the Slack token you received during the previous step
~$ export APPROUVEZ_SLACK_TOKEN=xoxb-1234567893-xxxxxxxxxxxx-lHLknn4jxufxghpXtVGe63H4

# Let's get this thing a bit more verbose than necessary
~$ export APPROUVEZ_LOG_LEVEL=debug

# Start the server, you should get something like this:
~$ approuvez serve
INFO[2021-01-22T15:51:20Z] started multiplexed HTTP/gRPC server on :8443

# Now, whether open another terminal or put the previous command in the background
# In case of a new terminal, disable mTLS and make it verbose for the client as well
~$ export APPROUVEZ_TLS_DISABLE=true; export APPROUVEZ_LOG_LEVEL=debug

# Then ask yourself a message :)
~$ approuvez ask -m "coucou" -u "foo@bar.baz"
DEBU[2021-01-22T15:57:23Z] establishing gRPC connection to the server..  endpoint="127.0.0.1:8443"
DEBU[2021-01-22T15:57:23Z] gRPC connection established                  
INFO[2021-01-22T15:57:23Z] session initiated successfully                session_id=c6db27ab-c530-4838-b90f-5f20ec9d9123
INFO[2021-01-22T15:57:23Z] message sent, waiting for user's decision    

# As the log suggests, the process should now be blocking, waiting for an interaction
# In the server's logs you should see something like this:
INFO[2021-01-22T15:57:23Z] new session initialized                       client_endpoint="127.0.0.1:64449" session_id=c6db27ab-c530-4838-b90f-5f20ec9d9123
DEBU[2021-01-22T15:57:23Z] fetching slack user                           user=foo@bar.baz
DEBU[2021-01-22T15:57:23Z] looking up slack user                         used-method=email user-string=foo@bar.baz
DEBU[2021-01-22T15:57:23Z] prompting slack user                          session_id=c6db27ab-c530-4838-b90f-5f20ec9d9123 user_id=UXX0X0X0
```

Checkout slack and make a decision, you should see a message like this:

![default_message](./images/default_message.png)

If you click on any button you should see the message evolving:

![denied_message](./images/denied_message.png)

Looking back at the client logs, you should see the process being exited with an `exit 1` as the denied button got clicked.

```bash
# Full the client logs
~$ approuvez ask -m "coucou" -u "foo@bar.baz"
DEBU[2021-01-22T15:57:23Z] establishing gRPC connection to the server..  endpoint="127.0.0.1:8443"
DEBU[2021-01-22T15:57:23Z] gRPC connection established                  
INFO[2021-01-22T15:57:23Z] session initiated successfully                session_id=c6db27ab-c530-4838-b90f-5f20ec9d9123
INFO[2021-01-22T15:57:23Z] message sent, waiting for user's decision    
INFO[2021-01-22T15:59:58Z] received response                             decision=DENY user_id=UXX0X0X0 user_name=foo
DEBU[2021-01-22T15:59:58Z] exited..                                      execution-time=2m34.907287907s
exit status 1
```

That is it, now feel free to play with the additionnal flags like `--link-{name,url}` or add markdown to your message to customize them a bit further. eg:

```bash
~$ approuvez ask \
     -m "*terraform/project* | _aws/dev_\n>About to update \`x\` resources" \
     -u "foo@bar.baz" \
     --link-name ":gitlab: :terraform: #1523" \
     --link-url "https://gitlab.com/terraform/project/-/jobs/1523"
```

![enhanced_message](./images/enhanced_message.png)

If you are sastisfied with the tool and emphaze a more serious usage of it I would recommend to look onto how to secure it [using mutual TLS](./mutual_tls.md).
