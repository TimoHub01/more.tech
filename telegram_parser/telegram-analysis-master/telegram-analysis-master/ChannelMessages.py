import configparser
import json
import asyncio
from datetime import date, datetime
import codecs
import pandas as pd
import numpy as np
from telethon import TelegramClient
from telethon.errors import SessionPasswordNeededError
from telethon.tl.functions.messages import (GetHistoryRequest, GetMessagesRequest)
from telethon.tl.types import (
    PeerChannel
)
from csv_to_sql import csv_to_sql
from csv_to_sql import get_data_for_model


# some functions to parse json date
class DateTimeEncoder(json.JSONEncoder):
    def default(self, o):
        if isinstance(o, datetime):
            return o.isoformat()

        if isinstance(o, bytes):
            return list(o)

        return json.JSONEncoder.default(self, o)

def get_title(row):
    message = row['message']

    if type(row['entities']) != str:
        return row

    entities = row['entities'] \
        .replace('[', '') \
        .replace(']', '')
    if entities:
        markup = json.loads(entities \
                            .replace('}, {', '}|,| {') \
                            .split('|,|')[0].replace("'", '"')
                            )

        offset = markup.get('offset')
        length = markup.get('length')

        row['title'] = message[offset:offset + length]
    return row

config = configparser.ConfigParser()
config.read("config-sample.ini")

# Setting configuration values
api_id = config['Telegram']['api_id']
api_hash = config['Telegram']['api_hash']

api_hash = str(api_hash)

phone = config['Telegram']['phone']
username = config['Telegram']['username']

# Create the client and connect
client = TelegramClient(username, api_id, api_hash)


async def main(phone, url):
    # Reading Configs
    await client.start()
    print("Client Created")
    # Ensure you're authorized
    if await client.is_user_authorized() == False:
        await client.send_code_request(phone)
        try:
            await client.sign_in(phone, input('Enter the code: '))
        except SessionPasswordNeededError:
            await client.sign_in(password=input('Password: '))
    if(url==""):
        user_input_channel = input('enter entity(telegram URL or entity id):')
    else:
        user_input_channel = url

    if user_input_channel.isdigit():
        entity = PeerChannel(int(user_input_channel))
    else:
        entity = user_input_channel

    my_channel = await client.get_entity(entity)

    offset_id = 0
    limit = 100
    all_messages = []
    total_messages = 0
    total_count_limit = 0

    while True:
        print("Current Offset ID is:", offset_id, "; Total Messages:", total_messages)
        history = await client(GetHistoryRequest(
            peer=my_channel,
            offset_id=offset_id,
            offset_date=None,
            add_offset=0,
            limit=limit,
            max_id=0,
            min_id=0,
            hash=0
        ))
        if not history.messages:
            break
        messages = history.messages
        index = 0

        for message in messages:
            try:

                all_messages.append(message.to_dict())

            except Exception as e:
                print(e)

        offset_id = messages[len(messages) - 1].id
        total_messages = len(all_messages)
        if total_count_limit != 0 and total_messages >= total_count_limit:
            break
    with open('channel_messages.json', 'w') as outfile:
        json.dump(all_messages, outfile, cls=DateTimeEncoder)

# with client:
#     client.loop.run_until_complete(main(phone, url=""))

    df = pd.read_json('channel_messages.json')
    print(df.loc[0, :])

    df = df.loc[:, ['date', 'message', 'entities', 'views', 'forwards']]
    df['title'] = np.nan

    df = df.apply(get_title, axis=1).drop(['entities'], axis=1)
    print(df.loc[0, ['title']])
    # df.reactions = df.reactions.apply(get_reactions_cnt)
    df.to_csv('TEST.csv')
    csv_to_sql()





