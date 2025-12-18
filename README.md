# sysview

Command-line tool to list active network processes and listening ports. The tool uses `lsof` to query active network connections and filters for processes in the LISTEN state, showing which ports are actively listening for incoming connections.

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
sysview
```

Lists all active network processes that are listening on ports, displaying:

- **PID**: process ID
- **COMMAND**: process name
- **USER**: user running the process
- **PORT**: network port number
- **ADDRESS**: IP address (`0.0.0.0` means listening on all interfaces)
- **STATE**: connection state (`LISTEN`)

```
Active Network Processes (5)
---------------------------------------------------------------------
PID      COMMAND            USER       PORT      ADDRESS       STATE
---------------------------------------------------------------------
12500    node               juraj      3000      ::1           LISTEN
580      ControlCenter      juraj      5000      0.0.0.0       LISTEN
420      Brave Browser      juraj      7000      0.0.0.0       LISTEN
250      Raycast            juraj      1234      127.0.0.1     LISTEN
585      rapportd           juraj      52001     0.0.0.0       LISTEN
```
