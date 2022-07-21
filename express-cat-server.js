//create web server
var express = require('express');
var app = express();
var fs = require('fs');
var server = require('http').createServer(app);
var https = require('https');

//setup https server
var options = {
    key: fs.readFileSync('key.pem'),
    cert: fs.readFileSync('cert.pem')
};

var httpsServer = https.createServer(options, app);


//set port
server.listen(80);
httpsServer.listen(443);



//serve random image from folder on /images
app.use('/images', express.static(__dirname + '/images'));


app.get('/', function(req, res){
    //check if using https
    if (req.protocol !== 'https') {
        res.redirect('https://' + req.headers.host + req.url);
    }
    //get list of random images in images/ folder
    const files = fs.readdirSync('./images/');
    //get random image
    const randomImage = files[Math.floor(Math.random() * files.length)];
    //set source header
    res.setHeader('Image-Source', 'images/' + randomImage);
    res.setHeader('X-Powered-By', 'KittyPower');
    //send image
    res.sendFile(__dirname + '/images/' + randomImage);
    console.log('Request received from ' + req.headers.host + ' for ' + randomImage);
}
);
