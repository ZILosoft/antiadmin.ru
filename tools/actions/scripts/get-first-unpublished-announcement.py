# set timezone to moscow

import os
import sys
import datetime
import time
import json

os.environ['TZ'] = 'Europe/Moscow'
time.tzset()

# get current date
now = datetime.datetime.now()

with open('../db/announcements.json', 'r') as f:
    announcements = json.load(f)
    for announcement in announcements:
        if announcement['published'] == False:
            if now < datetime.datetime.strptime(announcement['publishAfter'], '%d.%m.%Y %H:%M'):
                print(announcement['message'])
                announcement['published'] = True
                with open('../db/announcements.json', 'w') as f:
                    json.dump(announcements, f, indent=2)

                sys.exit(0)


