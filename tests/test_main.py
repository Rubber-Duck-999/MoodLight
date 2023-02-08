'''
Testing logic for Mood Light
'''

import unittest
from src.main import Colours

class TestColours(unittest.TestCase):
    '''Testing colours class'''
    def get_colours(self) -> None:
        '''Ensure colours are set'''
        colours = Colours()
        self.assertEqual(colours.blue, 50)
