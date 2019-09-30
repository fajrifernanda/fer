package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/kumparan/fer/util"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Generate(servicename string) {
	fmt.Println(`
   ________________     __              __
   / ____/ ____/ __ \   / /_____  ____  / /
  / /_  / __/ / /_/ /  / __/ __ \/ __ \/ / 
 / __/ / /___/ _, _/  / /_/ /_/ / /b_/ / /  
/_/   /_____/_/ |_|   \__/\____/\____/_/  

`)
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
	GetTemplates(servicename)
	os.Mkdir(servicename+"/service", os.ModePerm)
	os.Mkdir(servicename+"/pb", os.ModePerm)
	os.Mkdir(servicename+"/pb/content", os.ModePerm)

	util.CopyFolder("pb/content/", servicename+"/pb/content/")
	util.CopyFileContents("templates/service/service.go", servicename+"/service/service.go")

	os.Mkdir(servicename+"/client", os.ModePerm)
	GenerateProto2Go()
	CreateScaffoldScript()
	fmt.Println(servicename, "Scaffolding ...")
	RunScaffold(servicename)
	fmt.Println(servicename, "Generating client ...")
	GenerateRPCClient(protoCoreFile, servicename, serviceUrl)
	time.Sleep(1500 * time.Millisecond)
	fmt.Println(servicename, "Generating service&test ...")
	ReadServiceServer(servicename, serviceUrl, protoCoreFile)
	time.Sleep(1500 * time.Millisecond)
	fmt.Println(servicename, "Created")

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
git remote add origin "git@gitlab.kumparan.com:yowez/$servicename.git";
echo "git initialized";
echo "Oke";
echo "finish";
`

	bt := []byte(contents)
	ioutil.WriteFile("scaffold.sh", bt, 0644)

}

func GetTemplates(serviceName string) {
	contents := `#!/usr/bin/env bash
	servicename=$1;
	cd $servicename;
	git init;
	git remote add origin git@gitlab.kumparan.com:yowez/skeleton-service.git;
	git remote -v;
	git pull origin master;
	rm -rf .git;
	cd ..;

`

	bt := []byte(contents)
	err := ioutil.WriteFile("gettemplate.sh", bt, 0644)
	if err != nil {
		fmt.Println(err)
	}
	cmd := exec.Command("bash", "gettemplate.sh", serviceName)
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	_ = os.Remove("gettemplate.sh")
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
	os.Remove("scaffold.sh")
}
