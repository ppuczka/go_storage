package frontend

import (
	"fmt"
)

func NewFrontEnd(properties string) (FrontEnd, error) {
	switch properties {
		
	case "rest":
		return &restFrontEnd{}, nil

	default:
		return nil, fmt.Errorf("no such frontend %s", properties)
	}
}