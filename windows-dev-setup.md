Setting up a development environment on Windows
===============================================

In order to create an oci-connection to an Oracle database, we have to setup a development on windows, 
installing some extra tooling so we can compile go-code which uses the mattn/oci8 code:  





In summary
----------

Set following environment variables

    set GOPATH=Z:\dev
    set GOROOT=D:\Go\
    set Path=%PATH%;D:\Go\bin;D:\msys64\mingw64\bin
    set PKG_CONFIG_PATH=D:\oracle\instantclient_12_1


