# FakeSSH

A dockerized honeypot SSH server written in Go to log login attempts.
Password authentications always fail so no terminal access is given to the attacker.

[![](http://dockeri.co/image/fffaraz/fakessh)](https://hub.docker.com/r/fffaraz/fakessh)

## Quick Start

```
go install github.com/fffaraz/fakessh@latest
sudo setcap 'cap_net_bind_service=+ep' ~/go/bin/fakessh
fakessh [optional-log-directory]
```

OR

```
docker run -it --rm -p 22:22 fffaraz/fakessh
```

OR

```
docker run -d --restart=always -p 22:22 --name fakessh fffaraz/fakessh
docker logs -f fakessh
```

## Fork Added Features:

1、Report OpenSSH version as SSH-2.0-OpenSSH_8.9p1 Ubuntu-3ubuntu0.10 (Ubuntu 22.04 LTS).

2、Add support for private key login and log the login IP and username.

3、Add MySQL record for username-password-password SHA256 functionality. The database environment variables must be set, otherwise the program will not start.

Environment variables：
```

DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=your_db_host
DB_PORT=your_db_port
DB_NAME=your_db_name

```

Database structure：
```

CREATE TABLE `ssh` (
  `user` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
  `password` text COLLATE utf8mb4_general_ci NOT NULL,
  `sha256` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`user`,`sha256`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

```

Fork quick start
```
docker run -d \
           -e DB_USER=your_db_user \
           -e DB_PASSWORD=your_db_password \
           -e DB_HOST=your_db_host \
           -e DB_PORT=your_db_port \
           -e DB_NAME=your_db_name \
           --restart=always -p 22:22
           --name fakessh senlin1hao/fakessh

docker logs -f fakessh
```

### See also

* [jaksi/sshesame](https://github.com/jaksi/sshesame) - A fake SSH server that lets everyone in and logs their activity.
* [shazow/ssh-chat](https://github.com/shazow/ssh-chat) - Custom SSH server written in Go. Instead of a shell, you get a chat prompt.
* [gliderlabs/ssh](https://github.com/gliderlabs/ssh) - Easy SSH servers in Golang.
* [gliderlabs/sshfront](https://github.com/gliderlabs/sshfront) - Programmable SSH frontend.
* [desaster/kippo](https://github.com/desaster/kippo) - Kippo - SSH Honeypot.
* [micheloosterhof/cowrie](https://github.com/micheloosterhof/cowrie) - Cowrie SSH/Telnet Honeypot.
* [fzerorubigd/go0r](https://github.com/fzerorubigd/go0r) - A simple ssh honeypot in golang.
* [droberson/ssh-honeypot](https://github.com/droberson/ssh-honeypot) - Fake sshd that logs ip addresses, usernames, and passwords.
* [x0rz/ssh-honeypot](https://github.com/x0rz/ssh-honeypot) - Fake sshd that logs ip addresses, usernames, and passwords.
* [tnich/honssh](https://github.com/tnich/honssh) - HonSSH is designed to log all SSH communications between a client and server.
* [Learn from your attackers - SSH HoneyPot](https://www.robertputt.co.uk/learn-from-your-attackers-ssh-honeypot.html)
* [cowrie](https://github.com/cowrie/cowrie) - Cowrie SSH/Telnet Honeypot.
