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
	config      *config.Config
}

// copy(Telegram.WebApp.initData)
// NewBlumfield constructor
func NewBlumfield(log *logrus.Logger) (*Blumfield, error) {
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
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	return &Blumfield{
		BaseHeaders: baseHeaders,
		client:      resty.New(),
		log:         log,
		tools:       tools.NewTools(),
		config:      cfg,
		Tokens:      &models.BlumTokens{},
	}, nil
}

func (b *Blumfield) LoadTokensFromFile() error {
	tokens, err := b.config.LoadTokens()
	if err != nil {
		if err := jwts.ParseAndCheckToken(tokens.Auth); err != nil {
			if err := b.RenewAccessToken(); err != nil {
				return err
			}
			return errors.New("renewed tokens")
		}
	}

	b.BaseHeaders["Authorization"] = "Bearer " + tokens.Auth
	b.Tokens.Token.Access = tokens.Auth
	b.Tokens.Token.Refresh = tokens.Refresh

	return nil
}

func (b *Blumfield) Start(ctx context.Context) error {
	// Loading auth tokens
	if err := b.LoadTokensFromFile(); err != nil {
		b.log.Error("Error loading tokens: ", err)
	}

	// Checking daily claim
	if err := b.LogBalance(); err != nil {
		b.log.Error("Error checking balance: ", err)
	}
	b.tools.Delay(2) // Time Delay
	if err := b.ClaimCheckIn(); err != nil {
		b.log.Error("Error claiming check-in reward: ", err)
	}

	// Claiming farming
	if err := b.LogBalance(); err != nil {
		b.log.Error("Error checking balance: ", err)
	}
	b.tools.Delay(3) // Time Delay
	farmClaim, err := b.ClaimFarming()
	if err != nil {
		b.log.Error("Error claiming farming: ", err)
	}
	if farmClaim.AvailableBalance == "" {
		b.log.Info("Already claimed.")
	} else {
		b.log.Info("Claimed! Balance: ", farmClaim.AvailableBalance)
	}

	// Start farming
	if err := b.LogBalance(); err != nil {
		b.log.Error("Error checking balance: ", err)
	}
	b.tools.Delay(2) // Time Delay
	farmStart, err := b.StartFarming()
	if err != nil {
		b.log.Error("Error starting farming: ", err)
	}
	b.log.Info("Farming: ", farmStart.EarningsRate)

	// Complete tasks
	tasks, err := b.GetTasks()
	if err != nil {
		b.log.Error("Error retrieving tasks: ", err)
	}
	b.CompleteTasks(ctx, tasks)

	b.tools.Delay(3) // Time Delay
	if err := b.LogBalance(); err != nil {
		b.log.Error("Error checking balance: ", err)
	}
	b.PlayGame(ctx)
	return nil
}
