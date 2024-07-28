package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {

	case "run":
		run()

	case "child":
		child()

	default:
		panic("Bad command")
	}
}

func run() {
	fmt.Printf("Running %v\n", os.Args[2:])
	//use with test1 & test2 cmd := exec.Command(os.Args[2], os.Args[3:]...)

	// generate a stringlist with intial element child
	args := []string{"child"}
	// For the argument include and after the first argument 'run', e.g. /bin/bash, we append them into the list
	args = append(args, os.Args[2:]...)
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	/**
	test1： Before Run()
	In this case, we will set our rootname with the name "container", before the /bin/bash child process being run.
	Since the processes are isolated in our OS. Even if we change the root hostname back, our child process will still not be visible.
	*/
	// syscall.Sethostname([]byte("container"))
	cmd.Run()
	/**
	test2: After Run()
	Once again, the processes are isolated. The inital hostname of our child process will still be our inital root hostname. No matter how we change our root hostname.
	After the /bin/bash child process is run, our hostname will be changed to "container". But the child process is not visible of this change.
	*/
	//syscall.Sethostname([]byte("container"))
}

/**
	This is the function that will be called by run(). The run will call this function if the input command is "child".
	But notice that this calling "child" procedure is done by run(), not us.
	That is, we will call go run main.go run <our command>, and the run will call itself, by using the linux symbolic link that points to the executable of the current process,
with the parameter child + <our command>
	Then we will be able to start the child process, and also set the name(hostname) of this child process.


	Question: Why we cannot done this in run()?
	Answer: Because you can put the `syscall.Sethostname([]byte("container"))` nowhere in the run() method. Either putting the code snipet before or after the code `cmd.Run()` is going to set
our root hostname to container but not our container's name
*/

// It will set the hostname of the child process. Then start the child process with command. So we can clearly distinguish this process with the root process.
func child() {
	fmt.Printf("Currently Running Child with Command: %v\n", os.Args[2:])
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	syscall.Sethostname([]byte("container"))
	cmd.Run()
}

// anicaa-rogzephyrusg14ga401iuga401iu
