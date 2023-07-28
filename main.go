package main

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"github.com/skoona/hubPower/commons"
	"github.com/skoona/hubPower/providers"
	"github.com/skoona/hubPower/services"
	"github.com/skoona/hubPower/ui"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var err error
	commons.ShutdownSignals = make(chan os.Signal, 1)

	ctx, cancelHub := context.WithCancel(context.Background())
	//defer cancelHub()

	gui := app.NewWithID("net.skoona.projects.ggapcmon")
	commons.DebugLog("main()::RootURI: ", gui.Storage().RootURI().Path())
	gui.SetIcon(commons.SknSelectThemedResource(commons.AppIcon))

	go func(stopFlag chan os.Signal, a fyne.App) {
		signal.Notify(stopFlag, syscall.SIGINT, syscall.SIGTERM)
		sig := <-stopFlag // wait on ctrl-c
		cancelHub()
		time.Sleep(5 * time.Second)
		err = fmt.Errorf("Shutdown Signal Received: %v", sig.String())
		a.Quit()
	}(commons.ShutdownSignals, gui)

	cfg, err := providers.NewConfig(gui.Preferences())
	if err != nil {
		dialog.ShowError(fmt.Errorf("main()::NewConfig(): %v", err), gui.NewWindow("ggapcmon Configuration Failed"))
		commons.ShutdownSignals <- syscall.SIGINT
		cfg.ResetConfig()
	}
	//cfg.ResetConfig()

	service, err := services.NewService(ctx, cfg)
	if err != nil {
		log.Panic("main()::Service startup() failed: ", err.Error())
	}
	defer service.Shutdown()

	vp := ui.NewViewProvider(ctx, cfg, service)
	defer vp.Shutdown()

	win := vp.ShowMainPage()
	win.SetOnClosed(func() {
		if ctx.Err() == nil {
			commons.DebugLog("main::OnClose() window exit and cancel triggered")
			cancelHub()
		}
	})
	gui.Run()
	commons.DebugLog("main::Shutdown Ended ")
}
