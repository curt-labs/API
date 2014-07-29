# API Design

This document will lay out of the features that need to go into v3 of the CURT
 API.

### Services

- [User](#user-service)
- [Customer](#customer-service)
- [Vehicles](#vehicles-service)
- [Product](#product-service)
- [Category](#category-service)
- more to come...

***

#### <a href="user-service"></a>User Service

Manage all customer user manipulation. This will include things like:

- sign up
- password reset
- get
- update
- delete
- etc.

***

#### <a href="customer-service"></a>Customer Service

Manage all customer manipulation. Some of the features this will
include:

- pricing
- content
- cart integration
- shopping cart integration
- add/update/delete locations
- add/update/delete web properties

***

#### <a href="vehicles-service"></a>Vehicles Service

Query ACES vehicle configurations. Just a series of aggregation depending on
currently provided attributes. Tricky part here is going to be flexibility and
performance.

***

#### <a href="product-service"></a>Product Service

A series of data getters for a single product object, a series of product objects,
or a segment of a product object (prices, attributes, categories).

This will also include a certain level of customer data injection, solely based off
the provided API key.

I view the hard part of this design being the filtration functionality, we'll have to
decide on whether the best place for filtration should be. (Should it be on an array of
product objects or should this reside in the category service?)


> David's Part Filter Concept:

> A function that takes in a list of PartID's, and Filteration Rules. This function would return a list of parts that meet the filteration conditions. The beauty of having a function that takes a list of parts, and returns a list of parts, is that it is very modular and extensible. For example, it could be used as a standalone (I have these parts and I want to apply filter rules to them to narrow down my results). 

> It could also be used in conjunction with the vehicle lookup. Say you wanted to do the vehicle lookup and filter rules at the same time via the same request. The vehicle lookup would accept the list of filter rules. Once the vehicle lookup grabs all the parts, it would then call this filter function passing in the list of PartIDs and fowarding the filteration rules. It would then return the list of parts, and use that list to return to the user. This also allows use to keep all the messy filteration code out of the vehicle lookup, and out of all the individual GetParts calls. Below is a quick diagram I made to show the flow:

> ![alt text](http://i.imgur.com/eYwjUHt.png "Really quick and bad flowchart")


***

### <a href="category-service"></a>Category Service

A series of data retrieval endpoints, also with a certain level of customer data
injection. This will need to maintain more of a top-level category structure, we don't
want to end up in the same nested tree boat that we ended up with in v2.

Filtration may end up in this section, which will significantly increase the complexity
of the service.

***
***

### MiddleWare

- [Authentication](#authentication)
- [Analytics](#analytics)

***

#### <a href="authentication"></a>Authentication

This version of the API will require every request to be signed by an API key. Furthermore,
any request that injects/modifies data, will need to be signed by a private key, thus
enabling us to cut down on any bad data or security intrusions.

### <a href="analytics"></a>Analytics

Analytics will need to be treated as a priority in v3 of the API. We need to be able to
collect data on what our streamed data size is, the paramters of the request, the time spent
processing, and the final response code. Liking the approach of integrating with segment.io
for aggregating these metrics to different services.
