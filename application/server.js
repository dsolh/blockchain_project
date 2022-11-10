// ExpressJS Setup
const express = require("express");
const app = express();
var bodyParser = require("body-parser");

// Hyperledger Bridge
const { Wallets, Gateway } = require("fabric-network");
const fs = require("fs");
const path = require("path");
const ccpPath = path.resolve(__dirname, "ccp", "connection-org1.json");
let ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));

// Constants
const PORT = 8080;
const HOST = "0.0.0.0";

// use static file
app.use(express.static(path.join(__dirname, "views")));

// configure app to use body-parser
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

// main page routing
app.get("/", (req, res) => {
    res.sendFile(__dirname + "/index.html");
});

// login 
app.post("/login", (req, res) => {
    const wizardName = req.body.name;
    const country = req.body.country;
    const password = req.body.password;

    console.log("request:"+wizardName+", "+country+", "+password);
})

app.get("/readSpell", (req, res) => {
    const wizardName = req.query.name;

    console.log("request:"+wizardName);
})

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
