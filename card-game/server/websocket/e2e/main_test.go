package e2e

import (
	"fmt"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"testing"
)

var (
	DefaultWssEndpoint = ""
)

//
// Setup
//

func chDirToProjectRoot() {
	p, _ := os.Getwd()
	fmt.Println(p)

	os.Chdir("../")

	p, _ = os.Getwd()
	fmt.Println(p)
}

func slsDeploy() {
	cmd := exec.Command("sls", "deploy", "--stage", "test")
	out, _ := cmd.CombinedOutput()

	// output log
	fmt.Println(string(out))

	// check exit code
	if cmd.ProcessState.ExitCode() != 0 {
		panic(fmt.Sprintf("cmd failed : `sls %v`", cmd.Args))
	}
}

func slsInfo() string {
	cmd := exec.Command("sls", "info", "--stage", "test")
	out, _ := cmd.CombinedOutput()

	// output log
	fmt.Println(string(out))

	// check exit code
	if cmd.ProcessState.ExitCode() != 0 {
		panic(fmt.Sprintf("cmd failed : `sls %v`", cmd.Args))
	}

	re := regexp.MustCompile(`wss://.*`)
	endpoint := re.Find(out)

	if endpoint == nil {
		panic(fmt.Sprintf("endpoint not found in\n%s", string(out)))
	}

	fmt.Println(string(endpoint))

	return string(endpoint)
}

func TestMain(m *testing.M) {
	// chdir to project root
	chDirToProjectRoot()

	// sls deploy
	slsDeploy()

	// sls info
	DefaultWssEndpoint = slsInfo()

	// run Tests
	status := m.Run()
	os.Exit(status)
}

//
// Helper
//

func newWssConnection() (*websocket.Conn, string) {
	playerId := uuid.NewV4().String()
	return newWssConnectionWithArgs(DefaultWssEndpoint, playerId), playerId
}

func newWssConnectionWithArgs(url string, playerId string) *websocket.Conn {
	h := http.Header{}
	h.Add("X-Pobo380-Network-Games-Player-Id", playerId)

	c, resp, err := websocket.DefaultDialer.Dial(url, h)
	if err != nil {
		b, _ := ioutil.ReadAll(resp.Body)
		panic(fmt.Sprintf("Dial failed : %s\n%+v\n%s", err, resp, string(b)))
	}

	return c
}
