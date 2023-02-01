#!/usr/bin/env python3

import time
import os
import requests
try:
    from unicornhatmini import UnicornHATMini
except ImportError:
    from mock_hat import UnicornHATMini
import utilities
import logging
from datetime import datetime
import random
import json


filename = '/home/{}/sync/MoodLight.log'

try:
    name = utilities.get_user()
    filename = filename.format(name)
    os.remove(filename)
except OSError as error:
    pass

# Add the log message handler to the logger
logging.basicConfig(filename=filename,
                    format='%(asctime)s - %(levelname)s - %(message)s', 
                    level=logging.INFO)

unicornhatmini = UnicornHATMini()
unicornhatmini.set_brightness(0.1)

class Colours:
    def __init__(self):
        self.red = 100
        self.green = 75
        self.blue = 50
        self.server = ''
        self.state = 'ON'
        self.get_env()

    def get_env(self):
        logging.info('get_env()')
        config_name = '/home/{}/sync/config.json'.format(utilities.get_user())
        try:
            with open(config_name) as file:
                data = json.load(file)
            self.server = data["server_address"] + '/mood'
        except KeyError:
            logging.error("Variables not set")
        except IOError:
            logging.error('Could not read file')

    def get_colours(self):
        colour = random.randint(1,3)
        value = random.randint(0,255)
        if colour == 1:
            self.red = value
        elif colour == 2:
            self.green = value
        elif colour == 3:
            self.blue = value
        else:
            logging.info('How did this happen?')

    def get_data(self):
        '''Get mood status'''
        logging.info('get_data()')
        try:
            response = requests.get(self.server, timeout=5)
            if response.status_code == 200:
                logging.info("Requests successful")
                if 'mood' in response.json():
                    logging.info(response.json())
                    mood = response.json()['mood'][0]
                    logging.info(mood)
                    self.state = mood['state']
                logging.info('Mood Light is set to: {}'.format(self.state))
            else:
                logging.error('Response: {}'.format(response))
        except requests.ConnectionError as error:
            logging.error("Connection error: {}".format(error))
        except requests.Timeout as error:
            logging.error("Timeout on server: {}".format(error))
        except KeyError as error:
            logging.error("Key error on data: {}".format(error))

    def start(self):
        logging.info('start()')
        if self.state == 'ON':
            self.get_colours()
            unicornhatmini.set_all(self.red, self.green, self.blue)
            unicornhatmini.show()
            logging.info('Displaying colours: {} {} {}'.format(self.red, self.green, self.blue))
            time.sleep(30)
        else:
            logging.info('Night time')
            unicornhatmini.set_all(0, 0, 0)
            unicornhatmini.show()
            time.sleep(30)
        self.get_data()

if __name__ == "__main__":
    logging.info('Starting Program')
    colours = Colours()
    while True:
        colours.start()