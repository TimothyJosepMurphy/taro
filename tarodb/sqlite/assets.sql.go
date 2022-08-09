// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: assets.sql

package sqlite

import (
	"context"
	"database/sql"
	"time"
)

const allAssets = `-- name: AllAssets :many
SELECT asset_id, version, script_key_id, asset_family_sig_id, script_version, amount, lock_time, relative_lock_time, split_commitment_root_hash, split_commitment_root_value, anchor_utxo_id 
FROM assets
`

func (q *Queries) AllAssets(ctx context.Context) ([]Asset, error) {
	rows, err := q.db.QueryContext(ctx, allAssets)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Asset
	for rows.Next() {
		var i Asset
		if err := rows.Scan(
			&i.AssetID,
			&i.Version,
			&i.ScriptKeyID,
			&i.AssetFamilySigID,
			&i.ScriptVersion,
			&i.Amount,
			&i.LockTime,
			&i.RelativeLockTime,
			&i.SplitCommitmentRootHash,
			&i.SplitCommitmentRootValue,
			&i.AnchorUtxoID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const allInternalKeys = `-- name: AllInternalKeys :many
SELECT key_id, raw_key, key_family, key_index 
FROM internal_keys
`

func (q *Queries) AllInternalKeys(ctx context.Context) ([]InternalKey, error) {
	rows, err := q.db.QueryContext(ctx, allInternalKeys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []InternalKey
	for rows.Next() {
		var i InternalKey
		if err := rows.Scan(
			&i.KeyID,
			&i.RawKey,
			&i.KeyFamily,
			&i.KeyIndex,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const allMintingBatches = `-- name: AllMintingBatches :many
SELECT batch_id, batch_state, minting_tx_psbt, minting_output_index, genesis_id, creation_time_unix, key_id, raw_key, key_family, key_index 
FROM asset_minting_batches
JOIN internal_keys 
ON asset_minting_batches.batch_id = internal_keys.key_id
`

type AllMintingBatchesRow struct {
	BatchID            int32
	BatchState         int16
	MintingTxPsbt      []byte
	MintingOutputIndex sql.NullInt16
	GenesisID          sql.NullInt32
	CreationTimeUnix   time.Time
	KeyID              int32
	RawKey             []byte
	KeyFamily          int32
	KeyIndex           int32
}

func (q *Queries) AllMintingBatches(ctx context.Context) ([]AllMintingBatchesRow, error) {
	rows, err := q.db.QueryContext(ctx, allMintingBatches)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AllMintingBatchesRow
	for rows.Next() {
		var i AllMintingBatchesRow
		if err := rows.Scan(
			&i.BatchID,
			&i.BatchState,
			&i.MintingTxPsbt,
			&i.MintingOutputIndex,
			&i.GenesisID,
			&i.CreationTimeUnix,
			&i.KeyID,
			&i.RawKey,
			&i.KeyFamily,
			&i.KeyIndex,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const anchorGenesisPoint = `-- name: AnchorGenesisPoint :exec
WITH target_point(genesis_id) AS (
    SELECT genesis_id
    FROM genesis_points
    WHERE genesis_points.prev_out = ?
)
UPDATE genesis_points
SET anchor_tx_id = ?
WHERE genesis_id in (SELECT genesis_id FROM target_point)
`

type AnchorGenesisPointParams struct {
	PrevOut    []byte
	AnchorTxID sql.NullInt32
}

func (q *Queries) AnchorGenesisPoint(ctx context.Context, arg AnchorGenesisPointParams) error {
	_, err := q.db.ExecContext(ctx, anchorGenesisPoint, arg.PrevOut, arg.AnchorTxID)
	return err
}

const anchorPendingAssets = `-- name: AnchorPendingAssets :exec
WITH assets_to_update AS (
    SELECT script_key_id
    FROM assets 
    JOIN genesis_assets 
        ON assets.asset_id = genesis_assets.gen_asset_id
    JOIN genesis_points
        ON genesis_points.genesis_id = genesis_assets.genesis_point_id
    WHERE prev_out = ?
)
UPDATE assets
SET anchor_utxo_id = ?
WHERE script_key_id in (SELECT script_key_id FROM assets_to_update)
`

type AnchorPendingAssetsParams struct {
	PrevOut      []byte
	AnchorUtxoID sql.NullInt32
}

func (q *Queries) AnchorPendingAssets(ctx context.Context, arg AnchorPendingAssetsParams) error {
	_, err := q.db.ExecContext(ctx, anchorPendingAssets, arg.PrevOut, arg.AnchorUtxoID)
	return err
}

const assetsByGenesisPoint = `-- name: AssetsByGenesisPoint :many
SELECT assets.asset_id, version, script_key_id, asset_family_sig_id, script_version, amount, lock_time, relative_lock_time, split_commitment_root_hash, split_commitment_root_value, anchor_utxo_id, gen_asset_id, genesis_assets.asset_id, asset_tag, meta_data, output_index, asset_type, genesis_point_id, genesis_id, prev_out, anchor_tx_id
FROM assets 
JOIN genesis_assets 
    ON assets.asset_id = genesis_assets.gen_asset_id
JOIN genesis_points
    ON genesis_points.genesis_id = genesis_assets.genesis_point_id
WHERE prev_out = ?
`

type AssetsByGenesisPointRow struct {
	AssetID                  int32
	Version                  int32
	ScriptKeyID              int32
	AssetFamilySigID         sql.NullInt32
	ScriptVersion            int32
	Amount                   int64
	LockTime                 sql.NullInt32
	RelativeLockTime         sql.NullInt32
	SplitCommitmentRootHash  []byte
	SplitCommitmentRootValue sql.NullInt64
	AnchorUtxoID             sql.NullInt32
	GenAssetID               int32
	AssetID_2                []byte
	AssetTag                 string
	MetaData                 []byte
	OutputIndex              int32
	AssetType                int16
	GenesisPointID           int32
	GenesisID                int32
	PrevOut                  []byte
	AnchorTxID               sql.NullInt32
}

func (q *Queries) AssetsByGenesisPoint(ctx context.Context, prevOut []byte) ([]AssetsByGenesisPointRow, error) {
	rows, err := q.db.QueryContext(ctx, assetsByGenesisPoint, prevOut)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AssetsByGenesisPointRow
	for rows.Next() {
		var i AssetsByGenesisPointRow
		if err := rows.Scan(
			&i.AssetID,
			&i.Version,
			&i.ScriptKeyID,
			&i.AssetFamilySigID,
			&i.ScriptVersion,
			&i.Amount,
			&i.LockTime,
			&i.RelativeLockTime,
			&i.SplitCommitmentRootHash,
			&i.SplitCommitmentRootValue,
			&i.AnchorUtxoID,
			&i.GenAssetID,
			&i.AssetID_2,
			&i.AssetTag,
			&i.MetaData,
			&i.OutputIndex,
			&i.AssetType,
			&i.GenesisPointID,
			&i.GenesisID,
			&i.PrevOut,
			&i.AnchorTxID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const assetsInBatch = `-- name: AssetsInBatch :many
SELECT
    gen_asset_id, asset_id, asset_tag, meta_data, output_index, asset_type,
    genesis_points.prev_out prev_out
FROM genesis_assets
JOIN genesis_points
    ON genesis_assets.genesis_point_id = genesis_points.genesis_id
JOIN asset_minting_batches batches
    ON genesis_points.genesis_id = batches.genesis_id
JOIN internal_keys keys
    ON keys.key_id = batches.batch_id
WHERE keys.raw_key = ?
`

type AssetsInBatchRow struct {
	GenAssetID  int32
	AssetID     []byte
	AssetTag    string
	MetaData    []byte
	OutputIndex int32
	AssetType   int16
	PrevOut     []byte
}

func (q *Queries) AssetsInBatch(ctx context.Context, rawKey []byte) ([]AssetsInBatchRow, error) {
	rows, err := q.db.QueryContext(ctx, assetsInBatch, rawKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AssetsInBatchRow
	for rows.Next() {
		var i AssetsInBatchRow
		if err := rows.Scan(
			&i.GenAssetID,
			&i.AssetID,
			&i.AssetTag,
			&i.MetaData,
			&i.OutputIndex,
			&i.AssetType,
			&i.PrevOut,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const bindMintingBatchWithTx = `-- name: BindMintingBatchWithTx :exec
WITH target_batch AS (
    SELECT batch_id
    FROM asset_minting_batches batches
    JOIN internal_keys keys
        ON batches.batch_id = keys.key_id
    WHERE keys.raw_key = ?
)
UPDATE asset_minting_batches 
SET minting_tx_psbt = ?, minting_output_index = ?, genesis_id = ?
WHERE batch_id IN (SELECT batch_id FROM target_batch)
`

type BindMintingBatchWithTxParams struct {
	RawKey             []byte
	MintingTxPsbt      []byte
	MintingOutputIndex sql.NullInt16
	GenesisID          sql.NullInt32
}

func (q *Queries) BindMintingBatchWithTx(ctx context.Context, arg BindMintingBatchWithTxParams) error {
	_, err := q.db.ExecContext(ctx, bindMintingBatchWithTx,
		arg.RawKey,
		arg.MintingTxPsbt,
		arg.MintingOutputIndex,
		arg.GenesisID,
	)
	return err
}

const confirmChainTx = `-- name: ConfirmChainTx :exec
WITH target_txn(txn_id) AS (
    SELECT anchor_tx_id
    FROM genesis_points points
    JOIN asset_minting_batches batches
        ON batches.genesis_id = points.genesis_id
    JOIN internal_keys keys
        ON batches.batch_id = keys.key_id
    WHERE keys.raw_key = ?
)
UPDATE chain_txns
SET block_height = ?, block_hash = ?, tx_index = ?
WHERE txn_id in (SELECT txn_id FROm target_txn)
`

type ConfirmChainTxParams struct {
	RawKey      []byte
	BlockHeight sql.NullInt32
	BlockHash   []byte
	TxIndex     sql.NullInt32
}

func (q *Queries) ConfirmChainTx(ctx context.Context, arg ConfirmChainTxParams) error {
	_, err := q.db.ExecContext(ctx, confirmChainTx,
		arg.RawKey,
		arg.BlockHeight,
		arg.BlockHash,
		arg.TxIndex,
	)
	return err
}

const fetchAllAssets = `-- name: FetchAllAssets :many
WITH genesis_info AS (
    -- This CTE is used to fetch the base asset information from disk based on
    -- the raw key of the batch that will ultimately create this set of assets.
    -- To do so, we'll need to traverse a few tables to join the set of assets
    -- with the genesis points, then with the batches that reference this
    -- points, to the internal key that reference the batch, then restricted
    -- for internal keys that match our main batch key.
    SELECT
        gen_asset_id, asset_id, asset_tag, meta_data, output_index, asset_type,
        genesis_points.prev_out prev_out
    FROM genesis_assets
    JOIN genesis_points
        ON genesis_assets.genesis_point_id = genesis_points.genesis_id
), key_fam_info AS (
    -- This CTE is used to perform a series of joins that allow us to extract
    -- the family key information, as well as the family sigs for the series of
    -- assets we care about. We obtain only the assets found in the batch
    -- above, with the WHERE query at the bottom.
    SELECT 
        sig_id, gen_asset_id, genesis_sig, tweaked_fam_key, raw_key, key_index, key_family
    FROM asset_family_sigs sigs
    JOIN asset_families fams
        ON sigs.key_fam_id = fams.family_id
    JOIN internal_keys keys
        ON keys.key_id = fams.internal_key_id
    -- TODO(roasbeef): or can join do this below?
    WHERE sigs.gen_asset_id IN (SELECT gen_asset_id FROM genesis_info)
)
SELECT 
    version, internal_keys.raw_key AS script_key_raw, 
    internal_keys.key_family AS script_key_fam,
    internal_keys.key_index AS script_key_index, key_fam_info.genesis_sig, 
    key_fam_info.tweaked_fam_key, key_fam_info.raw_key AS fam_key_raw,
    key_fam_info.key_family AS fam_key_family, key_fam_info.key_index AS fam_key_index,
    script_version, amount, lock_time, relative_lock_time, 
    genesis_info.asset_id, genesis_info.asset_tag, genesis_info.meta_data, 
    genesis_info.output_index AS genesis_output_index, genesis_info.asset_type,
    genesis_info.prev_out AS genesis_prev_out,
    txns.raw_tx AS anchor_tx, txns.txid AS anchor_txid, txns.block_hash AS anchor_block_hash,
    utxos.outpoint AS anchor_outpoint
FROM assets
JOIN genesis_info
    ON assets.asset_id = genesis_info.gen_asset_id
LEFT JOIN key_fam_info
    ON assets.asset_id = key_fam_info.gen_asset_id
JOIN internal_keys
    ON assets.script_key_id = internal_keys.key_id
JOIN managed_utxos utxos
    ON assets.anchor_utxo_id = utxos.utxo_id
JOIN chain_txns txns
    ON utxos.txn_id = txns.txn_id
`

type FetchAllAssetsRow struct {
	Version            int32
	ScriptKeyRaw       []byte
	ScriptKeyFam       int32
	ScriptKeyIndex     int32
	GenesisSig         []byte
	TweakedFamKey      []byte
	FamKeyRaw          []byte
	FamKeyFamily       sql.NullInt32
	FamKeyIndex        sql.NullInt32
	ScriptVersion      int32
	Amount             int64
	LockTime           sql.NullInt32
	RelativeLockTime   sql.NullInt32
	AssetID            []byte
	AssetTag           string
	MetaData           []byte
	GenesisOutputIndex int32
	AssetType          int16
	GenesisPrevOut     []byte
	AnchorTx           []byte
	AnchorTxid         []byte
	AnchorBlockHash    []byte
	AnchorOutpoint     []byte
}

// TODO(roasbeef): identical to the above but no batch, how to combine?
// We use a LEFT JOIN here as not every asset has a family key, so this'll
// generate rows that have NULL values for the faily key fields if an asset
// doesn't have a family key. See the comment in fetchAssetSprouts for a work
// around that needs to be used with this query until a sqlc bug is fixed.
func (q *Queries) FetchAllAssets(ctx context.Context) ([]FetchAllAssetsRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchAllAssets)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchAllAssetsRow
	for rows.Next() {
		var i FetchAllAssetsRow
		if err := rows.Scan(
			&i.Version,
			&i.ScriptKeyRaw,
			&i.ScriptKeyFam,
			&i.ScriptKeyIndex,
			&i.GenesisSig,
			&i.TweakedFamKey,
			&i.FamKeyRaw,
			&i.FamKeyFamily,
			&i.FamKeyIndex,
			&i.ScriptVersion,
			&i.Amount,
			&i.LockTime,
			&i.RelativeLockTime,
			&i.AssetID,
			&i.AssetTag,
			&i.MetaData,
			&i.GenesisOutputIndex,
			&i.AssetType,
			&i.GenesisPrevOut,
			&i.AnchorTx,
			&i.AnchorTxid,
			&i.AnchorBlockHash,
			&i.AnchorOutpoint,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchAssetProof = `-- name: FetchAssetProof :one
WITH asset_info AS (
    SELECT assets.asset_id, keys.raw_key
    FROM assets
    JOIN internal_keys keys
        ON keys.key_id = assets.script_key_id
    WHERE keys.raw_key = ?
)
SELECT asset_info.raw_key AS script_key, asset_proofs.proof_file
FROM asset_proofs
JOIN asset_info
    ON asset_info.asset_id = asset_proofs.asset_id
`

type FetchAssetProofRow struct {
	ScriptKey []byte
	ProofFile []byte
}

func (q *Queries) FetchAssetProof(ctx context.Context, rawKey []byte) (FetchAssetProofRow, error) {
	row := q.db.QueryRowContext(ctx, fetchAssetProof, rawKey)
	var i FetchAssetProofRow
	err := row.Scan(&i.ScriptKey, &i.ProofFile)
	return i, err
}

const fetchAssetProofs = `-- name: FetchAssetProofs :many
WITH asset_info AS (
    SELECT assets.asset_id, keys.raw_key
    FROM assets
    JOIN internal_keys keys
        ON keys.key_id = assets.script_key_id
)
SELECT asset_info.raw_key AS script_key, asset_proofs.proof_file
FROM asset_proofs
JOIN asset_info
    ON asset_info.asset_id = asset_proofs.asset_id
`

type FetchAssetProofsRow struct {
	ScriptKey []byte
	ProofFile []byte
}

func (q *Queries) FetchAssetProofs(ctx context.Context) ([]FetchAssetProofsRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchAssetProofs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchAssetProofsRow
	for rows.Next() {
		var i FetchAssetProofsRow
		if err := rows.Scan(&i.ScriptKey, &i.ProofFile); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchAssetsByAnchorTx = `-- name: FetchAssetsByAnchorTx :many
SELECT asset_id, version, script_key_id, asset_family_sig_id, script_version, amount, lock_time, relative_lock_time, split_commitment_root_hash, split_commitment_root_value, anchor_utxo_id
FROM assets
WHERE anchor_utxo_id = ?
`

func (q *Queries) FetchAssetsByAnchorTx(ctx context.Context, anchorUtxoID sql.NullInt32) ([]Asset, error) {
	rows, err := q.db.QueryContext(ctx, fetchAssetsByAnchorTx, anchorUtxoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Asset
	for rows.Next() {
		var i Asset
		if err := rows.Scan(
			&i.AssetID,
			&i.Version,
			&i.ScriptKeyID,
			&i.AssetFamilySigID,
			&i.ScriptVersion,
			&i.Amount,
			&i.LockTime,
			&i.RelativeLockTime,
			&i.SplitCommitmentRootHash,
			&i.SplitCommitmentRootValue,
			&i.AnchorUtxoID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchAssetsForBatch = `-- name: FetchAssetsForBatch :many
WITH genesis_info AS (
    -- This CTE is used to fetch the base asset information from disk based on
    -- the raw key of the batch that will ultimately create this set of assets.
    -- To do so, we'll need to traverse a few tables to join the set of assets
    -- with the genesis points, then with the batches that reference this
    -- points, to the internal key that reference the batch, then restricted
    -- for internal keys that match our main batch key.
    SELECT
        gen_asset_id, asset_id, asset_tag, meta_data, output_index, asset_type,
        genesis_points.prev_out prev_out
    FROM genesis_assets
    JOIN genesis_points
        ON genesis_assets.genesis_point_id = genesis_points.genesis_id
    JOIN asset_minting_batches batches
        ON genesis_points.genesis_id = batches.genesis_id
    JOIN internal_keys keys
        ON keys.key_id = batches.batch_id
    WHERE keys.raw_key = ?
), key_fam_info AS (
    -- This CTE is used to perform a series of joins that allow us to extract
    -- the family key information, as well as the family sigs for the series of
    -- assets we care about. We obtain only the assets found in the batch
    -- above, with the WHERE query at the bottom.
    SELECT 
        sig_id, gen_asset_id, genesis_sig, tweaked_fam_key, raw_key, key_index, key_family
    FROM asset_family_sigs sigs
    JOIN asset_families fams
        ON sigs.key_fam_id = fams.family_id
    JOIN internal_keys keys
        ON keys.key_id = fams.internal_key_id
    -- TODO(roasbeef): or can join do this below?
    WHERE sigs.gen_asset_id IN (SELECT gen_asset_id FROM genesis_info)
)
SELECT 
    version, internal_keys.raw_key AS script_key_raw, 
    internal_keys.key_family AS script_key_fam,
    internal_keys.key_index AS script_key_index, key_fam_info.genesis_sig, 
    key_fam_info.tweaked_fam_key, key_fam_info.raw_key AS fam_key_raw,
    key_fam_info.key_family AS fam_key_family, key_fam_info.key_index AS fam_key_index,
    script_version, amount, lock_time, relative_lock_time, 
    genesis_info.asset_id, genesis_info.asset_tag, genesis_info.meta_data, 
    genesis_info.output_index AS genesis_output_index, genesis_info.asset_type,
    genesis_info.prev_out AS genesis_prev_out
FROM assets
JOIN genesis_info
    ON assets.asset_id = genesis_info.gen_asset_id
LEFT JOIN key_fam_info
    ON assets.asset_id = key_fam_info.gen_asset_id
JOIN internal_keys
    ON assets.script_key_id = internal_keys.key_id
`

type FetchAssetsForBatchRow struct {
	Version            int32
	ScriptKeyRaw       []byte
	ScriptKeyFam       int32
	ScriptKeyIndex     int32
	GenesisSig         []byte
	TweakedFamKey      []byte
	FamKeyRaw          []byte
	FamKeyFamily       sql.NullInt32
	FamKeyIndex        sql.NullInt32
	ScriptVersion      int32
	Amount             int64
	LockTime           sql.NullInt32
	RelativeLockTime   sql.NullInt32
	AssetID            []byte
	AssetTag           string
	MetaData           []byte
	GenesisOutputIndex int32
	AssetType          int16
	GenesisPrevOut     []byte
}

// We use a LEFT JOIN here as not every asset has a family key, so this'll
// generate rows that have NULL values for the faily key fields if an asset
// doesn't have a family key. See the comment in fetchAssetSprouts for a work
// around that needs to be used with this query until a sqlc bug is fixed.
func (q *Queries) FetchAssetsForBatch(ctx context.Context, rawKey []byte) ([]FetchAssetsForBatchRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchAssetsForBatch, rawKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchAssetsForBatchRow
	for rows.Next() {
		var i FetchAssetsForBatchRow
		if err := rows.Scan(
			&i.Version,
			&i.ScriptKeyRaw,
			&i.ScriptKeyFam,
			&i.ScriptKeyIndex,
			&i.GenesisSig,
			&i.TweakedFamKey,
			&i.FamKeyRaw,
			&i.FamKeyFamily,
			&i.FamKeyIndex,
			&i.ScriptVersion,
			&i.Amount,
			&i.LockTime,
			&i.RelativeLockTime,
			&i.AssetID,
			&i.AssetTag,
			&i.MetaData,
			&i.GenesisOutputIndex,
			&i.AssetType,
			&i.GenesisPrevOut,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchChainTx = `-- name: FetchChainTx :one
SELECT txn_id, txid, raw_tx, block_height, block_hash, tx_index
FROM chain_txns
WHERE txid = ?
`

func (q *Queries) FetchChainTx(ctx context.Context, txid []byte) (ChainTxn, error) {
	row := q.db.QueryRowContext(ctx, fetchChainTx, txid)
	var i ChainTxn
	err := row.Scan(
		&i.TxnID,
		&i.Txid,
		&i.RawTx,
		&i.BlockHeight,
		&i.BlockHash,
		&i.TxIndex,
	)
	return i, err
}

const fetchGenesisPointByAnchorTx = `-- name: FetchGenesisPointByAnchorTx :one
SELECT genesis_id, prev_out, anchor_tx_id 
FROM genesis_points
WHERE anchor_tx_id = ?
`

func (q *Queries) FetchGenesisPointByAnchorTx(ctx context.Context, anchorTxID sql.NullInt32) (GenesisPoint, error) {
	row := q.db.QueryRowContext(ctx, fetchGenesisPointByAnchorTx, anchorTxID)
	var i GenesisPoint
	err := row.Scan(&i.GenesisID, &i.PrevOut, &i.AnchorTxID)
	return i, err
}

const fetchManagedUTXO = `-- name: FetchManagedUTXO :one
SELECT utxo_id, outpoint, amt_sats, internal_key_id, tapscript_sibling, taro_root, txn_id
from managed_utxos
WHERE txn_id = ?
`

func (q *Queries) FetchManagedUTXO(ctx context.Context, txnID int32) (ManagedUtxo, error) {
	row := q.db.QueryRowContext(ctx, fetchManagedUTXO, txnID)
	var i ManagedUtxo
	err := row.Scan(
		&i.UtxoID,
		&i.Outpoint,
		&i.AmtSats,
		&i.InternalKeyID,
		&i.TapscriptSibling,
		&i.TaroRoot,
		&i.TxnID,
	)
	return i, err
}

const fetchMintingBatch = `-- name: FetchMintingBatch :one
SELECT batch_id, batch_state, minting_tx_psbt, minting_output_index, genesis_id, creation_time_unix, key_id, raw_key, key_family, key_index
FROM asset_minting_batches batches
JOIN internal_keys keys
    ON batches.batch_id = keys.key_id
WHERE keys.raw_key = ?
`

type FetchMintingBatchRow struct {
	BatchID            int32
	BatchState         int16
	MintingTxPsbt      []byte
	MintingOutputIndex sql.NullInt16
	GenesisID          sql.NullInt32
	CreationTimeUnix   time.Time
	KeyID              int32
	RawKey             []byte
	KeyFamily          int32
	KeyIndex           int32
}

func (q *Queries) FetchMintingBatch(ctx context.Context, rawKey []byte) (FetchMintingBatchRow, error) {
	row := q.db.QueryRowContext(ctx, fetchMintingBatch, rawKey)
	var i FetchMintingBatchRow
	err := row.Scan(
		&i.BatchID,
		&i.BatchState,
		&i.MintingTxPsbt,
		&i.MintingOutputIndex,
		&i.GenesisID,
		&i.CreationTimeUnix,
		&i.KeyID,
		&i.RawKey,
		&i.KeyFamily,
		&i.KeyIndex,
	)
	return i, err
}

const fetchMintingBatchesByInverseState = `-- name: FetchMintingBatchesByInverseState :many
SELECT batch_id, batch_state, minting_tx_psbt, minting_output_index, genesis_id, creation_time_unix, key_id, raw_key, key_family, key_index
FROM asset_minting_batches batches
JOIN internal_keys keys
    ON batches.batch_id = keys.key_id
WHERE batches.batch_state != ?
`

type FetchMintingBatchesByInverseStateRow struct {
	BatchID            int32
	BatchState         int16
	MintingTxPsbt      []byte
	MintingOutputIndex sql.NullInt16
	GenesisID          sql.NullInt32
	CreationTimeUnix   time.Time
	KeyID              int32
	RawKey             []byte
	KeyFamily          int32
	KeyIndex           int32
}

func (q *Queries) FetchMintingBatchesByInverseState(ctx context.Context, batchState int16) ([]FetchMintingBatchesByInverseStateRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchMintingBatchesByInverseState, batchState)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchMintingBatchesByInverseStateRow
	for rows.Next() {
		var i FetchMintingBatchesByInverseStateRow
		if err := rows.Scan(
			&i.BatchID,
			&i.BatchState,
			&i.MintingTxPsbt,
			&i.MintingOutputIndex,
			&i.GenesisID,
			&i.CreationTimeUnix,
			&i.KeyID,
			&i.RawKey,
			&i.KeyFamily,
			&i.KeyIndex,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchMintingBatchesByState = `-- name: FetchMintingBatchesByState :many
SELECT batch_id, batch_state, minting_tx_psbt, minting_output_index, genesis_id, creation_time_unix, key_id, raw_key, key_family, key_index
FROM asset_minting_batches batches
JOIN internal_keys keys
    ON batches.batch_id = keys.key_id
WHERE batches.batch_state = ?
`

type FetchMintingBatchesByStateRow struct {
	BatchID            int32
	BatchState         int16
	MintingTxPsbt      []byte
	MintingOutputIndex sql.NullInt16
	GenesisID          sql.NullInt32
	CreationTimeUnix   time.Time
	KeyID              int32
	RawKey             []byte
	KeyFamily          int32
	KeyIndex           int32
}

func (q *Queries) FetchMintingBatchesByState(ctx context.Context, batchState int16) ([]FetchMintingBatchesByStateRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchMintingBatchesByState, batchState)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchMintingBatchesByStateRow
	for rows.Next() {
		var i FetchMintingBatchesByStateRow
		if err := rows.Scan(
			&i.BatchID,
			&i.BatchState,
			&i.MintingTxPsbt,
			&i.MintingOutputIndex,
			&i.GenesisID,
			&i.CreationTimeUnix,
			&i.KeyID,
			&i.RawKey,
			&i.KeyFamily,
			&i.KeyIndex,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchSeedlingsForBatch = `-- name: FetchSeedlingsForBatch :many
WITH target_batch(batch_id) AS (
    SELECT batch_id
    FROM asset_minting_batches batches
    JOIN internal_keys keys
        ON batches.batch_id = keys.key_id
    WHERE keys.raw_key = ?
)
SELECT seedling_id, asset_name, asset_type, asset_supply, asset_meta,
    emission_enabled, asset_id, batch_id
FROM asset_seedlings 
WHERE asset_seedlings.batch_id in (SELECT batch_id FROM target_batch)
`

func (q *Queries) FetchSeedlingsForBatch(ctx context.Context, rawKey []byte) ([]AssetSeedling, error) {
	rows, err := q.db.QueryContext(ctx, fetchSeedlingsForBatch, rawKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AssetSeedling
	for rows.Next() {
		var i AssetSeedling
		if err := rows.Scan(
			&i.SeedlingID,
			&i.AssetName,
			&i.AssetType,
			&i.AssetSupply,
			&i.AssetMeta,
			&i.EmissionEnabled,
			&i.AssetID,
			&i.BatchID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const genesisAssets = `-- name: GenesisAssets :many
SELECT gen_asset_id, asset_id, asset_tag, meta_data, output_index, asset_type, genesis_point_id 
FROM genesis_assets
`

func (q *Queries) GenesisAssets(ctx context.Context) ([]GenesisAsset, error) {
	rows, err := q.db.QueryContext(ctx, genesisAssets)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GenesisAsset
	for rows.Next() {
		var i GenesisAsset
		if err := rows.Scan(
			&i.GenAssetID,
			&i.AssetID,
			&i.AssetTag,
			&i.MetaData,
			&i.OutputIndex,
			&i.AssetType,
			&i.GenesisPointID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const genesisPoints = `-- name: GenesisPoints :many
SELECT genesis_id, prev_out, anchor_tx_id 
FROM genesis_points
`

func (q *Queries) GenesisPoints(ctx context.Context) ([]GenesisPoint, error) {
	rows, err := q.db.QueryContext(ctx, genesisPoints)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GenesisPoint
	for rows.Next() {
		var i GenesisPoint
		if err := rows.Scan(&i.GenesisID, &i.PrevOut, &i.AnchorTxID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertAssetFamilyKey = `-- name: InsertAssetFamilyKey :one
INSERT INTO asset_families (
    tweaked_fam_key, internal_key_id, genesis_point_id 
) VALUES (
    ?, ?, ?
) ON CONFLICT 
    DO UPDATE SET genesis_point_id = EXCLUDED.genesis_point_id
RETURNING family_id
`

type InsertAssetFamilyKeyParams struct {
	TweakedFamKey  []byte
	InternalKeyID  int32
	GenesisPointID int32
}

func (q *Queries) InsertAssetFamilyKey(ctx context.Context, arg InsertAssetFamilyKeyParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, insertAssetFamilyKey, arg.TweakedFamKey, arg.InternalKeyID, arg.GenesisPointID)
	var family_id int32
	err := row.Scan(&family_id)
	return family_id, err
}

const insertAssetFamilySig = `-- name: InsertAssetFamilySig :one
INSERT INTO asset_family_sigs (
    genesis_sig, gen_asset_id, key_fam_id
) VALUES (
    ?, ?, ?
) RETURNING sig_id
`

type InsertAssetFamilySigParams struct {
	GenesisSig []byte
	GenAssetID int32
	KeyFamID   int32
}

func (q *Queries) InsertAssetFamilySig(ctx context.Context, arg InsertAssetFamilySigParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, insertAssetFamilySig, arg.GenesisSig, arg.GenAssetID, arg.KeyFamID)
	var sig_id int32
	err := row.Scan(&sig_id)
	return sig_id, err
}

const insertAssetSeedling = `-- name: InsertAssetSeedling :exec
INSERT INTO asset_seedlings (
    asset_name, asset_type, asset_supply, asset_meta,
    emission_enabled, batch_id
) VALUES (
    ?, ?, ?, ?, ?, ?
)
`

type InsertAssetSeedlingParams struct {
	AssetName       string
	AssetType       int16
	AssetSupply     int64
	AssetMeta       []byte
	EmissionEnabled bool
	BatchID         int32
}

func (q *Queries) InsertAssetSeedling(ctx context.Context, arg InsertAssetSeedlingParams) error {
	_, err := q.db.ExecContext(ctx, insertAssetSeedling,
		arg.AssetName,
		arg.AssetType,
		arg.AssetSupply,
		arg.AssetMeta,
		arg.EmissionEnabled,
		arg.BatchID,
	)
	return err
}

const insertAssetSeedlingIntoBatch = `-- name: InsertAssetSeedlingIntoBatch :exec
WITH target_key_id AS (
    -- We use this CTE to fetch the key_id of the internal key that's
    -- associated with a given batch. This can only return one value in
    -- practice since raw_key is a unique field. We then use this value below
    -- to insert the seedling and point to the proper batch_id, which is a
    -- foreign key that references the key_id of the internal key.
    SELECT key_id 
    FROM internal_keys keys
    WHERE keys.raw_key = ?
)
INSERT INTO asset_seedlings(
    asset_name, asset_type, asset_supply, asset_meta,
    emission_enabled, batch_id
) VALUES (
    ?, ?, ?, ?, ?, (SELECT key_id FROM target_key_id)
)
`

type InsertAssetSeedlingIntoBatchParams struct {
	RawKey          []byte
	AssetName       string
	AssetType       int16
	AssetSupply     int64
	AssetMeta       []byte
	EmissionEnabled bool
}

func (q *Queries) InsertAssetSeedlingIntoBatch(ctx context.Context, arg InsertAssetSeedlingIntoBatchParams) error {
	_, err := q.db.ExecContext(ctx, insertAssetSeedlingIntoBatch,
		arg.RawKey,
		arg.AssetName,
		arg.AssetType,
		arg.AssetSupply,
		arg.AssetMeta,
		arg.EmissionEnabled,
	)
	return err
}

const insertChainTx = `-- name: InsertChainTx :one
INSERT INTO chain_txns (
    txid, raw_tx
) VALUES (
    ?, ?
)
RETURNING txn_id
`

type InsertChainTxParams struct {
	Txid  []byte
	RawTx []byte
}

func (q *Queries) InsertChainTx(ctx context.Context, arg InsertChainTxParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, insertChainTx, arg.Txid, arg.RawTx)
	var txn_id int32
	err := row.Scan(&txn_id)
	return txn_id, err
}

const insertGenesisAsset = `-- name: InsertGenesisAsset :one
INSERT INTO genesis_assets (
    asset_id, asset_tag, meta_data, output_index, asset_type, genesis_point_id
) VALUES (
    ?, ?, ?, ?, ?, ?
) RETURNING gen_asset_id
`

type InsertGenesisAssetParams struct {
	AssetID        []byte
	AssetTag       string
	MetaData       []byte
	OutputIndex    int32
	AssetType      int16
	GenesisPointID int32
}

func (q *Queries) InsertGenesisAsset(ctx context.Context, arg InsertGenesisAssetParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, insertGenesisAsset,
		arg.AssetID,
		arg.AssetTag,
		arg.MetaData,
		arg.OutputIndex,
		arg.AssetType,
		arg.GenesisPointID,
	)
	var gen_asset_id int32
	err := row.Scan(&gen_asset_id)
	return gen_asset_id, err
}

const insertGenesisPoint = `-- name: InsertGenesisPoint :one
INSERT INTO genesis_points(
    prev_out
) VALUES (
    ?
) RETURNING genesis_id
`

func (q *Queries) InsertGenesisPoint(ctx context.Context, prevOut []byte) (int32, error) {
	row := q.db.QueryRowContext(ctx, insertGenesisPoint, prevOut)
	var genesis_id int32
	err := row.Scan(&genesis_id)
	return genesis_id, err
}

const insertInternalKey = `-- name: InsertInternalKey :one
INSERT INTO internal_keys (
    raw_key, key_family, key_index
) VALUES (?, ?, ?) RETURNING key_id
`

type InsertInternalKeyParams struct {
	RawKey    []byte
	KeyFamily int32
	KeyIndex  int32
}

func (q *Queries) InsertInternalKey(ctx context.Context, arg InsertInternalKeyParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, insertInternalKey, arg.RawKey, arg.KeyFamily, arg.KeyIndex)
	var key_id int32
	err := row.Scan(&key_id)
	return key_id, err
}

const insertManagedUTXO = `-- name: InsertManagedUTXO :one
WITH target_key(key_id) AS (
    SELECT key_id
    FROM internal_keys
    WHERE raw_key = ?
)
INSERT INTO managed_utxos (
    outpoint, amt_sats, internal_key_id, tapscript_sibling, taro_root, txn_id
) VALUES (
    ?, ?, (SELECT key_id FROM target_key), ?, ?, ?
) RETURNING utxo_id
`

type InsertManagedUTXOParams struct {
	RawKey           []byte
	Outpoint         []byte
	AmtSats          int64
	TapscriptSibling []byte
	TaroRoot         []byte
	TxnID            int32
}

func (q *Queries) InsertManagedUTXO(ctx context.Context, arg InsertManagedUTXOParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, insertManagedUTXO,
		arg.RawKey,
		arg.Outpoint,
		arg.AmtSats,
		arg.TapscriptSibling,
		arg.TaroRoot,
		arg.TxnID,
	)
	var utxo_id int32
	err := row.Scan(&utxo_id)
	return utxo_id, err
}

const insertNewAsset = `-- name: InsertNewAsset :exec
INSERT INTO assets (
    version, script_key_id, asset_id, asset_family_sig_id, script_version, 
    amount, lock_time, relative_lock_time
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?
)
`

type InsertNewAssetParams struct {
	Version          int32
	ScriptKeyID      int32
	AssetID          int32
	AssetFamilySigID sql.NullInt32
	ScriptVersion    int32
	Amount           int64
	LockTime         sql.NullInt32
	RelativeLockTime sql.NullInt32
}

func (q *Queries) InsertNewAsset(ctx context.Context, arg InsertNewAssetParams) error {
	_, err := q.db.ExecContext(ctx, insertNewAsset,
		arg.Version,
		arg.ScriptKeyID,
		arg.AssetID,
		arg.AssetFamilySigID,
		arg.ScriptVersion,
		arg.Amount,
		arg.LockTime,
		arg.RelativeLockTime,
	)
	return err
}

const newMintingBatch = `-- name: NewMintingBatch :exec
INSERT INTO asset_minting_batches (
    batch_state, batch_id, creation_time_unix
) VALUES (0, ?, ?)
`

type NewMintingBatchParams struct {
	BatchID          int32
	CreationTimeUnix time.Time
}

func (q *Queries) NewMintingBatch(ctx context.Context, arg NewMintingBatchParams) error {
	_, err := q.db.ExecContext(ctx, newMintingBatch, arg.BatchID, arg.CreationTimeUnix)
	return err
}

const updateAssetProof = `-- name: UpdateAssetProof :exec
WITH target_asset(asset_id) AS (
    SELECT asset_id
    FROM assets
    JOIN internal_keys keys
        ON keys.key_id = assets.script_key_id
    WHERE keys.raw_key = ?
)
INSERT INTO asset_proofs (
    asset_id, proof_file
) VALUES (
    (SELECT asset_id FROM target_asset), ?
) ON CONFLICT 
    DO UPDATE SET proof_file = EXCLUDED.proof_file
`

type UpdateAssetProofParams struct {
	RawKey    []byte
	ProofFile []byte
}

func (q *Queries) UpdateAssetProof(ctx context.Context, arg UpdateAssetProofParams) error {
	_, err := q.db.ExecContext(ctx, updateAssetProof, arg.RawKey, arg.ProofFile)
	return err
}

const updateBatchGenesisTx = `-- name: UpdateBatchGenesisTx :exec
WITH target_batch AS (
    SELECT batch_id
    FROM asset_minting_batches batches
    JOIN internal_keys keys
        ON batches.batch_id = keys.key_id
    WHERE keys.raw_key = ?
)
UPDATE asset_minting_batches
SET minting_tx_psbt = ?
WHERE batch_id in (SELECT batch_id FROM target_batch)
`

type UpdateBatchGenesisTxParams struct {
	RawKey        []byte
	MintingTxPsbt []byte
}

func (q *Queries) UpdateBatchGenesisTx(ctx context.Context, arg UpdateBatchGenesisTxParams) error {
	_, err := q.db.ExecContext(ctx, updateBatchGenesisTx, arg.RawKey, arg.MintingTxPsbt)
	return err
}

const updateMintingBatchState = `-- name: UpdateMintingBatchState :exec
WITH target_batch AS (
    -- This CTE is used to fetch the ID of a batch, based on the serialized
    -- internal key associated with the batch. This internal key is as the
    -- actual Taproot internal key to ultimately mint the batch. This pattern
    -- is used in several other queries.
    SELECT batch_id
    FROM asset_minting_batches batches
    JOIN internal_keys keys
        ON batches.batch_id = keys.key_id
    WHERE keys.raw_key = ?
)
UPDATE asset_minting_batches 
SET batch_state = ? 
WHERE batch_id in (SELECT batch_id FROM target_batch)
`

type UpdateMintingBatchStateParams struct {
	RawKey     []byte
	BatchState int16
}

func (q *Queries) UpdateMintingBatchState(ctx context.Context, arg UpdateMintingBatchStateParams) error {
	_, err := q.db.ExecContext(ctx, updateMintingBatchState, arg.RawKey, arg.BatchState)
	return err
}