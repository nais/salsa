package commands

import (
    "bytes"
    "fmt"
    "log"
    "os/exec"
)

type CmdConfig struct {
    workDir string
    cmd string
    args []string
}

func GradleScan(workDir string){
    c := CmdConfig{
        workDir: workDir,
        cmd:     "./gradlew",
        args:    []string{"-q", "dependencies", "--configuration", "runtimeClasspath"},
    }
    c.ExecuteCommand()
}

func (c CmdConfig) ExecuteCommand() {
    cmd := exec.Command(c.cmd, c.args...)

    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    cmd.Dir = c.workDir
    err := cmd.Run()

    if err != nil {
        log.Printf("cmd.Run: %s failed: %v\n", cmd, err)
    }
    outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
    if len(errStr) > 1 {
        fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
    }
    fmt.Printf(outStr)
}
