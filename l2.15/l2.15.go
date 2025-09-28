package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	// Перехват Ctrl+C — чтобы не убивать шелл
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for range sigChan {
			fmt.Println("\n[!] Команда прервана (Ctrl+C)")
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		// Показываем текущую директорию как приглашение
		dir, _ := os.Getwd()
		fmt.Printf("%s$ ", filepath.Base(dir))

		line, err := reader.ReadString('\n')
		if err != nil { // Ctrl+D — выход
			fmt.Println("\n[+] Завершение работы шелла")
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Поддержка конвейеров
		commands := strings.Split(line, "|")
		if len(commands) > 1 {
			if err := runPipeline(commands); err != nil {
				fmt.Println("Ошибка:", err)
			}
			continue
		}

		// Обычная команда
		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "exit":
			fmt.Println("[+] Выход")
			return
		case "cd":
			runCd(args)
		case "pwd":
			runPwd()
		case "echo":
			runEcho(args)
		case "kill":
			runKill(args)
		case "ps":
			runPs()
		default:
			runExternal(args)
		}
	}
}

func runCd(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: cd <path>")
		return
	}
	if err := os.Chdir(args[1]); err != nil {
		fmt.Println("cd error:", err)
	}
}

func runPwd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("pwd error:", err)
		return
	}
	fmt.Println(dir)
}

func runEcho(args []string) {
	if len(args) > 1 {
		fmt.Println(strings.Join(args[1:], " "))
	}
}

func runKill(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: kill <pid>")
		return
	}
	pid, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Invalid PID:", args[1])
		return
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("Process not found:", err)
		return
	}
	if err := proc.Kill(); err != nil {
		fmt.Println("Failed to kill process:", err)
	} else {
		fmt.Printf("Process %d terminated\n", pid)
	}
}

func runPs() {
	cmd := exec.Command("ps", "aux")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("ps error:", err)
	}
}

func runExternal(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Command error:", err)
	}
}

func runPipeline(commands []string) error {
	var cmds []*exec.Cmd
	for _, cmdStr := range commands {
		args := strings.Fields(strings.TrimSpace(cmdStr))
		if len(args) == 0 {
			continue
		}
		cmds = append(cmds, exec.Command(args[0], args[1:]...))
	}

	for i := 0; i < len(cmds)-1; i++ {
		out, err := cmds[i].StdoutPipe()
		if err != nil {
			return err
		}
		cmds[i+1].Stdin = out
	}

	cmds[len(cmds)-1].Stdout = os.Stdout

	for _, cmd := range cmds {
		cmd.Stderr = os.Stderr
	}

	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return err
		}
	}

	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return err
		}
	}
	return nil
}
