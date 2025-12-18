# sysview

Command-line tool to list active network processes and connections. The tool uses `lsof` to query active network connections and displays all connections with their local and remote addresses.

```
Active Network Processes (48)
----------------------------------------------------------------------------------------------------
PID      COMMAND              PROTO    PORT     LOCAL ADDRESS        REMOTE ADDRESS       STATE
----------------------------------------------------------------------------------------------------
626      identitys            TCP      1024     fe80:1a::3437:40d... fe80:1a::f26f:da2... ESTABLISHED
626      identitys            TCP      1028     fe80:1a::3437:40d... fe80:1a::f26f:da2... ESTABLISHED
12602    node                 TCP      3000     ::1                  -                    LISTEN
12602    node                 TCP      3000     ::1                  ::1:52733            ESTABLISHED
68082    remotepai            UDP      3722     0.0.0.0              -                    -
585      rapportd             UDP      3722     0.0.0.0              -                    -
...
```

```
Total Connections: 44
TCP: 37
UDP: 7
IPv4: 34
IPv6: 10
Listening: 5
Established: 31
```

## Installation

### Build from Source

1. Clone the repository:

```bash
git clone https://github.com/stefanicjuraj/sysview
```

2. Build the binary:

```bash
go build -o sysview .
```

3. Install (choose one):

   **System-wide installation** (requires sudo):

   ```bash
   sudo cp sysview /usr/local/bin/
   ```

   **User installation** (no sudo required):

   ```bash
   cp sysview ~/go/bin/
   ```

   Make sure `~/go/bin` is in your PATH.

## Usage

```bash
sysview                    # Show all network connections (default)
sysview stats              # Show network statistics summary
sysview state ESTABLISHED  # Filter by connection state (LISTEN, ESTABLISHED, TIME_WAIT, CLOSE_WAIT, etc.)
```

### `stats`

Displays network statistics summary:

- Total connections
- TCP vs UDP counts
- IPv4 vs IPv6 counts
- Listening vs Established connections

### `state`

Filter connections by state:

- `sysview state LISTEN` - Show only listening ports
- `sysview state ESTABLISHED` - Show only established connections
- `sysview state TIME_WAIT` - Show connections in time wait state
- `sysview state CLOSE_WAIT` - Show connections in close wait state
- Or any other connection state
