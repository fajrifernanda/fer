package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"gitlab.com/fajrifernanda/kumpa-kit-go/util"
)

func Generate(servicename string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter proto path (Ex:\"pb/content/content.proto\"): ")
	text2, _ := reader.ReadString('\n')
	protoPath := strings.ReplaceAll(text2, "\n", "")
	protoPath = strings.ReplaceAll(protoPath, "proto", "pb.go")

	fmt.Println(">>>" + servicename + "<<<")
	fmt.Println("RPC Path : ", protoPath)

	protoCoreFile := protoPath
	serviceUrl := "gitlab.kumparan.com/yowez/" + servicename

	fmt.Println("Creating ", servicename)
	os.RemoveAll(servicename)
	os.Mkdir(servicename, os.ModePerm)
	os.Mkdir(servicename+"/service", os.ModePerm)
	os.Mkdir(servicename+"/pb", os.ModePerm)
	os.Mkdir(servicename+"/pb/content", os.ModePerm)

	util.CopyFolder("templates", servicename)
	util.CopyFolder("pb/content/", servicename+"/pb/content/")
	util.CopyFileContents("templates/service/service.go", servicename+"/service/service.go")

	os.Mkdir(servicename+"/client", os.ModePerm)
	GenerateProto2Go()
	CreateScaffoldScript()
	RunScaffold(servicename)
	fmt.Println(servicename, "Created")
	GenerateRPCClient(protoCoreFile, servicename, serviceUrl)
	ReadServiceServer(servicename, serviceUrl, protoCoreFile)

}

func CreateScaffoldScript() {
	contents := `#!/usr/bin/env bash

servicename=$1;
find $servicename -type f -exec sed -i '' "s/skeleton-service/$servicename/g" {} \;
cp $servicename/*.example $servicename/config.yml;
cd $servicename;
go mod tidy;
go get;
echo "finish scaffolding";

echo "running test";
make mockgen;
make test;
echo "finish test";

git init;
git remote add gitlab.kumparan.com/yowez/$servicename;
echo "git initialized";
echo "Oke";
echo "finish";
`

	bt := []byte(contents)
	ioutil.WriteFile("scaffold.sh", bt, 0644)

}

func GenerateProto2Go() {
	contents := `#!/usr/bin/env bash
	protoc --go_out=plugins=grpc:. pb/*.proto
	ls pb/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'
`

	bt := []byte(contents)
	ioutil.WriteFile("proto2go.sh", bt, 0644)
	exec.Command("bash", "proto2go.sh")
	os.Remove("proto2go.sh")
}

func RunScaffold(serviceName string) {
	cmd := exec.Command("bash", "scaffold.sh", serviceName)
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Output \n", out.String())
	os.Remove("scaffold.sh")
}
