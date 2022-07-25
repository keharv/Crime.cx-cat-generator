
#create a web crawler that uses selenium to scrape the data from the website

import requests
from bs4 import BeautifulSoup
from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException
import time
import hashlib
import psycopg2
import os
import random
import magic


#import env variables
from dotenv import load_dotenv




def generateRandomString(length):
    letters = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
    return ''.join(random.choice(letters) for i in range(length))

#create a function that will scrape the data from the website
def scrape_data(url):
    chrome_options = webdriver.ChromeOptions()
    chrome_path = './chromedriver'
    #create a webdriver linux chrome driver located at ./chromedriver
    driver = webdriver.Chrome(chrome_path, chrome_options=chrome_options)

    #navigate to the url
    driver.get(url)
    #wait for the page to load
    time.sleep(5)
    #create a beautiful soup object
    html = driver.page_source
    print("Scraping data from: " + url)
    soup = BeautifulSoup(html, 'html.parser')
    #close the driver
    driver.close()

    #return the data
    return soup

def save_images(data):
    #get all img's from html contained in data and save all the images to the scraped folder
    
    images = data.find_all('img')
    for image in images:
        try:
            #ensure valid image link
            if image['src'] is None:
                continue
            #get the image link
            image_link = image['src']
            #check if image link is valid
            if image_link.startswith('//'):
                image_link = 'http:' + image_link

            if not image_link.startswith('http'):
                image_link = 'http://' + image_link
            src = image_link
            #get the image data
            image_data = requests.get(src).content
            image_type = magic.from_buffer(image_data, mime=True)
            image_extension = image_type.split('/')[1]
            if '+' in image_extension:
                image_extension = image_extension.split('+')[0]
            image_name = generateRandomString(10) + '.' + image_extension
            
            #generate md5 of the image
            
            md5 = hashlib.md5(image_data).hexdigest()

            #check if image already exists in database
            cur.execute("SELECT * FROM images WHERE md5 = %s", (md5,))
            if cur.fetchone() is not None:
                continue
            #save the image
            with open('scraped/' + image_name, 'wb') as handler:
                handler.write(image_data)
            #insert the image into the database
            cur.execute("INSERT INTO images (name, md5) VALUES (%s, %s)", (image_name, md5))
            conn.commit()
        except Exception as e:
            print(image)
            print(e)
            continue
def insert_new_links(data):
    #get all links from the data
    links = data.find_all('a')
    for link in links:
        try:
            href = link['href']
            if not href.startswith('http'):
                href = 'http://' + href
            #check if link already exists in database
            cur.execute("SELECT * FROM links WHERE link = %s", (href,))
            if cur.fetchone() is not None:
                continue
            #insert the link into the database
            cur.execute("INSERT INTO links (link, scraped) VALUES (%s, FALSE)", (href,))
        except Exception as e:
            print(e)
            print(link)
            continue
            

def do_crawling():
    #get oldest link in database
    cur.execute("SELECT link FROM links WHERE SCRAPED <> TRUE ORDER BY id ASC LIMIT 1")
    oldest_link = cur.fetchone()[0]
    #get the url of the oldest link
    url = oldest_link
    #scrape the data
    data = scrape_data(url)
    #save the images
    save_images(data)
    #insert the new links
    insert_new_links(data)
    #update the link in the database
    cur.execute("UPDATE links SET scraped = TRUE WHERE link = %s", (url,))
    #commit the changes
    conn.commit()

if __name__ == '__main__':
    load_dotenv()

    database_url = os.getenv('DATABASE_URL')
    database_user = os.getenv('DATABASE_USER')
    database_password = os.getenv('DATABASE_PASSWORD')

    #connect to the database
    conn = psycopg2.connect(database_url, user=database_user, password=database_password)

    #create a cursor
    cur = conn.cursor()

    #create tables
    cur.execute("CREATE TABLE IF NOT EXISTS images (id SERIAL PRIMARY KEY, name VARCHAR(255), md5 VARCHAR(255))")
    cur.execute("CREATE TABLE IF NOT EXISTS links (id SERIAL PRIMARY KEY, link VARCHAR(512), scraped BOOLEAN)")
    
    #commit the changes
    conn.commit()
    #get number of links in database
    cur.execute("SELECT COUNT(*) FROM links")
    num_links = cur.fetchone()[0]
    #if there are no links in the database, do crawling
    if num_links == 0:
        #insert imgur home into database
        cur.execute("INSERT INTO links (link, scraped) VALUES ('https://imgur.com/', FALSE)")
        #commit the changes
        conn.commit()
    do_crawling()