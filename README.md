# Soviet Mainframe БЭСМ-6
An educational example program of Go programming 
based on having some fun creating an old school decoy program. 
This is an old hacker trick, to give people you password. 
No, you are not supposed to use this to actually trick people 
and you most likely can't because most login systems don't work like this anymore.

This is more for fun and historical perspective.

## How Does it Work?
When people use text based login systems like this:

    Login: erik
    Password: qwerty
    Большая Электронно-Счётная Машина 6
    Welcome to БЭСМ-6 comrade erik: March 29, 2021
    > 

The name of the computer is from the Soviet Mainframe BESM-6 from 1965

What a nefarious hacker could easily do 
is to make a program that looks exactly like this, 
but which instead of logging the user in, stores the user password. 
Thus the hacker can come back later and read a file containing
the passwords of users who have previously tried to log in.

## What is the Point?
The point of this program is to teach basic things like:

- How to work with input and output
- Launching external processes
- Socket based communication
- Asychronous programming with Go routines

## Building, Running and Using Programs
This module is setup as multiple packages which produce a number of command line tools such as client, server and shell.

You can build these tools in the following way:

   ❯ go build github.com/ordovician/decoy/cmd/client
   ❯ go build github.com/ordovician/decoy/cmd/server
   ❯ go build github.com/ordovician/decoy/cmd/shell
    
Why exactly this long build line? It stems from the go.mod file defining the root of the package like this:

    ❯ cat go.mod
    module github.com/ordovician/mainframe

    go 1.16

However you can use local paths by prefixing with `./` like this:

    ❯ go build ./cmd/shell
    
The build command will use this defined module name when building. You can also directly run commands this way:

    ❯ go run github.com/ordovician/decoy/cmd/shell
    ❯ go run ./cmd/shell
    
However after building you will get local executables which you can run like this:

    ❯ ./shell
    Login: thomas
    Password: qwerty
    Большая Электронно-Счётная Машина 6: June 6, 2022
    Welcome to БЭСМ-6 comrade thomas


## Running as Server and Connecting as Client
The login program can be run as a server which you cannot to using NetCat or Telnet

    ❯ go build ./cmd/server
    ❯ ./server

It will listen for connections on port 1234 (TODO Make port number user configurable).
You can connect with telnet or netcat like this:

    ❯ telnet localhost 1234
    Trying ::1...
    Connected to localhost.
    Escape character is '^]'.
    Login: thomas
    Password: qwerty
    Большая Электронно-Счётная Машина 6: June 14, 2022
    Welcome to БЭСМ-6 comrade thomas

Alternatively with netcat:

    ❯ nc localhost 1234

    
