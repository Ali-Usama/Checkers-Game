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
		Index:       newIndex,
		Board:       newGame.String(),
		Turn:        rules.PieceStrings[newGame.Turn],
		Black:       msg.Black,
		Red:         msg.Red,
		MoveCount:   0,
		BeforeIndex: types.NoFifoIndex,
		AfterIndex:  types.NoFifoIndex,
		Deadline:    types.FormatDeadline(types.GetNextDeadline(ctx)),
		Winner:      rules.PieceStrings[rules.NO_PLAYER],
		Wager:       msg.Wager,
		Denom:       msg.Denom,
	}

	// Confirm that the values in the object are correct by checking the validity of the players' addresses
	err := storedGame.Validate()
	if err != nil {
		return nil, err
	}

	// Send the new game to the tail because it's freshly created
	k.Keeper.SendToFifoTail(ctx, &storedGame, &systemInfo)
	// Save the `storedGame` object
	k.Keeper.SetStoredGame(ctx, storedGame)

	// prepare the ground for the next game
	systemInfo.NextId++
	k.Keeper.SetSystemInfo(ctx, systemInfo)
	ctx.GasMeter().ConsumeGas(types.CreateGameGas, "Create Game")

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.GameCreatedEventType,
			sdk.NewAttribute(types.GameCreatedEventCreator, msg.Creator),
			sdk.NewAttribute(types.GameCreatedEventGameIndex, newIndex),
			sdk.NewAttribute(types.GameCreatedEventBlack, msg.Black),
			sdk.NewAttribute(types.GameCreatedEventRed, msg.Red),
			sdk.NewAttribute(types.GameCreatedEventWager, strconv.FormatUint(msg.Wager, 10)),
			sdk.NewAttribute(types.GameCreatedEventDenom, msg.Denom),
		),
	)

	// return the newly created ID for reference
	return &types.MsgCreateGameResponse{
		GameIndex: newIndex,
	}, nil
}
