package blumfield

import (
	"blumfield/config"
	jwts "blumfield/internal/jwt"
	"blumfield/internal/models"
	"blumfield/internal/tools"
	"context"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type Blumfield struct {
	BaseHeaders map[string]string
	client      *resty.Client
	Tokens      *models.BlumTokens
	log         *logrus.Logger
	tools       *tools.Tools
	Config      *config.Config
}

// copy(Telegram.WebApp.initData)
// NewBlumfield constructor
func NewBlumfield(log *logrus.Logger, configName string) (*Blumfield, error) {
	// Define headers
	baseHeaders := map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0",
		"Origin":          "https://telegram.b.codes",
		"Sec-Fetch-Site":  "same-site",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Dest":  "empty",
		"Accept-Encoding": "gzip, deflate",
		"Accept-Language": "ru",
		"Lang":            "en",
	}
	cfg, err := config.LoadConfig(configName)
	if err != nil {
		return nil, err
	}

	return &Blumfield{
		BaseHeaders: baseHeaders,
		client:      resty.New(),
		log:         log,
		tools:       tools.NewTools(),
		Config:      cfg,
		Tokens:      &models.BlumTokens{},
	}, nil
}

func (b *Blumfield) LoadTokensFromFile() error {
	tokens, err := b.Config.LoadTokens()
	if err != nil || jwts.ParseAndCheckToken(tokens.Auth) != nil {
		if err := b.RenewAccessToken(); err != nil {
			return err
		}
		return errors.New("renewed tokens")
	}

	b.BaseHeaders["Authorization"] = "Bearer " + tokens.Auth
	b.Tokens.Token.Access = tokens.Auth
	b.Tokens.Token.Refresh = tokens.Refresh

	return nil
}

func (b *Blumfield) Start(ctx context.Context) error {
	// Loadig Random User-Agent
	if b.Config.Settings.RandomAgent {
		ua, err := b.tools.GetRandomUserAgent("config/user_agent.txt")
		if err != nil {
			return err
		}
		b.BaseHeaders["User-Agent"] = ua
	} else {
		b.log.Info("Using default User-Agent.")
	}

	// Loading auth tokens
	if err := b.LoadTokensFromFile(); err != nil {
		b.log.Error("Error loading tokens: ", err)
	}
	if ctx.Err() != nil {
		return nil
	}

	// Checking daily claim
	if err := b.LogBalance(ctx); err != nil {
		b.log.Error("Error checking balance: ", err)
	}
	b.tools.Delay(2) // Time Delay
	if err := b.ClaimCheckIn(); err != nil {
		b.log.Error("Error claiming check-in reward: ", err)
	}
	if ctx.Err() != nil {
		return nil
	}

	// Claiming farming
	if b.Config.Settings.Farming {
		b.farmWithContext(ctx)
		if ctx.Err() != nil {
			return nil
		}
		b.tools.Delay(2)
	}

	// Complete tasks
	if b.Config.Settings.Tasks {
		tasks, err := b.GetTasks(ctx)
		if err != nil {
			b.log.Error("Error retrieving tasks: ", err)
		}
		if ctx.Err() != nil {
			return nil
		}
		b.CompleteTasks(ctx, tasks)
		b.tools.Delay(3)
	}

	if b.Config.Settings.Gaming {
		if err := b.LogBalance(ctx); err != nil {
			b.log.Error("Error checking balance: ", err)
		}
		b.PlayGame(ctx)
	}
	return nil
}
