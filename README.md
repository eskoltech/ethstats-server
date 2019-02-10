# ethstats
[![License](https://img.shields.io/badge/License-GPLv3%202.0-brightgreen.svg?style=for-the-badge)](https://www.gnu.org/licenses/gpl-3.0)

>**Note**: in early development

Websocket server to report Ethereum node stats.


## Build

If you want to build the `ethstats` server by yourself or start hacking on it,
just install this repo using the command:

```bash
$ go get -v github.com/eskoltech/ethstats
```

Also, you can use the `make` utility to build an executable. Using `make b` will
compile it, and the resulting binary can be found in the `build/bin/` directory.

## Running

In order to start the `ethstats` server, you need to provide a `secret`. This secret will 
be used to authorize nodes to report stats to this server. Note that if the Ethereum node 
can't be logged into the server, the server can't receive any notifications from any node.
For example, to start a server with default network options and a weak secret, just execute:

```bash
$ ethstats --secret 1234
```
>**Note** that for default, the server is started at `localhost:3000` when is not started using `make start`

You can view the default network options using the `-h` flag, and customize it for
you requirements. Also, you can start the server using the `make start` command, and customize 
the make flags to adapt it to your needs. You can modify this flags:

| **Variable** 	| **Description** 	| **Default**         	|
|--------------	|-----------------	|---------------------	|
| `ADDR`       	| Server address  	| `192.168.0.20:3000` 	|
| `SECRET`     	| Server secret   	| `123456789`         	|
>If you want to update the `ADDR` variable, execute the make task like `make start ADDR=10.0.0.15:3000`

If all is right, you will see some like this:

```
        __  .__              __          __
  _____/  |_|  |__   _______/  |______ _/  |_  ______
_/ __ \   __\  |  \ /  ___/\   __\__  \\   __\/  ___/
\  ___/|  | |   Y  \\___ \  |  |  / __ \|  |  \___ \
 \___  >__| |___|  /____  > |__| (____  /__| /____  >
     \/          \/     \/            \/          \/  v0.1.0

INFO[2019-02-10T17:19:23+01:00] Starting websocket server in 192.168.0.15:3000
INFO[2019-02-10T17:19:23+01:00] Node relay started successfully
INFO[2019-02-10T17:19:23+01:00] Server started successfully

```

Now you can attach nodes to report stats to this server using the address and port where 
the server is listening.
