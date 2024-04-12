package requests

import (
	"fmt"
)

// Restaurant Request
type RestaurantRequest struct {
	Userid uint   `json:"userid"`
	Name   string `json:"name"`
	Format string `json:"format"`
}

func (r *RestaurantRequest) String() string {
	return fmt.Sprint(r.Userid, r.Name, r.Format)
}
