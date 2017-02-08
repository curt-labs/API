#MongoDB Indexing

Due to the default memory constraints in MongoDB, specific API calls may fail if too much data results from a query at once (specifically `/part` in this case).
The code to add the index into MongoDB is as such:

`db.products.createIndex({"brand.id": 1, id: 1}, {background: true})`

This command must be run in the `product_data` DB.
