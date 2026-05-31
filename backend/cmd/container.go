package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/hada-commerce/internal/agent"
	"github.com/Abraxas-365/hada-commerce/internal/agent/agentapi"
	"github.com/Abraxas-365/hada-commerce/internal/agentsession"
	"github.com/Abraxas-365/hada-commerce/internal/agentsession/agentsessioncontainer"
	"github.com/Abraxas-365/hada-commerce/internal/marketplace/marketplacesrv"
	"github.com/Abraxas-365/hada-commerce/internal/agentsession/agentsessionsrv"
	"github.com/Abraxas-365/hada-commerce/internal/containerx"
	"github.com/Abraxas-365/hada-commerce/internal/containerxdocker"
	"github.com/Abraxas-365/hada-commerce/internal/agentmemory/agentmemorycontainer"
	"github.com/Abraxas-365/hada-commerce/internal/agenttrigger/agenttriggercontainer"
	"github.com/Abraxas-365/hada-commerce/internal/approval/approvalcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/abtest/abtestcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/analytics/analyticscontainer"
	"github.com/Abraxas-365/hada-commerce/internal/bulkops/bulkopscontainer"
	"github.com/Abraxas-365/hada-commerce/internal/dashboard/dashboardcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/inventory/inventorycontainer"
	"github.com/Abraxas-365/hada-commerce/internal/audit/auditcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/cartrecovery/cartrecoverycontainer"
	"github.com/Abraxas-365/hada-commerce/internal/sitemap"
	"github.com/Abraxas-365/hada-commerce/internal/customer/customersrv"
	"github.com/Abraxas-365/hada-commerce/internal/emails"
	"github.com/Abraxas-365/hada-commerce/internal/order/ordersrv"
	"github.com/Abraxas-365/hada-commerce/internal/cart/cartcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/importexport"
	"github.com/Abraxas-365/hada-commerce/internal/customergroup/customergroupcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/checkout/checkoutcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/catalog/catalogcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/config"
	"github.com/Abraxas-365/hada-commerce/internal/customer/customercontainer"
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/fsx"
	"github.com/Abraxas-365/hada-commerce/internal/fsx/fsxlocal"
	"github.com/Abraxas-365/hada-commerce/internal/fsx/fsxs3"
	"github.com/Abraxas-365/hada-commerce/internal/iam/iamcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/jobx"
	"github.com/Abraxas-365/hada-commerce/internal/jobx/jobxredis"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/logx"
	"github.com/Abraxas-365/hada-commerce/internal/marketplace/marketplacecontainer"
	"github.com/Abraxas-365/hada-commerce/internal/media/mediacontainer"
	"github.com/Abraxas-365/hada-commerce/internal/notifx"
	"github.com/Abraxas-365/hada-commerce/internal/notifx/notifxconsole"
	"github.com/Abraxas-365/hada-commerce/internal/notifx/notifxses"
	"github.com/Abraxas-365/hada-commerce/internal/order/ordercontainer"
	"github.com/Abraxas-365/hada-commerce/internal/payment/paymentcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/product/productcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/giftcard/giftcardcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/promo/promocontainer"
	"github.com/Abraxas-365/hada-commerce/internal/search/searchcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/settings/settingscontainer"
	"github.com/Abraxas-365/hada-commerce/internal/plugin/plugincontainer"
	"github.com/Abraxas-365/hada-commerce/internal/shipping/shippingcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/storefrontcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/currency/currencycontainer"
	"github.com/Abraxas-365/hada-commerce/internal/i18n/i18ncontainer"
	"github.com/Abraxas-365/hada-commerce/internal/subscription/subscriptioncontainer"
	"github.com/Abraxas-365/hada-commerce/internal/review/reviewcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/returns/returnscontainer"
	"github.com/Abraxas-365/hada-commerce/internal/tax/taxcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/webhook/webhookcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/theme/themecontainer"
	"github.com/Abraxas-365/hada-commerce/internal/socialauth/socialauthcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/wishlist/wishlistcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/loyalty/loyaltycontainer"
	"github.com/Abraxas-365/hada-commerce/internal/blog/blogcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/bundle/bundlecontainer"
	"github.com/Abraxas-365/hada-commerce/internal/multistore/multistorecontainer"
	"github.com/Abraxas-365/hada-commerce/internal/notification/notificationcontainer"
	"github.com/Abraxas-365/hada-commerce/internal/collection/collectioncontainer"
	"github.com/Abraxas-365/hada-commerce/internal/recommendation/recommendationcontainer"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// Container holds shared infrastructure and composed module containers.
type Container struct {
	Config *config.Config

	// Infrastructure
	DB    *sqlx.DB
	Redis *redis.Client

	// File storage
	FileSystem fsx.FileSystem
	S3Client   *s3.Client

	// Background services
	JobClient    *jobx.Client
	NotifxClient *notifx.Client

	// Event bus
	EventBus eventbus.Bus

	// IAM
	IAM *iamcontainer.Container

	// Audit log
	Audit *auditcontainer.Container

	// Commerce domains
	Cart           *cartcontainer.Container
	Product        *productcontainer.Container
	Order          *ordercontainer.Container
	Payment        *paymentcontainer.Container
	Customer       *customercontainer.Container
	CustomerGroup  *customergroupcontainer.Container
	Catalog        *catalogcontainer.Container
	Storefront     *storefrontcontainer.Container
	GiftCard       *giftcardcontainer.Container
	Promo          *promocontainer.Container
	Media          *mediacontainer.Container
	Marketplace    *marketplacecontainer.Container
	Analytics      *analyticscontainer.Container
	Dashboard      *dashboardcontainer.Container
	Settings       *settingscontainer.Container
	Theme          *themecontainer.Container
	Plugin         *plugincontainer.Container
	Search         *searchcontainer.Container
	Shipping       *shippingcontainer.Container
	Tax            *taxcontainer.Container
	Currency       *currencycontainer.Container
	I18n           *i18ncontainer.Container
	Checkout       *checkoutcontainer.Container
	ImportExport   *importexport.Container
	Sitemap        *sitemap.Container
	Wishlist       *wishlistcontainer.Container
	CartRecovery   *cartrecoverycontainer.Container
	Subscription   *subscriptioncontainer.Container
	Inventory      *inventorycontainer.Container
	Review         *reviewcontainer.Container
	Returns        *returnscontainer.Container
	Webhook        *webhookcontainer.Container
	Loyalty        *loyaltycontainer.Container
	Bundle         *bundlecontainer.Container
	SocialAuth     *socialauthcontainer.Container
	Notification   *notificationcontainer.Container
	MultiStore     *multistorecontainer.Container
	BulkOps        *bulkopscontainer.Container
	Blog           *blogcontainer.Container
	Collection     *collectioncontainer.Container
	ABTest         *abtestcontainer.Container
	Recommendation  *recommendationcontainer.Container

	// AI Agent
	Agent        *agentapi.Handler
	AgentSession *agentsessioncontainer.Container
	AgentMemory  *agentmemorycontainer.Container
	AgentTrigger *agenttriggercontainer.Container
	Approval     *approvalcontainer.Container
}

func NewContainer(cfg *config.Config) *Container {
	logx.Info("Initializing application container...")

	c := &Container{Config: cfg}

	c.initInfrastructure()
	c.initModules()

	logx.Info("Application container initialized")
	return c
}

// ---------------------------------------------------------------------------
// Infrastructure
// ---------------------------------------------------------------------------

func (c *Container) initInfrastructure() {
	logx.Info("Initializing infrastructure...")

	// Database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Config.Database.Host,
		c.Config.Database.Port,
		c.Config.Database.User,
		c.Config.Database.Password,
		c.Config.Database.Name,
		c.Config.Database.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logx.Fatalf("Failed to connect to database: %v", err)
	}
	db.SetMaxOpenConns(c.Config.Database.MaxOpenConns)
	db.SetMaxIdleConns(c.Config.Database.MaxIdleConns)
	db.SetConnMaxLifetime(c.Config.Database.ConnMaxLifetime)
	c.DB = db
	logx.Info("  Database connected")

	// Redis
	c.Redis = redis.NewClient(&redis.Options{
		Addr:     c.Config.Redis.Address(),
		Password: c.Config.Redis.Password,
		DB:       c.Config.Redis.DB,
	})
	if _, err := c.Redis.Ping(context.Background()).Result(); err != nil {
		logx.Fatalf("Failed to connect to Redis: %v (Redis is required)", err)
	}
	logx.Info("  Redis connected")

	logx.Info("Infrastructure initialized")
}

// ---------------------------------------------------------------------------
// Module composition
// ---------------------------------------------------------------------------

func (c *Container) initModules() {
	logx.Info("Initializing modules...")

	c.initFileStorage()
	c.initJobx()
	c.initNotifx()

	// IAM — uses notifx for OTP and invitation emails
	c.IAM = iamcontainer.New(iamcontainer.Deps{
		DB:                 c.DB,
		Redis:              c.Redis,
		Cfg:                c.Config,
		OTPNotifier:        NewNotifxOTPNotifier(c.NotifxClient),
		InvitationNotifier: NewNotifxInvitationNotifier(c.NotifxClient),
	})

	// Event bus
	bus := eventbus.NewInMemoryBus()
	bus.SubscribeAll(func(ctx context.Context, event eventbus.Event) error {
		logx.WithFields(logx.Fields{
			"event_id":   event.ID,
			"event_type": string(event.Type),
			"tenant_id":  string(event.TenantID),
		}).Info("domain event published")
		return nil
	})
	c.EventBus = bus

	// Audit log
	c.Audit = auditcontainer.New(c.DB)

	// Commerce domains
	c.Cart = cartcontainer.New(c.DB, bus)
	c.Product = productcontainer.New(c.DB, bus)
	c.Order = ordercontainer.New(c.DB, bus)
	c.Payment = paymentcontainer.New(c.DB, bus)
	c.Customer = customercontainer.New(c.DB, bus, c.Config.Auth.JWT.SecretKey, c.Order.Service)
	c.CustomerGroup = customergroupcontainer.New(c.DB)
	c.Catalog = catalogcontainer.New(c.DB, bus)
	c.Theme = themecontainer.New(c.DB, bus)
	// Settings must be initialized before Storefront so the renderer can fetch store branding.
	c.Settings = settingscontainer.New(c.DB, bus)
	c.Storefront = storefrontcontainer.New(c.DB, bus, c.Theme.Service, storefrontcontainer.Deps{
		ProductLister:    c.Product.Service,
		CollectionGetter: c.Catalog.Service,
		SettingsGetter:   c.Settings.Service,
	})
	c.GiftCard = giftcardcontainer.New(c.DB)
	c.Promo = promocontainer.New(c.DB)
	mediaCont, err := mediacontainer.New(c.DB, "./uploads", fmt.Sprintf("http://localhost:%d/uploads", c.Config.Server.Port))
	if err != nil {
		logx.Fatal(fmt.Sprintf("Failed to initialize media module: %v", err))
	}
	c.Media = mediaCont
	c.Marketplace = marketplacecontainer.New(c.DB)
	c.Analytics = analyticscontainer.New(c.DB)
	c.Dashboard = dashboardcontainer.New(c.DB)
	c.Plugin = plugincontainer.New(c.DB, bus)
	c.Search = searchcontainer.New(c.DB)
	c.Shipping = shippingcontainer.New(c.DB, bus)
	c.Tax = taxcontainer.New(c.DB, bus)
	c.Currency = currencycontainer.New(c.DB)
	c.I18n = i18ncontainer.New(c.DB)
	c.Checkout = checkoutcontainer.New(
		c.DB, bus,
		c.Cart.Service,
		c.Order.Service,
		c.Shipping.Service,
		c.Tax.Service,
		c.Payment.Service,
		c.Promo.Service,
	)
	c.Wishlist = wishlistcontainer.New(c.DB)
	c.CartRecovery = cartrecoverycontainer.New(c.DB)
	c.Subscription = subscriptioncontainer.New(c.DB, bus)
	c.Inventory = inventorycontainer.New(c.DB, bus)
	c.Review = reviewcontainer.New(c.DB, bus)
	c.Returns = returnscontainer.New(c.DB, bus)
	c.Webhook = webhookcontainer.New(c.DB, bus)
	c.Loyalty = loyaltycontainer.New(c.DB, bus)
	c.Bundle = bundlecontainer.New(c.DB, bus)
	c.SocialAuth = socialauthcontainer.New(c.DB)
	c.Notification = notificationcontainer.New(c.DB, bus)
	c.MultiStore = multistorecontainer.New(c.DB, bus)
	c.BulkOps = bulkopscontainer.New(c.DB, bus)
	c.Blog = blogcontainer.New(c.DB, bus)
	c.Collection = collectioncontainer.New(c.DB, bus)
	c.ABTest = abtestcontainer.New(c.DB, bus)
	c.Recommendation = recommendationcontainer.New(c.DB)

	// Import/Export — depends on Product, Order, and Customer services.
	c.ImportExport = importexport.New(
		c.Product.Service,
		c.Order.Service,
		c.Customer.Service,
		c.Product.Service,
	)

	// Sitemap — reads products, categories/collections, and published pages.
	c.Sitemap = sitemap.New(c.Product.Service, c.Catalog.Service, c.Storefront.Service)

	// Transactional email notifications — wired last so all domain containers exist.
	emailHandler := emails.New(
		c.NotifxClient,
		c.Config.Notifx.FromAddress,
		c.Config.Notifx.FromName,
		newCustomerEmailResolver(c.Customer.Service),
		newOrderCustomerResolver(c.Order.Service),
	)
	emailHandler.RegisterSubscriptions(bus)
	logx.Info("  Email notifications wired to event bus")

	// Approval workflow — human-in-the-loop for sensitive agent actions.
	c.Approval = approvalcontainer.New(c.DB)
	logx.Info("  Approval workflow initialized")

	// Agent memory — persistent knowledge base.
	c.AgentMemory = agentmemorycontainer.New(c.DB)
	logx.Info("  Agent memory initialized")

	// Agent triggers — event-driven agent actions.
	c.AgentTrigger = agenttriggercontainer.New(c.DB)
	logx.Info("  Agent triggers initialized")

	// Agent Sessions — workspace management via Docker containers.
	// Initialized before the AI Agent so the workspace provider can be passed in.
	var containerMgr containerx.Manager
	if mgr, err := containerxdocker.New(); err != nil {
		logx.Warnf("Docker unavailable, agent sessions disabled: %v", err)
	} else {
		containerMgr = mgr
		c.AgentSession = agentsessioncontainer.New(agentsessioncontainer.Deps{
			DB:        c.DB,
			Manager:   containerMgr,
			PresetSvc: c.Marketplace.PresetService,
		})
		logx.Info("  Agent session manager initialized")
	}

	// AI Agent — optional, only initialized when API key is configured.
	if c.Config.Agent.APIKey != "" {
		agentSvc := c.BuildAgentServices()
		var wp agentapi.WorkspaceProvider
		var cp agentapi.ChatPersister
		if c.AgentSession != nil {
			wp = &workspaceProviderAdapter{svc: c.AgentSession.Service}
			cp = &chatPersisterAdapter{repo: c.AgentSession.ChatRepo}
		}
		c.Agent = agentapi.NewHandler(
			c.Config.Agent.APIKey,
			c.Config.Agent.Model,
			"", // use default system prompt
			agentSvc,
			&presetProviderAdapter{svc: c.Marketplace.PresetService},
			cp,
			wp,
			containerMgr,
		)
		logx.Info("  AI Agent chat handler initialized")
	}

	logx.Info("All modules initialized")
}

// ---------------------------------------------------------------------------
// Lifecycle
// ---------------------------------------------------------------------------

func (c *Container) StartBackgroundServices(ctx context.Context) {
	logx.Info("Starting background services...")
	go c.JobClient.Start(ctx)
	c.IAM.StartBackgroundServices(ctx)
}

func (c *Container) Cleanup() {
	logx.Info("Cleaning up resources...")

	if c.DB != nil {
		if err := c.DB.Close(); err != nil {
			logx.Errorf("Error closing database: %v", err)
		} else {
			logx.Info("  Database connection closed")
		}
	}

	if c.Redis != nil {
		if err := c.Redis.Close(); err != nil {
			logx.Errorf("Error closing Redis: %v", err)
		} else {
			logx.Info("  Redis connection closed")
		}
	}

	logx.Info("Cleanup complete")
}

// ---------------------------------------------------------------------------
// File storage
// ---------------------------------------------------------------------------

func (c *Container) initFileStorage() {
	storageMode := getEnv("STORAGE_MODE", "local")

	switch storageMode {
	case "s3":
		awsRegion := getEnv("AWS_REGION", "us-east-1")
		awsBucket := getEnv("AWS_BUCKET", "hada-uploads")

		cfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(awsRegion))
		if err != nil {
			logx.Fatalf("Unable to load AWS SDK config: %v", err)
		}
		c.S3Client = s3.NewFromConfig(cfg)
		c.FileSystem = fsxs3.NewS3FileSystem(c.S3Client, awsBucket, "")
		logx.Infof("  S3 file system configured (bucket: %s, region: %s)", awsBucket, awsRegion)

	case "local":
		uploadDir := getEnv("UPLOAD_DIR", "./uploads")
		localFS, err := fsxlocal.NewLocalFileSystem(uploadDir)
		if err != nil {
			logx.Fatalf("Failed to initialize local file system: %v", err)
		}
		c.FileSystem = localFS
		logx.Infof("  Local file system configured (path: %s)", localFS.GetBasePath())

	default:
		logx.Fatalf("Unknown STORAGE_MODE: %s (use 'local' or 's3')", storageMode)
	}
}

// ---------------------------------------------------------------------------
// Job queue
// ---------------------------------------------------------------------------

func (c *Container) initJobx() {
	queue := jobxredis.NewRedisQueue(c.Redis)
	c.JobClient = jobx.NewClient(queue,
		jobx.WithConcurrency(c.Config.Jobx.Concurrency),
		jobx.WithQueues(c.Config.Jobx.Queues...),
		jobx.WithPollInterval(c.Config.Jobx.PollInterval),
		jobx.WithShutdownTimeout(c.Config.Jobx.ShutdownTimeout),
		jobx.WithDequeueTimeout(c.Config.Jobx.DequeueTimeout),
		jobx.WithDefaultRetryDelay(c.Config.Jobx.DefaultRetryDelay),
	)
	logx.Info("  Job queue configured")
}

// ---------------------------------------------------------------------------
// Notifications
// ---------------------------------------------------------------------------

func (c *Container) initNotifx() {
	var provider notifx.EmailSender

	switch c.Config.Notifx.Provider {
	case "ses":
		awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
			awsConfig.WithRegion(c.Config.Notifx.AWSRegion))
		if err != nil {
			logx.Fatalf("Unable to load AWS config for notifx: %v", err)
		}
		sesClient := ses.NewFromConfig(awsCfg)
		provider = notifxses.NewSESProvider(sesClient, c.Config.Notifx.FromAddress)
		logx.Infof("  Notifx: SES provider (region: %s)", c.Config.Notifx.AWSRegion)

	default:
		provider = notifxconsole.NewConsoleProvider()
		logx.Info("  Notifx: console provider (dev mode)")
	}

	c.NotifxClient = notifx.NewClient(provider)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// ---------------------------------------------------------------------------
// Notification adapters for IAM
// ---------------------------------------------------------------------------

// NotifxOTPNotifier implements otp.NotificationService using notifx
type NotifxOTPNotifier struct {
	client *notifx.Client
}

func NewNotifxOTPNotifier(client *notifx.Client) *NotifxOTPNotifier {
	return &NotifxOTPNotifier{client: client}
}

func (n *NotifxOTPNotifier) SendOTP(ctx context.Context, contact string, code string) error {
	return n.client.SendEmail(ctx, notifx.EmailMessage{
		To:       []string{contact},
		Subject:  "Your verification code",
		HTMLBody: fmt.Sprintf("<h2>Your verification code is: <strong>%s</strong></h2><p>This code will expire shortly.</p>", code),
		TextBody: fmt.Sprintf("Your verification code is: %s", code),
	})
}

// NotifxInvitationNotifier implements invitation.NotificationService using notifx
type NotifxInvitationNotifier struct {
	client *notifx.Client
}

func NewNotifxInvitationNotifier(client *notifx.Client) *NotifxInvitationNotifier {
	return &NotifxInvitationNotifier{client: client}
}

func (n *NotifxInvitationNotifier) SendInvitation(ctx context.Context, email string, token string, tenantID kernel.TenantID, invitedBy kernel.UserID) error {
	return n.client.SendEmail(ctx, notifx.EmailMessage{
		To:       []string{email},
		Subject:  "You've been invited to join",
		HTMLBody: fmt.Sprintf("<h2>You've been invited!</h2><p>Use the following token to accept your invitation: <strong>%s</strong></p>", token),
		TextBody: fmt.Sprintf("You've been invited! Use the following token to accept your invitation: %s", token),
	})
}

// ---------------------------------------------------------------------------
// Email resolver adapters for the emails package
// ---------------------------------------------------------------------------

// customerEmailResolver implements emails.EmailResolver using the customer service.
type customerEmailResolver struct {
	customers *customersrv.Service
}

func newCustomerEmailResolver(svc *customersrv.Service) *customerEmailResolver {
	return &customerEmailResolver{customers: svc}
}

func (r *customerEmailResolver) ResolveEmail(ctx context.Context, tenantID string, customerID string) (string, error) {
	c, err := r.customers.GetByID(ctx, kernel.TenantID(tenantID), kernel.CustomerID(customerID))
	if err != nil {
		return "", err
	}
	return string(c.Email), nil
}

// orderCustomerResolver implements emails.OrderCustomerResolver using the order service.
type orderCustomerResolver struct {
	orders *ordersrv.Service
}

func newOrderCustomerResolver(svc *ordersrv.Service) *orderCustomerResolver {
	return &orderCustomerResolver{orders: svc}
}

func (r *orderCustomerResolver) ResolveOrderCustomerID(ctx context.Context, tenantID string, orderID string) (string, error) {
	o, err := r.orders.GetByID(ctx, kernel.TenantID(tenantID), kernel.OrderID(orderID))
	if err != nil {
		return "", err
	}
	return string(o.CustomerID), nil
}

// presetProviderAdapter adapts marketplacesrv.PresetService to agentapi.PresetProvider.
type presetProviderAdapter struct {
	svc *marketplacesrv.PresetService
}

func (a *presetProviderAdapter) GetPresetSystemPrompt(ctx context.Context, presetID string) (string, error) {
	p, err := a.svc.Get(ctx, kernel.PresetID(presetID))
	if err != nil {
		return "", err
	}
	return p.SystemPrompt, nil
}

func (a *presetProviderAdapter) GetPresetToolsManifest(ctx context.Context, presetID string) (json.RawMessage, error) {
	p, err := a.svc.Get(ctx, kernel.PresetID(presetID))
	if err != nil {
		return nil, err
	}
	return p.ToolsManifest, nil
}

// chatPersisterAdapter adapts agentsession.ChatRepository to agentapi.ChatPersister.
type chatPersisterAdapter struct {
	repo agentsession.ChatRepository
}

func (a *chatPersisterAdapter) SaveMessage(ctx context.Context, sessionID, role, content string, toolCalls json.RawMessage) error {
	msg := agentsession.ChatMessage{
		ID:        uuid.New().String(),
		SessionID: kernel.AgentSessionID(sessionID),
		Role:      role,
		Content:   content,
		CreatedAt: time.Now(),
	}
	return a.repo.SaveMessage(ctx, msg)
}

// workspaceProviderAdapter adapts agentsessionsrv.Service to agentapi.WorkspaceProvider.
type workspaceProviderAdapter struct {
	svc *agentsessionsrv.Service
}

func (a *workspaceProviderAdapter) GetActiveWorkspace(ctx context.Context, tenantID, sessionID string) (string, string, error) {
	sess, err := a.svc.GetSession(ctx, kernel.TenantID(tenantID), kernel.AgentSessionID(sessionID))
	if err != nil {
		return "", "", err
	}
	if sess.Status != agentsession.SessionStatusRunning || sess.ContainerID == "" {
		return "", "", nil
	}
	return sess.ContainerID, sess.FrontendURL, nil
}

// ---------------------------------------------------------------------------
// Agent tools wiring
// ---------------------------------------------------------------------------

// BuildAgentServices returns a fully populated agent.Services struct
// that maps every domain container's Service to the agent tool layer.
func (c *Container) BuildAgentServices() agent.Services {
	return agent.Services{
		Storefront:      c.Storefront.Service,
		Products:        c.Product.Service,
		Orders:          c.Order.Service,
		Promos:          c.Promo.Service,
		Catalog:         c.Catalog.Service,
		Themes:          c.Theme.Service,
		Shipping:        c.Shipping.Service,
		Tax:             c.Tax.Service,
		Payment:         c.Payment.Service,
		Search:          c.Search.Service,
		CustomerGroups:  c.CustomerGroup.Service,
		GiftCards:       c.GiftCard.Service,
		CartRecovery:    c.CartRecovery.Service,
		Currency:        c.Currency.Service,
		I18n:            c.I18n.Service,
		Subscriptions:   c.Subscription.Service,
		Inventory:       c.Inventory.Service,
		Reviews:         c.Review.Service,
		Returns:         c.Returns.Service,
		Webhooks:        c.Webhook.Service,
		Audit:           c.Audit.Service,
		Loyalty:         c.Loyalty.Service,
		Bundles:         c.Bundle.Service,
		Dashboard:       c.Dashboard.Service,
		Notifications:   c.Notification.Service,
		MultiStore:      c.MultiStore.Service,
		BulkOps:         c.BulkOps.Service,
		Blog:            c.Blog.Service,
		Collections:     c.Collection.Service,
		ABTest:          c.ABTest.Service,
		Recommendations: c.Recommendation.Service,
		Memory:          c.AgentMemory.Service,
		Approval:        c.Approval.Service,
	}
}
