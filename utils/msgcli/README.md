# msgcli tool.

the msgcli tool (for sending events to a running msgclient) is standard command line tool written in Go, so all you need is to build the binary (or go-run the code), but the binary gives you more flexibility as you can include it on your own scripts just as you will do with docker compose or git or else, just make sure you make the binary chmod executable (or if you prefer using a leading ./) and that it is path findable

some few points:
Sending events is mostly achieved under the msgcli events push commands. For example:
```
msgcli events push --org ukamatestorg --scope global --route inventory.accounting.accounting.sync -m '{
"Id":"2eb8d9e6-11f1-44ad-b018-5e41e5361588",
"Item": "Ukama saas fees",
"OpexFee": "199"
}'
```

Events structs fields are optional and when skipped,  random values will be assigned  to them automatically
All commands and sub-command do have self explanatory help commands, so most of the time you will find the answers by your own :stuck_out_tongue:
```
msgcli events push --help                                                                                                          <<<
The push command pushes events to a running service throught it's associated
message client.

Usage:
  msgcli events push [flags]

Aliases:
  push, p

Flags:
  -h, --help             help for push
  -m, --message string   message for the event (should be in json format)
  -o, --org string       name of the org to send the event to (default "ukamatestorg")
  -r, --route string     route for the event (should match "system.service.object.action")
  -s, --scope string     event scope. Must match one of the following: ["local" "global"] (default "local")

Global Flags:
      --config string   config file (default is $HOME/.msgcli.yaml)

```
Default values are configurable, through config file or environment variables
Right now, all events are not supported, as support for events are added on a needed basis. So if you want to support new events you will have to make some very minimalist and straightforward changes (one line of code + one struct to add overall)
Please add me as one of the reviewers for any changes you intend to make
If you have any question let me know
