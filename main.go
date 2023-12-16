package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

// SSH connection structure
type ConnectionDetails struct {
    Name string `json:"name"`
    User string `json:"user"`
    IP   string `json:"ip"`
    Key  string `json:"key"`
}

// Returns path of dotfile
func getDotFilePath() string {
    homeDir, _ := os.UserHomeDir()
    return filepath.Join(homeDir, ".conman")
}

// Loads the saved SSH connections from the dotfile
func LoadConnections() ([]ConnectionDetails, error) {
    dotFile := getDotFilePath()

    // Check if file exists
    if _, err := os.Stat(dotFile); os.IsNotExist(err) {
        return []ConnectionDetails{}, nil
    }
	
    // Read file contents
    file, err := ioutil.ReadFile(dotFile)
    if err != nil {
        return nil, err
    }

    // JSON to ConnectionDetails slice
    var connections []ConnectionDetails
    if err := json.Unmarshal(file, &connections); err != nil {
        return nil, err
    }

    return connections, nil
}

// Saves SSH connections to a file
func SaveConnections(connections []ConnectionDetails) error {
    data, err := json.Marshal(connections)
    if err != nil {
        return err
    }

    dotFile := getDotFilePath()

    return ioutil.WriteFile(dotFile, data, 0600) // Read-only perms
}

// Adds a new SSH connection
func AddConnection(name, user, ip, key string) error {
    connections, err := LoadConnections()
    if err != nil {
        return err
    }

    connections = append(connections, ConnectionDetails{
        Name: name,
        User: user,
        IP:   ip,
        Key:  key,
    })

    return SaveConnections(connections)
}

// Establish an SSH connection
func Connect(name string) error {
    connections, err := LoadConnections()
    if err != nil {
        return err
    }

    for _, conn := range connections {
        if conn.Name == name {
            cmd := exec.Command("ssh", "-i", conn.Key, fmt.Sprintf("%s@%s", conn.User, conn.IP))
            cmd.Stdin = os.Stdin
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr
            return cmd.Run()
        }
    }

    return fmt.Errorf("connection not found")
}

// Delete an SSH connection by its name
func DeleteConnection(name string) error {
    connections, err := LoadConnections()
    if err != nil {
        return err
    }

    // Filter out the connection to be deleted
    updatedConnections := []ConnectionDetails{}
    found := false
    for _, conn := range connections {
        if conn.Name != name {
            updatedConnections = append(updatedConnections, conn)
        } else {
            found = true
        }
    }

    if !found {
        return fmt.Errorf("connection %s not found", name)
    }

    // Save the updated list of connections
    return SaveConnections(updatedConnections)
}

// SplitAddress splits the address into user and IP
func SplitAddress(address string) (user, ip string, err error) {
    parts := strings.Split(address, "@")
    if len(parts) != 2 {
        return "", "", fmt.Errorf("invalid address format")
    }
    return parts[0], parts[1], nil
}

// Flags
func printHelp() {
    fmt.Println(`Usage of conman:
    -k string
        Path to the SSH key
    -a string
        SSH address in the format user@ip
    -n string
        Custom name for the connection
    -c string
        Connect using a saved connection
    -ls
        List saved connections
    -v
        Verbose output for listing
    -d string
        Delete a saved connection
    -h
        Show this help message`)
}

func main() {
    var (
        keyPath   = flag.String("k", "", "Path to the SSH key")
        address   = flag.String("a", "", "SSH address in the format user@ip")
        name      = flag.String("n", "", "Custom name for the connection")
        connect   = flag.String("c", "", "Connect using a saved connection")
        list      = flag.Bool("ls", false, "List saved connections")
        verbose   = flag.Bool("v", false, "Verbose output for listing")
        deleteCon = flag.String("d", "", "Delete a saved connection")
	help      = flag.Bool("h", false, "Show help message")
    )

    flag.Parse()

    // Adding a new connection
    if *keyPath != "" && *address != "" && *name != "" {
        user, ip, err := SplitAddress(*address)
        if err != nil {
            fmt.Println("Error parsing address:", err)
            return
        }

        if err := AddConnection(*name, user, ip, *keyPath); err != nil {
            fmt.Println("Error adding connection:", err)
            return
        }
        fmt.Println(*name, "has been created and added to list")
    }

    // Connecting to a saved connection
    if *connect != "" {
        if err := Connect(*connect); err != nil {
            fmt.Println("Error connecting:", err)
            return
        }
    }

    // Listing connections
    if *list {
        connections, err := LoadConnections()
        if err != nil {
            fmt.Println("Error loading connections:", err)
            return
        }

        fmt.Println("Saved entries:")
        for _, conn := range connections {
            if *verbose {
                fmt.Printf("%s %s@%s %s\n", conn.Name, conn.User, conn.IP, conn.Key)
            } else {
                fmt.Println(conn.Name)
            }
        }
    }

    // Deleting a connection
    if *deleteCon != "" {
        if err := DeleteConnection(*deleteCon); err != nil {
            fmt.Println("Error deleting connection:", err)
            return
        }

        fmt.Println(*deleteCon, "has been deleted")
    }

    // Print flags
    if *help {
        printHelp()
        return
    }
}

