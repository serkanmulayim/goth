#"KTIDRUSEAP4BYFYL" // binance Google Auth recovery key

from binance.client import Client
import json
import os
import errno
import sys 
import time
import pandas as pd
import numpy as np

api_key="Oos1jUZpXSEmIxxTF5V0rFMmMRpSsopzWhCjobiLPkyzpJ36WDHYwlwzGPWfEskJ";
api_secret="kE5ORsdAyyA6wNdm85fiHPwhCt51iLsQWXkxvqR71lI6maf9AFEB4H3C7rJZ8uuG";
DATA_FOLDER = "../data/";

symbol = "BTCUSDT"
interval = Client.KLINE_INTERVAL_1DAY

if not os.path.exists(os.path.dirname(DATA_FOLDER)):
    try:
        os.makedirs(os.path.dirname(DATA_FOLDER))
    except OSError as exc: # Guard against race condition
        if exc.errno != errno.EEXIST:
                    raise

start =  "August 01, 2017"
end = "February 01, 2018"

client = Client(api_key, api_secret);
res = client.get_historical_klines(symbol, interval, start, end);


np_train = np.array(res)

floats = np_train.astype(np.float)

index = np.transpose(floats)[0].astype(np.int)
data = np.delete(floats,0, axis=1) #delete opentime col
data = np.delete(data,5, axis=1) #delete closetime col
data = np.delete(data,9, axis=1) #can be ingnored
#qav = quote asset volume
#tbbav = taker buy base asset volume
#tbqav = taker buy quote asset volume
cols = ['open', 'high', 'low', 'close', 'volume', 'qav', 'numtrades','tbbav', 'tbqav']
df = pd.DataFrame(data, index=index, columns=cols)
df.to_csv(DATA_FOLDER+'data.csv')


