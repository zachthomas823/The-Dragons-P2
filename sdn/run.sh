#! /bin/bash -xe

cd sdn_Command
./sdn_Comman_Exec &
cd ..

cd sdn_Proxy
./sdn_Proxy &
cd ..

cd sdn_Reasource
./sdn_Reasource &
cd ..

cd ..

cd dashboard
./dashboard &