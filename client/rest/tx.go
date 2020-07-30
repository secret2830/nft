package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"

	"github.com/irismod/nft/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	// Mint an NFT
	r.HandleFunc(
		"/nft/nfts/denoms/issue",
		issueDenomHandlerFn(cliCtx),
	).Methods("POST")

	// Mint an NFT
	r.HandleFunc(
		"/nft/nfts/mint",
		mintNFTHandlerFn(cliCtx),
	).Methods("POST")

	// Update an NFT tokenData
	r.HandleFunc(
		fmt.Sprintf("/nft/nfts/{%s}/{%s}", RestParamDenom, RestParamTokenID),
		editNFTHandlerFn(cliCtx),
	).Methods("PUT")

	// Transfer an NFT to an address
	r.HandleFunc(
		fmt.Sprintf("/nft/nfts/{%s}/{%s}/transfer", RestParamDenom, RestParamTokenID),
		transferNFTHandlerFn(cliCtx),
	).Methods("POST")

	// Burn an NFT
	r.HandleFunc(
		fmt.Sprintf("/nft/nfts/{%s}/{%s}/burn", RestParamDenom, RestParamTokenID),
		burnNFTHandlerFn(cliCtx),
	).Methods("POST")
}

func issueDenomHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req issueDenomReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		// create the message
		msg := types.NewMsgIssueDenom(req.ID, req.Name, req.Schema, req.Owner)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		authclient.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func mintNFTHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req mintNFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		if req.Recipient.Empty() {
			req.Recipient = req.Owner
		}
		// create the message
		msg := types.NewMsgMintNFT(
			req.ID,
			req.Denom,
			req.Name,
			req.URI,
			req.Data,
			req.Owner,
			req.Recipient,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		authclient.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func editNFTHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req editNFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		vars := mux.Vars(r)
		// create the message
		msg := types.NewMsgEditNFT(
			vars[RestParamTokenID],
			vars[RestParamDenom],
			req.Name,
			req.URI,
			req.Data, req.Owner,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		authclient.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func transferNFTHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req transferNFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		recipient, err := sdk.AccAddressFromBech32(req.Recipient)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		vars := mux.Vars(r)
		// create the message
		msg := types.NewMsgTransferNFT(
			vars[RestParamTokenID],
			vars[RestParamDenom],
			req.Name,
			req.URI,
			req.Data,
			req.Owner,
			recipient,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		authclient.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

func burnNFTHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req burnNFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		vars := mux.Vars(r)

		// create the message
		msg := types.NewMsgBurnNFT(
			req.Owner,
			vars[RestParamTokenID],
			vars[RestParamDenom],
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		authclient.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
