# smurfin-checkout
The checkout service
exposes a single endpoint, POST /checkout, which takes account info and begins the checkout process.
It confirms item information using the catalog service and then fires off a validate-payment event
iT then sends response to client
Once payment is validated, it then sends payment success event which fires off an email service and updates to the cart and catalog service.
work in progress
