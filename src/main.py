''' Mood Light program
'''
#!/usr/bin/env python3
# pylint: disable=consider-using-f-string import-error logging-format-interpolation bare-except

import time
import os
import logging
import random
import json
import requests
try:
    from unicornhatmini import UnicornHATMini
except:
    from mock_hat import UnicornHATMini
import utilities


fileName: str = '/home/{}/sync/MoodLight.log'

try:
    name: str = utilities.get_user()
    fileName: str = fileName.format(name)
    os.remove(fileName)
except OSError as error:
    pass

# Add the log message handler to the logger
logging.basicConfig(filename=fileName,
                    format='%(asctime)s - %(levelname)s - %(message)s',
                    level=logging.INFO)

unicornhatmini = UnicornHATMini()
unicornhatmini.set_brightness(0.1)

class Colours:
    '''Class for managing the onboard hat leds'''
    def __init__(self) -> None:
        self.red = 100
        self.green = 75
        self.blue = 50
        self.server = ''
        self.state = 'ON'
        self.get_env()

    def get_env(self) -> None:
        '''
        Acquires variables for the object from the main config file
        '''
        logging.info('get_env()')
        config_name: str = '/home/{}/sync/config.json'.format(utilities.get_user())
        try:
            with open(config_name, encoding="utf-8") as file:
                data = json.load(file)
            self.server = data["server_address"] + '/mood'
        except KeyError:
            logging.error("Variables not set")
        except IOError:
            logging.error('Could not read file')

    def get_colours(self) -> None:
        '''
        Sets RGB colours for the leds
        '''
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

    def get_data(self) -> None:
        '''
        Get mood status
        '''
        logging.info('get_data()')
        try:
            response = requests.get(self.server, timeout=5)
            if response.status_code == 200:
                logging.info("Requests successful")
                if 'state' in response.json():
                    logging.info(response.json())
                    self.state = response.json()['state']
                logging.info('Mood Light is set to: {}'.format(self.state))
            else:
                logging.error('Response: {}'.format(response))
        except requests.ConnectionError as connection_error:
            logging.error("Connection error: {}".format(connection_error))
        except requests.Timeout as timeout_error:
            logging.error("Timeout on server: {}".format(timeout_error))
        except KeyError as key_error:
            logging.error("Key error on data: {}".format(key_error))

    def start(self) -> None:
        '''
        Swicthes between ON and OFF depending on mood setting
        '''
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
