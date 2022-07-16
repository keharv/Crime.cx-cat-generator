import requests


# https://api.imgur.com/post/v1/posts/t/dailykittens?client_id=546c25a59c58ad7&filter[window]=week&include=adtiles%2Cadconfig%2Ccover&page=6&sort=-viral


# go through each page of the daily kitten gallery
for page in range(1, 40000):
    # get the json data for the page
    url = "https://api.imgur.com/post/v1/posts/t/dailykittens?client_id=546c25a59c58ad7&filter[window]=week&include=adtiles%2Cadconfig%2Ccover&page=" + str(page) + "&sort=-viral"
    r = requests.get(url)
    data = r.json()
    # go through each image in the page
    for image in data['posts']:
        # get the image url
        url = image['cover']['url']
        # download the image
        r = requests.get(url)
        # save the image
        with open("./images/" + image['id'] + ".jpg", "wb") as f:
            f.write(r.content)
        print("Saved " + image['id'] + ".jpg")
    print("Page " + str(page) + " done")