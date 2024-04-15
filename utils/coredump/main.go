/*
 * Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
 * or more contributor license agreements. Licensed under the Apache License 2.0.
 * See the file "LICENSE" for details.
 */

// coredump provides a tool for extracting stack traces from coredumps.
// It also includes a test suite to unit test pf-host-agent components against
// a set of coredumps to validate stack extraction code.

package main

import (
	"context"
	"errors"
	"flag"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/elastic/otel-profiling-agent/utils/coredump/modulestore"
	"github.com/peterbourgon/ff/v3/ffcli"
	log "github.com/sirupsen/logrus"
)

// moduleStoreS3Bucket defines the S3 bucket used for the module store.
const moduleStoreS3Bucket = "optimyze-proc-mem-testdata"

func main() {
	log.SetReportCaller(false)
	log.SetFormatter(&log.TextFormatter{})

	store := initModuleStore()

	root := ffcli.Command{
		Name:       "coredump",
		ShortUsage: "coredump <subcommand> [flags]",
		ShortHelp:  "Tool for creating and managing coredump test cases",
		Subcommands: []*ffcli.Command{
			newAnalyzeCmd(store),
			newCleanCmd(store),
			newExportModuleCmd(store),
			newNewCmd(store),
			newRebaseCmd(store),
			newUploadCmd(store),
			newGdbCmd(store),
		},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			log.Fatalf("%v", err)
		}
	}
}

func initModuleStore() *modulestore.Store {
	cfg := aws.NewConfig().WithRegion("eu-central-1")
	sess := session.Must(session.NewSession(cfg))
	s3Client := s3.New(sess)
	return modulestore.New(s3Client, moduleStoreS3Bucket, "modulecache")
}
