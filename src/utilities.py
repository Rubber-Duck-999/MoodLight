'''Utilities File
- Common functions used
'''
import getpass

def get_user() -> str:
    '''Returns the username'''
    try:
        username = getpass.getuser()
    except OSError:
        username = 'pi'
    return username
