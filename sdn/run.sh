#! /bin/bash -xe

cd sdn_Command
chmod 777 sdn_Command
./sdn_Command &
sleep 2
cd ..

cd sdn_Proxy
chmod 777 sdn_Proxy
./sdn_Proxy &
sleep 2
cd ..

cd sdn_Reasource
chmod 777 sdn_Reasource
./sdn_Reasource &
sleep 2
cd ..

cd ..

cd dashboard
chmod 777 dashboard
./dashboard &
sleep 2