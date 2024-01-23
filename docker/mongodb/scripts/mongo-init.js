// Enable sharding for database `linkZapURL`
sh.enableSharding("linkZapURL")

// Setup shardingKey for collection `linkZapURL.url`
db.adminCommand({ shardCollection: "linkZapURL.url", key: { shardID: 1 } })

db = db.getSiblingDB("linkZapURL");

db.url.createIndex({ shardID: 1, ID: 1 }, { unique: true });


