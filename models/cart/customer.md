# Customer

A customer resource instance represents a customer account with the shop. Customer accounts store contact information for the customer, saving logged-in customers the trouble of having to provide it at every checkout. For security reasons, the customer resoruce instance does not store credit card information. Customers will always have to provide this information at checkout.

![customer diagram](customer.png)
The customer resource instance also stores some additional information about the customer for the benefit of the shop owner, including: the number of orders, the amount of money s/he has spent and number of orders s/he has made throughout his/her history with the shop as well as the shop owner's notes and tags for the customer.

The shop's use of customer accounts will depend on its Customer Checkout settings. For shop owners, this is located in the shop admin dashboard in the "Preferences" tab, on the "Checkout and Payment" page. There are three options for this setting:

  * __Guest checkout only:__ Customer accounts are disabled. meaning that customers can't log in and can only check out as guests.
  * __Guest checkout with optional sign-in:__ Customer accounts are optional. Customers have the choice of either signing into their account or simply checking out as a guest. Under this setting, customers can create accounts for themselves; the shop owner can also an account for a customer and then invite him/her by email to use it.
  * __Sign-in required:__ Customer accounts are required. Customers can't check out unless logged in, and under this setting, the shop owner must create accounts for them first.

#### What can you do with Customer?

The API lets you do the following with the Customer resource. More detailed versions of these general actions may be available:

_GET /cart/customers_
Receive a list of all Customers

_GET /cart/customers/search?query=Bob_
Search for customers matching supplied query

_GET /cart/customers/#{id}.json_
Receive a single Customer

_POST /cart/customers.json_
Create a new Customer

_PUT /cart/customers/#{id}.json_
Modify an existing Customer

_DELETE /cart/customers/#{id}.json_
Remove a Customer from the database

_GET /cart/customers/count.json_
Receive a count of all Customers

_GET /cart/orders.json?customer_id=207119551_
Find orders belonging to this customer
