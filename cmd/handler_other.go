package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func generateGRPCHandler(entity string, fileNamingConvention string, sm ...*SafetyManager) {
	// Create gRPC directory
	grpcDir := filepath.Join(DirInternal, DirHandler, DirGRPC)
	_ = os.MkdirAll(grpcDir, 0755)

	generateProtoFile(grpcDir, entity, fileNamingConvention, sm...)
	generateGRPCServerFile(grpcDir, entity, fileNamingConvention, sm...)
}

func generateProtoFile(dir, entity string, fileNamingConvention string, sm ...*SafetyManager) {
	entityLower := strings.ToLower(entity)

	// Apply naming convention to filename
	var filename string
	if fileNamingConvention == "snake_case" {
		filename = filepath.Join(dir, toSnakeCase(entity)+".proto")
	} else if fileNamingConvention == "kebab-case" {
		filename = filepath.Join(dir, toKebabCase(entity)+".proto")
	} else {
		filename = filepath.Join(dir, entityLower+".proto")
	}

	var content strings.Builder
	content.WriteString("syntax = \"proto3\";\n\n")
	content.WriteString(fmt.Sprintf("package %s;\n\n", entityLower))
	content.WriteString(fmt.Sprintf("option go_package = \"./%s\";\n\n", entityLower))

	content.WriteString(fmt.Sprintf("service %sService {\n", entity))
	content.WriteString(fmt.Sprintf("  rpc Create%s(Create%sRequest) returns (Create%sResponse);\n", entity, entity, entity))
	content.WriteString(fmt.Sprintf("  rpc Get%s(Get%sRequest) returns (%sResponse);\n", entity, entity, entity))
	content.WriteString(fmt.Sprintf("  rpc Update%s(Update%sRequest) returns (Update%sResponse);\n", entity, entity, entity))
	content.WriteString(fmt.Sprintf("  rpc Delete%s(Delete%sRequest) returns (Delete%sResponse);\n", entity, entity, entity))
	content.WriteString(fmt.Sprintf("  rpc List%ss(List%ssRequest) returns (List%ssResponse);\n", entity, entity, entity))
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message %s {\n", entity))
	content.WriteString("  int32 id = 1;\n")
	content.WriteString("  string name = 2;\n")
	content.WriteString("  string email = 3;\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message Create%sRequest {\n", entity))
	content.WriteString("  string name = 1;\n")
	content.WriteString("  string email = 2;\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message Create%sResponse {\n", entity))
	content.WriteString(fmt.Sprintf("  %s %s = 1;\n", entity, entityLower))
	content.WriteString("  string message = 2;\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message Get%sRequest {\n", entity))
	content.WriteString("  int32 id = 1;\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message %sResponse {\n", entity))
	content.WriteString(fmt.Sprintf("  %s %s = 1;\n", entity, entityLower))
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message Update%sRequest {\n", entity))
	content.WriteString("  int32 id = 1;\n")
	content.WriteString("  string name = 2;\n")
	content.WriteString("  string email = 3;\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message Update%sResponse {\n", entity))
	content.WriteString("  string message = 1;\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message Delete%sRequest {\n", entity))
	content.WriteString("  int32 id = 1;\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message Delete%sResponse {\n", entity))
	content.WriteString("  string message = 1;\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message List%ssRequest {\n", entity))
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("message List%ssResponse {\n", entity))
	content.WriteString(fmt.Sprintf("  repeated %s %ss = 1;\n", entity, entityLower))
	content.WriteString("  int32 total = 2;\n")
	content.WriteString("}\n")

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		ui.Error(fmt.Sprintf("Error writing proto file: %v", err))
		return
	}
}

func generateGRPCServerFile(dir, entity string, fileNamingConvention string, sm ...*SafetyManager) {
	// Get the module name from go.mod
	moduleName := getModuleName()
	importPath := getImportPath(moduleName)

	entityLower := strings.ToLower(entity)

	// Apply naming convention to filename
	var filename string
	if fileNamingConvention == "snake_case" {
		filename = filepath.Join(dir, toSnakeCase(entity)+"_server.go")
	} else if fileNamingConvention == "kebab-case" {
		filename = filepath.Join(dir, toKebabCase(entity)+"-server.go")
	} else {
		filename = filepath.Join(dir, entityLower+"_server.go")
	}

	var content strings.Builder
	content.WriteString("package grpc\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"context\"\n\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/usecase\"\n", importPath))
	content.WriteString(fmt.Sprintf("\tpb \"%s/internal/handler/grpc/%s\"\n", importPath, entityLower))
	content.WriteString(")\n\n")

	content.WriteString(fmt.Sprintf("type %sServer struct {\n", entity))
	content.WriteString(fmt.Sprintf("\tpb.Unimplemented%sServiceServer\n", entity))
	content.WriteString(fmt.Sprintf("\tusecase usecase.%sUseCase\n", entity))
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("func New%sServer(uc usecase.%sUseCase) *%sServer {\n", entity, entity, entity))
	content.WriteString(fmt.Sprintf("\treturn &%sServer{usecase: uc}\n", entity))
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("func (s *%sServer) Create%s(ctx context.Context, req *pb.Create%sRequest) (*pb.Create%sResponse, error) {\n", entity, entity, entity, entity))
	content.WriteString(fmt.Sprintf("\tinput := usecase.Create%sInput{\n", entity))
	content.WriteString("\t\tName:  req.Name,\n")
	content.WriteString("\t\tEmail: req.Email,\n")
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\toutput, err := s.usecase.Create%s(input)\n", entity))
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn nil, err\n")
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\treturn &pb.Create%sResponse{\n", entity))
	content.WriteString(fmt.Sprintf("\t\t%s: &pb.%s{\n", entity, entity))
	content.WriteString(fmt.Sprintf("\t\t\tId:    int32(output.%s.ID),\n", entity))
	content.WriteString(fmt.Sprintf("\t\t\tName:  output.%s.Name,\n", entity))
	content.WriteString(fmt.Sprintf("\t\t\tEmail: output.%s.Email,\n", entity))
	content.WriteString("\t\t},\n")
	content.WriteString("\t\tMessage: output.Message,\n")
	content.WriteString("\t}, nil\n")
	content.WriteString("}\n\n")

	content.WriteString(fmt.Sprintf("func (s *%sServer) Get%s(ctx context.Context, req *pb.Get%sRequest) (*pb.%sResponse, error) {\n", entity, entity, entity, entity))
	content.WriteString(fmt.Sprintf("\t%s, err := s.usecase.Get%s(int(req.Id))\n", entityLower, entity))
	content.WriteString("\tif err != nil {\n")
	content.WriteString("\t\treturn nil, err\n")
	content.WriteString("\t}\n\n")

	content.WriteString(fmt.Sprintf("\treturn &pb.%sResponse{\n", entity))
	content.WriteString(fmt.Sprintf("\t\t%s: &pb.%s{\n", entity, entity))
	content.WriteString(fmt.Sprintf("\t\t\tId:    int32(%s.ID),\n", entityLower))
	content.WriteString(fmt.Sprintf("\t\t\tName:  %s.Name,\n", entityLower))
	content.WriteString(fmt.Sprintf("\t\t\tEmail: %s.Email,\n", entityLower))
	content.WriteString("\t\t},\n")
	content.WriteString("\t}, nil\n")
	content.WriteString("}\n")

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		ui.Error(fmt.Sprintf("Error writing grpc server file: %v", err))
		return
	}
}

func generateCLIHandler(entity string, fileNamingConvention string, sm ...*SafetyManager) {
	// Create CLI directory
	cliDir := filepath.Join(DirInternal, DirHandler, DirCLI)
	_ = os.MkdirAll(cliDir, 0755)

	// Get the module name from go.mod
	moduleName := getModuleName()
	importPath := getImportPath(moduleName)

	entityLower := strings.ToLower(entity)

	// Apply naming convention to filename
	var filename string
	if fileNamingConvention == "snake_case" {
		filename = filepath.Join(cliDir, toSnakeCase(entity)+"_commands.go")
	} else if fileNamingConvention == "kebab-case" {
		filename = filepath.Join(cliDir, toKebabCase(entity)+"-commands.go")
	} else {
		filename = filepath.Join(cliDir, entityLower+"_commands.go")
	}

	var content strings.Builder
	content.WriteString("package cli\n\n")
	content.WriteString("import (\n")
	content.WriteString("\t\"fmt\"\n")
	content.WriteString("\t\"strconv\"\n\n")
	content.WriteString("\t\"github.com/spf13/cobra\"\n")
	content.WriteString(fmt.Sprintf("\t\"%s/internal/usecase\"\n", importPath))
	content.WriteString(")\n\n")

	// CLI struct
	content.WriteString(fmt.Sprintf("type %sCLI struct {\n", entity))
	content.WriteString(fmt.Sprintf("\tusecase usecase.%sUseCase\n", entity))
	content.WriteString("}\n\n")

	// Constructor
	content.WriteString(fmt.Sprintf("func New%sCLI(uc usecase.%sUseCase) *%sCLI {\n", entity, entity, entity))
	content.WriteString(fmt.Sprintf("\treturn &%sCLI{usecase: uc}\n", entity))
	content.WriteString("}\n\n")

	// Create command
	content.WriteString(fmt.Sprintf("func (c *%sCLI) Create%sCommand() *cobra.Command {\n", entity, entity))
	content.WriteString("\tcmd := &cobra.Command{\n")
	content.WriteString("\t\tUse:   \"crear\",\n")
	content.WriteString(fmt.Sprintf("\t\tShort: \"Crear un nuevo %s\",\n", entityLower))
	content.WriteString("\t\tRun: func(cmd *cobra.Command, args []string) {\n")
	content.WriteString("\t\t\tname, _ := cmd.Flags().GetString(\"name\")\n")
	content.WriteString("\t\t\temail, _ := cmd.Flags().GetString(\"email\")\n\n")

	content.WriteString(fmt.Sprintf("\t\t\tinput := usecase.Create%sInput{\n", entity))
	content.WriteString("\t\t\t\tName:  name,\n")
	content.WriteString("\t\t\t\tEmail: email,\n")
	content.WriteString("\t\t\t}\n\n")

	content.WriteString(fmt.Sprintf("\t\t\toutput, err := c.usecase.Create%s(input)\n", entity))
	content.WriteString("\t\t\tif err != nil {\n")
	content.WriteString("\t\t\t\tfmt.Printf(\"Error: %v\\n\", err)\n")
	content.WriteString("\t\t\t\treturn\n")
	content.WriteString("\t\t\t}\n\n")

	content.WriteString(fmt.Sprintf("\t\t\tfmt.Printf(\"%s created: %%+v\\n\", output.%s)\n", entity, entity))
	content.WriteString("\t\t},\n")
	content.WriteString("\t}\n\n")

	content.WriteString("\tcmd.Flags().StringP(\"name\", \"n\", \"\", \"Name of the ")
	content.WriteString(entityLower)
	content.WriteString("\")\n")
	content.WriteString("\tcmd.Flags().StringP(\"email\", \"e\", \"\", \"Email of the ")
	content.WriteString(entityLower)
	content.WriteString("\")\n")
	content.WriteString("\tcmd.MarkFlagRequired(\"name\")\n")
	content.WriteString("\tcmd.MarkFlagRequired(\"email\")\n\n")
	content.WriteString("\treturn cmd\n")
	content.WriteString("}\n\n")

	// Get command
	content.WriteString(fmt.Sprintf("func (c *%sCLI) Get%sCommand() *cobra.Command {\n", entity, entity))
	content.WriteString("\treturn &cobra.Command{\n")
	content.WriteString("\t\tUse:   \"obtener [id]\",\n")
	content.WriteString(fmt.Sprintf("\t\tShort: \"Obtener %s por ID\",\n", entityLower))
	content.WriteString("\t\tArgs:  cobra.ExactArgs(1),\n")
	content.WriteString("\t\tRun: func(cmd *cobra.Command, args []string) {\n")
	content.WriteString("\t\t\tid, err := strconv.Atoi(args[0])\n")
	content.WriteString("\t\t\tif err != nil {\n")
	content.WriteString("\t\t\t\tfmt.Printf(\"Invalid ID: %v\\n\", err)\n")
	content.WriteString("\t\t\t\treturn\n")
	content.WriteString("\t\t\t}\n\n")

	content.WriteString(fmt.Sprintf("\t\t\t%s, err := c.usecase.Get%s(id)\n", entityLower, entity))
	content.WriteString("\t\t\tif err != nil {\n")
	content.WriteString("\t\t\t\tfmt.Printf(\"Error: %v\\n\", err)\n")
	content.WriteString("\t\t\t\treturn\n")
	content.WriteString("\t\t\t}\n\n")

	content.WriteString(fmt.Sprintf("\t\t\tfmt.Printf(\"%s: %%+v\\n\", %s)\n", entity, entityLower))
	content.WriteString("\t\t},\n")
	content.WriteString("\t}\n")
	content.WriteString("}\n")

	if err := writeGoFile(filename, content.String(), sm...); err != nil {
		ui.Error(fmt.Sprintf("Error writing cli handler file: %v", err))
		return
	}
}

func generateWorkerHandler(entity string, fileNamingConvention string, sm ...*SafetyManager) {
	// Create worker directory
	workerDir := filepath.Join(DirInternal, DirHandler, DirWorker)
	_ = os.MkdirAll(workerDir, 0755)

	// Get the module name from go.mod
	moduleName := getModuleName()

	entityLower := strings.ToLower(entity)

	// Apply naming convention to filename
	var filename string
	if fileNamingConvention == "snake_case" {
		filename = filepath.Join(workerDir, toSnakeCase(entity)+"_worker.go")
	} else if fileNamingConvention == "kebab-case" {
		filename = filepath.Join(workerDir, toKebabCase(entity)+"-worker.go")
	} else {
		filename = filepath.Join(workerDir, entityLower+"_worker.go")
	}

	content := fmt.Sprintf(`package worker

import (
	"encoding/json"
	"log"
	
	"%s/internal/usecase"
)

type %sWorker struct {
	usecase usecase.%sUseCase
}

func New%sWorker(uc usecase.%sUseCase) *%sWorker {
	return &%sWorker{usecase: uc}
}

func (w *%sWorker) Process%sJob(jobData []byte) error {
	var input usecase.Create%sInput
	
	if err := json.Unmarshal(jobData, &input); err != nil {
		log.Printf("Failed to unmarshal job data: %%v", err)
		return err
	}
	
	output, err := w.usecase.Create%s(input)
	if err != nil {
		log.Printf("Failed to process %s job: %%v", err)
		return err
	}
	
	log.Printf("%s job completed: %%v", output.Message)
	return nil
}

func (w *%sWorker) ProcessBatch%sJob(jobData []byte) error {
	var inputs []usecase.Create%sInput
	
	if err := json.Unmarshal(jobData, &inputs); err != nil {
		log.Printf("Failed to unmarshal batch job data: %%v", err)
		return err
	}
	
	for _, input := range inputs {
		if _, err := w.usecase.Create%s(input); err != nil {
			log.Printf("Failed to process %s in batch: %%v", err)
			// Continue processing other items
		}
	}
	
	log.Printf("Batch %s job completed")
	return nil
}
`, moduleName, entity, entity, entity, entity, entity, entity, entity, entity, entity, entity, entityLower, entity, entity, entity, entity, entity, entityLower, entity)

	if err := writeFile(filename, content, sm...); err != nil {
		ui.Error(fmt.Sprintf("Error writing worker file: %v", err))
		return
	}
}

func generateSOAPHandler(entity string, fileNamingConvention string, sm ...*SafetyManager) {
	// Create SOAP directory
	soapDir := filepath.Join(DirInternal, DirHandler, DirSOAP)
	_ = os.MkdirAll(soapDir, 0755)

	// Get the module name from go.mod
	moduleName := getModuleName()

	entityLower := strings.ToLower(entity)

	// Apply naming convention to filename
	var filename string
	if fileNamingConvention == "snake_case" {
		filename = filepath.Join(soapDir, toSnakeCase(entity)+"_client.go")
	} else if fileNamingConvention == "kebab-case" {
		filename = filepath.Join(soapDir, toKebabCase(entity)+"-client.go")
	} else {
		filename = filepath.Join(soapDir, entityLower+"_client.go")
	}

	content := fmt.Sprintf(`package soap

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	
	"%s/internal/usecase"
)

type %sSOAPClient struct {
	usecase usecase.%sUseCase
	baseURL string
}

func New%sSOAPClient(uc usecase.%sUseCase, baseURL string) *%sSOAPClient {
	return &%sSOAPClient{
		usecase: uc,
		baseURL: baseURL,
	}
}

type SOAPEnvelope struct {
	XMLName xml.Name    `+"`xml:\"soap:Envelope\"`"+`
	Header  interface{} `+"`xml:\"soap:Header\"`"+`
	Body    SOAPBody    `+"`xml:\"soap:Body\"`"+`
}

type SOAPBody struct {
	XMLName xml.Name    `+"`xml:\"soap:Body\"`"+`
	Content interface{} `+"`xml:\",innerxml\"`"+`
}

type Create%sRequest struct {
	XMLName xml.Name `+"`xml:\"Create%sRequest\"`"+`
	Name    string   `+"`xml:\"name\"`"+`
	Email   string   `+"`xml:\"email\"`"+`
}

type Create%sResponse struct {
	XMLName xml.Name `+"`xml:\"Create%sResponse\"`"+`
	ID      int      `+"`xml:\"id\"`"+`
	Message string   `+"`xml:\"message\"`"+`
}

func (c *%sSOAPClient) Create%s(name, email string) (*Create%sResponse, error) {
	request := Create%sRequest{
		Name:  name,
		Email: email,
	}
	
	envelope := SOAPEnvelope{
		Body: SOAPBody{Content: request},
	}
	
	xmlData, err := xml.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SOAP request: %%v", err)
	}
	
	resp, err := http.Post(
		c.baseURL+"/Create%s",
		"text/xml; charset=utf-8",
		strings.NewReader(string(xmlData)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send SOAP request: %%v", err)
	}
	defer resp.Body.Close()
	
	var responseEnvelope SOAPEnvelope
	if err := xml.NewDecoder(resp.Body).Decode(&responseEnvelope); err != nil {
		return nil, fmt.Errorf("failed to decode SOAP response: %%v", err)
	}
	
	var createResponse Create%sResponse
	if err := xml.Unmarshal([]byte(responseEnvelope.Body.Content.(string)), &createResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %%v", err)
	}
	
	return &createResponse, nil
}
`, moduleName, entity, entity, entity, entity, entity, entity, entity, entity, entity, entity, entity, entity, entity, entity, entity, entity)

	if err := writeFile(filename, content, sm...); err != nil {
		ui.Error(fmt.Sprintf("Error writing soap file: %v", err))
		return
	}
}

