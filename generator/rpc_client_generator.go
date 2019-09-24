package generator

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/dave/jennifer/jen"
)

// Service :nodoc:
type Service struct {
	Name       string
	Parameters []string
	Returns    []string
}

var (
	serviceUrl        string
	serviceClientName string
	serviceUrlProto   string
)

// GenerateRPCClient ... :nodoc:
func GenerateRPCClient(path string, serviceName string, serviceRepo string) {
	f := NewFile("client")
	serviceUrl = serviceRepo
	serviceObj, _ := ParseProtoToArray(serviceName, path)
	serviceUrlProto = serviceUrl + "/pb/" + getSimpleNameFromProtoPath(path)
	f.ImportAlias(serviceUrlProto, "pb")
	for _, v := range serviceObj {
		var returns string
		var bodyReturn string

		jParam := "("
		jParam += strings.Join(v.Parameters, ",")
		jParam += ")"

		for _, aReturn := range v.Returns {
			returns = returns + "," + aReturn
		}
		for _, aReturn := range v.Returns {
			if strings.Contains(aReturn, "*") {
				aReturn = strings.Replace(aReturn, "*", "&", -1)
				aReturn += "{}"
			}
			if strings.Contains(aReturn, "error") {
				aReturn = strings.Replace(aReturn, "error", "nil", -1)
			}
			bodyReturn = bodyReturn + "," + aReturn

		}
		returns = returns[1:]
		bodyReturn = bodyReturn[1:]
		returns = "( " + returns + " )"
		f.Func().Params(
			Id("c").Op("*").Id("client"),
		).Id(v.Name).Params(getFunctionArgsClient(jParam)...).Id(returns).Block(
			AddConn(),
			AddErrChecker(),
			CloseConn(),
			AddNewClient(serviceClientName),
			Return(AddClientReturn(v.Name)),
		)
	}

	buf := &bytes.Buffer{}
	_ = f.Render(buf)
	splitPath := strings.Split(path, "/")
	savePath := splitPath[2]
	savePath = strings.Replace(savePath, ".pb.go", ".go", -1)
	ioutil.WriteFile(serviceName+"/"+"client"+"/"+savePath, buf.Bytes(), 0644)

}

func getFunctionArgsClient(in string) (args []Code) {
	strlong := SplitBetweenTwoChar(in, "(", ")")
	strs := strings.Split(strlong, ", ")
	strs = strs[1 : len(strs)-1]
	args = append(args, Code(Id("ctx").Qual("context", "Context")))
	for _, v := range strs {
		if strings.Contains(v, ".") {
			argItem := strings.Split(v, " ")
			argName := argItem[0]
			argPath := strings.Split(argItem[1], ".")[0]
			if strings.Contains(argPath, "pb") {
				argPath = serviceUrlProto
			}
			argType := strings.Split(argItem[1], ".")[1]
			args = append(args, Code(Id(argName).Op("*").Qual(argPath, argType)))
		}
	}
	args = append(args, Code(Id("opts").Op("...").Qual("google.golang.org/grpc", "CallOption")))

	return
}

// ParseProtoToArray ... :nodoc:
func ParseProtoToArray(serviceName string, path string) ([]Service, error) {
	interfaceName := "Client" + " interface"
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)

	isServiceClient := false
	services := []Service{}
	protoFunctions := []string{}
	for scanner.Scan() {
		var text string
		if strings.Contains(scanner.Text(), interfaceName) {
			isServiceClient = true
			serviceClientName = GetServiceClient(scanner.Text())
		}
		if isServiceClient == true {
			if strings.Contains(scanner.Text(), "//") {
				continue
			}
			if strings.Contains(scanner.Text(), "*") {
				text = scanner.Text()
				text = strings.Replace(text, "*", "*pb.", -1)
				protoFunctions = append(protoFunctions, text)
			}
		}
		if isServiceClient == true && strings.Contains(scanner.Text(), "}") {
			isServiceClient = false
		}
	}
	for _, v := range protoFunctions {
		params := splitFunctionParameters(v)
		services = append(services, params)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return services, nil
}

func splitFunctionParameters(function string) Service {
	splittedFunction := strings.Split(function, "(")
	for k := range splittedFunction {
		splittedFunction[k] = strings.Replace(splittedFunction[k], ")", "", -1)
	}
	params := strings.Split(splittedFunction[1], ",")
	returns := strings.Split(splittedFunction[2], ",")

	newService := Service{Name: splittedFunction[0], Parameters: params, Returns: returns}
	return newService
}

func AddConn() (s *Statement) {
	s = List(Id("conn"), Id("err")).Op(":=").Id("c").Dot("Conn").Dot("Get").Parens(Id("ctx"))
	return
}
func AddErrChecker() (s *Statement) {
	s = If(
		Err().Op("!=").Nil(),
	).Block(
		Return(Nil(), Err()),
	)
	return
}

func CloseConn() (s *Statement) {
	s = Defer().Id("conn").Dot("Close()")
	return
}

func AddNewClient(client string) (s *Statement) {
	s = Id("cli").Op(":=").Id("pb").Dot("New" + client).Parens(Id("conn").Dot("ClientConn"))
	return
}

func AddClientReturn(funcName string) (s *Statement) {
	s = Id("cli").Dot(funcName).Parens(List(Id("ctx"), Id("in"), Id("opts...")))
	return
}

func GetServiceClient(text string) (client string) {
	client = text
	client = strings.ReplaceAll(client, "type", "")
	client = strings.ReplaceAll(client, "interface", "")
	client = strings.ReplaceAll(client, "{", "")
	client = strings.ReplaceAll(client, "\n", "")
	client = strings.ReplaceAll(client, " ", "")
	return
}
