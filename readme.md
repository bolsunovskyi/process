# Process Manager

Library that provides APIs to run external processes

Main library interface is located at `Manager` structure. To init this structure use `CreateManager` function.
Main library methods:

- `AddProcess` - run new process with arguments
- `TerminateProcess` - kill process
- `GetProcesses` - get all processes
- `Shutdown` - kill all processes and flush logs

Library config is represented in `ManagerConfig` structure:

type ManagerConfig struct {
	LogsFolder        string
	RenewOldProcesses bool
	ProcessesListFile string
}

- `LogsFolder` - folder where individual process logs will be located
- `RenewOldProcesses` - if set to true process will be renewed after main daemon restart
- `ProcessesListFile` - location of file that contains list of processes that needs to be renew

See more details and example of usage in `cmd` package.  