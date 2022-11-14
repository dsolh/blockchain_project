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

// blockchain connection
async function ccCall(func, args){
    try {
        // load the network configuration
        const ccpPath = path.resolve(__dirname, "ccp", "connection-org1.json");
        let ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));
        
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), "wallet");
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get("appUser");
        if (!identity) {
            console.log(
                'An identity for the user "appUser" does not exist in the wallet'
            );
            console.log("Run the registerUser.js application before retrying");
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, {
            wallet,
            identity: "appUser",
            discovery: { enabled: true, asLocalhost: true },
        });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork("mychannel");

        // Get the contract from the network.
        const contract = network.getContract("whozardCC");

        // Submit the specified transaction.
        var result = "";
        if (func == "CreateSpell") {
            result = await contract.submitTransaction("CreateSpell", args[0], args[1]);
        } else if (func == "ReadSpell") {
            result = await contract.evaluateTransaction("ReadSpell", args[0]);
        } else if (func == "WizardExists"){
            result = await contract.evaluateTransaction("WizardExists", args[0]);
        } else if (func == "RegisterWizard"){
            result = await contract.submitTransaction("RegisterWizard", args[0], args[1]);
        } else if (func == "CastSpell") {
            result = await contract.submitTransaction("CastSpell", args[0], args[1], args[2]);
        } else if (func == "CheckIdentity") {
            result = await contract.evaluateTransaction("CheckIdentity", args[0], args[1]);
        } else if (func == "RecentSpell") {
            result = await contract.evaluateTransaction("RecentSpell", args[0], args[1]);
        } else {
            result = "Not supported fuction"
        }

        // Disconnect from the gateway.
        await gateway.disconnect();

        return result;
    } catch(error) {
        throw new Error(error);
    }
    
}

// login 
app.post("/login", async (req, res) => {
    const wizardName = req.body.wizardName;
    const country = req.body.country;
    const password = req.body.password;

    // check password
    var rst = true;
    switch (country) {
        case "British" :
            if (password != "British is the best") {
                rst = false;
            }
            break;
        case "USA" :
            if (password != "USA is the best") {
                rst = false;
            }
            break;
        case "France" :
            if (password != "France is the best") {
                rst = false;
            }
            break;
        case "German" :
            if (password != "German is the best") {
                rst = false;
            }
            break;
        case "China" :
            if (password != "China is the best") {
                rst = false;
            }
            break;
    }
    if (!rst) {
        res.status(200).json({ result: `wrong password`});
        console.log("wrong password");
        return;
    }

    // check existance of wizard
    try {
        result = await ccCall("WizardExists", [wizardName]);
        result = result.toString();
        
        if (result == "true") {
            // when the wizard is registered
            // If you want to run it by multi channel, then modify here.
            res.status(200).json({ result: "login success"});
        } else {
            // when the wizard is not registered
            res.status(200).json({ result: "not registered"});
        }
    } catch (error) {
        res.status(400).json({ result: "Login has been failed"});
    }
})

// register wizard
app.post("/register_wizard", async (req, res) => {
    const wizardName = req.body.wizardName;
    const country = req.body.country;
    const args = [wizardName, country];

    // register
    try {
        result = await ccCall("RegisterWizard", args);

        res.status(200).json({ result: "You're a wizard now!"});
        console.log("registered wizard:"+wizardName+", "+country);
    } catch (error) {
        res.status(200).json({ result: "Register process has been failed"});
    }
})

// cast spell
app.post("/cast_spell", async (req, res) => {
    const wizardName = req.body.wizardName;
    const spell = req.body.spell;
    // instead of setting multichennel
    const country = req.body.country;
    const args = [wizardName, spell, country];

    try {
        result = await ccCall("CastSpell", args);

        res.status(200).json({ result: "The magic worked well"});
        console.log("result of CastSpell:"+wizardName+", "+spell);
    } catch (error) {
        var msg = error.message;
        if (msg.includes("not registered")) {
            res.status(200).json({ result: `You are not registered`});
        } else if (msg.includes("not existed")) {
            res.status(200).json({ result: `Spell \'${spell}\' is not existed`});
        } else if (msg.includes("country")) {
            // instead of setting multichennel
            res.status(200).json({ result: `Country isn't correct`});
        } else {
            res.status(200).json({ result: "Casting spell has been failed"});
        }
        console.log("Transaction has been failed");
    }
})

// create spell
app.post("/create_spell", async (req, res) => {
    const spell = req.body.spell
    const category = req.body.category

    var args = [spell, category]

    try {
        //register
        result = await ccCall("CreateSpell", args);

        res.status(200).json({ result: `${spell} is created` });
        console.log("Transaction has been submitted");
    } catch (error) {
        var msg = error.message;
        if (msg.includes("already existed")) {
            res.status(200).json({ result: `spell \'${spell}\' is already existed`});
        } else {
            res.status(400).json({ result: "Transaction has been failed"});
        }
        console.log("Transaction has been failed");
    }
})

// check identity
app.get("/check_identity", async (req, res) => {
    const wizardName = req.query.wizardName;
    // instead of setting multichennel
    const country = req.query.country;
    const args = [wizardName, country];

    try {
        result = await ccCall("CheckIdentity", args);
        result = result.toString();

        res.status(200).json({ result: `${wizardName}'s identity is \'${result}\'`});
        console.log("result of Check Identity:"+wizardName+" ,"+result);
    } catch (error) {
        var msg = error.message;
        if (msg.includes("not registered")) {
            res.status(200).json({ result: `${wizardName} is not registered`});
        } else if (msg.includes("country")) {
            // instead of setting multichennel
            res.status(200).json({ result: `Country isn't correct`});
        } else {
            res.status(200).json({ result: "Checking identity has been failed"});
        }
        console.log("Transaction has been failed");
    }
})

// check recent spell
app.get("/recent_spell", async (req, res) => {
    const wizardName = req.query.wizardName;
    // instead of setting multichennel
    const country = req.query.country;
    const args = [wizardName, country];

    try {
        result = await ccCall("RecentSpell", args);
        result = result.toString();

        res.status(200).json({ result: `${wizardName} has casted \'${result}\' recently`});
        console.log("result of Recent Spell:"+wizardName+" ,"+result);
    } catch (error) {
        var msg = error.message;
        if (msg.includes("not registered")) {
            res.status(200).json({ result: `${wizardName} is not registered`});
        } else if (msg.includes("country")) {
            // instead of setting multichennel
            res.status(200).json({ result: `Country isn't correct`});
        } else {
            res.status(200).json({ result: "Checking recent spell has been failed"});
        }
        console.log("Transaction has been failed");
    }
})

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
