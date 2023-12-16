
# ConMan - SSH Connection Manager

ConMan is a simple yet effective CLI tool for managing SSH connections. It allows you to save SSH connection details under a user-defined name, quickly establish these connections, and manage your saved entries with ease.

## Features

- **Save SSH Connections**: Store your SSH connection details with a custom name for easy access.
- **Connect Quickly**: Connect to your saved SSH servers using just the custom name.
- **List and Manage**: View all saved connections and delete them as needed.

## Installation

Clone the repository and build the project:

```bash
git clone https://github.com/Bryce-Moore/Conman
cd conman
go build
```

## Usage

### Adding a Connection

```bash
./conman -k [path/to/key.pem] -a [user@ip] -n [custom-name]
```

- `-k`: Path to your SSH key.
- `-a`: SSH address (format: user@ip).
- `-n`: Custom name for the connection.

### Connecting to a Server

```bash
./conman -c [custom-name]
```

- `-c`: Connect using a saved connection name.

### Listing Connections

```bash
./conman -ls
```

- `-ls`: List all saved connections.
- `-v`: Verbose output (optional).

### Deleting a Connection

```bash
./conman -d [custom-name]
```

- `-d`: Delete a saved connection.

### Help

```bash
./conman -h
```

- `-h`: Show help message.

## Configuration

The tool saves connection details in a hidden JSON file located at `~/.conman`
