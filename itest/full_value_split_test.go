package itest

import (
	"bytes"
	"context"

	"github.com/lightninglabs/taro/chanutils"
	"github.com/lightninglabs/taro/tarorpc"
	"github.com/lightningnetwork/lnd/lntest/wait"
	"github.com/stretchr/testify/require"
)

// testFullValueSplitSend tests that we can properly send the full value of a
// normal asset.
func testFullValueSend(t *harnessTest) {
	// First, we'll make an normal assets with enough units to allow us to
	// send it around a few times.
	rpcAssets := mintAssetsConfirmBatch(
		t, t.tarod, []*tarorpc.MintAssetRequest{simpleAssets[0]},
	)

	genInfo := rpcAssets[0].AssetGenesis
	genBootstrap := rpcAssets[0].AssetGenesis.GenesisBootstrapInfo

	ctxb := context.Background()

	// Now that we have the asset created, we'll make a new node that'll
	// serve as the node which'll receive the assets.
	secondTarod := setupTarodHarness(
		t.t, t, t.lndHarness.BackendCfg, t.lndHarness.Bob, t.universeServer,
	)
	defer func() {
		require.NoError(t.t, secondTarod.stop(true))
	}()

	// Next, we'll attempt to complete three transfers of the full value of
	// the asset between our main node and Bob.
	var (
		numSends     = 3
		fullAmount   = rpcAssets[0].Amount
		receiverAddr *tarorpc.Addr
		err          error
		sendResp     *tarorpc.SendAssetResponse
		importResp   *tarorpc.ImportProofResponse
	)

	for i := 0; i < numSends; i++ {
		// Create an address for the receiver and send the asset. We
		// start with Bob receiving the asset, then sending it back
		// to the main node, and so on.
		if i%2 == 0 {
			receiverAddr, err = secondTarod.NewAddr(
				ctxb, &tarorpc.NewAddrRequest{
					GenesisBootstrapInfo: genBootstrap,
					Amt:                  fullAmount,
				},
			)
			require.NoError(t.t, err)

			assertAddrCreated(
				t.t, secondTarod, rpcAssets[0], receiverAddr,
			)

			sendResp, err = sendAssetsToAddr(t.tarod, receiverAddr)
		} else {
			receiverAddr, err = t.tarod.NewAddr(
				ctxb, &tarorpc.NewAddrRequest{
					GenesisBootstrapInfo: genBootstrap,
					Amt:                  fullAmount,
				},
			)
			require.NoError(t.t, err)

			assertAddrCreated(
				t.t, t.tarod, rpcAssets[0], receiverAddr,
			)

			sendResp, err = sendAssetsToAddr(secondTarod, receiverAddr)
		}

		require.NoError(t.t, err)
		sendRespJSON, err := formatProtoJSON(sendResp)
		require.NoError(t.t, err)
		t.Logf("Got response from sending assets: %v", sendRespJSON)

		// Mine a block to force the send we created above to confirm.
		_ = mineBlocks(t, t.lndHarness, 1, len(rpcAssets))

		// Import the proof for the send on the receiving node.
		if i%2 == 0 {
			importResp, err = sendProof(
				t.tarod, secondTarod, receiverAddr, genInfo,
			)
		} else {
			importResp, err = sendProof(
				secondTarod, t.tarod, receiverAddr, genInfo,
			)
		}

		require.NoError(t.t, err)
		importRespJSON, err := formatProtoJSON(importResp)
		require.NoError(t.t, err)
		t.Logf("Got response from importing transfer proof: %v", importRespJSON)

		// TODO(jhb): update checks to include amounts
		err = wait.Predicate(func() bool {
			resp, err := t.tarod.ListTransfers(
				ctxb, &tarorpc.ListTransfersRequest{},
			)
			require.NoError(t.t, err)
			require.Len(t.t, resp.Transfers, i+1)

			sameAssetID := func(xfer *tarorpc.AssetTransfer) bool {
				return bytes.Equal(xfer.AssetSpendDeltas[0].AssetId,
					rpcAssets[0].AssetGenesis.AssetId)
			}

			return chanutils.All(resp.Transfers, sameAssetID)
		}, defaultTimeout/2)
		require.NoError(t.t, err)
	}
}
