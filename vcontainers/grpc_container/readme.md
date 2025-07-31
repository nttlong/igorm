#$env:PATH += ";D:\protoc-31.1-win64\bin"
2- Run
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

%USERPROFILE%\go\bin\protoc-gen-go.exe
%USERPROFILE%\go\bin\protoc-gen-go-grpc.exe

#gen code cd D:\mygrpc\proto
protoc --go_out=. --go-grpc_out=. invoker.proto
#protoc   --go_out=.   --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative  invoker.proto
#protoc --go_out=. --go-grpc_out=.  invoker/invoker.proto