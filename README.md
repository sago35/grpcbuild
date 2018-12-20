# grpcbuild

Sample for gRPC build tool for C project.


# Usage

Windows 環境で作っているので、それ以外の環境は適宜読み替えてください。

片方の cmd.exe は以下のように立ち上げます。

    $ cd .\cmd\server
    $ set PATH=%GOPATH%\src\github.com\sago35\grpcbuild\cmd\builder;%PATH%
    $ go build
    $ server.exe

もう片方の cmd.exe は以下のように立ち上げます。

    $ cd .\cmd\builder
    $ go build
    $ go build -o dummyld.exe ..\dummycc
    $ go build -o dummycc.exe ..\dummycc

以下のようにしてそれぞれの実装を実行します。

    $ builder.exe -m 1
    $ builder.exe -m 2
    $ builder.exe -m 3 -threads 8
    $ builder.exe -m 4 -threads 8
    $ builder.exe -m 5 -threads 8
    $ builder.exe -m 6 -threads 8
    $ builder.exe -m 7 -threads 8

