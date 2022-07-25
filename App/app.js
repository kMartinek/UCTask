/*
TO DO:

- get FilePath via command line parameter
    - read from cmd
    - if no args then print error msg || interface
    - check if file path is correct (file exists)
- connect to service
- send file for processing
- recieve and print info

*/

import {connect, StringCodec} from "nats"
import fs from "fs"
//import process from 'minimist'

//Read arguments


const args = process.argv.slice(2);
console.log (args);
if(args == undefined)
{
    console.log ("You haven't entered path argument!");
    process.exit();
}

//Check for file 

function checkFileExists(path)
{
    console.log("Path: ", path)
    let flag = true;
    try
    {
        fs.accessSync(path);
    }
    catch(e)
    {
        flag = false;
    }
    console.log("Flag", flag)
    return flag;
}

if(checkFileExists(args[0]) == false)
{
    console.log("The file you are looking for does not exist!")
    //process.exit();
}
else 
{
    console.log("File exists")
}


// Connect to server
const server = {servers : "localhost:4222"}
/*try
{
    const nc = await connect(server)
    console.log(`connected to ${nc.getServer()}`);

    const done = nc.closed();

    const err = await done;
    if (err){
        console.log(`error closing:`, err);
    }
}
catch(err)
{
    console.log(`error connecting to ${JSON.stringify(server)}`)
}
*/
const nc = await connect(server)

const sc = StringCodec();

// the client makes a request and receives a promise for a message
// by default the request times out after 1s (1000 millis) and has
// no payload.
await nc.request("mp4InitSegment", new Uint8Array(args[0]), { timeout: 30000 })
    .then((m) => {
    console.log(`got response: ${sc.decode(m.data)}`);
    })
    .catch((err) => {
    console.log(`problem with request: ${err.message}`);
    });

await nc.request()

await nc.close();


