package keeper

import (
	"context"
	"strconv"

	"github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateGame(goCtx context.Context, msg *types.MsgCreateGame) (*types.MsgCreateGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the new game's ID
	systemInfo, found := k.Keeper.GetSystemInfo(ctx)
	if !found {
		panic("SystemInfo not found")
	}
	newIndex := strconv.FormatUint(systemInfo.NextId, 10)

	// Create the object to be stored
	newGame := rules.New()
	storedGame := types.StoredGame{
		Index: newIndex,
		Board: newGame.String(),
		Turn:  rules.PieceStrings[newGame.Turn],
		Black: msg.Black,
		Red:   msg.Red,
	}

	// Confirm that the values in the object are correct by checking the validity of the players' addresses
	err := storedGame.Validate()
	if err != nil {
		return nil, err
	}

	// Save the `storedGame` object
	k.Keeper.SetStoredGame(ctx, storedGame)

	// prepare the ground for the next game
	systemInfo.NextId++
	k.Keeper.SetSystemInfo(ctx, systemInfo)

	// return the newly created ID for reference
	return &types.MsgCreateGameResponse{
		GameIndex: newIndex,
	}, nil
}
