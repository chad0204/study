#编译go、java文件

##1. 安装proto解压, mac brew install protobuf

##2. 编译成java、python直接
./protoc.exe --java_out=./ ./demo.proto
./protoc.exe --python_out=./ ./demo.proto

##3. 编译成go需要在protoc的bin目录下添加protoc-gen-go.exe
   进入protoc-gen-go编译 go build main.go 重命名main.exe为protoc-gen-go.exe
   将protoc-gen-go.exe移到protoc的bin目录下
      
   mac的话把构建的protoc-gen-go二进制文件放到gopath下



   demo.proto文件需要添加package和option
   执行 ./protoc.exe --go_out=./ ./demo.proto
   ./protoc.exe --go_out=. demo.proto


##4. 生成grpc代码
protoc --go_out=plugins=grpc:. hello.proto
还可以自己重写protoc-gen-go插件, 自定义代码生成插件


go get -v -u github.com/golang/protobuf/proto(没有protobuf需要下载protobuf包)
