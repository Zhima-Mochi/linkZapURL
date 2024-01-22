#!/bin/bash

mongosh << EOF
use admin;
db.createUser({ user: "root", pwd: "password", roles: [{ role: "root", db: "admin" }] });
exit;
EOF