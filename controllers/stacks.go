package controllers
import (
	"io/ioutil"
	"os"
	"fmt"
	"os/exec"
	"log"
	"github.com/HariniGB/openstack-api/models"
	"encoding/json"
	"strings"
)

func DeployTemplate(name, tempStr string, params map[string]string) error {
	bytes := []byte(tempStr)
	id, _ := flake.NextID()
	file := fmt.Sprintf("/tmp/%d", id)
	err := ioutil.WriteFile(file, bytes, os.ModeAppend)

	paramArr := []string{"stack", "create",  name,  "-t", file}
	for k,v := range params {
		paramArr = append(paramArr,"--parameter", fmt.Sprintf("%s=%s", k, v))
	}

	cmd := exec.Command("openstack",  paramArr...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(string(b))
		return fmt.Errorf("unable to execute template file")
	}

	return nil
}

func UndeployTemplate(name string) error {
	paramArr := []string{"stack", "delete",  name}
	cmd := exec.Command("openstack",  paramArr...)
	err := cmd.Run()
	if err != nil {
		log.Print(err)
		return fmt.Errorf("unable to delete stack")
	}

	return nil
}

func GetStacks() []models.StackResponse {
	arr := []string {"stack", "list", "-f", "json", "-c", "Stack Name"}
	cmd := exec.Command("openstack",  arr...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(err)
		return nil
	}

	entries := []map[string]interface{}{}
	err = json.Unmarshal(b, &entries)
	if err != nil {
		log.Print(err)
		return nil
	}

	stacks := []models.StackResponse{}
	for _, entry := range entries {
		if stack, ok := entry["Stack Name"]; ok && strings.Index(stack.(string), "-network") == -1 {
			st := GetStack(stack.(string))
			if st != nil {
				stacks = append(stacks, *st)
			}
		}
	}
	return stacks
}

func GetStack(name string) *models.StackResponse {
	arr := []string {"stack", "show", name,  "-f", "json", "-c", "stack_status", "-c", "stack_name", "-c", "creation_time",
		"-c", "stack_status_reason", "-c", "outputs"}

	cmd := exec.Command("openstack",  arr...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(err)
		return nil
	}

	out := &models.Stack{}
	err = json.Unmarshal(b, out)
	if err != nil {
		log.Print(err)
		return nil
	}

	return out.ToStackResponse()
}
