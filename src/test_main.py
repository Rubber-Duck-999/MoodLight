'''
Testing logic for Mood Light
'''

# pylint: disable=import-error

import unittest
import main

class TestColours(unittest.TestCase):
    '''Testing colours class'''
    def get_colours(self) -> None:
        '''Ensure colours are set'''
        colours = main.Colours()
        self.assertEqual(colours.blue, 50)
