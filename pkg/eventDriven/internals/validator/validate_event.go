package validator

import (
	"fmt"
	"strings"

	"github.com/triasbrata/golibs/pkg/eventDriven/internals/types"
)

func ValidateEvent(handler interface{}) error {
	if _, safe := handler.(types.ClientHandler); safe {
		return nil
	}
	if _, safe := handler.(types.ClientHandlerDataOnly); safe {
		return nil
	}
	if _, safe := handler.(types.ClientHandlerNoParams); safe {
		return nil
	}
	return fmt.Errorf("cant register handler, must apply with this types %v", strings.Join([]string{
		fmt.Sprintf("(%T)", types.ClientHandler(nil)),
		fmt.Sprintf("(%T)", types.ClientHandlerDataOnly(nil)),
		fmt.Sprintf("(%T)", types.ClientHandler(nil)),
	}, " or "))
}
