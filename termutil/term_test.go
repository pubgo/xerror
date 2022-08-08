package termutil

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"testing"

	"github.com/mattn/go-isatty"
)

func TestName(t *testing.T) {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		// don't log to terminal
		fmt.Println("don't log to terminal")
	}
}

func TestHomeDir(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	//defaultConf := filepath.Join(usr.HomeDir, ".config", "photon", "config")
	//if _, err := os.Stat(defaultConf); os.IsNotExist(err) {
	//	log.Fatal(err)
	//}
	fmt.Println(usr.HomeDir)
	fmt.Println(usr.Name)
	fmt.Println(usr.Username)
}
