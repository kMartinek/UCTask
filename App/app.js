import { connect, StringCodec } from "nats"
import fs from "fs"
import minimist from 'minimist'

const connServer = "localhost:4222";


//Read path argument
const args = minimist(process.argv.slice(2))
const filePath = args.path;

if (filePath == undefined) {
    console.log("You haven't entered --path argument!\n");
    process.exit();
}

//Check if file exists 
function checkFileExists(path) {
    let flag = true;
    try {
        fs.accessSync(path);
    }
    catch (e) {
        flag = false;
    }
    return flag;
}


if (!checkFileExists(filePath)) {
    console.log("The file you are looking for does not exist!")
    process.exit();
}
else {
    console.log("File exists")
}

function checkMP4Extension(filePath) {
    let splitArr = filePath.split('.');
    if (splitArr[splitArr.length - 1] == 'mp4')
        return true;
    else
        return false;
}

if (!checkMP4Extension(filePath)) {
    console.log("Given file does not have .mp4 extension");
    process.exit();
}


// Connect to server
const server = { servers: connServer }
try {
    const nc = await connect(server);
    console.log("Succesfully connected to:", connServer)
    const sc = StringCodec();

    //Send request - await response
    let enc = new TextEncoder();
    const response = await nc.request("mp4InitSegment", enc.encode(filePath), { timeout: 30000 })
    console.log(`Response: ${sc.decode(response.data)}`);
    await nc.close();
}
catch (exception) {
    console.log("Connection to server failed!");
}

