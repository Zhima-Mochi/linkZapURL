// Enable sharding for database `linkZapURL`
sh.enableSharding("linkZapURL")

// Setup shardingKey for collection `linkZapURL.url`
db.adminCommand({ shardCollection: "linkZapURL.url", key: { seq: 1 } })

db = db.getSiblingDB("linkZapURL");

db.url.createIndex({ seq: 1, ID: 1 }, { unique: true });


