# LUGNotify
Sends meeting information to a specifed Matrix room

## Usage
Put script to call it in a Cron job or Systemd timer to run weekly on Wednesdays
Included in matrixbot.nix is my Systemd timer and service config.

### Using the command

1. Build binary

```bash
go build src/matrixbot.go
```

2. Create a matrix account for sending messages.
3. Add the account to either a private chat with yourself or a groupchat.
4. Note the internal room ID formatted as !UPPERCASEandLOWERCASEletters:matrix.org
5. Add to config as such
```json
{
    "username": "<account_username>",
    "password": "<account_password>",
    "internalRoomID": "<the_whole_room_id_you_noted>"
}

```
6. Specify the config file path and optionally redirect stderr to a log file
``` bash
matrixbot --config /etc/matrixbot/config.json &>> /etc/matrixbot/errors.log
```
If not specified, the config file will be sourced from the directory of invocation.

7. Profit

## Details
During Summer break, the program doesn't send a message. When there isn't a meeting during the school year it sends.

~~~
No Meeting Today
~~~

When there is a meeting that day, it sends a message such as the following

~~~
Topic :  Give thanks 
Presenter :  - 
Notes :  Fall break 
~~~
