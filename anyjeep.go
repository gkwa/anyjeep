package anyjeep

import (
	"fmt"
	"log/slog"
	"os/exec"
	"sync"
)

func Main() int {
	slog.Debug("anyjeep", "test", true)
	test()

	return 0
}

func test() {
	commands := []string{
		"echo Command 1",
		"echo Command 2",
		"echo Command 3",
		"echo Command 4",
		"echo Command 5",
		"echo Command 6",
		"echo Command 7",
		"echo Command 8",
		"echo Command 9",
		"echo Command 10",
	}

	numCommands := len(commands)

	concurrencyLimit := 1
	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, concurrencyLimit)
	results := make(chan string, numCommands)

	for _, cmdStr := range commands {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore
		go func(cmd string) {
			defer func() {
				<-semaphore // Release semaphore
				wg.Done()
			}()

			output, err := runCommand(cmd)
			if err != nil {
				results <- fmt.Sprintf("Error executing command: %s", err)
				return
			}
			results <- output
		}(cmdStr)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}

func runCommand(cmdStr string) (string, error) {
	cmd := exec.Command("bash", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
