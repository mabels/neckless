package main

import (
	"testing"
)

func TestCrazyBee(t *testing.T) {
	args := CrazyBeeArgs{}
	buildArgs([]string{}, &args)
}

func TestCrazyBeeCreatePipeLine(t *testing.T) {
	args := CrazyBeeArgs{}
	buildArgs([]string{"pipeline"}, &args)
	if len(args.pipeline.name) == 48 {
		t.Errorf("pipeline.name not ok: %s", args.pipeline.name)
	}
}

func TestCrazyBeeCreatePipeLineSet(t *testing.T) {
	args := CrazyBeeArgs{}
	buildArgs([]string{
		"pipeline",
		"--name", "Test.Name",
	}, &args)
	if args.pipeline.name != "Test.Name" {
		t.Errorf("pipeline.name not ok: %s", args.pipeline.name)
	}
}
