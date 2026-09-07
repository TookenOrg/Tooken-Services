package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/TookenOrg/tooken-services/internal/api/server"
	blockchainService "github.com/TookenOrg/tooken-services/internal/blockchain/services"
	"github.com/TookenOrg/tooken-services/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *Handler) DeployAllImplementationsAsync(gCtx *gin.Context) {

	go func() {
		logger.LogInfo("🚀 Starting asynchronous deployment of all implementations")

		_, err := h.blockchainSvc.DeployAllImplementations(context.Background())
		if err != nil {
			logger.LogError("Failed to deploy all implementations: %v", err.Error())
			return
		}
		logger.LogInfo("🆗 Asynchronous deployment of all implementations completed successfully.")
	}()

	gCtx.JSON(http.StatusAccepted, server.APIResponse{
		Message: "The deployment of all implementations has been started asynchronously.",
	})
}

func (h *Handler) DeployIdentityFactory(gCtx *gin.Context) {

	logger.LogInfo("🚀 Starting deployment of Identity Factory")

	deploymentContractDetails, err := h.blockchainSvc.DeployIdentityFactory(gCtx.Request.Context())
	if err != nil {
		logger.LogError("Failed to deploy Identity Factory: %v", err)
		gCtx.JSON(http.StatusInternalServerError, server.InitializeIdentityFactoryResponse{
			Message: "Failed to deploy Identity Factory: " + err.Error(),
		})
	}

	logger.LogInfo("🆗 Identity Factory deployed successfully")

	gCtx.JSON(http.StatusOK, server.InitializeIdentityFactoryResponse{
		Message: "Identity Factory deployed successfully",
		Data:    &deploymentContractDetails,
	})
}

func (h *Handler) ConfigureAuthority(gCtx *gin.Context) {

	logger.LogInfo("🚀 Starting deployment of configure Authority")

	txDetails, err := h.blockchainSvc.ConfigureAuthority(gCtx.Request.Context())
	if err != nil {
		gCtx.JSON(http.StatusInternalServerError, server.ConfigureAuthorityResponse{
			Message: "Failed to configure Authority: " + err.Error(),
		})
		return
	}

	logger.LogInfo("🆗 Authority configured successfully")

	gCtx.JSON(http.StatusOK, server.ConfigureAuthorityResponse{
		Message: "The Authority has been configured successfully.",
		Data:    &txDetails,
	})
}

func (h *Handler) DeployAndInitTrexFactoryAsync(gCtx *gin.Context) {
	go func() {
		logger.LogInfo("🚀 Starting asynchronous deployment of TREX Factory")
		err := h.blockchainSvc.DeployAndInitTrexFactory(context.Background())
		if err != nil {
			logger.LogError("Failed to deploy and init TREX Factory: %v", err.Error())
			return
		}
		logger.LogInfo("🆗 Asynchronous deployment of TREX Factory completed successfully.")
	}()

	gCtx.JSON(http.StatusAccepted, server.APIResponse{
		Message: "The deployment of TREX Factory has been started asynchronously.",
	})
}

func (h *Handler) DeployTrexSuiteAsync(gCtx *gin.Context) {
	go func() {
		logger.LogInfo("🚀 Starting asynchronous deployment of TREX Suite")
		_, err := h.blockchainSvc.DeployTrexSuite(context.Background())
		if err != nil {
			logger.LogError("Failed to deploy TREX Suite: %v", err.Error())
			return
		}
		logger.LogInfo("🆗 Asynchronous deployment of TREX Suite completed successfully.")
	}()

	gCtx.JSON(http.StatusAccepted, server.APIResponse{
		Message: "The deployment of TREX Suite has been started asynchronously.",
	})
}

func (h *Handler) GetTrexSuiteInfos(gCtx *gin.Context, params server.GetTrexSuiteInfosParams) {
	gCtx.JSON(http.StatusOK, server.APIResponse{
		Message: "The Trex Suite information has been retrieved successfully.",
	})
}

func (h *Handler) CreateIdentity(gCtx *gin.Context) {

	logger.LogInfo("🚀 Starting deployment of Identity")

	var req server.CreateIdentityRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusInternalServerError, server.APIResponse{
			Message: "Failed to parse request body: " + err.Error(),
		})
		return
	}

	newIdentity, err := h.blockchainSvc.CreateIdentity(gCtx.Request.Context(), req)
	if err != nil {
		logger.LogError("Failed to deploy Identity: %v", err.Error())
		gCtx.JSON(http.StatusInternalServerError, server.APIResponse{
			Message: "Failed to deploy Identity: " + err.Error(),
		})
		return
	}
	logger.LogInfo("🆗 Deployment of Identity completed successfully.")

	gCtx.JSON(http.StatusOK, server.CreateIdentityResponse{
		Message: "The Identity has been created successfully.",
		Data:    &newIdentity,
	})
}

func (h *Handler) AddClaimToIdentity(gCtx *gin.Context) {

	logger.LogInfo("🚀 Starting adding claims to Identity")

	var req server.AddClaimRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusInternalServerError, server.AddClaimResponse{
			Message: "Failed to parse request body: " + err.Error(),
		})
		return
	}

	tx, err := h.blockchainSvc.AddClaimToIdentity(gCtx.Request.Context(), req)
	if err != nil {
		gCtx.JSON(http.StatusBadRequest, server.AddClaimResponse{
			Message: err.Error(),
		})
		return
	}

	logger.LogInfo("🆗 Add claim completed successfully.")

	gCtx.JSON(http.StatusOK, server.AddClaimResponse{
		Message: "The claim has been added to the Identity successfully.",
		Data:    &server.TxHashName{TransactionHash: tx.Hash().Hex(), OperationName: "AddClaimToIdentity"},
	})
}

func (h *Handler) CreateTokenContract(gCtx *gin.Context) {

	logger.LogInfo("🚀 Starting creation of Token Contract")

	var req server.CreateTokenRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusInternalServerError, server.CreateTokenResponse{
			Message: "Failed to parse request body: " + err.Error(),
			Data:    &server.TokenInfos{},
		})
		return
	}

	tokenInfos, err := h.blockchainSvc.CreateToken(gCtx.Request.Context(), req)
	if err != nil {
		gCtx.JSON(http.StatusBadRequest, server.CreateTokenResponse{
			Message: "The token creation failed." + err.Error(),
		})
		return
	}

	logger.LogInfo("🆗 Token Contract created successfully.")

	gCtx.JSON(http.StatusOK, server.CreateTokenResponse{
		Message: "The Token Contract has been created successfully.",
		Data:    &tokenInfos,
	})
}

func (h *Handler) GetTokenInfos(gCtx *gin.Context, tokenAddress string) {
	gCtx.JSON(http.StatusOK, server.APIResponse{
		Data:    "",
		Message: "The Token information has been retrieved successfully.",
	})
}

func (h *Handler) MintTokenAsync(gCtx *gin.Context) {

	logger.LogInfo("🚀 Starting mint on token")

	var req server.MintTokenRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusInternalServerError, server.MintTokenResponse{
			Message: "Failed to parse request body: " + err.Error(),
		})
		return
	}

	reqCtx := gCtx.Request.Context()

	// Validate before answering 202: the work runs in a goroutine, so this is
	// the only moment where the caller can still be told the request is wrong.
	if err := h.blockchainSvc.ValidateMintBurnInput(reqCtx, req.TokenContractAddress, req.To, req.Amount); err != nil {
		if errors.Is(err, blockchainService.ErrInvalidMintBurnInput) {
			logger.LogWarn("Refused mint request: %s", err.Error())
			gCtx.JSON(http.StatusBadRequest, server.MintTokenResponse{Message: err.Error()})
			return
		}

		logger.LogError("Failed to validate the mint request: %s", err.Error())
		gCtx.JSON(http.StatusInternalServerError, server.MintTokenResponse{
			Message: "Unable to validate the request.",
		})
		return
	}

	ctx := context.WithoutCancel(reqCtx)

	go func(ctx context.Context) {
		_, err := h.blockchainSvc.Mint(ctx, req.TokenContractAddress, req.To, req.Amount)
		if err != nil {
			logger.LogError("Failed to Mint: %s", err.Error())
			return
		}
	}(ctx)

	gCtx.JSON(http.StatusAccepted, server.MintTokenResponse{
		Message: "Mint started successfully.",
	})
}

func (h *Handler) BurnTokenAsync(gCtx *gin.Context) {
	logger.LogInfo("🚀 Starting burn on token")

	var req server.BurnTokenRequest
	if err := gCtx.ShouldBindJSON(&req); err != nil {
		gCtx.JSON(http.StatusInternalServerError, server.BurnTokenResponse{
			Message: "Failed to parse request body: " + err.Error(),
		})
		return
	}

	reqCtx := gCtx.Request.Context()

	// Validate before answering 202: the work runs in a goroutine, so this is
	// the only moment where the caller can still be told the request is wrong.
	if err := h.blockchainSvc.ValidateMintBurnInput(reqCtx, req.TokenContractAddress, req.To, req.Amount); err != nil {
		if errors.Is(err, blockchainService.ErrInvalidMintBurnInput) {
			logger.LogWarn("Refused burn request: %s", err.Error())
			gCtx.JSON(http.StatusBadRequest, server.BurnTokenResponse{Message: err.Error()})
			return
		}

		logger.LogError("Failed to validate the burn request: %s", err.Error())
		gCtx.JSON(http.StatusInternalServerError, server.BurnTokenResponse{
			Message: "Unable to validate the request.",
		})
		return
	}

	ctx := context.WithoutCancel(reqCtx)

	go func(ctx context.Context) {
		_, err := h.blockchainSvc.Burn(ctx, req.TokenContractAddress, req.To, req.Amount)
		if err != nil {
			logger.LogError("Failed to Burn: %s", err.Error())
			return
		}
	}(ctx)

	gCtx.JSON(http.StatusAccepted, server.BurnTokenResponse{
		Message: "Burn started successfully.",
	})
}
