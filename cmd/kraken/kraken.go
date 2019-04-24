package main

import (
	"github.com/russell-kvashnin/kraken/internal/app/kraken/ioc"
	"github.com/russell-kvashnin/kraken/internal/app/kraken/startup"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		ioc.ConfigModule,
		ioc.ErrorModule,
		ioc.LoggerModule,
		ioc.FSModule,
		ioc.PersistenceModule,
		ioc.MessagingModule,
		ioc.BusModule,
		ioc.FileModule,
		ioc.HealthModule,
		ioc.StartupModule,
		ioc.ApiModule,
		ioc.RpcModule,

		fx.Invoke(func(s *startup.Scenario) error {
			err := s.Execute()
			if err != nil {
				return err
			}

			return nil
		}),
	)

	app.Run()
}
