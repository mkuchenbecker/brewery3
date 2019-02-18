go get -u github.com/golang/dep/cmd/dep
go get -u github.com/kyoh86/richgo
go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
go get github.com/mattn/goveralls
go get -d -u github.com/golang/protobuf/protoc-gen-go

PROTOC_ZIP=protoc-3.3.0-linux-x86_64.zip
curl -OL https://github.com/google/protobuf/releases/download/v3.3.0/$PROTOC_ZIP
sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
rm -f $PROTOC_ZIP

dep ensure

GIT_TAG="v1.2.0" # change as needed
go get -d -u github.com/golang/protobuf/protoc-gen-go
git -C $GOPATH/src/github.com/golang/protobuf checkout $GIT_TAG
go install github.com/golang/protobuf/protoc-gen-go

go get github.com/golang/mock/gomock
go install github.com/golang/mock/mockgen

