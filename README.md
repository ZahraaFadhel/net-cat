# 🚀 Net-Cat Project

## 📝 Overview

Welcome to the TCP-Chat Project! This project is a recreation of NetCat in a Server-Client Architecture using Go. The goal is to create a group chat application that mimics the functionality of NetCat, enabling communication over TCP connections.

NetCat (`nc`), a command-line utility, reads and writes data across network connections using TCP or UDP. It is widely used for anything involving TCP, UDP, or UNIX-domain sockets. It can open TCP connections, send UDP packages, listen on arbitrary TCP and UDP ports, and more.

This project aims to replicate the functionality of NetCat with the following features:

## 🌟 Features

- **📡 TCP Connection**: Establish TCP connections between the server and multiple clients (1 to many relationship).
- **👤 Client Name Requirement**: Clients must provide a name upon connection.
- **🔢 Connection Control**: Limit the number of concurrent connections to a maximum of 10.
- **💬 Messaging**: Clients can send messages to the chat. Messages will not be broadcasted if they are empty.
- **⏰ Timestamped Messages**: Messages are timestamped and include the sender's username.
- **🕰 Message History**: New clients receive the complete message history upon joining.
- **🔔 Join/Leave Notifications**: All clients are notified when a new client joins or leaves the chat.
- **🔄 Continuous Operation**: The server remains operational even if a client disconnects.
- **⚠️ Error Handling**: Robust error handling on both the server and client sides.

## Usage/Examples

```bash
git clone https://learn.reboot01.com/git/ok/net-cat.git

> cd /net-cat
```
```go
go run tcp-server.go
#Listening on port 8989 by default
Listening on the port :8989

# to listen on custom port simply do:
go run tcp-server.go <port>

```



## Authors

- [@ok](https://learn.reboot01.com/git/ok)
- [@zfadhel](https://learn.reboot01.com/git/zfadhel)
- [@mohani](https://learn.reboot01.com/git/mohani)

