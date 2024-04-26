# FrontDoor

FrontDoor is a reverse proxy server designed to redirect traffic from HTTP to HTTPS while serving as a proxy to a
backend service. It ensures secure communication by upgrading all incoming HTTP requests to HTTPS and proxies these
requests to a specified backend server. This approach enhances security and ensures that sensitive data is transmitted
over encrypted connections.

## Features

- HTTP to HTTPS Redirection: Automatically redirects all HTTP traffic to HTTPS.
- Reverse Proxy: Proxies HTTPS requests to an internal backend, allowing for centralized handling of SSL/TLS
  termination.
- Graceful Shutdown: Properly handles SIGINT and SIGTERM signals for graceful shutdown of the server.
- Logging: Utilizes zerolog for structured and console-friendly logging.
- HTTP/2 Support: Ready for HTTP/2 communications with backend services.

## Requirements

- Go 1.15 or higher.
- Systemd (for service management on Linux).
- SSL/TLS certificates.

## Installation

### 1. Clone the repository:

```bash
git clone https://github.com/yourusername/frontdoor.git
cd frontdoor
```

### 2.Build the binary:

```bash
go build -ldflags "-s -w" -o frontdoor
sudo cp -f frontdoor /usr/local/bin/frontdoor
```

### 3. Set permissions to allow binding to privileged ports (for non-root users):

If at all possible please make sure to run this as a non root user. There is really no reason for you to run it as a
root user as you can allow the binary to bind to privileged ports, without the need of full root capabilities, like
below:

```bash
sudo setcap 'cap_net_bind_service=+eip' /usr/local/bin/frontdoor
```

## Usage

Start FrontDoor with necessary flags specifying the paths to your SSL/TLS certificates and the address of your backend
server:

```bash
/usr/local/bin/frontdoor -key path/to/server.key -cert path/to/server.crt -backend 127.0.0.1:8080
```

### Flags

- -key: Path to the SSL/TLS private key file.
- -cert: Path to the SSL/TLS certificate file.
- -backend: Address of the backend server.

### Systemd Integration

To manage FrontDoor as a systemd service, create a service file named frontdoor.service in /etc/systemd/system/ with the
following content:

```systemd
[Unit]
Description=FrontDoor Service
After=network.target

[Service]
Type=simple
User=frontdoor
ExecStart=/usr/local/bin/frontdoor -key /path/to/server.key -cert /path/to/server.crt -backend 127.0.0.1:8080
Restart=on-failure

[Install]
WantedBy=multi-user.target
Enable and start the service:

```

```bash
sudo systemctl daemon-reload
sudo systemctl enable frontdoor.service
sudo systemctl start frontdoor.service
```

Contributing
Contributions are welcome! Please feel free to submit pull requests or create issues for bugs, questions, and feature
requests.

Thank you for using FrontDoor!

## License
- MIT License
- See LICENSE.md
