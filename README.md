# Spoliteware - a polite and considerate spyware

A polite spyware that asks permission to collect data beforehand.

## Usage

Start server side executable with needed configuration and put thanksdir in the same directory. Received data will be put in the created directory beside the executable.

Launch client executable for a one-time data collection and transfer with a pre-defined server address (recompile with the serverAddr changed to your server's address beforehand). If permission has been granted - the data will be POSTed to the server and a thanks message will be delivered and printed.

## Features

### v0.1.1

#### Client
- Display thanks message from the server

#### Server
- Thanks system and thanksfiles

### v0.1.0

#### Client
- System information

#### Server
- Collected data is sorted per machine in a corresponding folder (hostname_username_IP)
- TLS support

## TODO
- More green telemetry options
- Be more open on what the program does

## License 
MIT