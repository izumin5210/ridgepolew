package ridgepolew

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
)

type Ridgepole interface {
	Exec(context.Context, []string) error
}

type RidgepoleDocker struct {
	workDir string
	inR     io.Reader
	outW    io.Writer
	errW    io.Writer
}

func NewRidgepole(workDir string, inR io.Reader, outW, errW io.Writer) *RidgepoleDocker {
	return &RidgepoleDocker{
		workDir: workDir,
		inR:     inR,
		outW:    outW,
		errW:    errW,
	}
}

func NewDefaultRidgepole() *RidgepoleDocker {
	return NewRidgepole("", os.Stdin, os.Stdout, os.Stderr)
}

var _ Ridgepole = (*RidgepoleDocker)(nil)

func (r *RidgepoleDocker) Exec(ctx context.Context, args []string) error {
	args, err := r.buildRunArgs(args)
	if err != nil {
		return fmt.Errorf("failed to build args: %w", err)
	}

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdin = r.inR
	cmd.Stdout = r.outW
	cmd.Stderr = r.errW

	return cmd.Run()
}

func (r *RidgepoleDocker) buildRunArgs(args []string) ([]string, error) {
	ip, err := getLoopbackIP()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve loopback IP: %w", err)
	}

	wd := r.workDir
	if wd == "" {
		wd, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current working directory: %w", err)
		}
	}

	return append([]string{
		"run",
		"--rm",
		"--add-host", "localhost:" + ip.String(),
		"--volume", wd + ":/workdir",
		"izumin5210/ridgepole",
	}, args...), nil
}

func getLoopbackIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	localhost := net.IPv4(127, 0, 0, 1)

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		default:
			continue
		}
		if !ip.IsLoopback() && ip.To4() != nil {
			localhost = ip
		}
	}

	return localhost, nil
}
