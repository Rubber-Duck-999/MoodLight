''' Mock unicornhat library to stop issues with
testing the program
'''

class UnicornHATMini:
    '''Mock class'''
    def set_brightness(self, value: int) -> None:
        '''Set brightness value'''
        print(value)

    def set_all(self, red: int, green: int, blue: int) -> None:
        '''Set all leds'''
        print(red, green, blue)

    def show(self) -> None:
        '''Show colours'''
