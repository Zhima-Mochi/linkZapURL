db = db.getSiblingDB("linkZapURL");

var collections = db.getCollectionNames();
if (!collections.includes("url")) {
    db.createCollection("url");
    print("Collection 'url' created.");
} else {
    print("Collection 'url' already exists.");
}

var indexes = db.url.getIndexes();
var idIndexExists = indexes.some(function(index) {
    return index.key.hasOwnProperty("ID") && index.unique;
});

if (!idIndexExists) {
    db.url.createIndex({ "ID": 1 }, { unique: true });
    print("Unique index on 'ID' created.");
} else {
    print("Unique index on 'ID' already exists.");
}
