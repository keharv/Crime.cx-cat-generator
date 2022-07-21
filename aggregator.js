var Crawler = require("simplecrawler");
var fs = require('fs');

var crawler = new Crawler("https://imgur.com/gallery/BYp0Naz/");
crawler.respectRobotsTxt = false;
crawler.allowInitialDomainChange = true;
crawler.decodeResponses = true;
crawler.maxConcurrency = 25;
crawler.on("fetchcomplete", function(queueItem, responseBuffer, response) {
    console.log("I just received %s (%d bytes)", queueItem.url, responseBuffer.length);
    console.log("It was a resource of type %s", response.headers['content-type']);
    if(response.headers['content-type'].indexOf('image') > -1) {
        //save to from-crawler folder
        fs.writeFile('./from-crawler/' + queueItem.url.split('/').pop(), responseBuffer, function(err) {
            if(err) {
                console.log(err);
            }
        }
        );
}
});


crawler.start();