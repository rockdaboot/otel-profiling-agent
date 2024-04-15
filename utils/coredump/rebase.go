/*
 * Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
 * or more contributor license agreements. Licensed under the Apache License 2.0.
 * See the file "LICENSE" for details.
 */

package main

import (
	"context"
	"flag"
	"fmt"
	"os/exec"

	"github.com/elastic/otel-profiling-agent/utils/coredump/modulestore"
	"github.com/peterbourgon/ff/v3/ffcli"
)

type rebaseCmd struct {
	store *modulestore.Store

	allowDirty bool
}

func newRebaseCmd(store *modulestore.Store) *ffcli.Command {
	args := &rebaseCmd{store: store}

	set := flag.NewFlagSet("rebase", flag.ExitOnError)
	set.BoolVar(&args.allowDirty, "allow-dirty", false, "Allow uncommitted changes in git")

	return &ffcli.Command{
		Name:       "rebase",
		Exec:       args.exec,
		ShortUsage: "rebase",
		ShortHelp:  "Update all test cases by running them and saving the current unwinding",
		FlagSet:    set,
	}
}

func (cmd *rebaseCmd) exec(context.Context, []string) (err error) {
	cases, err := findTestCases(true)
	if err != nil {
		return fmt.Errorf("failed to find test cases: %w", err)
	}

	if !cmd.allowDirty {
		if err = exec.Command("git", "diff", "--quiet").Run(); err != nil {
			return fmt.Errorf("refusing to work on a dirty source tree. " +
				"please commit your changes first or pass `-allow-dirty` to ignore")
		}
	}

	for _, testCasePath := range cases {
		var testCase *CoredumpTestCase
		testCase, err = readTestCase(testCasePath)
		if err != nil {
			return fmt.Errorf("failed to read test case: %w", err)
		}

		core, err := OpenStoreCoredump(cmd.store, testCase.CoredumpRef, testCase.Modules)
		if err != nil {
			return fmt.Errorf("failed to open coredump: %w", err)
		}

		testCase.Threads, err = ExtractTraces(context.Background(), core, false, nil)
		core.Close()
		if err != nil {
			return fmt.Errorf("failed to extract traces: %w", err)
		}

		if err = writeTestCase(testCasePath, testCase, true); err != nil {
			return fmt.Errorf("failed to write test case: %w", err)
		}
	}

	return nil
}
