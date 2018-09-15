#!/bin/bash

# download localforwarder: TODO

cd /lf/
ls -la ./Microsoft.LocalForwarder.ConsoleHost 

chmod +x ./Microsoft.LocalForwarder.ConsoleHost
./Microsoft.LocalForwarder.ConsoleHost noninteractive
