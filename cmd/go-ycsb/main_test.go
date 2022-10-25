package main

import (
	"github.com/spf13/cobra"
	"reflect"
	"testing"
)

func Test_initClientCommand(t *testing.T) {
	type args struct {
		m *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initClientCommand(tt.args.m)
		})
	}
}

func Test_initialGlobal(t *testing.T) {
	type args struct {
		dbName       string
		onProperties func()
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialGlobal(tt.args.dbName, tt.args.onProperties)
		})
	}
}

func Test_newLoadCommand(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newLoadCommand(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newLoadCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newRunCommand(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newRunCommand(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRunCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newShellCommand(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newShellCommand(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newShellCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runClientCommandFunc(t *testing.T) {
	type args struct {
		cmd            *cobra.Command
		args           []string
		doTransactions bool
		command        string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"1", args{
			cmd:            nil,
			args:           []string{"httpdb"},
			doTransactions: true,
			command:        "run",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runClientCommandFunc(tt.args.cmd, tt.args.args, tt.args.doTransactions, tt.args.command)
		})
	}
}

func Test_runLoadCommandFunc(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runLoadCommandFunc(tt.args.cmd, tt.args.args)
		})
	}
}

func Test_runShellCommand(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runShellCommand(tt.args.args)
		})
	}
}

func Test_runShellCommandFunc(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runShellCommandFunc(tt.args.cmd, tt.args.args)
		})
	}
}

func Test_runShellDeleteCommand(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runShellDeleteCommand(tt.args.cmd, tt.args.args)
		})
	}
}

func Test_runShellInsertCommand(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runShellInsertCommand(tt.args.cmd, tt.args.args)
		})
	}
}

func Test_runShellReadCommand(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runShellReadCommand(tt.args.cmd, tt.args.args)
		})
	}
}

func Test_runShellScanCommand(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runShellScanCommand(tt.args.cmd, tt.args.args)
		})
	}
}

func Test_runShellTableCommand(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runShellTableCommand(tt.args.cmd, tt.args.args)
		})
	}
}

func Test_runShellUpdateCommand(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runShellUpdateCommand(tt.args.cmd, tt.args.args)
		})
	}
}

func Test_runTransCommandFunc(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTransCommandFunc(tt.args.cmd, tt.args.args)
		})
	}
}

func Test_shellLoop(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shellLoop()
		})
	}
}
