package applier

import (
	"io"
	"os"

	"github.com/helmfile/helmfile/pkg/app"
	hconfig "github.com/helmfile/helmfile/pkg/config"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(writer io.Writer, logLevel string) *zap.SugaredLogger {
	var cfg zapcore.EncoderConfig
	cfg.MessageKey = "message"
	out := zapcore.AddSync(writer)
	var level zapcore.Level
	err := level.Set(logLevel)
	if err != nil {
		panic(err)
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		out,
		level,
	)
	return zap.New(core).Sugar()
}

func ApplyHelmfile() error {
	//environment := "default"
	source := "/home/vishal/cdrive/work/git/ukama/infra-as-code/releases/init-helmfile.yaml"
	globalConfig := new(hconfig.GlobalOptions)

	valueset := []string{"lookupImageTag", "1ab632f"}
	set := map[string]interface{}{}
	set[valueset[0]] = valueset[1]
	applyOptions := &hconfig.ApplyOptions{}
	//diffOptions := &hconfig.DiffOptions{}

	//globalConfig.Environment = environment
	globalConfig.StateValuesSet = valueset
	globalConfig.File = source
	globalConfig.Debug = true
	globalConfig.EnableLiveOutput = true
	globalConfig.LogLevel = "debug"
	globalConfig.SetLogger(NewLogger(os.Stdout, "debug"))

	globalImpl := hconfig.NewGlobalImpl(globalConfig)
	globalImpl.SetSet(set)

	applyImpl := hconfig.NewApplyImpl(globalImpl, applyOptions)
	//diffImpl := hconfig.NewDiffImpl(globalImpl, diffOptions)

	if err := applyImpl.ValidateConfig(); err != nil {
		log.Errorf("ApplyImpl failed %s", err.Error())
		return err
	}

	//log.Infof("Applying Globalconfig:\n %+v, \n applyConfig:\n %+v diffconfig:\n %+v", globalConfig, applyOptions, diffOptions)
	log.Infof("Applying NewApplyGlobalImpl:\n %+v, \n ApplyOption:\n %+v", applyImpl.GlobalImpl, applyImpl.ApplyOptions)

	//a := app.New(applyImpl)
	helmfile := app.New(applyImpl)
	log.Infof("Applying helmfile:\n %+v", helmfile)
	if err2 := helmfile.Apply(applyImpl); err2 != nil {
		log.Errorf("Error: %v", err2)
		return err2
	}

	// if err := helmfile.Diff(diffImpl); err != nil {
	// 	switch e := err.(type) {
	// 	case *app.Error:
	// 		if e.Code() == 2 {
	// 			log.Info("Changes detected. Applying...")

	// 			if err2 := helmfile.Apply(applyImpl); err2 != nil {
	// 				return err2
	// 			}

	// 			log.Infof("Changes applied.")
	// 		} else {
	// 			log.Errorf("error on diff %s", err.Error())
	// 			return err
	// 		}
	// 	default:
	// 		log.Errorf("Error: %v", err)
	// 		if strings.HasSuffix(err.Error(), "no state file found") {
	// 			return err
	// 		}
	// 		return nil
	// 	}
	// } else {
	// 	log.Infof("No changes detected.")
	// }

	// err := a.Apply(applyImpl)
	// if err != nil {
	// 	switch e := err.(type) {
	// 	case *app.NoMatchingHelmfileError:
	// 		noMatchingExitCode := 3
	// 		if globalImpl.AllowNoMatchingRelease {
	// 			noMatchingExitCode = 0
	// 		}
	// 		return errors.NewExitError(e.Error(), noMatchingExitCode)
	// 	case *app.MultiError:
	// 		return errors.NewExitError(e.Error(), 1)
	// 	case *app.Error:
	// 		return errors.NewExitError(e.Error(), e.Code())
	// 	default:
	// 		panic(fmt.Errorf("BUG: please file an github issue for this unhandled error: %T: %v", e, e))
	// 	}
	// }

	return nil

	// r, err := applier.New(
	// 	nil,
	// 	applier.Source(source),
	// 	applier.Once(true),
	// 	applier.Environment(environment),
	// 	applier.Values(m),
	// )
	// if err != nil {
	// 	return err
	// }

}

// type Runner struct {
// 	logger *zap.SugaredLogger

// 	config    *HelmFileConfig
// 	diffConf  HelmFileDiffConfig
// 	applyConf HelmFileApplyConfig

// 	assetsDir string
// 	interval  time.Duration
// 	once      bool
// 	synced    bool
// }

// func New(box *packr.Box, opts ...Option) (*Runner, error) {
// 	//l := apputil.NewLogger(os.Stderr, "debug")

// 	r := &Runner{
// 		interval: 10 * time.Second,
// 		diffConf: HelmFileDiffConfig{
// 			detailedExitcode: true,
// 		},
// 		applyConf: HelmFileApplyConfig{
// 			//logger: l,
// 		},
// 		config: &HelmFileConfig{
// 			//logger: l,
// 			env: "",
// 		},
// 	}

// 	for i := range opts {
// 		if err := opts[i](r); err != nil {
// 			return nil, err
// 		}
// 	}

// 	if r.config.env == "" {
// 		r.config.env = "default"
// 	}

// 	if r.config.fileOrDir == "" {
// 		r.config.fileOrDir = fmt.Sprintf("%s/helmfile.yaml", r.assetsDir)
// 	}

// 	//r.logger = l

// 	return r, nil
// }

// var DefaultAssetsDir = "assets"

// func (r *Runner) RunOnce() error {

// 	logger := r.logger

// 	helmfile := app.New(r.config)
// 	if err := helmfile.Diff(r.diffConf); err != nil {
// 		switch e := err.(type) {
// 		case *app.Error:
// 			if e.Code() == 2 {
// 				logger.Info("Changes detected. Applying...")

// 				if err2 := helmfile.Apply(r.applyConf); err2 != nil {
// 					return err2
// 				}

// 				logger.Infof("Changes applied.")
// 			} else {
// 				return err
// 			}
// 		default:
// 			r.logger.Errorf("Error: %v", err)
// 			if strings.HasSuffix(err.Error(), "no state file found") {
// 				return err
// 			}
// 			return nil
// 		}
// 	} else {
// 		logger.Infof("No changes detected.")
// 	}

// 	return nil
// }

// func (r *Runner) Run() error {
// 	stopSig := signals.SetupSignalHandler()

// 	if r.once {
// 		return r.RunOnce()
// 	}

// 	stop := make(chan struct{}, 0)
// 	errs := make(chan error, 0)

// 	go func() {
// 		for {
// 			if err := r.RunOnce(); err != nil {
// 				errs <- err
// 				return
// 			}

// 			r.logger.Infof("Waiting for %-8v", r.interval)
// 			nextTime := time.Now()
// 			nextTime = nextTime.Add(r.interval)
// 			time.Sleep(time.Until(nextTime))

// 			select {
// 			case <-stop:
// 				r.logger.Info("Gracefully stopped the run loop")
// 				return
// 			default:
// 			}
// 		}
// 	}()

// 	select {
// 	case <-stopSig:
// 		// TODO Immediately cancel the RunOnce call on SIGTERM
// 		r.logger.Info("Stopping the run loop")
// 		stop <- struct{}{}
// 		return nil
// 	case err := <-errs:
// 		return err
// 	}

// 	return nil
// }
