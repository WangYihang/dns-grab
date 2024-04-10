package main

import (
	"os"
	"time"

	"github.com/WangYihang/dns-grab/pkg/loaders"
	"github.com/WangYihang/dns-grab/pkg/model"
	"github.com/WangYihang/dns-grab/pkg/option"
	"github.com/WangYihang/gojob"
	"github.com/WangYihang/gojob/pkg/runner"
	"github.com/WangYihang/gojob/pkg/utils"
	"github.com/jessevdk/go-flags"
)

var Opt option.Option

func init() {
	Opt.Version = model.PrintVersion
	if _, err := flags.Parse(&Opt); err != nil {
		os.Exit(1)
	}
}

func main() {
	scheduler := gojob.NewScheduler().
		SetNumWorkers(Opt.NumWorkers).
		SetMaxRetries(Opt.MaxTries).
		SetMaxRuntimePerTaskSeconds(Opt.MaxRuntimePerTaskSeconds).
		SetNumShards(int64(Opt.NumShards)).
		SetShard(int64(Opt.Shard)).
		SetTotalTasks(utils.Count(loaders.Get(Opt.InputFilePath, "txt"))).
		SetOutputFilePath(Opt.OutputFilePath).
		SetMetadata("build", map[string]string{
			"version": model.Version,
			"commit":  model.Commit,
			"date":    model.Date,
		}).
		SetMetadata("runner", runner.Runner).
		SetMetadata("arguments", Opt).
		SetMetadata("started_at", time.Now().Format(time.RFC3339)).
		Start()
	for line := range utils.Cat(Opt.InputFilePath) {
		scheduler.Submit(
			model.NewTask(
				model.WithQNAME(line),
				model.WithQTYPE(Opt.QType),
				model.WithResolver(Opt.Resolver),
			),
		)
	}
	scheduler.Wait()
}
