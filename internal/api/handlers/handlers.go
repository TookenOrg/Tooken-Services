package handlers

import (
	"context"
	"sync"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	assetsService "github.com/TookenOrg/tooken-services/internal/assets_managements/services"
	authService "github.com/TookenOrg/tooken-services/internal/auth/services"
	blockchainService "github.com/TookenOrg/tooken-services/internal/blockchain/services"
	deliveryService "github.com/TookenOrg/tooken-services/internal/delivery/services"
	ordersService "github.com/TookenOrg/tooken-services/internal/orders/services"
	paymentsService "github.com/TookenOrg/tooken-services/internal/payments/services"
	usersService "github.com/TookenOrg/tooken-services/internal/users/services"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
)

// orderDeliverer is the delivery seen from the handler: one call, no answer.
// Declaring it here rather than taking the concrete service lets a test put a
// no-op in its place, so the handler suite can check the HTTP contract without
// a chain running behind it.
type orderDeliverer interface {
	Deliver(ctx context.Context, orderRef string)
}

var _ server.ServerInterface = (*Handler)(nil)

type Handler struct {
	blockchainSvc *blockchainService.Service
	paymentsSvc   *paymentsService.Service
	authSvc       *authService.Service
	userSvc       *usersService.Service
	realEstateSvc *assetsService.Service
	orderSvc      *ordersService.Service
	deliverySvc   orderDeliverer

	// background counts the jobs a request started and left running after
	// its answer: deployments, burns, deliveries. Start them with
	// startBackground, never with a bare go statement, or neither
	// WaitBackground nor CloseAndWait can see them.
	background sync.WaitGroup

	// closed is set by CloseAndWait, under closing: from then on no job
	// starts, which is what makes waiting for the running ones safe.
	closing sync.Mutex
	closed  bool
}

// startBackground runs job in the background and tracks it.
//
// Once CloseAndWait was called the job is dropped instead, and logged under
// what. The request asking for it gets no answer anyway: the server closes
// the connections before it closes the handler.
func (h *Handler) startBackground(what string, job func()) {
	h.closing.Lock()
	defer h.closing.Unlock()

	if h.closed {
		logger.LogError("❌ Not started, the server is shutting down: %s", what)
		return
	}
	h.background.Go(job)
}

// WaitBackground blocks until every job started in the background by a
// request has returned.
//
// Those jobs outlive the request on purpose, and keep writing to the database
// and to the chain after the answer was sent. Whoever tears the process down
// has to wait for them first: a test before it drops its database, a server
// before it exits.
//
// It only waits for the jobs already started. The requests that start them
// must have returned, or stopped arriving, before it is called: a job started
// while it waits on an empty set may be missed, which is why a server has to
// stop accepting requests first.
func (h *Handler) WaitBackground() {
	h.background.Wait()
}

// CloseAndWait stops the handler from starting background jobs, then blocks
// until the running ones have returned. It is the last step of a shutdown,
// once the server no longer serves requests.
//
// Unlike WaitBackground it is safe while requests are still running, as when
// the server gave up waiting for them: a job they try to start is refused
// rather than missed by the wait. It is meant to be called once; the handler
// starts no job afterwards.
func (h *Handler) CloseAndWait() {
	h.closing.Lock()
	h.closed = true
	h.closing.Unlock()

	h.background.Wait()
}

func NewHandler() *Handler {

	// Share the same blockchain service instance across all dependent services.
	blockchainSvc := blockchainService.NewService()

	return &Handler{
		blockchainSvc: blockchainSvc,
		paymentsSvc:   paymentsService.NewService(),
		authSvc:       authService.NewService(),
		realEstateSvc: assetsService.NewService(blockchainSvc),
		orderSvc:      ordersService.NewService(),
		userSvc:       usersService.NewService(),
		deliverySvc:   deliveryService.NewService(blockchainSvc),
	}
}

func (h *Handler) GetHealth(gCtx *gin.Context) {
	gCtx.JSON(200, gin.H{"status": "ok", "version": "v1.0.0"})
}
