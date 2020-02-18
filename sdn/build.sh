#!/bin/bash

cd sdn_Command
go build .
cd ..

cd sdn_Controller
go build .
cd ..

cd sdn_Launcher
go build .
cd ..

cd sdn_Reasource
go build .
cd ..

cd sdn_Proxy
go build .
cd ..