// package main
//
// import (
//
//	"fmt"
//	"log"
//	"os/exec"
//
// )
//
//	func main() {
//		cmd := exec.Command(
//			"poetry",
//			"run",
//			"python",
//			"/Users/khang.nguyen/Documents/hfsd-go/server/d3po/scripts/test.py",
//		)
//		if err := cmd.Run(); err != nil {
//			output, _ := cmd.CombinedOutput()
//			fmt.Println(string(output))
//			log.Fatal(err)
//		}
//
//		output, _ := cmd.CombinedOutput()
//		fmt.Println(string(output))
//	}
package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command(
		"poetry",
		"run", "python", "./d3po/scripts/test.py", "--train.json_path", "hehe",
	)

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	fmt.Println(string(stdout))
}
