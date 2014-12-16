package cart_ctlr

import (
	"encoding/json"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/models/cart"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	InventoryBehaviorBypass                  = "bypass"                    // default - do not claim inventory
	InventoryBehaviorDecrementIgnoringPolicy = "decrement_ignoring_policy" // ignore the product's inventory policy and claim amounts no matter what.
	InventoryBehaviorDecrementObeyingPolicy  = "decrement_obeying_policy"  // Obey the product's inventory policy
)

func CreateOrder(w http.ResponseWriter, req *http.Request, params martini.Params, enc encoding.Encoder, shop *cart.Shop) string {

	var o cart.Order
	qs := req.URL.Query()
	defer req.Body.Close()

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	if err := json.Unmarshal(data, &o); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	o.ShopId = shop.Id

	if hooks := qs.Get("send_webhooks"); hooks != "" {
		if wb_hooks, err := strconv.ParseBool(hooks); err == nil {
			o.SendWebhooks = wb_hooks
		}
	}
	if receipt := qs.Get("send_receipt"); receipt != "" {
		if wb_receipt, err := strconv.ParseBool(receipt); err == nil {
			o.SendReceipt = wb_receipt
		}
	}
	if fulfillment := qs.Get("send_fulfillment_receipt"); fulfillment != "" {
		if wb_fulfillment, err := strconv.ParseBool(fulfillment); err == nil {
			o.SendFulfillmentReceipt = wb_fulfillment
		}
	}
	if behav := qs.Get("inventory_behavior"); behav != "" {
		switch behav {
		case InventoryBehaviorDecrementIgnoringPolicy:
			o.InventoryBehavior = InventoryBehaviorDecrementIgnoringPolicy
		case InventoryBehaviorDecrementObeyingPolicy:
			o.InventoryBehavior = InventoryBehaviorDecrementObeyingPolicy
		default:
			o.InventoryBehavior = InventoryBehaviorBypass
		}
	}

	if err := o.Create(); err != nil {
		apierror.GenerateError(err.Error(), err, w, req)
		return ""
	}

	return encoding.Must(enc.Encode(o))
}
